package rboot

import "sync"

type Element struct {
	prev *Element
	incoming Message
	outgoing []Message
}

func (e *Element) Incoming() Message {
	return e.incoming
}

func (e *Element) Outgoing() []Message {
	return e.outgoing
}

// 获取上一条历史
func (e *Element) Prev() *Element {
	return e.prev
}

type history struct {
	root Element
	len int
	m sync.Mutex
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
func (h *history) insert(e Element) *history {
	e.prev = &h.root
	h.root = e
	h.len++
	return h
}

// 获取当前历史信息
func (h *history) current() *Element {
	h.m.Lock()
	defer h.m.Unlock()

	return &h.root
}

// 写入
func (h *history) push(in Message, out []Message) *history {
	h.m.Lock()
	defer h.m.Unlock()

	e := Element{incoming:in, outgoing:out}

	return h.insert(e)
}

// 清空历史记录
func (h *history) clear() *history {
	h.m.Lock()
	defer h.m.Unlock()

	return h.init()
}

type History map[string]*history

// 将用户操作写入历史，每位用户有一个 History 实例，当消息来源(用户)未知时将消息写入键值为 other 的 History 中，其他写入对应用户 History 中
func (h History) Push(in Message, out []Message) {
	u := "other"
	if in.From.ID != "" {
		u = in.From.ID
	}

	var uh *history
	var ok bool

	if uh, ok = h[u]; !ok {
		uh = newHistory()
	}

	uh.push(in, out)
}

// 用户历史信息
func (h History) Current(uid string) *Element {
	if _, ok := h[uid]; !ok {
		return nil
	}

	return h[uid].current()
}

// 用户上一条历史信息
func (h History) Prev(uid string) *Element {
	if _, ok := h[uid]; !ok {
		return nil
	}

	return h[uid].root.Prev()
}

// 用户前几条历史信息
func (h History) History(uid string, n int) []*Element {
	var uh *history
	var ok bool

	if uh, ok = h[uid]; !ok {
		return nil
	}

	var es = make([]*Element, n)

	root := &uh.root
	for i := 0; i < n; i ++ {
		root = root.prev
		if root == nil {
			break
		}
		es = append(es, root)
	}

	return es
}

// 清空历史记录
func (h History) Clear(uid string) {
	if uh, ok := h[uid]; ok {
		uh.clear()
	}
}





