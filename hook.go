package rboot

// Hook 应在 robot 启动之后，消息处理之前触发，
// 可以认为 Hook 是一个独立于聊天转接器和脚本之外的插件
type Hook interface {
	Fire(*Robot) error
}

// Hooks 是注册的钩子集合
type Hooks []Hook

// Add 为 Hooks 增加一个钩子
func (hooks Hooks) Add(hook Hook) {
	hooks = append(hooks, hook)
}

// Fire 执行所有注册的钩子
func (hooks Hooks) Fire(bot *Robot) error {
	for _, hook := range hooks {
		if err := hook.Fire(bot); err != nil {
			return err
		}
	}

	return nil
}
