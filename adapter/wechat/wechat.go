package wechat

import (
	"strings"

	"github.com/ghaoo/rboot"
	"github.com/ghaoo/rboot/adapter/wechat/sdk"
)

type wx struct {
	in  chan rboot.Message
	out chan rboot.Message

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
		in:     make(chan rboot.Message),
		out:    make(chan rboot.Message),
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

func (w *wx) Sync() {
	contact := w.client.AllContacts()
	users := make([]rboot.User, len(contact))

	for k, v := range contact {
		users[k] = rboot.User{
			ID:   v.UserName,
			Name: v.NickName,
			Data: map[string]interface{}{"wechatFriend": v},
		}
	}

	w.bot.SyncUsers(users)
}

func (w *wx) Incoming() chan rboot.Message {
	return w.in
}

func (w *wx) Outgoing() chan rboot.Message {
	return w.out
}

func (w *wx) run() {

	go func() {
		for msg := range w.out {
			if len(msg.Attachments) > 0 {
				for _, p := range msg.Attachments {
					w.client.SendFile(p.Path, msg.To.ID)
				}
			}
			w.client.SendTextMsg(msg.Content, msg.To.ID)

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
			friend := false
			if c := w.client.ContactByUserName(msg.SenderUserName); c != nil {
				senderName = c.NickName
				if c.Type == sdk.Friend || c.Type == sdk.FriendAndMember {
					friend = true
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
				w.in <- rboot.Message{
					To:      to,
					From:    from,
					Sender:  sender,
					Content: content,
					Mate: map[string]interface{}{
						"AtMe":         msg.AtMe,
						"SendByMySelf": msg.IsSendedByMySelf,
						"SendByFriend": friend,
						"GroupMsg":     msg.IsGroupMsg,
					},
				}
			}

		}

		go func() {
			w.Sync()
			for sync := range w.client.SyncNotify {
				if sync > 0 {
					w.Sync()
				}
			}
		}()

		if es.Hook != nil {
			es.Hook(e)
		}
	}
}

func init() {
	rboot.RegisterAdapter(`wechat`, New)
}
