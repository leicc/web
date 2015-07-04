package core

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	_ = iota
	LOG_FATAL
	LOG_ERROR
	LOG_DEBUG
	LOG_INFO
)

type CoreLoger struct {
	*log.Logger
	dir  string
	mask int8
}

func NewLoger(dir string, mask int8) *CoreLoger {
	if mask < 1 {
		mask = 0
	}
	if dir == "" {
		dir = "./log"
	}
	os.MkdirAll(dir, 0644)
	file := fmt.Sprintf("%s/%s.log", dir, time.Now().Format("20060102"))
	fs, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fs = os.Stdout
	}
	loger := &CoreLoger{log.New(fs, "", log.Lshortfile|log.Ldate|log.Ltime), dir, mask}
	return loger
}

func (this *CoreLoger) Write(mask int8, v ...interface{}) {
	if mask <= this.mask {
		this.Print(v...)
	}
}

func (this *CoreLoger) Writef(mask int8, format string, v ...interface{}) {
	if mask <= this.mask {
		this.Printf(format, v...)
	}
}

var std = NewLoger("./cache/log", 8)

func Log(mask int8, v ...interface{}) {
	std.Write(mask, v...)
}

func Logf(mask int8, format string, v ...interface{}) {
	std.Writef(mask, format, v...)
}
