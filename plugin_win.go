// +build windows

package rboot

import (
	"strings"
)

func setupPlugin(bot *Robot, in *Message) []*Message {
	rule := in.Header.Get("rule")

	plug := bot.plugins[rule]

	for _, pc := range plug.Command {
		for _, c := range pc.Cmd {
			args := strings.Split(c, " ")

			out, err := runCommand(pc.Dir, args[0], args[1:]...)
			if err != nil {
				return NewMessages(err.Error())
			}

			bot.Outgoing(NewMessage(out, in.From))
		}

	}

	return nil
}
