package rboot

type MESSAGE_TYPE_ENUM int

const (
	// 系统消息
	SYSTEM_MESSAGE_TYPE MESSAGE_TYPE_ENUM = iota
	// 广播消息
	BROADCAST_MESSAGE_TYPE
	// 心跳消息
	HEART_BEAT_MESSAGE_TYPE
	// 上线通知
	CONNECTED_MESSAGE_TYPE
	// 下线通知
	DISCONNECTED_MESSAGE_TYPE
	// 服务断开链接通知(服务端关闭)
	BREAK_MESSAGE_TYPE
)

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
	Type       MESSAGE_TYPE_ENUM      `json:"type"`       // 类型
}

type Location struct {
	Coordinates Coordinates
}

type Coordinates struct {
	Lat  float64
	Long float64
}
