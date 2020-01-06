package wxwork

import (
	"fmt"
	"github.com/ghaoo/rboot"
	"github.com/ghaoo/wxwork"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type workwx struct {
	in  chan rboot.Message
	out chan rboot.Message

	agent *wxwork.Agent
}

func newWework(bot *rboot.Robot) rboot.Adapter {
	wx := &workwx{
		in:    make(chan rboot.Message),
		out:   make(chan rboot.Message),
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

func (wx *workwx) Name() string {
	return "wework"
}

func (wx *workwx) Incoming() chan rboot.Message {
	return wx.in
}

func (wx *workwx) Outgoing() chan rboot.Message {
	return wx.out
}

func (wx *workwx) parseRecvHandle(w http.ResponseWriter, r *http.Request) {
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
		logrus.Errorln("parse receive msg err:", err)
	}

	msg := rboot.Message{
		Channel: "wxwork",
		From:    rboot.User{ID: recv.FromUsername, Name: recv.FromUsername},
		Content: recv.Content,
		Location: rboot.Location{
			Lat:  recv.LocationX,
			Long: recv.LocationY,
		},
		Mate: map[string]interface{}{"originMsg": recv},
	}

	wx.in <- msg

}

// 监听 rboot Outgoing
func (wx *workwx) listenOutgoing() {
	for msg := range wx.out {
		var wmsg *wxwork.Message
		switch msg.Mate["msgtype"] {
		case MSG_TYPE_MARKDOWN:
			wmsg = wxwork.NewMarkdownMessage(msg.Content)
		default:
			wmsg = wxwork.NewTextMessage(msg.Content)
		}

		wmsg.SetUser(msg.To.ID)

		_, err := wx.agent.SendMessage(wmsg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func init() {
	rboot.RegisterAdapter("wxwork", newWework)
}
