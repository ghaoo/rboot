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
