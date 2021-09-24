package dolphin

import (
	"log"
	"net/http"
	"sync"
)

type App struct {
	handlers HandlerChain

	logger *log.Logger

	pool sync.Pool

	port *int
}

// Run starts the app and listens on the given port.
func (app *App) Run() {
	addr := resolveListenAddr(app.port)

	app.log("Server running at %s.\n", addr)

	err := http.ListenAndServe(addr, app)
	if err != nil {
		app.log("Failed to run server: %v\n", err)
	}
}

// ServeHTTP implements the http.Handler interface.
func (app *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := app.pool.Get().(*Context)
	ctx.reset(req)

	ctx.app = app

	ctx.Use(app.handlers...)

	ctx.Next()
	ctx.writeResponse(rw)

	ctx.finalize()
}

// Use registers one or more middlewares or request handlers to the app.
func (app *App) Use(handlers ...HandlerFunc) {
	app.handlers = append(app.handlers, handlers...)
}

func (app *App) log(fmt string, args ...interface{}) {
	handler := log.Printf
	if app.logger != nil {
		handler = app.logger.Printf
	}

	handler(fmt, args...)
}
