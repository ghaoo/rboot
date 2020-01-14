package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ghaoo/rboot"
)

func signature(ts int64, secret string) string {
	strToSign := fmt.Sprintf("%d\n%s", ts, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToSign))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

type incoming struct {
	MsgType           string            `json:"msgtype"`
	Text              *Text             `json:"text"`
	MsgId             string            `json:"msgId"`
	CreateAt          int64             `json:"createAt"`
	ConversationType  string            `json:"conversationType"` // 1-单聊、2-群聊
	ConversationId    string            `json:"conversationId"`   // // 加密的会话ID
	ConversationTitle string            `json:"conversationId"`   // 会话标题（群聊时才有）
	SenderId          string            `json:"senderId"`
	SenderNick        string            `json:"senderNick"`
	SenderCorpId      string            `json:"senderCorpId"`
	SenderStaffId     string            `json:"senderStaffId"`
	ChatbotUserId     string            `json:"chatbotUserId"`
	AtUsers           map[string]string `json:"atUsers"`

	SessionWebhook string `json:"sessionWebhook"`
	IsAdmin        bool   `json:"isAdmin"`
}

func (ding *dingtalk) listenIncoming(w http.ResponseWriter, r *http.Request) {
	ts := r.Header.Get("Timestamp")
	//token := r.Header.Get("Token")
	sign := r.Header.Get("Sign")

	// timestamp 与系统当前时间戳如果相差1小时以上，则认为是非法的请求。
	tsi, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		log.Printf("请求头可能未附加时间戳信息!!")
	}

	if int64(time.Now().UnixNano()-time.Hour.Nanoseconds()) < tsi*1000 {
		log.Printf("超时访问!!")
		return
	}

	if sign != signature(tsi, ding.secret) {
		log.Printf("非法访问!!")
		return
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	var in incoming
	_ = json.Unmarshal(body, &in)

	msg := rboot.NewMessage(in.Text.Content)
	msg.Sender = in.SenderNick
	msg.Header.Set("sender", in.SenderId)
	msg.Header.Set("dinghook", in.SessionWebhook)
	msg.Header.Set("ding", "1")

	ding.in <- msg

	out := <-ding.out

	// 非当前返回不做反应
	if out.Header.Get("ding") != "1" {
		return
	}

	dmsg := ding.buildMessage(out)
	result, err := json.Marshal(dmsg)
	if err != nil {
		log.Println("marshal outgoing message error: ", err)
		return
	}
	w.Write(result)
}
