// +build darwin dragonfly freebsd linux netbsd openbsd

package rboot

import (
	"strings"
)

func setupPlugin(bot *Robot, in *Message) []*Message {
	rule := in.Header.Get("rule")

	plug := bot.plugins[rule]

	for _, pc := range plug.Command {
		for _, c := range pc.Cmd {
			out, err := runCommand(pc.Dir, "/bin/sh", "-c", c)
			if err != nil {
				return NewMessages(err.Error())
			}

			if len(strings.TrimSpace(out)) > 0 {
				bot.Outgoing(NewMessage(out, in.From))
			}
		}

	}

	return nil
}
