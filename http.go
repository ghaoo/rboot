package rboot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 超时时间(秒)
var timeout = 60

// 计算消息签名
// 1. 将参数按照 时间\n加密秘钥\n消息内容 排列，获取需要加密的字符串
// 2. 使用 sha256 将上面的字符串加密获取加密后的字符串
// 3. 将 sha256 加密后的字符串使用 base64 编码获取最终的签名值
func signature(datetime, secret, content string) string {
	strToSign := fmt.Sprintf("%s\n%s\n%s", datetime, secret, content)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToSign))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

// verifySign 验证签名
func verifySign(sign, secret, content, datetime string) error {
	dt, err := time.Parse("2006-01-02 15:04:05", datetime)
	if err != nil {
		return fmt.Errorf("datetime format error, should 2006-01-02 15:04:05")
	}

	if time.Now().Sub(dt).Seconds() > float64(timeout) {
		return fmt.Errorf("timeout! the request time is long ago, please try again")
	}

	if sign != signature(datetime, secret, content) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// listenIncoming 用于传入消息，为保证消息的安全性，消息应该进行签名加密
func (bot *Robot) listenIncoming(w http.ResponseWriter, r *http.Request) {
	listenMessage(w, r, bot.inputChan)
}

// listenOutgoing 用于传出消息，为保证消息的安全性，消息应该进行签名加密
func (bot *Robot) listenOutgoing(w http.ResponseWriter, r *http.Request) {
	listenMessage(w, r, bot.outputChan)
}

// listenMessage 监听消息
func listenMessage(w http.ResponseWriter, r *http.Request, msgChan chan *Message) {
	sign := r.Header.Get("sign")
	datetime := r.Header.Get("datetime")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("message read failed, errmsg: " + err.Error()))
		return
	}
	defer r.Body.Close()

	secret := os.Getenv("ROBOT_SECRET")

	if err = verifySign(sign, secret, string(content), datetime); err != nil {
		w.WriteHeader(403)
		w.Write([]byte(err.Error()))
		return
	}

	var msg = NewMessage(string(content), r.Header.Get("to"))

	msg.From = r.Header.Get("from")
	msg.Sender = r.Header.Get("sender")
	msg.Header = Header(r.Header)

	msgChan <- msg

	w.WriteHeader(200)
	w.Write([]byte("发送成功"))
}
