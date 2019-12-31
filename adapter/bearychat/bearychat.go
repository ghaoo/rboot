package bearychat

import (
	"encoding/json"
	"fmt"
	"github.com/ghaoo/rboot"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// bearychat adapter
type beary struct {
	in  chan rboot.Message
	out chan rboot.Message
}

func newBeary(bot *rboot.Robot) rboot.Adapter {
	beary := &beary{
		in:  make(chan rboot.Message),
		out: make(chan rboot.Message),
	}

	beary.run(bot)

	return beary
}

func (b *beary) Name() string {
	return `bearychat`
}

func (b *beary) Incoming() chan rboot.Message {
	return b.in
}

func (b *beary) Outgoing() chan rboot.Message {
	return b.out
}

// 监听 rboot 需要发送给 bearychat 的消息
func (b *beary) listenOutgoing() {
	for msg := range b.out {
		res := Response{
			Text:     msg.Content,
			Markdown: true,
			Channel:  msg.Channel,
			User:     msg.To.ID,
		}

		fmt.Println("out channel----", msg.Channel)
		fmt.Println("out user----", res.User)

		var attCount = len(msg.Attachments)
		var atts = make([]Attachment, len(msg.Attachments))
		if attCount > 0 {
			for _, matt := range msg.Attachments {
				att := Attachment{
					Title:  matt.Title,
					Text:   matt.Text,
					Color:  matt.Color,
					Images: matt.Images,
				}

				atts = append(atts, att)
			}
		}

		res.Attachments = atts

		if err := sendMessage(res); err != nil {
			fmt.Println(err)
		}
	}
}

// 监听 bearychat 传入 rboot 的消息
func (b *beary) listenIncoming(w http.ResponseWriter, r *http.Request) {

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		logrus.WithField("func", "bearychat listenIncoming").Errorf("listen incoming message err: %v", err)
		return
	}

	// 验证token
	if req.Token != os.Getenv("BEARYCHAT_OUT_TOKEN") {
		w.WriteHeader(http.StatusNotExtended)
		return
	}

	// 是否需要删除 bearychat 设置的 TRIGGER_WORD（和 scripts 相关）
	// req.Text = strings.TrimPrefix(req.Text, os.Getenv("TRIGGER_WORD"))

	msg := rboot.Message{
		Channel:     req.ChannelName,
		From:        rboot.User{ID: req.UserName, Name: req.UserName},
		Sender:      rboot.User{ID: req.UserName, Name: req.UserName},
		Content:     req.Text,
		Attachments: nil,
	}

	b.in <- msg
}

func (b *beary) run(bot *rboot.Robot) {
	go b.listenOutgoing()

	bot.Router.HandleFunc("/beary", b.listenIncoming).Methods("GET").Name("beary_listen_message")
}

func init() {
	rboot.RegisterAdapter(`bearychat`, newBeary)
}
