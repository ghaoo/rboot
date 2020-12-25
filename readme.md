
# Rboot 

[![Build Status](https://travis-ci.org/ghaoo/rboot.svg?branch=master)](https://travis-ci.org/ghaoo/rboot) [![Go Report Card](https://goreportcard.com/badge/github.com/ghaoo/rboot)](https://goreportcard.com/report/github.com/ghaoo/rboot) [![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ghaoo/rboot?color=%2B&style=flat-square)](https://golang.org/) [![GoDoc](http://godoc.org/github.com/ghaoo/rboot?status.svg)](http://godoc.org/github.com/ghaoo/rboot)


`rboot` 是一个使用 `golang` 写的，简单、高效的聊天机器人框架，易于扩展，它可以工作在不同的聊天服务上，并通过扩展脚本可实现 `聊天`、`工作助手`、`服务监控`、`警报触发` 等功能。

## golang版本需求

golang `v1.13+`

## 快速创建自己的机器人

```shell script
$ go get github.com/ghaoo/rboot
$ cd $GOPATH/github.com/ghaoo/rboot/robot
$ go build
$ ./robot
```

## 关于消息转接器

消息转接器是用来监听消息的传入和传出，通过消息转接器可以将聊天客户端的消息发送到机器人，经过脚本处理后返回消息发送给客户端。

`rboot` 提供了 `命令行cli` `微信网页版` `企业微信` `钉钉` `倍洽` 聊天转接器的简单实现。

## 关于插件

`Plugin` 并没有提供太多开箱即用的插件，除了一个`help`插件，其他的需要开发者根据自己的需求去开发。

**help插件用法**：

`!help <plugin>`：查看插件帮助信息，当命令不带插件名称时会列出所有插件帮助信息，带插件名称只列出此插件的帮助信息。

### 使用golang编写插件

> 在文件夹 `robot/plugins` 下有简单的插件案例，开发者可查看插件编写方法。

### 使用其他语言编写插件

`Plugin` 不仅可以使用golang编写插件，也可以使用脚本插件来执行系统命令或使用脚本语言编写的插件文件。

脚本插件是用来解析脚本语言的`Plugin`插件，它是rboot插件的一个扩展。通过`yaml`配置文件来执行系统命令或脚本。

> 因为脚本插件是建立在`Plugin`基础之上的，每个脚本都会被注册到`Plugin`之中，所以确保插件之间名称不要重叠，否则可能先注册的插件会被后注册的插件替换！

#### 如何编写脚本插件

##### 配置

- `PLUGIN_DIR`：脚本插件配置文件存放的文件夹，若不配置默认为`scripts`

##### 快速开始

我们可以通过创建一个`yaml`文件来创建一个脚本插件，通过文件中的配置选项来实现对脚本插件的配置。比如我们创建一个 `hello.yml` 文件，它的内容如下：

```yaml
name: hello
version: 0.1.0
ruleset:
    hello: "^hello"
usage:
    hi: echo hello world and 你好
description: 脚本插件示例
command:
    -
        cmd:
            - echo hi
            - echo hello world
    -
        dir: plugins
        cmd:
            - echo 你好
```

这个插件使用的是系统命令 `echo`。它的意思是：当我们输入“hello”后，脚本会返回 `hi`，`hello world` 和 `你好` 三条信息。

配置中各个字段的含义：

配置|必须|意义
---|:---:|:---:
name|是|插件名称 
ruleset|是|规则集合
version|否|插件版本
usage|否|插件用法
description|否|插件简介
command|是|插件命令集
---|---|---
dir|否|命令执行文件夹
cmd|是|插件命令

> `command`可配置多个命令集，执行顺序为从上到下依次执行
>
> `cmd`可配置多条命令，执行顺序为从上到下依次执行

脚本插件支持`系统命令`和`脚本语言`。系统命令模式如上面的`hello.yml`，只需在文件中填写文件夹和系统命令，当你发出命令后，机器人就会从上到下依次执行。

脚本语言是建立在系统命令模式之上的执行方式，我们可以使用系统命令调用语言脚本，从而执行比较复杂的脚本。比如我们使用python输出“hello robot”。

我们的python脚本如下：

```python
#!/usr/bin/env python

print("Hello, robot! i am a python script")
```

我们的配置文件如下：

```yaml
name: pyscript
version: 0.1.0
ruleset:
    py: "^hello python"
usage:
    py: execute python script
description: python插件示例
command:
        dir: script
        cmd:
            - ./hello.py
```

当我们输入 `hello python` 时，机器人会调用 `hello.py` 脚本，脚本输出"Hello, robot! i am a python script"并通过机器人展示给我们。

> 在不同操作系统下请确认 `目录分隔符` 是否符合当前系统设置。
> `windows` 下请使用 `\`
> `unix` 下请使用 `/`

## 文档

[![Rboot](https://img.shields.io/badge/%E4%B8%AD%E6%96%87%E6%96%87%E6%A1%A3-rboot1.2.0-green)](https://www.kancloud.cn/ghaoo/rboot/1476883)

[![GoDoc](http://godoc.org/github.com/ghaoo/rboot?status.svg)](http://godoc.org/github.com/ghaoo/rboot)

### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。




