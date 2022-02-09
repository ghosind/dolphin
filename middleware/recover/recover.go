package recover

import (
	"net/http"

	"github.com/ghosind/dolphin"
)

// Recover returns a middleware that recovers from panics, and it'll return 500 error default.
func Recover(config ...Config) dolphin.HandlerFunc {
	cfg := getConfig(config...)

	return func(ctx *dolphin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if cfg.Handler != nil {
					cfg.Handler(ctx, err.(error))
				}
			}
		}()

		ctx.Next()
	}
}

func defaultHandler(ctx *dolphin.Context, err error) {
	ctx.String("Internal Server Error", http.StatusInternalServerError)
	ctx.Abort()
}
