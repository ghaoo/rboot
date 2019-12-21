package richman

import (
	"math/rand"
	"time"
)

//摇骰子，生成小于等于6的随机数
func ShakeDice() (dice int) {
	dice = RandNumber() % 6
	return dice + 1
}

//生成随机数
func RandNumber() int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(99999)
}
