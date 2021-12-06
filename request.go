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

	bodyOnce *sync.Once

	request *http.Request
}

// reset resets request object to the initial state.
func (req *Request) reset() {
	req.body = nil
	req.bodyOnce = &sync.Once{}
	req.request = nil
}

// BasicAuth returns the username and password that provided in the request
// header 'Authorization' field.
func (req *Request) BasicAuth() (username, password string, ok bool) {
	return req.request.BasicAuth()
}

// Body returns the body of the request.
func (req *Request) Body() string {
	if req.body == nil {
		req.bodyOnce.Do(func() {
			buf := new(strings.Builder)
			io.Copy(buf, req.request.Body)
			body := buf.String()

			req.body = &body
		})
	}

	return *req.body
}

// Cookie returns the cookie by the specific name.
func (req *Request) Cookie(key string) (cookie *http.Cookie, err error) {
	return req.request.Cookie(key)
}

// File returns the file from multipart form by the specific key.
func (req *Request) File(key string) (file multipart.File, fileHeader *multipart.FileHeader, err error) {
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
	return req.request.Host
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

// Post returns the body of the request, it's the alias of Request.Body().
func (req *Request) Post() string {
	return req.Body()
}

// PostForm returns the form data from the request by the specific key.
func (req *Request) PostForm(key string) string {
	return req.request.FormValue(key)
}

// Query returns the query string value from the request by the specific key.
func (req *Request) Query(key string) string {
	return req.request.FormValue(key)
}

// MultiValuesQuery returns the string array type values from the request query string by the
// specific key.
func (req *Request) MultiValuesQuery(key string) []string {
	if req.request.Form == nil {
		err := req.request.ParseForm()
		debugPrintf("Failed to parse request form: %v", err)
	}

	return req.request.Form[key]
}

// RawQuery returns raw query string (withoud ?).
func (req *Request) RawQuery() string {
	return req.request.URL.RawQuery
}

// Referrer returns the referring URL of the request.
func (req *Request) Referer() string {
	return req.request.Referer()
}

// UserAgent returns the client user agent of the request.
func (req *Request) UserAgent() string {
	return req.request.UserAgent()
}
