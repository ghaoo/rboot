package wework

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

	agent *wxwork.Agent
}

func newWework(bot *rboot.Robot) rboot.Adapter {
	wx := &wework{
		in:    make(chan *rboot.Message),
		out:   make(chan *rboot.Message),
		agent: agent(),
	}

	bot.Router.HandleFunc("/wxwork", wx.agent.CallbackVerify).Methods("GET")
	bot.Router.HandleFunc("/wxwork", wx.parseRecvHandle).Methods("POST")

	go wx.listenOutgoing()

	return wx
}

func agent() *wxwork.Agent {
	corpid := os.Getenv("WORKWX_CORP_ID")
	secret := os.Getenv("WORKWX_SECRET")
	agentid, err := strconv.Atoi(os.Getenv("WORKWX_AGENT_ID"))
	if err != nil {
		panic(err)
	}
	a := wxwork.NewAgent(corpid, secret, agentid)

	token := os.Getenv("WORKWX_RECV_TOKEN")
	encodingAESKey := os.Getenv("WORKWX_RECV_AES_KEY")
	a.SetCallback(token, encodingAESKey)

	return a
}

func (wx *wework) Name() string {
	return "wework"
}

func (wx *wework) Incoming() chan *rboot.Message {
	return wx.in
}

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

	recv, err := wx.agent.ParseRecvMessage(signature, timestamp, nonce, data)
	if err != nil {
		logrus.WithField("func", "parseRecvHandle.ParseRecvMessage").Error(err)
	}

	msg := rboot.NewMessage(recv.Content)
	msg.From = rboot.User{ID: recv.FromUsername, Name: recv.FromUsername}
	msg.Sender = rboot.User{ID: recv.FromUsername, Name: recv.FromUsername}
	msg.Header.Set("AgentId", strconv.Itoa(wx.agent.AgentID))

	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(&recv); err != nil {
		logrus.WithField("func", "parseRecvHandle.gob.NewEncoder").Error(err)
	}
	msg.Header.Set("Data", buf.String())

	fmt.Println(buf.String())

	wx.in <- msg

}

// 监听 rboot Outgoing
func (wx *wework) listenOutgoing() {
	for msg := range wx.out {
		var wmsg *wxwork.Message

		msgtype := msg.MsgType()
		title := msg.Header.Get("title")
		desc := msg.Header.Get("description")
		mediaid := msg.Header.Get("mediaid")
		url := msg.Header.Get("url")

		switch msgtype {
		case MSG_TYPE_MARKDOWN:
			wmsg = wxwork.NewMarkdownMessage(msg.String())

		case MSG_TYPE_IMAGE, MSG_TYPE_VOICE, MSG_TYPE_FILE:
			wmsg = wxwork.NewMediaMessage(msgtype, mediaid)

		case MSG_TYPE_VIDEO:
			wmsg = wxwork.NewVideoMessage(title, desc, mediaid)

		case MSG_TYPE_TEXTCARD:
			btntxt := msg.Header.Get("btntxt")
			wmsg = wxwork.NewTextCardMessage(title, desc, url, btntxt)

		default:
			wmsg = wxwork.NewTextMessage(msg.String())
		}

		wmsg.SetUser(msg.To.ID)

		_, err := wx.agent.SendMessage(wmsg)
		if err != nil {
			logrus.WithField("func", "wxwork listenOutgoing").Errorf("listen outgoing message err: %v", err)
		}
	}
}

func init() {
	rboot.RegisterAdapter("wework", newWework)
}
