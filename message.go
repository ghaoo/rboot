package rboot

import (
	"bufio"
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/textproto"
	"strings"
)

type Message struct {
	To     []string  // 接收者
	From   string    // 发送者
	Header Header    // 头信息
	Body   io.Reader // 消息主体
}

// ReadMessage 读取消息
func ReadMessage(r io.Reader) (msg *Message, err error) {
	tp := textproto.NewReader(bufio.NewReader(r))

	hdr, err := tp.ReadMIMEHeader()
	msg = &Message{
		Header: Header(hdr),
		Body:   tp.R,
	}

	return msg, err
}

func NewMessage(content string) *Message {
	return &Message{
		Header: Header{},
		Body:   strings.NewReader(content),
	}
}

func NewMessageWithReader(body io.Reader) *Message {
	return &Message{
		Header: Header{},
		Body:   body,
	}
}

func (m *Message) String() string {
	content, err := ioutil.ReadAll(m.Body)
	if err != nil {
		logrus.Error(err)
	}

	m.Body = bytes.NewBuffer(content)

	return string(content)
}

// MsgType 消息类型，类型名称会转换成小写
func (m *Message) MsgType() string {
	return strings.ToLower(m.Header.Get("MsgType"))
}

// File 获取消息中的附件存放位置
func (m *Message) FilePath() string {
	return m.Header.Get("File")
}

// SetAttachment 为消息设置附件，多个附件以
func (m *Message) AddFile(contentType, filepath string) {
	m.Header.Set("MsgType", contentType)
	m.Header.Set("File", filepath)
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
