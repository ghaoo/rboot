package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/ghaoo/rboot"
)

func TestCli(t *testing.T) {
	rawMsg := "hello world!!"

	stdin = strings.NewReader(rawMsg)
	var buf bytes.Buffer
	stdout = &buf

	cli := NewCli()

	in := <-cli.Incoming()
	if in.Content != rawMsg {
		t.Error("cli provider not ingesting incoming messages")
	}

	out := cli.Outgoing()
	out <- rboot.Message{Content: rawMsg}
	close(out)

	to := time.After(5 * time.Second)
	for buf.Len() == 0 {
		select {
		case <-to:
			t.Fatal("could not read output buffer")
		default:
		}
	}

	const expectedStdout = "hello world!!\n"
	gotOut := buf.String()

	if expectedStdout != gotOut {
		t.Errorf("wrong output prompt. Expected output:%v. Got: %v", expectedStdout, gotOut)
	}
}
