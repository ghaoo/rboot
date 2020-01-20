package rboot

type plugin struct {
	Script
	language string
}

func (sf ScriptFunc) Handle(bot *Robot, incoming *Message) []*Message {
	return sf(bot, incoming)
}
