package pressure

import "testing"

func TestLogger(t *testing.T) {
	l := NewLogger(DEBUG)
	l.LogDebug("Hi There", "this is a debug log message.")
	l.LogWarning("Hi There", "this is a warning log message.")
	l.LogError("Hi There", "this is a error log message.")
}
