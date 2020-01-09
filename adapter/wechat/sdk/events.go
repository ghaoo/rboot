package sdk

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	EVENT_STOP_LOOP    int = 500
	EVENT_CLIENT_START int = iota + 1
	EVENT_NEW_MESSAGE
	EVENT_CONTACT_CHANGE
)

// Event ...
type Event struct {
	Type int
	From string
	To   string
	Data interface{}
	Time int64
}

// EventContactData 通讯录中删人 或者有人修改资料的时候
type ContactData struct {
	ChangeType int
	Contact    Contact
}

// EventMsgData 新消息
type MsgData struct {
	IsGroupMsg       bool
	IsMediaMsg       bool
	IsSendedByMySelf bool
	MsgType          int64
	AtMe             bool
	MediaURL         string
	Content          string
	FromUserName     string
	SenderUserName   string
	ToUserName       string
	OriginalMsg      map[string]interface{}
}

type Stream struct {
	sync.RWMutex
	Event chan Event
	Hook  func(Event)
}

func newStream() *Stream {
	return &Stream{
		Event: make(chan Event),
	}
}

func (es *Stream) emitContactChangeEvent(c Contact, ct int) {
	data := ContactData{
		ChangeType: ct,
		Contact:    c,
	}

	event := Event{
		Type: EVENT_CONTACT_CHANGE,
		From: `Server`,
		To:   `End`,
		Time: time.Now().Unix(),
		Data: data,
	}
	es.Event <- event
}

func (es *Stream) emitClientStartEvent() {

	event := Event{
		Type: EVENT_CLIENT_START,
		From: `Server`,
		To:   `End`,
		Time: time.Now().Unix(),
	}
	es.Event <- event
}

// Go 皮皮虾我们走
func (wechat *WeChat) Go() {
	es := wechat.Stream

	for e := range es.Event {
		switch e.Type {
		case EVENT_STOP_LOOP:
			return
		}

		if es.Hook != nil {
			es.Hook(e)
		}
	}
}

// Hook modify event on fly
func (wechat *WeChat) Hook(f func(Event)) {
	es := wechat.Stream
	es.Hook = f
}

func (wechat *WeChat) emitNewMessageEvent(m map[string]interface{}) {

	fromUserName := m[`FromUserName`].(string)
	senderUserName := fromUserName
	toUserName := m[`ToUserName`].(string)
	content := m[`Content`].(string)
	isSendedByMySelf := fromUserName == wechat.MySelf.UserName
	var groupUserName string
	if strings.HasPrefix(fromUserName, `@@`) {
		groupUserName = fromUserName
	} else if strings.HasPrefix(toUserName, `@@`) {
		groupUserName = toUserName
	}
	isGroupMsg := false
	if len(groupUserName) > 0 {
		isGroupMsg = true
		wechat.UpdateGroupIfNeeded(groupUserName)
	}
	msgType := m[`MsgType`].(float64)
	mid := m[`MsgId`].(string)

	isMediaMsg := false
	mediaURL := ``
	route := ``

	switch msgType {
	case 3: // 图片消息
		route = `webwxgetmsgimg`
	case 47: // 动画表情
		pid, _ := m[`HasProductId`].(float64)
		if pid == 0 {
			route = `webwxgetmsgimg`
		}
	case 34: // 语音消息
		route = `webwxgetvoice`
	case 43: // 视频通话消息
		route = `webwxgetvideo`
	case 49: // 分享链接
		route = `webwxgetmedia`
	case 37: // VERIFYMSG 好友验证消息
		route = `webwxverifyuser`
	case 42:
		// 分享名片
	case 10002:
		// 撤回消息
	}
	if len(route) > 0 {
		isMediaMsg = true
		mediaURL = fmt.Sprintf(`%v/%s?msgid=%v&%v`, wechat.BaseURL, route, mid, wechat.SkeyKV())
	}
	isAtMe := false
	if isGroupMsg && !isSendedByMySelf {
		atme := `@`
		if len(wechat.MySelf.DisplayName) > 0 {
			atme += wechat.MySelf.DisplayName
		} else {
			atme += wechat.MySelf.NickName
		}
		isAtMe = strings.Contains(content, atme)

		infos := strings.Split(content, `:<br/>`)
		if len(infos) != 2 {
			return
		}

		contact := wechat.ContactByUserName(infos[0])
		if contact == nil {
			wechat.ForceUpdateGroup(groupUserName)
			log.Errorf(`can't find contact info, so ignore this message %s`, m)
			return
		}

		senderUserName = contact.UserName
		content = infos[1]
	}

	data := MsgData{
		IsGroupMsg:       isGroupMsg,
		IsMediaMsg:       isMediaMsg,
		IsSendedByMySelf: isSendedByMySelf,
		MsgType:          int64(msgType),
		AtMe:             isAtMe,
		MediaURL:         mediaURL,
		Content:          content,
		FromUserName:     fromUserName,
		SenderUserName:   senderUserName,
		ToUserName:       toUserName,
		OriginalMsg:      m,
	}

	event := Event{
		Type: EVENT_NEW_MESSAGE,
		From: `Server`,
		To:   `End`,
		Time: time.Now().Unix(),
		Data: data,
	}
	wechat.Stream.Event <- event
}

func (wechat *WeChat) handleServerEvent(resp *syncMessageResponse) {

	es := wechat.Stream

	if resp.DelContactCount > 0 {
		for _, v := range resp.DelContactList {
			go es.emitContactChangeEvent(Contact{UserName: v[`UserName`].(string)}, Delete) // 已经删除的联系人这里构造一个
		}
	}

	if resp.ModContactCount > 0 {
		for _, v := range resp.ModContactList {
			contact := wechat.ContactByUserName(v[`UserName`].(string))
			if contact != nil {
				go es.emitContactChangeEvent(*contact, Modify)
			}
		}
	}

	if resp.AddMsgCount > 0 {
		for _, v := range resp.AddMsgList {
			go wechat.emitNewMessageEvent(v)
		}
	}
}
