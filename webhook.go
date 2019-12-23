package rboot

// Header webhook请求客户端时的头信息，key-value结构
type Header map[string][]string

//
type webhook struct {
	Header Header // 头信息
	Body []Message // 消息实体
}

//

