package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCli(t *testing.T) {
	rawMsg := "hello world!!"

	stdin = strings.NewReader(rawMsg)
	var buf bytes.Buffer
	stdout = &buf
}
