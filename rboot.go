package rboot

type Rboot struct {
	plugs map[string]Plugin
}

// 适配连接器
type Connecter interface {
	Incoming() chan Message
	Outgoing() chan Message
}
