package rboot

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"sync"
)

type Response struct {
	*Rboot
	Matcher string

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

	scrName, ok := res.matchRuleset(text)

	if ok {

		action, err := DirectiveAction(scrName)

		if err != nil {
			return err
		}

		return action(res)
	}

	return fmt.Errorf(`Response Receive: no matching scripts... `)
}

func (res *Response) ReceiveWithReader(in io.Reader) error {
	msg, err := ReadMessage(in)

	if err != nil {
		return err
	}

	return res.Receive(msg)
}

func (res *Response) matchRuleset(msg string) (string, bool) {
	for scr, rules := range rulesets {
		for matcher, rule := range rules {
			if res.match(rule, msg) {
				res.Matcher = matcher
				return scr, true
			}
		}
	}

	log.Printf(`no match script`)
	return ``, false
}

func (res *Response) match(pattern, msg string) bool {

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
