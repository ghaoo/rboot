package bearychat

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
