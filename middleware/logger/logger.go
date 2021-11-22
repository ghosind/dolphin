package logger

import (
	"io"
	"text/template"
	"time"

	"github.com/ghosind/dolphin"
)

// LoggerCongfig is the config for Logger middleware.
type LoggerCongfig struct {
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
	// Path is the request path.
	Path string
	// StatusCode is the response status code.
	StatusCode int
}

// Logger is the builtin logger middleware for log requests' information.
func Logger(config ...LoggerCongfig) dolphin.HandlerFunc {
	cfg := getConfig(config...)
	tpl := template.Must(template.New("LoggerFormat").Parse(*cfg.Format))

	return func(ctx *dolphin.Context) {
		start := time.Now()
		data := requestLogData{
			IP:     ctx.IP(),
			Method: ctx.Method(),
			Path:   ctx.Path(),
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

func getConfig(cfg ...LoggerCongfig) *LoggerCongfig {
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

func getDefaultConfig() *LoggerCongfig {
	defaultFormat := "{{.Method}} {{.Path}} {{.StatusCode}} {{.MilLatency}}ms"

	return &LoggerCongfig{
		Format: &defaultFormat,
		Output: nil,
	}
}
