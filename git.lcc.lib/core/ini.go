package core

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type iniNodeItem map[string]string
type iniNodeEntery map[string]iniNodeItem
type iniNodeFile map[string]iniNodeEntery

var localIniStore iniNodeFile
var mutex sync.RWMutex

func init() {
	localIniStore = make(iniNodeFile)
}

func readFromCache(hash, entery, key string) (val string, isok bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	val, isok = localIniStore[hash][entery][key]
	return
}

func pushIntoCache(hash, entery, key, str string) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, isok := localIniStore[hash][entery]; !isok {
		localIniStore[hash][entery] = make(iniNodeItem)
	}
	localIniStore[hash][entery][key] = str
}

func cleanFromCache(hash string) {
	if _, isok := localIniStore[hash]; isok {
		mutex.Lock()
		delete(localIniStore, hash)
		localIniStore[hash] = make(iniNodeEntery)
		mutex.Unlock()
	}
}

type IniConfig struct {
	file string
	hash string
}

func NewIni(file string) *IniConfig {
	iniconf := &IniConfig{file, fmt.Sprintf("%x", md5.Sum([]byte(file)))}
	if _, isok := localIniStore[iniconf.hash]; !isok {
		localIniStore[iniconf.hash] = make(iniNodeEntery)
	}
	return iniconf
}

func (this *IniConfig) readIniItem(entery, key string) (string, error) {
	rd, err := os.Open(this.file)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer rd.Close()
	dr := bufio.NewReader(rd)
	entery = fmt.Sprintf("[%s]", entery)
	line, str, isfix := "", "", false

	for {
		line, err = dr.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == entery {
			isfix = true
			continue
		}
		if isfix && strings.HasPrefix(line, key) {
			segs := strings.SplitN(line, "=", 2)
			if len(segs) == 2 {
				str = strings.TrimSpace(segs[1])
			}
			break
		}
	}
	return str, err
}

func (this *IniConfig) GetItem(entery, key string) string {
	str, isok := readFromCache(this.hash, entery, key)
	if !isok {
		str, _ = this.readIniItem(entery, key)
		pushIntoCache(this.hash, entery, key, str)
	}
	return str
}

func (this *IniConfig) PutItem(entery, key, value string) {
	pushIntoCache(this.hash, entery, key, value)
}

func (this *IniConfig) Format2File() {
	absfile, _ := filepath.Abs(this.file)
	dir := filepath.Dir(absfile)
	os.MkdirAll(dir, 0777)
	fp, err := os.OpenFile(absfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fp.Sync()
		fp.Close()
	}()
	for entry, items := range localIniStore[this.hash] {
		fmt.Fprintf(fp, "[%s]\r\n", entry)
		for key, value := range items {
			fmt.Fprintf(fp, "%s=%s\r\n", key, value)
		}
		fmt.Fprint(fp, "\r\n")
	}
}

func (this *IniConfig) ReLoad() {
	cleanFromCache(this.hash)
}
