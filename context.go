package dolphin

import "net/http"

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

func (ctx *Context) Next() {
	ctx.index++

	for !ctx.isAbort && ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}
