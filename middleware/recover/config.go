package recover

// Config is the recover middleware config.
type Config struct {
	// Handler is the handler that will be trigged when catching some panics.
	// It'll return a 500 error if the handler is not set.
	Handler RecoverHandler
}

// DefaultConfig is the default recover middleware config.
var DefaultConfig = Config{
	Handler: defaultHandler,
}

func getConfig(config ...Config) Config {
	if len(config) < 1 {
		return DefaultConfig
	}

	return config[0]
}
