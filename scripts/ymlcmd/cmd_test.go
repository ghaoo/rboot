package ymlcmd

import (
	"testing"
)

func TestCmd_Load(t *testing.T) {
	b, err := load("./test/composer.yml")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(b))
}

func TestCmd_AllCmd(t *testing.T) {
	cmds, err := allCmd("./test")
	if err != nil {
		t.Error(err)
	}

	t.Log(cmds)
}
