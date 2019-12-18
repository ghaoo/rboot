package plug

// 插件
type Plug struct {
	Name string // 插件名称
	Ruleset map[string]string // 脚本规则集合
	Usage       string // 帮助信息
	Description string // 简介
}
