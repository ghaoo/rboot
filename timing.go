package rboot

import (
	"fmt"
	"time"
)

// 运行中的计时器，包含了名称，运行次数，方法和下次执行时间
type runTiming struct {
	name string
	n    int
	f    func()
	next time.Time
}

// timer 定义了定时器的基本结构，包含了定时器的实体，结束时间，名称和需要执行的命令
type timer struct {
	t *time.Timer
	r runTiming
}

// NewTimer 新建一个定时器
func NewTimer(d time.Duration, name string, f func()) *timer {
	t := &timer{t: time.NewTimer(d)}

	t.r = runTiming{
		name: name,
		n:    0,
		f:    f,
		next: time.Now().Add(d),
	}

	go func() {
		for {
			select {
			case <-t.t.C:
				// 定时器到期，执行
				t.r.f()
			}
		}
	}()

	return t
}

// Timer 返回 *time.Timer
func (t *timer) Timer() *time.Timer {
	return t.t
}

// Stop 停止定时器定时
func (t *timer) Stop() {
	if !t.t.Stop() {
		<-t.t.C
	}
}

// Name 返回定时器名称
func (t *timer) Name() string {
	return t.r.name
}

// EndTime 返回定时器结束时间
func (t *timer) EndTime() time.Time {
	return t.r.next
}

// Reset 将计时器更改为在持续时间 d 之后过期
func (t *timer) Reset(d time.Duration) bool {
	t.r.next = time.Now().Add(d)
	return t.t.Reset(d)
}

// ticker 续断器，包含了续断器实体，名称，已经执行的次数，下一次执行时间和需要执行的脚本
type ticker struct {
	t *time.Ticker // 续断器实体
	r runTiming
}

// NewTicker 新建一个续断器实体
func NewTicker(d time.Duration, name string, f func()) *ticker {
	t := &ticker{t: time.NewTicker(d)}

	t.r = runTiming{
		name: name,
		n:    0,
		f:    f,
		next: time.Now().Add(d),
	}

	go func() {
		for {
			select {
			case <-t.t.C:

				t.r.n += 1

				t.r.next = time.Now().Add(d)

				t.r.f()
			}
		}
	}()

	return t
}

// Ticker 返回一个续断器
func (t *ticker) Ticker() *time.Ticker {
	return t.t
}

// Name 返回续断器名称
func (t *ticker) Name() string {
	return t.r.name
}

// Count 返回续断器执行的次数
func (t *ticker) Count() int {
	return t.r.n
}

// NextTime 返回下一次执行的时间
func (t *ticker) NextTime() time.Time {
	return t.r.next
}

// Stop 停止计时
func (t *ticker) Stop() {
	t.t.Stop()
}

// StrToDuration 将小时、分、秒转换为time.Duration
func StrToDuration(t int, arc string) (time.Duration, error) {

	switch arc {
	case "时", "小时", "H", "h":
		return time.Duration(t) * time.Hour, nil
	case "分", "分钟", "M", "m":
		return time.Duration(t) * time.Minute, nil
	case "秒", "秒钟", "S", "s":
		return time.Duration(t) * time.Second, nil
	}

	return 0, fmt.Errorf("非法字符串")
}
