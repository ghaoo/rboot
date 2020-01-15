package command

import (
	"github.com/go-yaml/yaml"
	"testing"
)

func TestCmd_Load(t *testing.T) {
	b, err := load("./test/echo.yml")
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

func TestCmd_Run(t *testing.T) {
	b, err := load("./test/echo.yml")
	if err != nil {
		t.Error(err)
	}

	var cmd = Cmd{}
	err = yaml.Unmarshal(b, &cmd)
	if err != nil {
		t.Error(err)
	}

	for _, c := range cmd.Cmd {
		out, err := runCommand("/bin/sh", c)

		if err != nil {
			t.Error(err)
		}

		t.Log(out)
	}
}
