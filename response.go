package dolphin

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"sync"
)

type Response struct {
	body *bytes.Buffer

	header http.Header

	statusCode int
}

var responsePool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return &Response{}
	},
}

func (resp *Response) reset() {
	resp.body = &bytes.Buffer{}
	resp.header = make(http.Header)
	resp.statusCode = http.StatusOK
}

func (resp *Response) write(rw http.ResponseWriter) {
	for key, val := range resp.header {
		rw.Header()[key] = val
	}

	if resp.statusCode <= 0 || resp.statusCode > 999 {
		resp.statusCode = http.StatusOK
	}

	rw.WriteHeader(resp.statusCode)

	io.Copy(rw, resp.body)
}

func (resp *Response) SetBody(data []byte) (int, error) {
	return resp.body.Write(data)
}

func (resp *Response) SetStatusCode(code int) error {
	if code <= 0 || code > 999 {
		return errors.New("invalid status code")
	}

	resp.statusCode = code

	return nil
}

func (resp *Response) SetContentType(contentType string) {
	resp.header.Set(HTTPHeaderContentType, contentType)
}

func (resp *Response) SetHeader(key, val string) {
	resp.header.Set(key, val)
}

func (resp *Response) AddHeader(key, val string) {
	resp.header.Add(key, val)
}
