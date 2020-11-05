package wework

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/ghaoo/rboot"
	"github.com/ghaoo/wxwork"
	"github.com/sirupsen/logrus"
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

// Name 转接器名称
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
		logrus.WithFields(logrus.Fields{
			"mod": `wework`,
		}).Errorln("read callback msg err:", err)
	}

	recv, err := wx.client.ParseRecvMessage(signature, timestamp, nonce, data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"mod": `wework`,
		}).Error(err)
	}

	msg := rboot.NewMessage(recv.Content)
	msg.From = recv.FromUsername
	msg.Sender = recv.FromUsername
	msg.Header.Set("AgentId", strconv.Itoa(wx.client.AgentID))
	msg.Header.Set("msgtype", "markdown")

	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(&recv); err != nil {
		logrus.WithFields(logrus.Fields{
			"mod": `rboot`,
		}).Error(err)
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
		url := out.Header.Get("url")

		switch out.Header.Get("msgtype") {
		case MSG_TYPE_TEXT:
			msg = wxwork.NewTextMessage(out.String())

		case MSG_TYPE_IMAGE, MSG_TYPE_VOICE, MSG_TYPE_FILE:
			if out.Header.Get("file") != "" {
				media, err := wx.client.MediaUpload(out.Header.Get("file"))
				if err != nil {
					msg = wxwork.NewTextMessage("上传附件失败：" + err.Error())
				} else {
					msg = wxwork.NewMediaMessage(out.Header.Get("msgtype"), media.MediaId)
				}
			}

		case MSG_TYPE_VIDEO:
			if out.Header.Get("file") != "" {
				media, err := wx.client.MediaUpload(out.Header.Get("file"))
				if err != nil {
					msg = wxwork.NewTextMessage("上传视频失败：" + err.Error())
				} else {
					msg = wxwork.NewVideoMessage(title, desc, media.MediaId)
				}
			}

		case MSG_TYPE_TEXTCARD:
			btntxt := out.Header.Get("btntxt")
			msg = wxwork.NewTextCardMessage(title, desc, url, btntxt)

		default:
			msg = wxwork.NewMarkdownMessage(out.String())
		}

		msg.SetUser(out.To)

		_, err := wx.client.SendMessage(msg)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"mod": `wework`,
			}).Errorf("listen outgoing message err: %v", err)
		}
	}
}

func init() {
	rboot.RegisterAdapter("wework", newWework)
}
