package dolphin

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
)

type Context struct {
	Request *Request

	Response *Response

	app *App

	handlers HandlerChain

	index int

	isAbort bool

	state map[string]interface{}
}

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

func (ctx *Context) finalize() {
	requestPool.Put(ctx.Request)
	responsePool.Put(ctx.Response)

	ctx.app.pool.Put(ctx)
}

func (ctx *Context) writeResponse(rw http.ResponseWriter) {
	ctx.Response.write(rw)
}

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

func (ctx *Context) Next() {
	ctx.index++

	for !ctx.isAbort && ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

func (ctx *Context) Use(handlers ...HandlerFunc) {
	ctx.handlers = append(ctx.handlers, handlers...)
}

func (ctx *Context) Log(fmt string, args ...interface{}) {
	ctx.app.log(fmt, args...)
}

func (ctx *Context) Get(key string) (interface{}, bool) {
	val, ok := ctx.state[key]

	return val, ok
}

func (ctx *Context) Set(key string, val interface{}) {
	ctx.state[key] = val
}

func (ctx *Context) Cookie(key string) (*http.Cookie, error) {
	return ctx.Request.Cookie(key)
}

func (ctx *Context) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.Request.File(key)
}

func (ctx *Context) Header(key string) string {
	return ctx.Request.Header(key)
}

func (ctx *Context) MultiValuesHeader(key string) []string {
	return ctx.Request.MultiValuesHeader(key)
}

func (ctx *Context) Method() string {
	return ctx.Request.Method()
}

func (ctx *Context) Path() string {
	return ctx.Request.Path()
}

func (ctx *Context) Post() string {
	return ctx.Request.Post()
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Request.PostForm(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Request.Query(key)
}

func (ctx *Context) QueryDefault(key, defaultValue string) string {
	val := ctx.Query(key)

	if val == "" {
		return defaultValue
	}
	return val
}

func (ctx *Context) MultiValuesQuery(key string) []string {
	return ctx.Request.MultiValuesQuery(key)
}

func (ctx *Context) MultiValuesQueryDefault(key string, defaultValues []string) []string {
	val := ctx.MultiValuesQuery(key)

	if len(val) == 0 {
		return defaultValues
	}
	return val
}

func (ctx *Context) Success(data interface{}) error {
	return ctx.JSON(data, http.StatusOK)
}

func (ctx *Context) Fail(data interface{}) error {
	return ctx.JSON(data, http.StatusBadRequest)
}

func (ctx *Context) JSON(data interface{}, statusCode ...int) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ctx.send(payload, MIMETypeJSON, statusCode...)
}

func (ctx *Context) HTML(html string, statusCode ...int) error {
	return ctx.send([]byte(html), MIMETypeHTML, statusCode...)
}

func (ctx *Context) String(data string, statusCode ...int) error {
	return ctx.send([]byte(data), MIMETypeText, statusCode...)
}

func (ctx *Context) Redirect(url string, statusCode ...int) {
	if url == "back" {
		url = ctx.Request.Header(HTTPHeaderReferrer)
	}

	code := http.StatusFound

	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	ctx.Response.SetStatusCode(code)
	ctx.Response.SetHeader(HTTPHeaderLocation, url)
}

func (ctx *Context) AddCookies(cookies ...*http.Cookie) {
	ctx.Response.AddCookies(cookies...)
}

func (ctx *Context) SetHeader(key, val string) {
	ctx.Response.SetHeader(key, val)
}
