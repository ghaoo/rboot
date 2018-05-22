package bot

// 消息元数据
type Message struct {
	Room         string
	FromUserID   string
	FromUserName string
	ToUserID     string
	ToUserName   string
	Message      string
	AtMe         bool
}
