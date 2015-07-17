package core

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	maxFormSize = int64(10 << 20)
)

type Input struct {
	req *http.Request
}

func NewInput(req *http.Request) *Input {
	return &Input{req: req}
}

func (this *Input) Get(key string) string {
	return this.req.FormValue(key)
}

func (this *Input) Post(key string) string {
	return this.req.PostFormValue(key)
}

func (this *Input) Cookie(name string) string {
	val, err := this.req.Cookie(name)
	if err != nil {
		return ""
	}
	return val.Value
}

func (this *Input) RawData() (buf []byte, err error) {
	reader := io.LimitReader(this.req.Body, maxFormSize+1)
	buf, err = ioutil.ReadAll(reader)
	if err != nil {
		return buf, err
	}
	if int64(len(buf)) > maxFormSize {
		err = errors.New("http: POST too large")
	}
	return
}
