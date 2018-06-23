package rboot

import (
	"sync"
	"path"
)

type Plugin struct {
	Type string      // 插件类型
	Path string      // 插件路由
	Data interface{} // 插件数据
}

type plugStream struct {
	sync.RWMutex
	stream      chan Plugin
	wg          sync.WaitGroup
	hook        func(Plugin)

	Handlers    map[string]func(Plugin)
}

func newStream() *plugStream {
	return &plugStream{
		stream:      make(chan Plugin),
		Handlers:    make(map[string]func(Plugin)),
	}
}

func (ps *plugStream) init() {
	go func() {
		ps.wg.Wait()
		close(ps.stream)
	}()
}

func (ps *plugStream) merge(name string, plug chan Plugin) {
	ps.Lock()
	defer ps.Unlock()

	ps.wg.Add(1)

	go func(a chan Plugin) {
		for n := range a {
			ps.stream <- n
		}
		ps.wg.Done()
	}(plug)
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	return path.Clean(p)
}

func isPathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		return false
	}
	n := len(pattern)
	return len(path) >= n && path[0:n] == pattern
}

func findMatch(mux map[string]func(Plugin), path string) string {
	n := -1
	pattern := ""
	for m := range mux {
		if !isPathMatch(m, path) {
			continue
		}
		if len(m) > n {
			pattern = m
			n = len(m)
		}
	}
	return pattern
}

func (ps *plugStream) match(path string) string {
	return findMatch(ps.Handlers, path)
}

// 插件信息处理
func (ps *plugStream) Handle(path string, handler func(Plugin)) {
	ps.Handlers[cleanPath(path)] = handler
}

// Hook modify event on fly
func (ps *plugStream) Hook(f func(Plugin)) {
	ps.hook = f
}

func (ps *plugStream) Loop() {
	for e := range ps.stream {
		func(p Plugin) {
			ps.RLock()
			defer ps.RUnlock()
			if pattern := ps.match(p.Path); pattern != "" {
				ps.Handlers[pattern](p)
			}
		}(e)
		if ps.hook != nil {
			ps.hook(e)
		}
	}
}
