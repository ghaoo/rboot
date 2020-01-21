// +build darwin dragonfly freebsd linux netbsd openbsd

package ymlcmd

import "github.com/ghaoo/rboot"

func setup(bot *rboot.Robot, in *rboot.Message) []*rboot.Message {
	rule := in.Header.Get("rule")

	cmd := command[rule]

	for _, cs := range cmd.Command {
		for _, c := range cs.Cmd {
			out, err := runCommand(cs.Dir, "/bin/sh", "-c", c)
			if err != nil {
				return rboot.NewMessages(err.Error())
			}

			bot.Outgoing(rboot.NewMessage(out, in.From))
		}

	}

	return nil
}
