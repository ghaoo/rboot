package rboot

type Message struct {
	Channel    string                 `json:"channel"`    // 通道
	To         User                   `json:"to"`         // 发给的用户
	From       User                   `json:"from"`       // 来源(群组或个人)
	Sender     User                   `json:"sender"`     // 发送者(个人)
	Content    string                 `json:"content"`    // 内容
	Broadcast  bool                   `json:"broadcast"`  // 广播消息
	Data       interface{}            `json:"data"`       // 源内容
	Mate       map[string]interface{} `json:"mate"`       // 附加信息
	Attachment []string               `json:"attachment"` // 附件位置
}

type Location struct {
	Coordinates Coordinates
}

type Coordinates struct {
	Lat  float64
	Long float64
}
