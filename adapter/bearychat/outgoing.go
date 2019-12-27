package bearychat

// 消息传入请求
type Request struct {
	Token       string `json:"token,omitempty"`
	TimeStamp   int64  `json:"ts,omitempty"`
	Text        string `json:"text,omitempty"`
	TriggerWord string `json:"trigger_word,omitempty"`
	Subdomain   string `json:"subdomain,omitempty"`
	ChannelName string `json:"channel_name,omitempty"`
	UserName    string `json:"user_name,omitempty"`
}

func (req Request) build() {

}

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
