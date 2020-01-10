package wechat

// 微信助手

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	sdk "github.com/ghaoo/wechat"
	"github.com/sirupsen/logrus"
)

func (w *wx) Assisant(e sdk.Event) {

	if e.Type == sdk.EVENT_NEW_MESSAGE {
		msg := e.Data.(sdk.MsgData)

		//w.CreateChatRoomRequest(msg)

		if msg.IsGroupMsg {
			if msg.AtMe {
				realcontent := strings.TrimSpace(strings.Replace(msg.Content, "@"+w.client.MySelf.NickName, "", 1))
				if realcontent == "统计人数" {
					stat, err := w.ChatRoomMember(msg.FromUserName)
					if err == nil {
						ans := fmt.Sprintf("据统计群里男生 %d 人， 女生 %d 人，未知性别者 %d 人 (ó-ò) ", stat["man"], stat["woman"], stat["none"])

						w.client.SendTextMsg(ans, msg.FromUserName)
					} else {
						w.client.SendTextMsg(err.Error(), msg.FromUserName)
					}
				}
			} else if msg.MsgType == 10000 && strings.Contains(msg.Content, `加入了群聊`) {
				nn, err := search(msg.Content, `"`, `"通过`)
				if err != nil {
					logrus.Errorf(`发送欢迎消息失败 %s`, msg.Content)
				}
				w.client.SendTextMsg(`欢迎【`+nn+`】加入群聊`, msg.FromUserName)
			}
		}

		if msg.MsgType == 10002 {
			// 消息撤回
		}

		if msg.MsgType == 37 {
			w.AutoAcceptAddFirendRequest(msg)
		}
	}
}

// 统计群里男生和女生数量
func (w *wx) ChatRoomMember(room_name string) (map[string]int, error) {

	stats := make(map[string]int)

	RoomContactList, err := w.client.MembersOfGroup(room_name)
	if err != nil {
		return nil, err
	}

	man := 0
	woman := 0
	none := 0
	for _, v := range RoomContactList {

		member := w.client.ContactByUserName(v.UserName)

		if member.Sex == 1 {
			man++
		} else if member.Sex == 2 {
			woman++
		} else {
			none++
		}

	}

	stats = map[string]int{
		"woman": woman,
		"man":   man,
		"none":  none,
	}

	return stats, nil
}

func (w *wx) addFriend(data map[string]interface{}) error {
	return w.verifyUser(data, 2)
}

// AcceptAddFriend ...
func (w *wx) acceptAddFriend(data map[string]interface{}) error {
	return w.verifyUser(data, 3)
}

// 添加好友和通过好友验证
func (w *wx) verifyUser(rinfo map[string]interface{}, status int) error {

	url := fmt.Sprintf(`%s/webwxverifyuser?r=%s&%s`, w.client.BaseURL, strconv.FormatInt(time.Now().Unix(), 10), w.client.PassTicketKV())

	data := map[string]interface{}{
		`BaseRequest`:        w.client.BaseRequest,
		`Opcode`:             status,
		`VerifyUserListSize`: 1,
		`VerifyUserList`: []map[string]string{
			{
				`Value`:            rinfo["UserName"].(string),
				`VerifyUserTicket`: rinfo["Ticket"].(string),
			},
		},
		`VerifyContent`:  ``,
		`SceneListCount`: 1,
		`SceneList`:      []int{33},
		`skey`:           w.client.BaseRequest.Skey,
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

// 自动添加好友
func (w *wx) AutoAcceptAddFirendRequest(msg sdk.MsgData) {
	rInfo := msg.OriginalMsg[`RecommendInfo`].(map[string]interface{})

	err := w.acceptAddFriend(rInfo)
	if err != nil {
		logrus.Error(err)
	}
	err = w.client.SendTextMsg(`添加好友`+msg.FromUserName, `filehelper`)
	if err != nil {
		logrus.Error(err)
	}
}

func search(source, prefix, suffix string) (string, error) {

	index := strings.Index(source, prefix)
	if index == -1 {
		err := fmt.Errorf("can't find [%s] in [%s]", prefix, source)
		return ``, err
	}
	index += len(prefix)

	end := strings.Index(source[index:], suffix)
	if end == -1 {
		err := fmt.Errorf("can't find [%s] in [%s]", suffix, source)
		return ``, err
	}

	result := source[index : index+end]

	return result, nil
}

func Match(pattern, msg string) ([]string, bool) {
	r := regexp.MustCompile(pattern)

	if submatch := r.FindStringSubmatch(msg); submatch != nil {
		return submatch, true
	}

	return nil, false
}
