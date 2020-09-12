package rboot

const (
	HOOK_BEFORE_INCOMING int = iota + 1
	HOOK_AFTER_OUTGOING
)

var hookType = map[int]string{
	HOOK_BEFORE_INCOMING: "传入消息前",
	HOOK_AFTER_OUTGOING:  "消息传出后",
}

// Hook
type Hook interface {
	Types() []int
	Fire(*Robot, *Message) error
}

type MsgHooks map[int][]Hook

func (hooks MsgHooks) Add(hook Hook) {
	for _, typ := range hook.Types() {
		hooks[typ] = append(hooks[typ], hook)
	}
}

func (hooks MsgHooks) Fire(typ int, bot *Robot, msg *Message) error {
	for _, hook := range hooks[typ] {
		if err := hook.Fire(bot, msg); err != nil {
			return err
		}
	}
	return nil
}
