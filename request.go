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

// reset resets request object to the initial state.
func (req *Request) reset() {
	req.body = nil
	req.jsonBody = nil
	req.request = nil
}

// Cookie returns the cookie by the specific name.
func (req *Request) Cookie(key string) (*http.Cookie, error) {
	return req.request.Cookie(key)
}

// File returns the file from multipart form by the specific key.
func (req *Request) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return req.request.FormFile(key)
}

// Header returns the value from the request header by the specific key.
func (req *Request) Header(key string) string {
	return req.request.Header.Get(key)
}

// MultiValuesHeader returns the string array type values from the request header by the
// specific key.
func (req *Request) MultiValuesHeader(key string) []string {
	return req.request.Header.Values(key)
}

// Methods returns the request method.
func (req *Request) Method() string {
	return req.request.Method
}

// Path returns the request path.
func (req *Request) Path() string {
	return req.request.URL.Path
}

// Query returns the query string value from the request by the specific key.
func (req *Request) Query(key string) string {
	return req.request.FormValue(key)
}

// MultiValuesQuery returns the string array type values from the request query string by the
// specific key.
func (req *Request) MultiValuesQuery(key string) []string {
	if req.request.PostForm == nil {
		req.request.ParseForm()
	}

	return req.request.PostForm[key]
}

// Post returns the body of the request.
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

// PostJSON returns the request body and parses it to the interface{} object.
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

// PostFrom returns the form data from the request by the specific key.
func (req *Request) PostForm(key string) string {
	return req.request.FormValue(key)
}
