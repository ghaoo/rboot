package rboot

const (
	HOOK_BEFORE_INCOMING int = iota + 1
	HOOK_AFTER_OUTGOING
)

var hookType = map[int]string{
	HOOK_BEFORE_INCOMING: "传入消息前",
	HOOK_AFTER_OUTGOING:  "消息传出后",
}

// Hook interface
type Hook interface {
	Types() []int
	Fire(*Robot, *Message) error
}

type MsgHooks map[int][]Hook

// Add 新增钩子
func (hooks MsgHooks) Add(hook Hook) {
	for _, typ := range hook.Types() {
		hooks[typ] = append(hooks[typ], hook)
	}
}

// Fire 执行
func (hooks MsgHooks) Fire(typ int, bot *Robot, msg *Message) error {
	for _, hook := range hooks[typ] {
		if err := hook.Fire(bot, msg); err != nil {
			return err
		}
	}
	return nil
}
