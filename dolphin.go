package dolphin

import (
	"log"
)

type HandlerFunc func(*Context)

type HandlerChain []HandlerFunc

type O map[string]interface{}

type Config struct {
	Logger *log.Logger

	Port *int
}

var debugMode bool = false

func init() {
	setDebugMode()
}

func New(config *Config) *App {
	return &App{
		logger:   config.Logger,
		port:     config.Port,
		handlers: HandlerChain{},
	}
}

func Default() *App {
	defaultPort := 8080

	app := New(&Config{
		Port: &defaultPort,
	})

	return app
}
