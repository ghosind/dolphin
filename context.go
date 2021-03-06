package dolphin

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
)

// Context is the context instance for the request.
type Context struct {
	// Request is the wrapped HTTP request.
	Request *Request
	// Response is the wrapped HTTP response.
	Response *Response
	// app is the framework application instance.
	app *App
	// handlers is the handler chain.
	handlers HandlerChain
	// index is the current handler index.
	index int
	// isAbort is the flag to indicate if the current handler chain should be
	// aborted.
	isAbort bool
	// pathVariables stores current request path variables.
	pathVariables map[string]string
	// sm is the mutex for protecting the context state.
	sm sync.RWMutex
	// state is the context state, it can be used to store any data and pass to
	// the next handlers.
	state map[string]any
}

// allocateContext returns a new context instance.
func allocateContext() *Context {
	return &Context{}
}

// reset the context instance to initial state.
func (ctx *Context) reset(app *App, req *http.Request) {
	ctx.Request = app.reqPool.Get().(*Request)
	ctx.Request.reset()
	ctx.Request.request = req

	ctx.Response = app.resPool.Get().(*Response)
	ctx.Response.reset()

	ctx.app = app
	ctx.handlers = HandlerChain{}
	ctx.index = -1
	ctx.isAbort = false
	ctx.pathVariables = make(map[string]string)
	ctx.state = make(map[string]any)

	ctx.Use(app.handlers...)
}

// finalize releases the context, request, and response resources.
func (ctx *Context) finalize() {
	ctx.app.reqPool.Put(ctx.Request)
	ctx.app.resPool.Put(ctx.Response)

	ctx.app.pool.Put(ctx)
}

// writeResponse writes data from context to the response.
func (ctx *Context) writeResponse(rw http.ResponseWriter) {
	ctx.Response.write(rw)
}

// send writes the response data to the body buffer, sets the contentType and the status code
// if it's set.
func (ctx *Context) send(data []byte, contentType string, statusCode ...int) error {
	if len(statusCode) >= 1 {
		err := ctx.Response.SetStatusCode(statusCode[0])
		if err != nil {
			debugPrintf("Failed to set HTTP status code: %v", err)
			return err
		}
	} else {
		ctx.Response.SetStatusCode(http.StatusOK)
	}

	ctx.SetContentType(contentType)
	_, err := ctx.Response.SetBody(data)
	if err != nil {
		debugPrintf("Failed to set response body: %v", err)
		return err
	}

	return nil
}

// addPathVariable adds the given path variable to the context.
func (ctx *Context) addPathVariable(key string, value string) {
	ctx.sm.Lock()
	defer ctx.sm.Unlock()

	ctx.pathVariables[key] = value
}

// Abort stops the current handler chain.
func (ctx *Context) Abort() {
	ctx.isAbort = true
}

