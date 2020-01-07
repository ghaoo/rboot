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
				w.client.SendFile(msg.Header.Get("file"), msg.To.ID)
			}

			w.client.SendTextMsg(msg.String(), msg.To.ID)
		}
	}()

	es := w.client.Stream

	for e := range es.Event {
		switch e.Type {
		case sdk.EVENT_STOP_LOOP:
			return
		case sdk.EVENT_NEW_MESSAGE:
			msg := e.Data.(sdk.MsgData)

			toName := w.client.MySelf.NickName
			if msg.IsGroupMsg {
				if c := w.client.ContactByUserName(msg.ToUserName); c != nil {
					toName = c.NickName
				} else {
					toName = `无名`
				}
			}

			to := rboot.User{
				ID:   msg.ToUserName,
				Name: toName,
			}

			fromName := ``
			if c := w.client.ContactByUserName(msg.FromUserName); c != nil {
				fromName = c.NickName
			}

			from := rboot.User{
				ID:   msg.FromUserName,
				Name: fromName,
			}

			senderName := ``
			isFriend := false
			if c := w.client.ContactByUserName(msg.SenderUserName); c != nil {
				senderName = c.NickName
				if c.Type == sdk.Friend || c.Type == sdk.FriendAndMember {
					isFriend = true
				}
			}
			sender := rboot.User{
				ID:   msg.SenderUserName,
				Name: senderName,
			}

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

			if !msg.IsGroupMsg || msg.AtMe {
				rmsg := rboot.NewMessage(content)
				rmsg.To = to
				rmsg.From = from
				rmsg.Sender = sender
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
