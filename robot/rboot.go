package main

import (
	"github.com/ghaoo/rboot"
	_ "github.com/ghaoo/rboot/adapter"
	_ "github.com/ghaoo/rboot/robot/plugins"
	"github.com/sirupsen/logrus"
)

func main() {

	bot := rboot.New()

	sm := rboot.NewMsgHook(bot)

	bot.AddHook(sm)

	bot.Go()
}

func init() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	/*logfile := filepath.Join(os.Getenv("CACHE_PATH"), "log/go.log")

	writer, _ := rotatelogs.New(
		logfile+".%Y%m%d",
		rotatelogs.WithLinkName(logfile),
		rotatelogs.WithRotationCount(1000),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(writer)

	logrus.SetLevel(logrus.TraceLevel)*/
}
