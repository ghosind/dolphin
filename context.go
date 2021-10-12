package dolphin

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
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
	// state is the context state, it can be used to store any data and pass to
	// the next handlers.
	state map[string]interface{}
}

// allocateContext returns a new context instance.
func allocateContext() *Context {
	return &Context{}
}

// reset the context instance to initial state.
func (ctx *Context) reset(req *http.Request) {
	ctx.Request = requestPool.Get().(*Request)
	ctx.Request.reset()
	ctx.Request.request = req

	ctx.Response = responsePool.Get().(*Response)
	ctx.Response.reset()

	ctx.handlers = HandlerChain{}
	ctx.index = -1
	ctx.isAbort = false
	ctx.state = map[string]interface{}{}
}

// finalize releases the context, request, and response resources.
func (ctx *Context) finalize() {
	requestPool.Put(ctx.Request)
	responsePool.Put(ctx.Response)

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

	ctx.Response.SetContentType(contentType)
	_, err := ctx.Response.SetBody(data)
	if err != nil {
		debugPrintf("Failed to set response body: %v", err)
		return err
	}

	return nil
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
func (ctx *Context) Use(handlers ...HandlerFunc) {
	ctx.handlers = append(ctx.handlers, handlers...)
}

// Log call the app logger with the given format and args.
func (ctx *Context) Log(fmt string, args ...interface{}) {
	ctx.app.log(fmt, args...)
}

// Get retrieves the value of the given key from the context state.
func (ctx *Context) Get(key string) (interface{}, bool) {
	val, ok := ctx.state[key]

	return val, ok
}

// Set sets the value of the given key to the context state.
func (ctx *Context) Set(key string, val interface{}) {
	ctx.state[key] = val
}

// Cookie returns the named cookie provided in the request.
func (ctx *Context) Cookie(key string) (*http.Cookie, error) {
	return ctx.Request.Cookie(key)
}

// File returns the named file provided in the request.
func (ctx *Context) File(key string) (multipart.File, *multipart.FileHeader, error) {
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

// Method returns the request method.
func (ctx *Context) Method() string {
	return ctx.Request.Method()
}

// Path returns the request path.
func (ctx *Context) Path() string {
	return ctx.Request.Path()
}

// Post returns the request post data.
func (ctx *Context) Post() string {
	return ctx.Request.Post()
}

// PostJSON returns the request post data as JSON object.
func (ctx *Context) PostJSON() (interface{}, error) {
	return ctx.Request.PostJSON()
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

// Success writes the given data to the response body as JSON, and set the status code to 200 (OK).
func (ctx *Context) Success(data interface{}) error {
	return ctx.JSON(data, http.StatusOK)
}

// Fail writes the given data to the response body as JSON, and set the status code
// to 400 (Bad Request).
func (ctx *Context) Fail(data interface{}) error {
	return ctx.JSON(data, http.StatusBadRequest)
}

// JSON stringifies and writes the given data to the response body, and set the content type
// to "application/json".
func (ctx *Context) JSON(data interface{}, statusCode ...int) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ctx.send(payload, MIMETypeJSON, statusCode...)
}

// HTML writes the given data to the response body as HTML, and set the content type
// to "text/html".
func (ctx *Context) HTML(html string, statusCode ...int) error {
	return ctx.send([]byte(html), MIMETypeHTML, statusCode...)
}

// String writes the given string to the response body, and sets the context type
// to "text/plain".
func (ctx *Context) String(data string, statusCode ...int) error {
	return ctx.send([]byte(data), MIMETypeText, statusCode...)
}

// Redirect redirects the request to the given URL, it'll set status code to 302 (Found) as
// default status code if it isn't set.
func (ctx *Context) Redirect(url string, statusCode ...int) {
	if url == "back" {
		url = ctx.Request.Header(HeaderReferrer)
	}

	code := http.StatusFound

	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	ctx.Response.SetStatusCode(code)
	ctx.Response.SetHeader(HeaderLocation, url)
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

// SetContentType sets the content type of the response.
func (ctx *Context) SetContentType(contentType string) {
	ctx.Response.SetContentType(contentType)
}
