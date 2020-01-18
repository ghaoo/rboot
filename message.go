package rboot

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/textproto"
	"strings"

	"github.com/sirupsen/logrus"
)

// Message 表示一个消息的结构
type Message struct {
	To         string    // 消息接收者
	From       string    // 消息来源
	Sender     string    // 发送者
	Header     Header    // 头信息
	KeepHeader bool      // 如果为true则传入消息的Header在一次会话结束之前不会清除
	Body       io.Reader // 消息主体
}

// NewMessages 新建一组消息
func NewMessages(content string, to ...string) []*Message {
	msg := make([]*Message, 1)
	msg[0] = NewMessage(content, to...)

	return msg
}

// NewMessage 新建一条消息
func NewMessage(content string, to ...string) *Message {
	msg := &Message{
		Header: Header{},
		Body:   strings.NewReader(content),
	}

	if len(to) > 0 {
		msg.To = to[0]
	}

	return msg
}

// String 读取消息内容为 string
func (m *Message) String() string {
	content, err := ioutil.ReadAll(m.Body)
	if err != nil {
		logrus.Error(err)
	}

	m.Body = bytes.NewBuffer(content)

	return string(content)
}

// Bytes 读取消息内容为 []byte
func (m *Message) Bytes() []byte {
	content, err := ioutil.ReadAll(m.Body)
	if err != nil {
		logrus.Error(err)
	}

	m.Body = bytes.NewBuffer(content)

	return content
}

// SetCc 为消息设置抄送
func (m *Message) SetCc(to ...string) {
	m.Header.Set("Cc", strings.Join(to, ","))
}

// Cc 返回消息抄送信息
func (m *Message) Cc() []string {
	cc := m.Header.Get("Cc")

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
