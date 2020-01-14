package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/ghaoo/wechat"
)

type groupUser struct {
	UserName string
}

// 创建群聊
func (w *wx) createChatRoom(users []groupUser) error {

	url := fmt.Sprintf(`%s/webwxcreatechatroom?r=%s&%s&lang=zh_CN`, w.client.BaseURL, strconv.FormatInt(time.Now().Unix(), 10), w.client.PassTicketKV())

	data := map[string]interface{}{
		`BaseRequest`: w.client.BaseRequest,
		`MemberCount`: len(users),
		`MemberList`:  users,
		`Topic`:       ``,
	}

	bs, _ := json.Marshal(data)

	var resp sdk.Response

	err := w.client.Execute(url, bytes.NewReader(bs), &resp)
	if err != nil {
		return err
	}
	if resp.IsSuccess() {
		return nil
	}
	return resp.Error()
}
