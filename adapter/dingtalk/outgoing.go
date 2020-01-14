package dingtalk

import (
	"bytes"
	"encoding/json"
	"github.com/ghaoo/rboot"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var baseHookUrl = "https://oapi.dingtalk.com/robot/send"

func (ding *dingtalk) listenOutgoing() {
	for out := range ding.out {

		accessToken := os.Getenv("DING_ROBOT_HOOK_ACCESS_TOKEN")
		secret := os.Getenv("DING_ROBOT_HOOK_SECRET")
		timestamp := time.Now().UnixNano() / int64(time.Millisecond)

		sign := signature(timestamp, secret)

		query := url.Values{}

		query.Set("access_token", accessToken)
		query.Set("timestamp", strconv.FormatInt(timestamp, 10))
		query.Set("sign", sign)

		hookUrl, _ := url.Parse(baseHookUrl)
		hookUrl.RawQuery = query.Encode()

		dmsg := ding.buildMessage(out)
		dimsg, _ := json.Marshal(dmsg)

		req, err := http.NewRequest("POST", hookUrl.String(), bytes.NewBuffer(dimsg))
		if err != nil {
			log.Printf("create request failed: %v", err)
			return
		}

		req.Header.Add("Accept-Charset", "utf8")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		var caller Result
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("call the webhook failed: %v", err)
			return
		}

		if err = json.NewDecoder(resp.Body).Decode(caller); err != nil {
			log.Printf("decode response body error: %v", err)
			return
		}

		if caller.ErrCode != 0 {
			log.Printf("response error message: %s", caller.ErrMsg)
			return
		}
	}
}

func (ding *dingtalk) buildMessage(msg *rboot.Message) *Message {
	var dmsg *Message

	title := msg.Header.Get("title")
	if title == "" {
		title = "消息"
	}

	switch msg.Header.Get("msgtype") {
	case "text":
		dmsg = NewTextMessage(msg.String())
	case "link":
		dmsg = NewLinkMessage(title, msg.String(), msg.Header.Get("msgUrl"), msg.Header.Get("picUrl"))
	case "actionCard":
		var card *ActionCard

		err := json.Unmarshal(msg.Bytes(), &card)
		if err != nil {
			dmsg = NewEmptyMessage()
		} else {
			dmsg = NewActionCardMessage(card)
		}

	case "feedCard":
		var links []Link

		err := json.Unmarshal(msg.Bytes(), &links)
		if err != nil {
			dmsg = NewEmptyMessage()
		} else {
			dmsg = NewFeedCardMessage(links)
		}
	default:
		dmsg = NewMarkdownMessage(title, msg.String())
	}

	atMobiles := msg.Header["atMobiles"]

	atAll, _ := strconv.ParseBool(msg.Header.Get("atAll"))

	dmsg.At = At{
		AtMobiles: atMobiles,
		IsAtAll:   atAll,
	}

	return dmsg
}

type Result struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
