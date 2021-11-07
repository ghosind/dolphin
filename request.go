package dolphin

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
)

// Request is the wrapped HTTP request object.
type Request struct {
	body *string

	locker sync.Mutex

	request *http.Request
}

// reset resets request object to the initial state.
func (req *Request) reset() {
	req.body = nil
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

// Host reads and returns the request "Host" header.
func (req *Request) Host() string {
	return req.Header(HeaderHost)
}

// IP returns the request client ip.
func (req *Request) IP() string {
	return req.request.RemoteAddr
}

// Method returns the request method.
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
	if req.request.Form == nil {
		req.request.ParseForm()
	}

	return req.request.Form[key]
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

// PostForm returns the form data from the request by the specific key.
func (req *Request) PostForm(key string) string {
	return req.request.FormValue(key)
}
