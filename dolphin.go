package dolphin

import (
	"log"
	"net/http"
	"sync"
)

// HandlerFunc is the function that register as a handler to the app.
type HandlerFunc func(*Context)

// HandlerChain is a chain of handlers.
type HandlerChain []HandlerFunc

// O is an alias for map that contains string key and interface{} value.
type O map[string]interface{}

// Config is the configuration for the dolphin web application.
type Config struct {
	// Logger is the logger used by the app, dolphin will use log.Printf if this
	// have not set.
	Logger *log.Logger
	// Port is the port to listen on.
	Port *int
}

// debugMode indicates the enable/disable status of debug mode.
var debugMode bool = false

func init() {
	// Load debug mode setting from environment variable.
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
		server: &http.Server{},
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
