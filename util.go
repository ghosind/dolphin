package dolphin

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// defaultPort is the default listening port of the dolphin framework.
var defaultPort = ":8080"

// DebugWriter is the writter used to print debug messages, default `os.Stdout`.
var DebugWriter io.Writer = os.Stdout

// setDebugMode sets global variable `debugMode` and enables debug mode if
// environment variable `DOLPHIN_DEBUG` is set to "true".
func setDebugMode() {
	debugMode = strings.ToLower(os.Getenv("DOLPHIN_DEBUG")) == "true"
	if debugMode {
		debugPrintf("Debug mode enabled.")
	}
}

// debugPrintf prints message if debug mode is enabled.
func debugPrintf(format string, args ...interface{}) {
	if !debugMode {
		return
	}

	format = "[DOLPHIN] " + format

	fmt.Fprintf(DebugWriter, format, args...)
}

// resolveListenAddr resolves the listening address from parameter or
// environment variable `DOLPHIN_PORT`. Port number of listening should
// be greater than 0 and less than 65536.
func resolveListenAddr(port *int) string {
	if port != nil && *port > 0 && *port < 65536 {
		return fmt.Sprintf(":%d", *port)
	}

	envPort := os.Getenv("DOLPHIN_PORT")
	if envPort != "" {
		port, err := strconv.Atoi(envPort)
		if err != nil || port <= 0 || port > 65535 {
			debugPrintf("Environment variable \"DOLPHIN_PORT\"(%s) is invalid: %s", envPort)
			return defaultPort
		}

		return ":" + envPort
	}

	debugPrintf("No port setting, use default port 8080.")
	return defaultPort
}
