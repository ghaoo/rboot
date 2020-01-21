// +build windows

package ymlcmd

import (
	"github.com/ghaoo/rboot"
	"strings"
)

func setup(bot *rboot.Robot, in *rboot.Message) []*rboot.Message {
	rule := in.Header.Get("rule")

	cmd := command[rule]

	for _, cs := range cmd.Command {
		for _, c := range cs.Cmd {
			args := strings.Split(c, " ")

			out, err := runCommand(cs.Dir, args[0], args[1:]...)
			if err != nil {
				return rboot.NewMessages(err.Error())
			}

			bot.Outgoing(rboot.NewMessage(out, in.From))
		}

	}

	return nil
}
