package utils

import (
	"fmt"
	"testing"
	"time"
)

var testTimer = NewTimer(1*time.Hour, "test", func() {
	fmt.Println("test timer")
})

func TestTiming_Timer(t *testing.T) {
	t.Log("endtime1", testTimer.EndTime())

	testTimer.Reset(1 * time.Second)

	t.Log("endtime2", testTimer.EndTime())

	t.Log(testTimer.Name())

	time.Sleep(2 * time.Second)
}

var testTicker = NewTicker(time.Second, "test", func() {
	fmt.Println("test ticker")
})

func TestTiming_Ticker(t *testing.T) {
	t.Log("name", testTicker.Name())

	t.Logf("next %d -- %v", testTicker.Count(), testTicker.NextTime())

	time.Sleep(2 * time.Second)

	t.Logf("next %d -- %v", testTicker.Count(), testTicker.NextTime())

	testTicker.Stop()

	t.Logf("next %d -- %v", testTicker.Count(), testTicker.NextTime())

	time.Sleep(time.Second)
}

func TestStrToDuration(t *testing.T) {
	d, err := StrToDuration(1, "小时")
	if err != nil {
		t.Error(err)
	}
	t.Log("小时", int64(d))

	d, err = StrToDuration(1, "分钟")
	if err != nil {
		t.Error(err)
	}
	t.Log("分钟", int64(d))

	d, err = StrToDuration(1, "秒钟")
	if err != nil {
		t.Error(err)
	}
	t.Log("秒钟", int64(d))
}
