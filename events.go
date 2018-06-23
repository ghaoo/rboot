package rboot

import (
	"path"
	"sync"
	"time"
	"strings"
	"strconv"
)

type Event struct {
	Type string      // 插件类型
	Path string      // 插件路由
	Data interface{} // 插件数据
	Time int64
}

type eventStream struct {
	sync.RWMutex
	stream chan Event
	wg     sync.WaitGroup
	hook   func(Event)

	Handlers map[string]func(Event)
}

func newStream() *eventStream {
	return &eventStream{
		stream:   make(chan Event),
		Handlers: make(map[string]func(Event)),
	}
}

func (es *eventStream) init() {
	go func() {
		es.wg.Wait()
		close(es.stream)
	}()
}

func (es *eventStream) merge(name string, plug chan Event) {
	es.Lock()
	defer es.Unlock()

	es.wg.Add(1)

	go func(a chan Event) {
		for n := range a {
			es.stream <- n
		}
		es.wg.Done()
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

func findMatch(mux map[string]func(Event), path string) string {
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

func (es *eventStream) match(path string) string {
	return findMatch(es.Handlers, path)
}

// 插件信息处理
func (es *eventStream) Handle(path string, handler func(Event)) {
	es.Handlers[cleanPath(path)] = handler
}

// Hook modify event on fly
func (es *eventStream) Hook(f func(Event)) {
	es.hook = f
}

func (es *eventStream) Loop() {
	for e := range es.stream {
		func(a Event) {
			es.RLock()
			defer es.RUnlock()
			if pattern := es.match(a.Path); pattern != "" {
				es.Handlers[pattern](a)
			}
		}(e)
		if es.hook != nil {
			es.hook(e)
		}
	}
}

// Timer ...
type TimerData struct {
	Duration time.Duration
	Count    uint64
}

func newTimerCh(du time.Duration) chan Event {
	t := make(chan Event)

	go func(a chan Event) {
		n := uint64(0)
		for {
			n++
			time.After(du)
			e := Event{}
			e.Type = "timer"
			e.Path = "/timer/" + du.String()
			e.Data = TimerData{
				Duration: du,
				Count:    n,
			}
			e.Time = time.Now().Unix()
			t <- e
		}
	}(t)
	return t
}

// Timing ...
type TimingtData struct {
	Count uint64
}

func newTimingCh(hm string) chan Event {

	infos := strings.Split(hm, `:`)
	if len(infos) != 2 {
		panic(`hm incorrect`)
	}
	hour, _ := strconv.Atoi(infos[0])
	minute, _ := strconv.Atoi(infos[1])

	t := make(chan Event)

	go func(a chan Event) {
		n := uint64(0)
		for {
			now := time.Now()
			nh, nm, _ := now.Clock()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
			if n > 0 || hour > nh || (hour == nh && minute < nm) {
				next = next.Add(time.Hour * 24)
			}
			println(`next timing `, next)
			n++
			time.Sleep(next.Sub(now))
			e := Event{}
			e.Path = `/timing/` + hm
			e.Data = TimingtData{
				Count: n,
			}
			e.Time = time.Now().Unix()
			t <- e
		}
	}(t)
	return t
}

var usrEvent = make(chan Event)

func SendCustomEvent(path string, data interface{}) {
	e := Event{}
	e.Path = path
	e.Data = data
	e.Time = time.Now().Unix()
	usrEvent <- e
}

// 注册计时器
func (bot *Rboot) Timer(du time.Duration) {
	bot.es.merge(`timer`, newTimerCh(du))
}

// 注册定时发送事件
func (bot *Rboot) Timing(hm string) {
	bot.es.merge(`timing`, newTimingCh(hm))
}