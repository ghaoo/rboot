package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
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

		/*if out.Header.Get("ding") != "" {
			continue
		}*/

		var msg *Message

		title := out.Header.Get("title")
		if title == "" {
			title = "消息"
		}

		switch out.MsgType() {
		case "text":
			msg = NewTextMessage(out.String())
		case "link":
			msg = NewLinkMessage(title, out.String(), out.Header.Get("msgUrl"), out.Header.Get("picUrl"))
		case "actionCard":
			var card *ActionCard

			err := json.Unmarshal(out.Bytes(), &card)
			if err != nil {
				msg = NewEmptyMessage()
			} else {
				msg = NewActionCardMessage(card)
			}

		case "feedCard":
			var links []Link

			err := json.Unmarshal(out.Bytes(), &links)
			if err != nil {
				msg = NewEmptyMessage()
			} else {
				msg = NewFeedCardMessage(links)
			}
		default:
			msg = NewMarkdownMessage(title, out.String())
		}

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

		fmt.Println(hookUrl)
		dmsg, _ := json.Marshal(msg)

		fmt.Println("\n", string(dmsg))

		req, err := http.NewRequest("POST", hookUrl.String(), bytes.NewBuffer(dmsg))
		if err != nil {
			log.Printf("create request failed: %v", err)
		}

		req.Header.Add("Accept-Charset", "utf8")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		var caller Result
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("call the webhook failed: %v", err)
		}

		if err = json.NewDecoder(resp.Body).Decode(caller); err != nil {
			log.Printf("decode response body error: %v", err)
		}

		if caller.ErrCode != 0 {
			log.Printf("response error message: %s", caller.ErrMsg)
		}

	}
}

type Result struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
