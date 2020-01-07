package bearychat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	baseIncomingUrl = `https://hook.bearychat.com`
)

// 发送消息
type Response struct {
	Text         string       `json:"text"`
	Notification string       `json:"notification,omitempty"`
	Markdown     bool         `json:"markdown,omitempty"`
	Channel      string       `json:"channel,omitempty"`
	User         string       `json:"user,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Title  string   `json:"title,omitempty"`
	Text   string   `json:"text,omitempty"`
	Color  string   `json:"color,omitempty"`
	Images []string `json:"images,omitempty"`
}

type Res struct {
	Code    int `json:"code,omitempty"`
	Request int `json:"request,omitempty"`
}

func sendMessage(res Response) error {
	url := os.Getenv("BEARYCHAT_WEBHOOK")
	if !strings.HasPrefix(url, baseIncomingUrl) {
		url = path.Join(baseIncomingUrl, url)
	}

	msg, _ := json.Marshal(res)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
