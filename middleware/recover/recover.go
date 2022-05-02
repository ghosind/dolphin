package recover

import (
	"net/http"
	"runtime/debug"

	"github.com/ghosind/dolphin"
)

type PanicError struct {
	// Err is the error that caused the panic.
	Err error
	// StackTrace is the stack trace of the panic.
	StackTrace string
}

type RecoverHandler func(*dolphin.Context, *PanicError)

// Recover returns a middleware that recovers from panics, and it'll return 500 error default.
func Recover(config ...Config) dolphin.HandlerFunc {
	cfg := getConfig(config...)

	return func(ctx *dolphin.Context) {
		defer func() {
			if err := recover(); err != nil {
				e := PanicError{
					Err:        err.(error),
					StackTrace: string(debug.Stack()),
				}

				if cfg.Handler != nil {
					cfg.Handler(ctx, &e)
				}
			}
		}()

		ctx.Next()
	}
}

func defaultHandler(ctx *dolphin.Context, err *PanicError) {
	ctx.String("Internal Server Error", http.StatusInternalServerError)
	ctx.Abort()
}
