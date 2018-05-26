package bot

// 事件，发送接收消息、用户处理、插件系统都属于事件
type Event interface {
	SourceData() interface{}
}

type eventStream struct {
}
