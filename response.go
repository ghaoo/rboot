package rboot

import (
	"fmt"
	"io"
	"sync"
	"log"
)

type Response struct {
	*Rboot
	msg *Message

	sync.Mutex
}

func NewResponse(bot *Rboot) *Response {
	res := new(Response)
	res.Rboot = bot

	return res
}

func (res *Response) Receive(msg *Message) error {
	res.Lock()
	defer res.Unlock()

	_, err := msg.Read()
	if err != nil {
		return fmt.Errorf(`Response Receive: message read error %v `, err)
	}

	if msg.Header.From() == `` {
		msg.Header.Set(`From`, `System`)
	}

	if msg.Header.To() == `` {
		msg.Header.Set(`To`, `Nil`)
	}

	res.msg = msg

	res.executeDirectives()

	return nil

}

func (res *Response) ReceiveWithReader(in io.Reader) error {
	msg, err := ReadMessage(in)

	if err != nil {
		return err
	}

	return res.Receive(msg)
}

func (res *Response) Send(strs ...string) error {
	return res.connecter.Send(strs...)
}

func (res *Response) Reply(strs ...string) error {
	return res.connecter.Reply(strs...)
}

func (res *Response) executeDirectives() {
	for _, name := range res.Rboot.conf.Plugins {

		action, err := DirectiveAction(name)
		if err != nil {
			log.Print(err)
		}

		if err = action(res); err != nil {
			log.Print(err)
		}
	}
}
