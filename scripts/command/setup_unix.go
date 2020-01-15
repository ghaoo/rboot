// +build darwin dragonfly freebsd linux netbsd openbsd

package command

import (
	"fmt"
	"github.com/ghaoo/rboot"
	"os/exec"
)

func setup(bot *rboot.Robot, in *rboot.Message) []*rboot.Message {
	rule := in.Header.Get("rule")

	cmd := command[rule]

	for _, c := range cmd.Cmd {
		runCmd := exec.Command("/bin/sh", c)

		output, err := runCmd.CombinedOutput()
		if err != nil {
			bot.Outgoing(rboot.NewMessage(fmt.Sprintf("`error running command`: %v: %q", err, string(output)), in.From))
		}

		bot.Outgoing(rboot.NewMessage(string(output), in.From))
	}

	return nil
}
