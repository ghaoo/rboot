package wechat

import (
	"strconv"
	"strings"

	"github.com/ghaoo/rboot"
	"github.com/ghaoo/rboot/adapter/wechat/sdk"
)

type wx struct {
	in  chan *rboot.Message
	out chan *rboot.Message

	bot    *rboot.Robot
	client *sdk.WeChat
}

func New(bot *rboot.Robot) rboot.Adapter {

	// 初始化微信
	client, err := sdk.NewBot(nil)
	if err != nil {
		panic(err)
	}

	w := &wx{
		in:     make(chan *rboot.Message),
		out:    make(chan *rboot.Message),
		bot:    bot,
		client: client,
	}

	client.Hook(w.Assisant)

	go w.run()

	return w
}

func (w *wx) Name() string {
	return "wechat"
}

func (w *wx) Incoming() chan *rboot.Message {
	return w.in
}

func (w *wx) Outgoing() chan *rboot.Message {
	return w.out
}

func (w *wx) run() {

	go func() {
		for msg := range w.out {
			if msg.Header.Get("file") != "" {
				for _, to := range msg.To {
					w.client.SendFile(msg.Header.Get("file"), to)
				}
			}

			for _, to := range msg.To {
				w.client.SendTextMsg(msg.String(), to)
			}
		}
	}()

	es := w.client.Stream

	for e := range es.Event {
		switch e.Type {
		case sdk.EVENT_STOP_LOOP:
			return
		case sdk.EVENT_NEW_MESSAGE:
			msg := e.Data.(sdk.MsgData)

			to := []string{msg.ToUserName}

			from := msg.FromUserName

			sender := msg.SenderUserName

			content := msg.Content

			if msg.AtMe {
				atme := `@`
				if len(w.client.MySelf.DisplayName) > 0 {
					atme += w.client.MySelf.DisplayName
				} else {
					atme += w.client.MySelf.NickName
				}
				content = strings.TrimSpace(strings.TrimPrefix(content, atme))
			}

			isFriend := false
			if c := w.client.ContactByUserName(msg.SenderUserName); c != nil {
				if c.Type == sdk.Friend || c.Type == sdk.FriendAndMember {
					isFriend = true
				}
			}

			if !msg.IsGroupMsg || msg.AtMe {
				rmsg := rboot.NewMessage(content)
				rmsg.To = to
				rmsg.From = from
				rmsg.Header.Set("Sender", sender)
				rmsg.Header.Set("AtMe", strconv.FormatBool(msg.AtMe))
				rmsg.Header.Set("SendByMySelf", strconv.FormatBool(msg.IsSendedByMySelf))
				rmsg.Header.Set("GroupMsg", strconv.FormatBool(msg.IsGroupMsg))
				rmsg.Header.Set("IsFriend", strconv.FormatBool(isFriend))

				w.in <- rmsg
			}
		}

	}
}

func init() {
	rboot.RegisterAdapter(`wechat`, New)
}
