package rboot

import (
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"
)

type Response struct {
	*Rboot
	executePlugin *Plugin

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

	b, err := msg.Read()
	if err != nil {
		return fmt.Errorf(`Response Receive: message read error %v `, err)
	}

	if msg.Header.From() == `` {
		msg.Header.Set(`From`, `System`)
	}

	if msg.Header.To() == `` {
		msg.Header.Set(`To`, `Nil`)
	}

	text := string(b)

	plug_name, ok := checkRuleset(text)

	if !ok {
		return fmt.Errorf(`Response Receive: no matching plugin... `)
	}

	if ok {
		res.executePlugin, err = getPlugin(plug_name)

		if err != nil {
			return fmt.Errorf(`Response Receive: get plugin error %v `, err)
		}

		return nil
	}

	return nil

}

func (res *Response) ReceiveWithReader(in io.Reader) error {
	msg, err := ReadMessage(in)

	if err != nil {
		return err
	}

	return res.Receive(msg)
}

func checkRuleset(msg string) (string, bool) {
	for plug, rules := range rulesets {
		for _, rule := range rules {
			if match(rule, msg) {
				return plug, true
			}
		}
	}

	return ``, false
}

func match(pattern, msg string) bool {

	reg := regexp.MustCompile(pattern)

	if reg.MatchString(msg) {
		return true
	}

	return false
}

func (res *Response) Send(strs ...string) error {
	return res.connecter.Send(strs...)
}

func (res *Response) Reply(strs ...string) error {
	return res.connecter.Reply(strs...)
}

func newHeader(from, to, contentType string) Header {
	header := make(Header)

	now := time.Now().Local().String()

	header.Add(`From`, from)
	header.Add(`To`, to)
	header.Add(`ContentType`, contentType)
	header.Add(`Date`, now)

	return header
}
