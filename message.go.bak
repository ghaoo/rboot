package rboot

import "errors"

type Message struct {
	Channel   string      `json:"channel,omitempty"`   // 通道
	To        User        `json:"to"`                  // 发给的用户
	From      User        `json:"from"`                // 来源(群组或个人)
	Sender    User        `json:"sender"`              // 发送者(个人)
	Content   string      `json:"content"`             // 内容
	Broadcast bool        `json:"broadcast,omitempty"` // 广播消息
	Mate      Mate        `json:"mate,omitempty"`      // 附加信息
	Data      interface{} `json:"data,omitempty"`      // 源消息
}

func NewMessage(content string) Message {
	return Message{
		Channel:   "",
		To:        User{},
		From:      User{},
		Sender:    User{},
		Content:   content,
		Broadcast: false,
		Mate:      make(Mate),
	}
}

func (msg Message) SetTo(id, name string) {
	msg.To = User{ID: id, Name: name}
}

func (msg Message) SetFrom(id, name string) {
	msg.From = User{ID: id, Name: name}
}

func (msg Message) SetSender(id, name string) {
	msg.Sender = User{ID: id, Name: name}
}

type Mate map[string]interface{}

var defaultMate = Mate{}

func (m Mate) Has(key string) bool {
	_, ok := m[key]
	return ok
}

func (m Mate) Get(key string) interface{} {
	return m[key]
}

func (m Mate) GetString(key string) string {
	if value, ok := m[key].(string); ok {
		return value
	}
	return ""
}

func (m Mate) GetBool(key string) bool {
	if value, ok := m[key].(bool); ok {
		return value
	}
	return false
}

func (m Mate) GetInt(key string) (int, error) {
	if value, ok := m[key].(int); ok {
		return value, nil
	}
	return 0, errors.New("not int")
}

func (m Mate) Set(key string, value interface{}) {
	m[key] = value
}

func (m Mate) Del(key string) {
	delete(m, key)
}
