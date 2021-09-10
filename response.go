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

	cookies []*http.Cookie

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
	resp.cookies = make([]*http.Cookie, 0)
	resp.header = make(http.Header)
	resp.statusCode = http.StatusOK
}

func (resp *Response) write(rw http.ResponseWriter) {
	if len(resp.cookies) > 0 {
		for _, cookie := range resp.cookies {
			if cookie == nil {
				continue
			}
			resp.AddHeader(HTTPHeaderSetCookie, cookie.String())
		}
	}

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

func (resp *Response) AddCookies(cookies ...*http.Cookie) {
	resp.cookies = append(resp.cookies, cookies...)
}

func (resp *Response) SetContentType(contentType string) {
	resp.header.Set(HTTPHeaderContentType, contentType)
}

func (resp *Response) AddHeader(key, val string) {
	resp.header.Add(key, val)
}

func (resp *Response) SetHeader(key, val string) {
	resp.header.Set(key, val)
}

func (resp *Response) SetStatusCode(code int) error {
	if code <= 0 || code > 999 {
		return errors.New("invalid status code")
	}

	resp.statusCode = code

	return nil
}
