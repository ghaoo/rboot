package rboot

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/textproto"
	"strings"
)

type Message struct {
	Header Header    // 头信息
	Body   io.Reader // 消息主体
}

// ReadMessage 读取消息
func ReadMessage(r io.Reader) (msg Message, err error) {
	tp := textproto.NewReader(bufio.NewReader(r))

	hdr, err := tp.ReadMIMEHeader()
	msg = Message{
		Header: Header(hdr),
		Body:   tp.R,
	}

	return msg, err
}

func NewMessage(content string) Message {
	return Message{
		Header: Header{},
		Body:   strings.NewReader(content),
	}
}

func NewMessageWithReader(body io.Reader) Message {
	return Message{
		Header: Header{},
		Body:   body,
	}
}

func (m Message) String() string {
	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		logrus.Error(err)
	}
	return string(body)
}

// ToUser 返回接收者
// 当接收者为多个时以 , 隔开
func (m Message) ToUser() []string {
	to := m.Header.Get("to")
	return strings.Split(to, ",")

}

// FromUser 返回发送者
func (m Message) FromUser() string {
	return m.Header.Get("From")
}

// AddTo 设置接收者
func (m Message) AddTo(user ...string) {
	to := strings.Join(user, ",")
	m.Header.Add("to", to)
}

// SetFrom 设置发送者
func (m Message) SetFrom(user string) {
	m.Header.Set("from", user)
}

// SetAttachment 为消息设置附件
func (m Message) AddFile(contentType, filepath string) {
	m.Header.Add("Content-Type", contentType)
	m.Header.Add("file", filepath)
}

type Header map[string][]string

func (h Header) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}

func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

func (h Header) Del(key string) {
	textproto.MIMEHeader(h).Del(key)
}
