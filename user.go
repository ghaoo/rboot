package rboot

type User struct {
	ID   string                 // 用户唯一标识
	Name string                 // 用户名称
	Data map[string]interface{} // 附加信息
}
