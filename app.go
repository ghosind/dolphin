package dolphin

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// App is the dolphin web server engine.
type App struct {
	certFile *string

	handlers HandlerChain

	keyFile *string

	logger *log.Logger

	pool *sync.Pool

	port int

	reqPool *sync.Pool

	resPool *sync.Pool

	server *http.Server
}

// Run starts the app and listens on the given port.
func (app *App) Run() {
	app.initServer()

	err := app.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		app.log("Failed to run server: %v\n", err)
	}
}

// RunTLS starts the app to provide HTTPS service and listens on the given port.
func (app *App) RunTLS() error {
	if app.certFile == nil || len(*app.certFile) <= 0 {
		return ErrNoTLSCert
	}
	if app.keyFile == nil || len(*app.keyFile) <= 0 {
		return ErrNoTLSKey
	}

	app.initServer()

	err := app.server.ListenAndServeTLS(*app.certFile, *app.keyFile)

	if err != nil && err != http.ErrServerClosed {
		app.log("Failed to run server: %v\n", err)
		return err
	}

	return nil
}

// Shutdown tries to close active connections and stops the server.
func (app *App) Shutdown(ctx ...context.Context) error {
	if len(ctx) == 0 {
		ctx = append(ctx, context.Background())
	}

	return app.server.Shutdown(ctx[0])
}

// ServeHTTP implements the http.Handler interface.
func (app *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := app.pool.Get().(*Context)
	ctx.reset(app, req)

	ctx.Next()
	ctx.writeResponse(rw)

	ctx.finalize()
}

// Use registers one or more middlewares or request handlers to the app.
func (app *App) Use(handlers ...HandlerFunc) {
	app.handlers = append(app.handlers, handlers...)
}

// initServer gets listenning port and initialize HTTP server.
func (app *App) initServer() {
	addr := resolveListenAddr(&app.port)
	app.log("Server running at %s.\n", addr)

	app.server.Addr = addr
	app.server.Handler = app
}

// LoggerWriter returns the app logger's writer, or os.Stderr if the app logger is not set.
func (app *App) LoggerWriter() io.Writer {
	if app.logger != nil {
		return app.logger.Writer()
	}

	return os.Stderr
}

// log logs a message to the app's logger or log.Printf.
func (app *App) log(fmt string, args ...interface{}) {
	logger := log.Printf
	if app.logger != nil {
		logger = app.logger.Printf
	}

	logger(fmt, args...)
}
