package logger

import (
	"io"
	"text/template"
	"time"

	"github.com/ghosind/dolphin"
)

// Congfig is the config for Logger middleware.
type Congfig struct {
	// Format is the log format.
	Format *string
	Output *io.Writer
}

type requestLogData struct {
	// IP is the request client ip.
	IP string
	// MilLatency is the latency in milliseconds.
	MilLatency int64
	// Method is the request method.
	Method string
	// Path is the request path with the query string.
	Path string
	// StatusCode is the response status code.
	StatusCode int
}

// Logger is the builtin logger middleware for log requests' information.
func Logger(config ...Congfig) dolphin.HandlerFunc {
	cfg := getConfig(config...)
	tpl := template.Must(template.New("LoggerFormat").Parse(*cfg.Format))

	return func(ctx *dolphin.Context) {
		start := time.Now()
		data := requestLogData{
			IP:     ctx.IP(),
			Method: ctx.Method(),
			Path:   ctx.Path(),
		}

		if query := ctx.RawQuery(); query != "" {
			data.Path += "?" + query
		}

		ctx.Next()

		latency := time.Since(start)
		data.MilLatency = latency.Milliseconds()
		data.StatusCode = ctx.Response.StatusCode()

		output := cfg.Output
		if output == nil {
			appWriter := ctx.LoggerWriter()
			output = &appWriter
		}
		tpl.Execute(*output, data)
	}
}

func getConfig(cfg ...Congfig) *Congfig {
	config := getDefaultConfig()

	if len(cfg) > 0 {
		userConfig := cfg[0]

		if userConfig.Format != nil {
			config.Format = userConfig.Format
		}
		if userConfig.Output != nil {
			config.Output = userConfig.Output
		}
	}

	return config
}

func getDefaultConfig() *Congfig {
	defaultFormat := "{{.Method}} {{.Path}} {{.StatusCode}} {{.MilLatency}}ms"

	return &Congfig{
		Format: &defaultFormat,
		Output: nil,
	}
}
