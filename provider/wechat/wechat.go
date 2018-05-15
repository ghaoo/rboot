package wechat

import (
	"Rboot"
	"github.com/KevinGong2013/wechat"
	"math/rand"
	"strconv"
	"time"
)

type weixin struct {
	Rboot.BasicProvider
	wechat *wechat.WeChat
}

func NewWeChat(r *Rboot.Robot) (Rboot.Provider, error) {
	wx := new(weixin)

	wx.SetRobot(r)

	wc, err := wechat.NewBot(nil)
	if err != nil {
		panic(err)
	}

	wx.wechat = wc

	return wx, nil
}

func (wx *weixin) Name() string {
	return "web wechat"
}

func clientMsgID() string {
	return strconv.FormatInt(time.Now().Unix()*1000, 10) + strconv.Itoa(rand.Intn(10000))
}

func (wx *weixin) Run() error {
	wx.wechat.Handle(`/login`, func(evt wechat.Event) {
		isSuccess := evt.Data.(int) == 1
		if isSuccess {
			log.Info(`登录成功......`)
		} else {
			log.Error(`登录失败......`)
		}
	})

	// 私聊
	wx.wechat.Handle(`/msg/solo`, func(evt wechat.Event) {
		wmsg := evt.Data.(wechat.EventMsgData)

		var msg = new(Rboot.Message)

		msg.FromUser = wmsg.FromUserName
		msg.Text = wmsg.Content

		err := wx.Receive(msg)

		if err != nil {
			log.Errorf(`receive msg error: %v`, err)
		}

	})

	// 微信群
	wx.wechat.Handle(`/msg/group`, func(evt wechat.Event) {
		wmsg := evt.Data.(wechat.EventMsgData)

		if wmsg.AtMe {

			msg := &Rboot.Message{
				ID:       clientMsgID(),
				FromUser: wmsg.FromUserName,
				Room:     wmsg.FromUserName,
				Text:     wmsg.Content,
			}

			err := wx.Receive(msg)

			if err != nil {
				log.Errorf(`receive msg error: %v`, err)
			}
		}

	})

	wx.wechat.Go()

	return nil
}

func (wx *weixin) Close() error {
	return nil
}

func (wx *weixin) Receive(msg *Rboot.Message) error {

	return wx.Robot.Receive(msg)
}

func (wx *weixin) Send(res *Rboot.Response, strings ...string) error {
	for _, str := range strings {
		err := wx.wechat.SendTextMsg(str, res.FromUser())

		if err != nil {
			log.Errorf(`Send msg error: %v`, err)
		}
	}

	return nil
}

func (wx *weixin) Reply(res *Rboot.Response, strings ...string) error {
	for _, str := range strings {
		err := wx.wechat.SendTextMsg(str, res.FromUser())

		if err != nil {
			log.Errorf(`Send msg error: %v`, err)
		}
	}

	return nil
}

func (wx *weixin) chatRoomMember(room_name string) (map[string]int, error) {

	stats := make(map[string]int)

	RoomContactList, err := wx.wechat.MembersOfGroup(room_name)
	if err != nil {
		return nil, err
	}

	man := 0
	woman := 0
	none := 0
	for _, v := range RoomContactList {

		member := wx.wechat.ContactByUserName(v.UserName)

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

func init() {
	Rboot.RegisterProvider(`wechat`, NewWeChat)
}
