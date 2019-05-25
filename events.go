package rboot

import (
	"github.com/sirupsen/logrus"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Event ...
type Event struct {
	Type string
	Path string
	From string
	To   string
	Data interface{}
	Time int64
}

// EventTimerData ...
type EventTimerData struct {
	Duration time.Duration
	Count    uint64
}

// EventTimingData ...
type EventTimingtData struct {
	Count uint64
}

type evtStream struct {
	sync.RWMutex
	srcMap      map[string]chan Event
	stream      chan Event
	wg          sync.WaitGroup
	sigStopLoop chan Event
	Handlers    map[string]func(Event)
	hook        func(Event)
	serverEvt   chan Event
}

func newEvtStream() *evtStream {
	return &evtStream{
		srcMap:      make(map[string]chan Event),
		stream:      make(chan Event),
		Handlers:    make(map[string]func(Event)),
		sigStopLoop: make(chan Event),
		serverEvt:   make(chan Event, 10),
	}
}

func (es *evtStream) init() {
	es.merge("internal", es.sigStopLoop)
	es.merge(`serverEvent`, es.serverEvt)

	go func() {
		es.wg.Wait()
		close(es.stream)
	}()
}

func (es *evtStream) merge(name string, ec chan Event) {
	es.Lock()
	defer es.Unlock()

	es.wg.Add(1)
	es.srcMap[name] = ec

	go func(a chan Event) {
		for n := range a {
			n.From = name
			es.stream <- n
		}
		es.wg.Done()
	}(ec)
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

func (es *evtStream) match(path string) string {
	return findMatch(es.Handlers, path)
}

// Stop 皮皮虾快停下
func (bot *Rboot) Stop() {
	es := bot.evtStream
	go func() {
		e := Event{
			Path: "/sig/stoploop",
		}
		es.sigStopLoop <- e
	}()
}

// Handle 处理消息，联系人，登录态 等等 所有东西
func (bot *Rboot) Handle(path string, handler func(Event)) {
	bot.evtStream.Handlers[cleanPath(path)] = handler
}

// Hook modify event on fly
func (bot *Rboot) Hook(f func(Event)) {
	es := bot.evtStream
	es.hook = f
}

// ResetHandlers remove all regeisted handler
func (bot *Rboot) ResetHandlers() {
	for Path := range bot.evtStream.Handlers {
		delete(bot.evtStream.Handlers, Path)
	}
	return
}

// NewTimerCh ...
func newTimerCh(du time.Duration) chan Event {
	t := make(chan Event)

	go func(a chan Event) {
		n := uint64(0)
		for {
			n++
			time.Sleep(du)
			e := Event{}
			e.Path = "/timer/" + du.String()
			e.Time = time.Now().Unix()
			e.Data = EventTimerData{
				Duration: du,
				Count:    n,
			}
			t <- e

		}
	}(t)
	return t
}

// AddTimer ..
func (bot *Rboot) AddTimer(du time.Duration) {
	bot.evtStream.merge(`timer`, newTimerCh(du))
}

// NewTimingCh ...
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
			logrus.Debugf(`next timing %v`, next)
			n++
			time.Sleep(next.Sub(now))
			e := Event{}
			e.Path = `/timing/` + hm
			e.Time = time.Now().Unix()
			e.Data = EventTimingtData{
				Count: n,
			}
			t <- e
		}
	}(t)
	return t
}

// AddTiming ...
func (bot *Rboot) AddTiming(hm string) {
	bot.evtStream.merge(`timing`, newTimingCh(hm))
}

func (bot *Rboot) emitMessageEvent(m Message) {

	event := Event{
		Type: `SendMessage`,
		From: `Server`,
		Path: `/msg`,
		To:   `End`,
		Time: time.Now().Unix(),
		Data: m,
	}
	bot.evtStream.serverEvt <- event
}