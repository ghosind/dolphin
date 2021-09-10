package dolphin

import (
	"fmt"
	"log"
	"os"
)

// setDebugMode Set global variable debugMode and enable debug mode if
// environment variable "DOLPHIN_DEBUG" is set.
func setDebugMode() {
	debug := os.Getenv("DOLPHIN_DEBUG")
	if debug != "" {
		debugMode = true
		debugPrintf("Debug mode enabled.")
	} else {
		debugMode = false
	}
}

// debugPrintf Print message if debug mode is enabled.
func debugPrintf(format string, args ...interface{}) {
	if !debugMode {
		return
	}

	format = "[DOLPHIN] " + format

	log.Printf(format, args...)
}

// resolveListenAddr Resolve listen address by parameter or environment
// variable "DOLPHIN_PORT". Port number should be greater than 0 and less
// than 65535.
func resolveListenAddr(port *int) string {
	if port == nil || *port < 0 || *port > 65535 {
		if port := os.Getenv("DOLPHIN_PORT"); port != "" {
			debugPrintf("Get server port from env: %s", port)
			return ":" + port
		}

		debugPrintf("No port setting, use default port 8080.")
		return ":8080"
	}

	return fmt.Sprintf(":%d", *port)
}
