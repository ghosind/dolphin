package recover

import (
	"net/http"

	"github.com/ghosind/dolphin"
)

// Config is the recover middleware config.
type Config struct {
	// Handler is the handler that will be called when a panic happens.
	Handler func(ctx *dolphin.Context, err error)
}

// Recover returns a middleware that recovers from panics, and it'll return 500 error default.
func Recover(config *Config) dolphin.HandlerFunc {
	cfg := getConfig()

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

func getConfig(config ...Config) *Config {
	cfg := Config{
		Handler: defaultHandler,
	}

	if len(config) > 0 {
		if config[0].Handler != nil {
			cfg.Handler = config[0].Handler
		}
	}

	return &cfg
}

func defaultHandler(ctx *dolphin.Context, err error) {
	ctx.String("Internal Server Error", http.StatusInternalServerError)
}
