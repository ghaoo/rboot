package dingtalk

type At struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

// Message 基础消息结构
type Message struct {
	MsgType string `json:"msgtype"`
	At      At     `json:"at,omitempty"`

	Text       *Text       `json:"text,omitempty"`
	Markdown   *Markdown   `json:"markdown,omitempty"`
	Link       *Link       `json:"link,omitempty"`
	ActionCard *ActionCard `json:"actionCard,omitempty"`
	FeedCard   *FeedCard   `json:"feedCard,omitempty"`
}

func (msg *Message) AtAll() {
	msg.At.IsAtAll = true
}

func (msg *Message) AtMobiles(mobiles []string) {
	msg.At.AtMobiles = mobiles
}

// Text text类型
type Text struct {
	Content string `json:"content,omitempty"`
}

type Markdown struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// ActionCard 整体跳转actionCard类型
type ActionCard struct {
	Title          string            `json:"title,omitempty"`
	Text           string            `json:"text,omitempty"`
	SingleTitle    string            `json:"singleTitle,omitempty"`
	SingleURL      string            `json:"singleURL,omitempty"`
	HideAvatar     string            `json:"hideAvatar,omitempty"`
	BtnOrientation string            `json:"btnOrientation,omitempty"`
	Btns           map[string]string `json:"btns,omitempty"`
}

// Link feedCard类型 links 参数
type Link struct {
	Title      string `json:"title,omitempty"`
	Text       string `json:"title,omitempty"`
	MessageURL string `json:"messageURL,omitempty"`
	PicURL     string `json:"picURL,omitempty"`
}

// FeedCard feedCard类型
type FeedCard struct {
	Links []Link `json:"links,omitempty"`
}

// NewTextMessage 新建 text 类型消息
func NewTextMessage(content string) *Message {
	return &Message{
		MsgType: "text",
		Text:    &Text{Content: content},
	}
}

// NewMarkdownMessage 新建 markdown 类型消息
func NewMarkdownMessage(title, text string) *Message {
	return &Message{
		MsgType:  "markdown",
		Markdown: &Markdown{Title: title, Text: text},
	}
}

// NewLinkMessage 新建 link 类型消息
func NewLinkMessage(title, text, msgUrl, picUrl string) *Message {
	return &Message{
		MsgType: "link",
		Link: &Link{
			Title:      title,
			Text:       text,
			MessageURL: msgUrl,
			PicURL:     picUrl,
		},
	}
}

// NewActionCardMessage 新建 actionCard 类型消息
func NewActionCardMessage(card *ActionCard) *Message {
	return &Message{
		MsgType:    "actionCard",
		ActionCard: card,
	}
}

// 新建 feedCard 类型
func NewFeedCardMessage(links []Link) *Message {
	return &Message{
		MsgType:  "feedCard",
		FeedCard: &FeedCard{Links: links},
	}
}

// 新建 empty 类型
func NewEmptyMessage() *Message {
	return &Message{MsgType: "empty"}
}
