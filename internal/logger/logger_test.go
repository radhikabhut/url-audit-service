package logger

import (
	"testing"
)

func TestInitLogger(t *testing.T) {
	// Simply call it with text format and JSON format to ensure no panics
	InitLogger("debug", "json")
	InitLogger("info", "text")
	InitLogger("invalid", "invalid")
}
