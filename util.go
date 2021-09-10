package dolphin

import (
	"fmt"
	"log"
	"os"
)

func setDebugMode() {
	debug := os.Getenv("DOLPHIN_DEBUG")
	if debug != "" {
		debugMode = true
		debugPrintf("Debug mode enabled.")
	}
}

func debugPrintf(format string, args ...interface{}) {
	if !debugMode {
		return
	}

	format = "[DOLPHIN] " + format

	log.Printf(format, args...)
}

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
