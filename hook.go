package rboot

const (
	HOOK_BEFORE_INCOMING int = iota + 1
	HOOK_AFTER_OUTGOING
)

var hookType = map[int]string{
	HOOK_BEFORE_INCOMING: "传入消息前",
	HOOK_AFTER_OUTGOING:  "消息传出后",
}

var AllHookTypes = []int{
	HOOK_BEFORE_INCOMING,
	HOOK_AFTER_OUTGOING,
}

// Hook interface
type Hook interface {
	Types() []int
	Fire(*Message) error
}

type Hooks map[int][]Hook

// Add 新增钩子
func (hooks Hooks) Add(hook Hook) {
	for _, typ := range hook.Types() {
		hooks[typ] = append(hooks[typ], hook)
	}
}

// Fire 执行
func (hooks Hooks) Fire(typ int, bot *Robot, msg *Message) error {
	for _, hook := range hooks[typ] {
		if err := hook.Fire(msg); err != nil {
			return err
		}
	}
	return nil
}
