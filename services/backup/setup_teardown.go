package backup

import (
	"os"
	"testing"
)

// Setup setup
func Setup(t *testing.T) {
	os.MkdirAll("test/output", 0777)
	os.MkdirAll("test/output1", 0777) // to test dirhash
}

// Teardown setup
func Teardown(t *testing.T) {
	os.RemoveAll("test")
}
