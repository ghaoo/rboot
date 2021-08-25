package rboot

import (
	"crypto/md5"
	"fmt"
	"net/textproto"
	"sort"
	"strings"
	"time"
)

// Message 表示一个消息的结构
type Message struct {
	Channel    string    `json:"channel"` // 自定义通道
	To         string    `json:"to"`      // 消息接收者
	From       string    `json:"from"`    // 消息来源
	Sender     string    `json:"sender"`  // 发送者
	Header     Header    `json:"header"`  // 头信息
	Content    string    `json:"content"` // 消息文本内容
	Time       time.Time `json:"time"`    // 消息发送时间
	KeepHeader bool      `json:"-"`       // 如果为true则传入消息的Header在一次会话结束之前不会清除
}

// NewMessages 新建一组消息
func NewMessages(content string, to ...string) []*Message {
	msg := make([]*Message, 1)
	msg[0] = NewMessage(content, to...)

	return msg
}

// NewMessage 新建一条消息，支持多个接收人
func NewMessage(content string, to ...string) *Message {
	msg := &Message{
		Header:  Header{},
		Content: content,
		Time:    time.Now(),
	}

	if len(to) > 0 {
		msg.To = to[0]
		msg.SetCc(to[1:]...)
	}

	return msg
}

// String 读取消息内容为 string
func (m *Message) String() string {
	return m.Content
}

// SetCc 为消息设置抄送
func (m *Message) SetCc(to ...string) {
	m.Header.Set("Cc", strings.Join(to, ","))
}

// Cc 返回消息抄送信息
func (m *Message) Cc() []string {
	cc := m.Header.Get("Cc")

	if len(strings.TrimSpace(cc)) == 0 {
		return nil
	}

	return strings.Split(cc, ",")
}

// Header 消息附带的头信息，键-值对
type Header map[string][]string

// Add 将键、值对添加到Header，附加到与键关联的现有值
func (h Header) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}

// Set 将key设置为单个值，它替换与key的现有值
func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

// Get 从头信息中获取与给定键关联的第一个值
func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

// GetKey 从头信息中获取与给定键关联的多个值
func (h Header) GetKey(key string) []string {
	return h[textproto.CanonicalMIMEHeaderKey(key)]
}

// Del 删除与键关联的值
func (h Header) Del(key string) {
	textproto.MIMEHeader(h).Del(key)
}

// GetMsgChannel 获取消息通道
func GetMsgChannel(from, to string) string {
	uu := []string{from, to}
	sort.Strings(uu)

	has := md5.Sum([]byte(uu[0] + uu[1]))
	return fmt.Sprintf("%x", has)
}
