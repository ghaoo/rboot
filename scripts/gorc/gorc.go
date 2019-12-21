package gorc

import (
	"github.com/ghaoo/rboot"
	"os"
	goctx "context"
	"fmt"
)

var (
	root string // 下载文件存放位置
	thread int64 = 5 // 并发线程数
	manual = false // 是否手动指定线程数
	attempt = 3 // 下载失败重试次数
	blockSize int64 = 1 // 临时文件大小, 1是16m, 2是32m, 以此类推
)

func setup(ctx goctx.Context, bot *rboot.Robot) []rboot.Message {

	fmt.Println(bot.Match)
	switch bot.MatchRule {
	case `download`:
		return download(bot.Match[1])
	}

	return nil
}

func download(url string) []rboot.Message {
	err := Download(url)

	if err != nil {
		return []rboot.Message{{Content: "下载文件错误: " + err.Error()}}
	}

	return []rboot.Message{{Content: "开始下载"}}
}

func init() {
	root, _ = os.Getwd()

	root += "/download"

	if os.Getenv("DOWNLOAD_ROOT") != "" {
		root = os.Getenv("DOWNLOAD_ROOT")
	}

	rboot.RegisterScripts(`gorc`, rboot.Script{
		Action: setup,
		Ruleset: map[string]string{
			`download`: `^![down|download|下载]+ (.+)`,
		},
		Usage:       "!下载 <url>: 下载文件",
		Description: `多线程下载文件。`,
	})
}

