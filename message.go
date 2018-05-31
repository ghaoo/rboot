package rboot

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net/textproto"
	"time"
)

type Message struct {
	Header Header
	Body   io.Reader
}

func (msg *Message) Read() ([]byte, error) {
	return ioutil.ReadAll(msg.Body)
}

// 读消息
func ReadMessage(r io.Reader) (msg *Message, err error) {
	tp := textproto.NewReader(bufio.NewReader(r))

	hdr, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}

	return &Message{
		Header: Header(hdr),
		Body:   tp.R,
	}, nil
}

// 消息头
type Header map[string][]string

func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

var ErrHeaderNotPresent = errors.New("rboot: header not in message")

// 获取时间
func (h Header) Date() (time.Time, error) {
	hdr := h.Get("Time")

	if hdr == "" {
		return time.Now(), ErrHeaderNotPresent
	}
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02 15:04:05", hdr, loc)
}

// 获取消息类型
func (h Header) ContentType() (string, error) {
	hdr := h.Get("content-type")

	if hdr == "" {
		return ``, ErrHeaderNotPresent
	}

	return hdr, nil
}

// 获取消息来源
func (h Header) From() (string, error) {
	hdr := h.Get("From")

	if hdr == "" {
		return ``, ErrHeaderNotPresent
	}

	return hdr, nil
}

// 获取发送地址
func (h Header) To() (string, error) {
	hdr := h.Get("To")

	if hdr == "" {
		return ``, ErrHeaderNotPresent
	}

	return hdr, nil
}
