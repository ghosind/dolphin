package dolphin

import (
	"strings"
	"testing"
)

func TestDebugMode(t *testing.T) {
	// Disable debug mode
	t.Setenv("DOLPHIN_DEBUG", "")
	setDebugMode()
	if debugMode {
		t.Errorf("debug mode expect false, actual true.")
	}

	// Enable debug mode
	t.Setenv("DOLPHIN_DEBUG", "true")
	setDebugMode()
	if !debugMode {
		t.Errorf("debug mode expect true, actual false.")
	}

	// Clear debug mode
	t.Setenv("DOLPHIN_DEBUG", "")
	setDebugMode()
	if debugMode {
		t.Errorf("debug mode expect false, actual true.")
	}
}

func TestResolveListenAddr(t *testing.T) {
	var port *int = nil

	// Set addr as default port if argument is nil
	addr := resolveListenAddr(port)
	if strings.Compare(addr, ":8080") != 0 {
		t.Errorf("Listening address expect \"%s\", actual \"%s\",", ":8080", addr)
	}

	// Set addr as default port if argument is invalid.
	port = new(int)
	*port = -100
	addr = resolveListenAddr(port)
	if strings.Compare(addr, ":8080") != 0 {
		t.Errorf("Listening address expect \"%s\", actual \"%s\",", ":8080", addr)
	}

	*port = 65536
	addr = resolveListenAddr(port)
	if strings.Compare(addr, ":8080") != 0 {
		t.Errorf("Listening address expect \"%s\", actual \"%s\",", ":8080", addr)
	}

	// Get invalid port from environment
	t.Setenv("DOLPHIN_PORT", "0")
	addr = resolveListenAddr(nil)
	if strings.Compare(addr, ":8080") != 0 {
		t.Errorf("Listening address expect \"%s\", actual \"%s\",", ":8080", addr)
	}

	t.Setenv("DOLPHIN_PORT", "65536")
	addr = resolveListenAddr(nil)
	if strings.Compare(addr, ":8080") != 0 {
		t.Errorf("Listening address expect \"%s\", actual \"%s\",", ":8080", addr)
	}

	t.Setenv("DOLPHIN_PORT", "test")
	addr = resolveListenAddr(nil)
	if strings.Compare(addr, ":8080") != 0 {
		t.Errorf("Listening address expect \"%s\", actual \"%s\",", ":8080", addr)
	}

	t.Setenv("DOLPHIN_PORT", "5000-test")
	addr = resolveListenAddr(nil)
	if strings.Compare(addr, ":8080") != 0 {
		t.Errorf("Listening address expect \"%s\", actual \"%s\",", ":8080", addr)
	}

	// Get valid port from environment
	t.Setenv("DOLPHIN_PORT", "5000")
	addr = resolveListenAddr(nil)
	if strings.Compare(addr, ":5000") != 0 {
		t.Errorf("Listening address expect \"%s\", actual \"%s\",", ":5000", addr)
	}

	// Valid port number.
	*port = 4000
	addr = resolveListenAddr(port)
	if strings.Compare(addr, ":4000") != 0 {
		t.Errorf("Listening address expect \"%s\", actual \"%s\",", ":4000", addr)
	}

	// Clear environment
	t.Setenv("DOLPHIN_PORT", "")
}
