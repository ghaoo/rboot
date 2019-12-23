# Rboot

`rboot` 是一个使用 `golang` 写的，简单、高效的聊天机器人框架，易于扩展，它可以工作在不同的聊天服务上，并通过扩展脚本可实现 `聊天`、`工作助手`、`服务监控`、`警报触发` 等功能。


## 安装

```bash
go get -v github.com/ghaoo/rboot
```

### 快速开始

`rboot` 内置了 `cli` 和 `微信网页版` 的支持，`微信网页版` 使用的是 [KevinGong2013/wechat](https://github.com/KevinGong2013/wechat) 包，稍微做了修改。

#### 创建

1. 创建文件夹
2. 在文件夹下创建 `.env` 配置文件
3. 创建文件 main.go

创建文件夹和配置文件
```bash
mkdir rboot
cd rboot
touch .env
```

配置文件 .env 内容
```env
RBOOT_NAME=RBOOT
RBOOT_ALIAS=rboot
# 指定适配器 wechat cli
RBOOT_ADAPTER=wechat
# 指定储存器
RBOOT_MEMORIZER=memory
# boltdb 储存文件
BOLT_DB_FILE=.data/db/rboot.db
# 模拟大富翁游戏地图
MAP_FILE=maps/map.json
```
> 配置可自行添加，程序自动加载，使用时用 `os.Getenv()` 获取

创建 `main.go` 文件
```go
package main

import (
	_ "github.com/ghaoo/rboot/adapter"
	_ "github.com/ghaoo/rboot/memorizer"
	_ "github.com/ghaoo/rboot/scripts"

	"github.com/ghaoo/rboot"
)

func main() {
	// 创建 bot 实例
	bot := rboot.New()

    // 开始监听消息
	bot.Go()
}
```

编译并运行
```bash
# 编译
go build

# 运行
./rboot
```

查看脚本信息
```bash
!scripts // 所有已经加载的脚本及介绍
!help <script> // 查看脚本帮助信息
```

#### 聊天系统适配器 adapter

适配器编写也非常简单，只需要实现 `rboot.Adapter` 接口，并使用 `rboot.RegisterAdapter` 注册就可以了，具体实现可参考文件夹 `adapter`。

#### 脚本 script

脚本的编写非常简单，下面是一个 echo 的脚本：

当用户输入指令 `hello` 时，返回 `Hello World!`

```go
// 注册脚本
rboot.RegisterScripts(`echo`, rboot.Script{
	Action: func(ctx context.Context, bot *rboot.Robot) []rboot.Message {
		bot.SendText("Hello World!")
		return nil
    },
    Ruleset:     map[string]string{"hello": "hello"},
    Description: `返回 Hello World! `,
})
```

复杂应用请看脚本示例，脚本示例位置: `scripts` 文件夹。
- `ping`
- `timing`: 定时任务脚本
- `richman`: 大富翁模拟游戏，只支持微信，只用于学习

运行 `richman` 脚本需要配置地图文件位置 `MAP_FILE`, 默认地图文件在 `example` 文件夹下 `maps` 下

#### 储存器 memorizer

储存器默认实现了 `内存储存` 和 `bolt`, 也可自行编写，具体实现参考 `memorizer` 文件夹
指定 `储存器` 后可通过 `bot.Memory` 使用相关方法
- Save(bucket, key string, value []byte) error
- Find(bucket, key string) []byte
- FindAll(bucket string) map[string][]byte
- Update(bucket, key string, value []byte) error
- Delete(bucket, key string) error





