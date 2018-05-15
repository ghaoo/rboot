package wechat

import (
	"Rboot"
	"fmt"
	"github.com/KevinGong2013/wechat"
)

type rbootMsg struct {
	*Rboot.Message
	to string
}

func NewRbootMsg(wmsg wechat.EventMsgData) *rbootMsg {

	msg := &Rboot.Message{
		Text: wmsg.Content,
	}

	return &rbootMsg{
		msg,
		wmsg.FromUserName,
	}
}

func (am *rbootMsg) Path() string {
	return `rbootMsg`
}

func (am *rbootMsg) To() string {
	return am.to
}

func (am *rbootMsg) Content() map[string]interface{} {
	content := make(map[string]interface{}, 0)

	content["Type"] = 1
	content["Content"] = am.Text

	return content
}

func (am *rbootMsg) Description() string {
	return fmt.Sprintf(`[RbootTextMessage] %s`, am.Text)
}

func (am *rbootMsg) String() string {
	return am.Text
}
