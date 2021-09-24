package dolphin

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
)

// Request is the wrapped HTTP request object.
type Request struct {
	body *string

	jsonBody interface{}

	locker sync.Mutex

	request *http.Request
}

var requestPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return &Request{}
	},
}

// reset resets request object to initial state.
func (req *Request) reset() {
	req.body = nil
	req.jsonBody = nil
	req.request = nil
}

// Cookie returns the cookie by the specific name.
func (req *Request) Cookie(key string) (*http.Cookie, error) {
	return req.request.Cookie(key)
}

func (req *Request) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return req.request.FormFile(key)
}

func (req *Request) Header(key string) string {
	return req.request.Header.Get(key)
}

func (req *Request) MultiValuesHeader(key string) []string {
	return req.request.Header.Values(key)
}

func (req *Request) Method() string {
	return req.request.Method
}

func (req *Request) Path() string {
	return req.request.URL.Path
}

func (req *Request) Query(key string) string {
	return req.request.FormValue(key)
}

func (req *Request) MultiValuesQuery(key string) []string {
	return req.request.PostForm[key]
}

func (req *Request) Post() string {
	if req.body != nil {
		return *req.body
	}

	req.locker.Lock()
	defer req.locker.Unlock()

	buf := new(strings.Builder)
	io.Copy(buf, req.request.Body)
	body := buf.String()

	req.body = &body

	return body
}

func (req *Request) PostJSON() (interface{}, error) {
	if req.jsonBody != nil {
		return req.jsonBody, nil
	}

	body := req.Post()

	req.locker.Lock()
	defer req.locker.Unlock()

	var payload interface{}
	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return nil, err
	}

	req.jsonBody = payload

	return payload, nil
}

func (req *Request) PostForm(key string) string {
	return req.request.FormValue(key)
}
