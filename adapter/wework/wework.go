package wework

import (
	"bytes"
	"encoding/gob"
	"github.com/ghaoo/rboot"
	"github.com/ghaoo/wxwork"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type wework struct {
	in  chan *rboot.Message
	out chan *rboot.Message

	bot    *rboot.Robot
	client *wxwork.Agent
}

func newWework(bot *rboot.Robot) rboot.Adapter {
	wx := &wework{
		in:     make(chan *rboot.Message),
		out:    make(chan *rboot.Message),
		bot:    bot,
		client: newAgent(),
	}

	bot.Router.HandleFunc("/wework", wx.client.CallbackVerify).Methods("GET")
	bot.Router.HandleFunc("/wework", wx.parseRecvHandle).Methods("POST")

	go wx.listenOutgoing()

	contacts := wx.getContacts()
	if len(contacts) > 0 {
		wx.bot.SyncContacts(contacts)
	}

	return wx
}

func newAgent() *wxwork.Agent {
	corpid := os.Getenv("WORKWX_CORP_ID")
	secret := os.Getenv("WORKWX_SECRET")
	agentid, err := strconv.Atoi(os.Getenv("WORKWX_AGENT_ID"))
	if err != nil {
		panic(err)
	}
	a := wxwork.NewAgent(corpid, agentid)
	a = a.WithSecret(secret)

	token := os.Getenv("WORKWX_RECV_TOKEN")
	encodingAESKey := os.Getenv("WORKWX_RECV_AES_KEY")
	a.SetCallback(token, encodingAESKey)

	return a
}

// Name 适配器名称
func (wx *wework) Name() string {
	return "wework"
}

// Incoming 企业微信传入的消息
func (wx *wework) Incoming() chan *rboot.Message {
	return wx.in
}

// Outgoing 传入消息经脚本处理后向企业微信输入的消息
func (wx *wework) Outgoing() chan *rboot.Message {
	return wx.out
}

func (wx *wework) parseRecvHandle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	signature := query.Get("msg_signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorln("read callback msg err:", err)
	}

	recv, err := wx.client.ParseRecvMessage(signature, timestamp, nonce, data)
	if err != nil {
		logrus.WithField("func", "parseRecvHandle.ParseRecvMessage").Error(err)
	}

	msg := rboot.NewMessage(recv.Content)
	msg.From = recv.FromUsername
	msg.Sender = recv.FromUsername
	msg.Header.Set("AgentId", strconv.Itoa(wx.client.AgentID))

	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(&recv); err != nil {
		logrus.WithField("func", "parseRecvHandle.gob.NewEncoder").Error(err)
	}
	msg.Header.Set("Data", buf.String())

	wx.in <- msg

}

// listenOutgoing 监听 rboot Outgoing
func (wx *wework) listenOutgoing() {
	for out := range wx.out {
		var msg *wxwork.Message

		title := out.Header.Get("title")
		desc := out.Header.Get("description")
		mediaid := out.Header.Get("mediaid")
		url := out.Header.Get("url")

		switch out.MsgType() {
		case MSG_TYPE_TEXT:
			msg = wxwork.NewTextMessage(out.String())

		case MSG_TYPE_IMAGE, MSG_TYPE_VOICE, MSG_TYPE_FILE:
			msg = wxwork.NewMediaMessage(out.MsgType(), mediaid)

		case MSG_TYPE_VIDEO:
			msg = wxwork.NewVideoMessage(title, desc, mediaid)

		case MSG_TYPE_TEXTCARD:
			btntxt := out.Header.Get("btntxt")
			msg = wxwork.NewTextCardMessage(title, desc, url, btntxt)

		default:
			msg = wxwork.NewMarkdownMessage(out.String())
		}

		msg.SetUser(out.To)

		_, err := wx.client.SendMessage(msg)
		if err != nil {
			logrus.WithField("func", "wxwork listenOutgoing").Errorf("listen outgoing message err: %v", err)
		}
	}
}

func init() {
	rboot.RegisterAdapter("wework", newWework)
}
