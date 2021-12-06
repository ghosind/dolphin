package dolphin

import (
	"bytes"
	"io"
	"net/http"
)

// Response is the HTTP response wrapper.
type Response struct {
	body *bytes.Buffer

	cookies []*http.Cookie

	header http.Header

	statusCode int
}

// reset resets response object to initial state.
func (resp *Response) reset() {
	resp.body = &bytes.Buffer{}
	resp.cookies = make([]*http.Cookie, 0)
	resp.header = make(http.Header)
	resp.statusCode = http.StatusOK
}

// write writes response to the specific HTTP response writer.
func (resp *Response) write(rw http.ResponseWriter) {
	// Add cookies to response.
	if len(resp.cookies) > 0 {
		for _, cookie := range resp.cookies {
			if cookie == nil {
				continue
			}
			resp.AddHeader(HeaderSetCookie, cookie.String())
		}
	}

	// Write response header.
	for key, val := range resp.header {
		rw.Header()[key] = val
	}

	// Set response status code to OK if not set or it's invalid.
	if resp.statusCode <= 0 || resp.statusCode > 999 {
		resp.statusCode = http.StatusOK
	}

	// Set response status code
	rw.WriteHeader(resp.statusCode)

	// Write response body.
	io.Copy(rw, resp.body)
}

// SetBody sets response body.
func (resp *Response) SetBody(data []byte) (len int, err error) {
	return resp.body.Write(data)
}

// AddCookies adds cookies setting to response, it will set response HTTP
// header "Set-Cookie" field.
func (resp *Response) AddCookies(cookies ...*http.Cookie) {
	resp.cookies = append(resp.cookies, cookies...)
}

// AddHeader adds value to the specific response HTTP header field.
func (resp *Response) AddHeader(key, val string) {
	resp.header.Add(key, val)
}

// SetHeader sets the specific response HTTP header field.
func (resp *Response) SetHeader(key, val string) {
	resp.header.Set(key, val)
}

// SetStatusCode sets the status code of the response.
func (resp *Response) SetStatusCode(code int) error {
	if code <= 0 || code > 999 {
		return ErrInvalidStatusCode
	}

	resp.statusCode = code

	return nil
}

// StatusCode gets the response status code.
func (resp *Response) StatusCode() int {
	return resp.statusCode
}