// Next calls the next handler in the handler chain.
func (ctx *Context) Next() {
	ctx.index++

	for !ctx.isAbort && ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

// Use registers one or more middlewares or request handlers to the context.
func (ctx *Context) Use(handler ...HandlerFunc) *Context {
	if len(handler) > 0 {
		ctx.handlers = append(ctx.handlers, handler...)
	}

	return ctx
}

// Log call the app logger with the given format and args.
func (ctx *Context) Log(fmt string, args ...any) {
	ctx.app.log(fmt, args...)
}

// LoggerWriter returns the app logger's writer, or os.Stderr if the app logger is not set.
func (ctx *Context) LoggerWriter() io.Writer {
	return ctx.app.LoggerWriter()
}

// Get retrieves the value of the given key from the context state.
func (ctx *Context) Get(key string) (any, bool) {
	ctx.sm.RLock()
	defer ctx.sm.RUnlock()

	val, ok := ctx.state[key]

	return val, ok
}

// Has returns true if the given key exists in the context state.
func (ctx *Context) Has(key string) bool {
	ctx.sm.Lock()
	defer ctx.sm.Unlock()

	_, ok := ctx.state[key]

	return ok
}

// Set sets the value of the given key to the context state.
func (ctx *Context) Set(key string, val any) any {
	ctx.sm.Lock()
	defer ctx.sm.Unlock()

	oldVal := ctx.state[key]
	ctx.state[key] = val

	return oldVal
}

// BasicAuth returns the username and password that provided in the request
// header 'Authorization' field.
func (ctx *Context) BasicAuth() (username, password string, ok bool) {
	return ctx.Request.BasicAuth()
}

// Body returns the body of the request.
func (ctx *Context) Body() string {
	return ctx.Request.Body()
}

// Cookie returns the named cookie provided in the request.
func (ctx *Context) Cookie(key string) (cookie *http.Cookie, err error) {
	return ctx.Request.Cookie(key)
}

// File returns the named file provided in the request.
func (ctx *Context) File(key string) (file multipart.File, fileHeader *multipart.FileHeader, err error) {
	return ctx.Request.File(key)
}

// Header returns the header value from the request by the given key.
func (ctx *Context) Header(key string) string {
	return ctx.Request.Header(key)
}

// MultiValuesHeader returns a string array value for the given key from the request header.
func (ctx *Context) MultiValuesHeader(key string) []string {
	return ctx.Request.MultiValuesHeader(key)
}

// IP returns the request client IP address.
func (ctx *Context) IP() string {
	return ctx.Request.IP()
}

// Method returns the request method.
func (ctx *Context) Method() string {
	return ctx.Request.Method()
}

// Path returns the request path.
func (ctx *Context) Path() string {
	return ctx.Request.Path()
}

// PathVariable returns the path variable by the given key.
func (ctx *Context) PathVariable(key string) string {
	return ctx.pathVariables[key]
}

// Post returns the request post data.
func (ctx *Context) Post() string {
	return ctx.Request.Body()
}

// PostJSON gets request body and parses to the given struct.
func (ctx *Context) PostJSON(payload any) error {
	body := ctx.Request.Body()

	err := json.Unmarshal([]byte(body), payload)
	if err != nil {
		return err
	}

	return nil
}

// PostForm returns the request form data.
func (ctx *Context) PostForm(key string) string {
	return ctx.Request.PostForm(key)
}

// Query returns the query value from the request by the given key.
func (ctx *Context) Query(key string) string {
	return ctx.Request.Query(key)
}

// QueryDefault returns the query value from the request by the given key, and returns the
// default value if the key is not found.
func (ctx *Context) QueryDefault(key, defaultValue string) string {
	val := ctx.Query(key)

	if val == "" {
		return defaultValue
	}
	return val
}

// MultiValuesQuery returns a string array value for the given key from the request query.
func (ctx *Context) MultiValuesQuery(key string) []string {
	return ctx.Request.MultiValuesQuery(key)
}

// MultiValuesQueryDefault returns a string array value for the given key from the request query
// string, and returns the default values if the key is not exists.
func (ctx *Context) MultiValuesQueryDefault(key string, defaultValues []string) []string {
	val := ctx.MultiValuesQuery(key)

	if len(val) == 0 {
		return defaultValues
	}
	return val
}

// RawQuery returns raw query string (withoud ?).
func (ctx *Context) RawQuery() string {
	return ctx.Request.RawQuery()
}

// Referer returns the referring URL of the request.
func (ctx *Context) Referer() string {
	return ctx.Request.Referer()
}

// UserAgent returns the client user agent of the request.
func (ctx *Context) UserAgent() string {
	return ctx.Request.UserAgent()
}

// Success writes the given data to the response body as JSON, and set the status code to 200 (OK).
func (ctx *Context) Success(data any) error {
	return ctx.JSON(data, http.StatusOK)
}

// Fail writes the given data to the response body as JSON, and set the status code
// to 400 (Bad Request).
func (ctx *Context) Fail(data any) error {
	return ctx.JSON(data, http.StatusBadRequest)
}

// JSON stringifies and writes the given data to the response body, and set the content type
// to "application/json".
func (ctx *Context) JSON(data any, statusCode ...int) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ctx.send(payload, "application/json", statusCode...)
}

// HTML writes the given data to the response body as HTML, and set the content type
// to "text/html".
func (ctx *Context) HTML(html string, statusCode ...int) error {
	return ctx.send([]byte(html), "text/html", statusCode...)
}

// String writes the given string to the response body, and sets the context type
// to "text/plain".
func (ctx *Context) String(data string, statusCode ...int) error {
	return ctx.send([]byte(data), "text/plain", statusCode...)
}

// Redirect redirects the request to the given URL, it'll set status code to 302 (Found) as
// default status code if it isn't set.
func (ctx *Context) Redirect(url string, statusCode ...int) {
	if url == "back" {
		url = ctx.Request.Referer()
	}

	code := http.StatusFound

	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	ctx.Response.SetStatusCode(code)
	ctx.Response.SetHeader("Location", url)
}

// Write writes the given data to the response body.
func (ctx *Context) Write(data []byte) (len int, err error) {
	return ctx.Response.SetBody(data)
}

// AddCookies adds one or more given cookies to the response.
func (ctx *Context) AddCookies(cookies ...*http.Cookie) {
	ctx.Response.AddCookies(cookies...)
}

// AddHeader appends the given header pair to the response.
func (ctx *Context) AddHeader(key, val string) {
	ctx.Response.AddHeader(key, val)
}

// SetHeader sets the value of the given header key.
func (ctx *Context) SetHeader(key, val string) {
	ctx.Response.SetHeader(key, val)
}

// SetStatusCode sets the status code of the response.
func (ctx *Context) SetStatusCode(code int) error {
	return ctx.Response.SetStatusCode(code)
}

// SetContentType sets the response HTTP header "Content-Type" field to the
// specific MIME type value.
func (ctx *Context) SetContentType(contentType string) {
	ctx.Response.SetHeader("Content-Type", contentType)
}
