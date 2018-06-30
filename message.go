package rboot

import (
	"bufio"
	"errors"
	"io"
	"net/textproto"
	"time"
	"bytes"
	"io/ioutil"
)

type Message struct {
	Header Header
	Content string
}

func (msg Message) String() string {

	return msg.Content
}

// 新建消息
func NewMessage(r io.Reader) Message {
	tp := textproto.NewReader(bufio.NewReader(r))

	hdr, _ := tp.ReadMIMEHeader()

	content, err := ioutil.ReadAll(tp.R)

	if err != nil {
		panic(err)
	}

	return Message{
		Header: Header(hdr),
		Content: string(content),
	}
}

func NewStringMessage(str string) Message {

	return NewMessage(bytes.NewBufferString(str))
}

// 消息头
type Header map[string][]string

func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

func (h Header) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}

func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

var ErrHeaderNotPresent = errors.New("rboot: header not in message")

// 获取时间
func (h Header) Date() (time.Time, error) {
	hdr := h.Get("Date")

	if hdr == "" {
		return time.Now(), ErrHeaderNotPresent
	}
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02 15:04:05", hdr, loc)
}

// 获取消息类型
func (h Header) ContentType() string {
	hdr := h.Get("content-type")

	return hdr
}

// 获取消息来源
func (h Header) From() string {
	hdr := h.Get("From")

	return hdr
}

// 获取发送地址
func (h Header) To() string {
	hdr := h.Get("To")

	return hdr
}
