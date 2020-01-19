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

// 超时时间
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

func (bot *Robot) listenIncoming(w http.ResponseWriter, r *http.Request) {
	sign := r.Header.Get("sign")
	datetime := r.Header.Get("datetime")
	dt, err := time.Parse("2006-01-02 15:04:05", datetime)
	if err != nil {
		w.WriteHeader(412)
		w.Write([]byte("datetime format is error, should 2006-01-02 15:04:05"))
		return
	}

	if time.Now().Sub(dt).Seconds() > float64(timeout) {
		w.WriteHeader(403)
		w.Write([]byte("timeout! the request time is long ago, please try again"))
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("the message read failed, errmsg: " + err.Error()))
		return
	}
	defer r.Body.Close()

	secret := os.Getenv("ROBOT_SECRET")

	if sign != signature(datetime, secret, string(content)) {
		w.WriteHeader(403)
		w.Write([]byte("bad request!"))
		return
	}

	//botId :=
}

/*func (bot *Robot) listenOutgoing(w http.ResponseWriter, r *http.Request) {

}*/