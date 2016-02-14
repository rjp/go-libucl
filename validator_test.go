package libucl

import (
	"testing"
)

func TestBufferToString(t *testing.T) {
	expected := "Hello World!"
	buffer := [128]byte{'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd', '!', 0}
	str := byteToBufferAdapter(buffer)
	if str != expected {
		t.Fatalf("bad: \"%s\", expected: \"%s\"", str, expected)
	}
}
