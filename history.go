package rboot

type History struct {
	prev     *History
	incoming Message
	outgoing []Message
}

func (h *History) Incoming() Message {
	return h.incoming
}

func (h *History) Outgoing() []Message {
	return h.outgoing
}

// 获取上一条历史
func (h *History) Prev() *History {
	return h.prev
}

type history struct {
	root History
	len  int
}

// 清空或初始化 history
func (h *history) init() *history {
	h.root.prev = &h.root
	h.len = 0
	return h
}

// 实例化 history
func newHistory() *history {
	return new(history).init()
}

// 向 history 中插入数据
func (h *history) insert(e History) *history {
	e.prev = &h.root
	h.root = e
	h.len++
	return h
}

// 获取当前历史信息
func (h *history) current() *History {

	return &h.root
}

// 写入
func (h *history) push(in Message, out []Message) *history {

	e := History{incoming: in, outgoing: out}

	return h.insert(e)
}

// 清空历史记录
func (h *history) clear() *history {
	return h.init()
}

// 用户历史消息计记录器
type Histories map[string]*history

// 将用户操作写入历史，每位用户有一个 History 实例，当消息来源(用户)未知时将消息写入键值为 other 的 History 中，其他写入对应用户 History 中
func (hs Histories) Push(in Message, out []Message) {
	u := "other"
	if in.From.ID != "" {
		u = in.From.ID
	}

	var uh *history
	var ok bool

	if uh, ok = hs[u]; !ok {
		uh = newHistory()
	}

	uh.push(in, out)

	hs[u] = uh
}

// 用户历史信息
func (hs Histories) Current(uid string) *History {
	if _, ok := hs[uid]; !ok {
		return hs[uid].current()
	}

	return nil
}

// 用户上一条历史信息
func (hs Histories) Prev(uid string) *History {
	if _, ok := hs[uid]; ok {
		return hs[uid].root.Prev()
	}

	return nil
}

// 用户前几条历史信息
func (hs Histories) PrevN(uid string, n int) []*History {
	var uh *history
	var ok bool

	if uh, ok = hs[uid]; !ok {
		return nil
	}

	var es = make([]*History, n)

	root := &uh.root
	for i := 0; i < n; i++ {
		root = root.prev
		if root == nil {
			break
		}
		es = append(es, root)
	}

	return es
}

// 清空历史记录
func (hs Histories) Clear(uid string) {
	if uh, ok := hs[uid]; ok {
		uh.clear()
	}
}

