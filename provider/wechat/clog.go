package wechat

import "github.com/num5/logger"

var log *logger.Log

func init() {
	// 初始化
	log = logger.NewLog(1000)

	// 设置log级别
	log.SetLevel("Debug")

	// 设置输出引擎
	log.SetEngine("file", `{"level":4, "spilt":"size", "filename":".logs/wechat.log", "maxsize":10}`)

	//log.DelEngine("console")

	// 设置是否输出行号
	//log.SetFuncCall(true)
}
