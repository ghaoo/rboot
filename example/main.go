package main

import (
	_ "Rboot/provider/wechat"
	"Rboot"
)

func main() {
	bot, err := Rboot.New()

	if err != nil {
		panic(err)
	}

	bot.Handle(
		&Rboot.Handler{Pattern: `11|22|12`, Run: testHandler},
	)

	bot.Run()
}

func testHandler(res *Rboot.Response) error {
	err := res.Send(`dasdasdasdasdas`)
	if err != nil {
		return err
	}

	return nil
}
