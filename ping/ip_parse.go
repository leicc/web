package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	ipsegment_file = "./ip.txt"
	ipdetail_file  = "./ipdetail.txt"
	ipping_file    = "./ipping.txt"
	ipserver_file  = "./ipserver.txt"
	ttlregex       = regexp.MustCompile(`(ttl|TTL)\=([\d]+)`)
	MaxCoRoutine   = 128
	pmux           sync.Mutex
	smux           sync.Mutex
	wg             sync.WaitGroup
	loger          *log.Logger
)

func ipByteToUInt32(bip []byte) uint32 {
	var nres uint32 = 0
	res := bytes.NewBuffer(bip)
	binary.Read(res, binary.BigEndian, &nres)
	return nres
}

func uint32ToIpByte(iip uint32) []byte {
	res := bytes.NewBuffer(nil)
	binary.Write(res, binary.BigEndian, iip)
	return res.Bytes()
}

//格式化解析IP
func ipSegmentParse() {
	rfp, err1 := os.Open(ipsegment_file)
	wfp, err2 := os.OpenFile(ipdetail_file, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err1 != nil || err2 != nil {
		loger.Println(err1, err2)
		os.Exit(0)
	}
	reader := bufio.NewReader(rfp)
	for {
		buf, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			break
		}
		astr := strings.Split(string(buf), "/")
		mark, err := strconv.ParseUint(astr[1], 10, 64)
		sip := net.ParseIP(astr[0]).To4()
		nlen := 32 - mark
		maxlen := uint32(^(0xffffff << nlen))
		start := ipByteToUInt32([]byte(sip))
		for i := maxlen; i > 0; i-- {
			bip := uint32ToIpByte(start + i)
			ip := net.IP(bip)
			fmt.Fprintf(wfp, "%s\r\n", ip)
		}
	}
	rfp.Close()
	wfp.Close()
}

func ipPingServer() {
	rfp, err1 := os.Open(ipdetail_file)
	pfp, err2 := os.OpenFile(ipping_file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	sfp, err3 := os.OpenFile(ipserver_file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err1 != nil || err2 != nil || err2 != nil {
		loger.Println(err1, err2, err3)
		os.Exit(0)
	}
	reader := bufio.NewReader(rfp)
	gochan := make(chan bool, MaxCoRoutine)
	for {
		buf, _, err := reader.ReadLine()
		if err != nil {
			loger.Println(err)
			break
		}
		ip := strings.TrimSpace(string(buf))
		if ip == "" {
			continue
		}

		gochan <- true
		wg.Add(1)
		go doPingAndServer(gochan, ip, pfp, sfp)
	}
	rfp.Close()
	pfp.Close()
	sfp.Close()
	wg.Wait()
}

func main() {
	lfp, _ := os.OpenFile("error.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	loger = log.New(lfp, "", log.LstdFlags)
	cmd := flag.String("command", "", "parseip or pingserver!")
	flag.Parse()
	switch *cmd {
	case "parseip":
		ipSegmentParse()
	case "pingserver":
		ipPingServer()
	default:
		flag.Usage()
	}
}

func doPingAndServer(gochan chan bool, ip string, pfp, sfp *os.File) {
	defer func() {
		<-gochan
		wg.Done()
	}()
	cmd := exec.Command("ping", ip, "-w", "1")
	str, err := cmd.Output()
	if err != nil {
		loger.Println(err, ip)
		return
	}
	res := ttlregex.FindSubmatch(str)
	if len(res) != 3 {
		loger.Println(string(str))
		return
	}
	pmux.Lock()
	fmt.Fprintf(pfp, "%s;%s\n", ip, string(res[2]))
	pmux.Unlock()

	rsp, err := http.Get("http://" + ip + ":80/")
	if err != nil {
		loger.Println(err)
	} else {
		smux.Lock()
		fmt.Fprintf(sfp, "%s;%s\r\n", ip, rsp.Header.Get("Server"))
		smux.Unlock()
		rsp.Body.Close()
	}
}
