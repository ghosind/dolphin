package dolphin

import (
	"log"
	"sync"
)

type HandlerFunc func(*Context)

type HandlerChain []HandlerFunc

type O map[string]interface{}

type Config struct {
	// Logger is the logger used by the app, dolphin will use log.Printf if this
	// have not set.
	Logger *log.Logger
	// Port is the port to listen on.
	Port *int
}

var debugMode bool = false

func init() {
	setDebugMode()
}

// New creates a new App instance.
func New(config *Config) *App {
	return &App{
		logger:   config.Logger,
		port:     config.Port,
		handlers: HandlerChain{},
		pool: sync.Pool{
			New: func() interface{} {
				return allocateContext()
			},
		},
	}
}

// Default creates a new App instance with default configuration.
func Default() *App {
	defaultPort := 8080

	app := New(&Config{
		Port: &defaultPort,
	})

	return app
}
