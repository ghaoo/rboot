[![Build Status](https://travis-ci.org/ghaoo/rboot.svg?branch=master)](https://travis-ci.org/ghaoo/rboot) [![Go Report Card](https://goreportcard.com/badge/github.com/ghaoo/rboot)](https://goreportcard.com/report/github.com/ghaoo/rboot) [![GoDoc](http://godoc.org/github.com/ghaoo/rboot?status.svg)](http://godoc.org/github.com/ghaoo/rboot)

# Rboot 

`rboot` 是一个使用 `golang` 写的，简单、高效的聊天机器人框架，易于扩展，它可以工作在不同的聊天服务上，并通过扩展脚本可实现 `聊天`、`工作助手`、`服务监控`、`警报触发` 等功能。

golang版本: 1.13+

## 快速创建自己的机器人

```shell script
$ go get github.com/ghaoo/rboot
$ cd $GOPATH/github.com/ghaoo/rboot/cmd
$ go build
$ ./cmd
```

## 关于消息转接器

消息转接器是用来监听消息的传入和传出，通过消息转接器可以将聊天客户端的消息发送到机器人，经过脚本处理后返回消息发送给客户端。

`rboot` 提供了 `命令行cli` `微信网页版` `企业微信` `钉钉` `倍洽` 聊天转接器的简单实现。

## 文档

中文文档：[Rboot](https://www.kancloud.cn/ghaoo/rboot/1476883)

开发请查看 [GoDoc](http://godoc.org/github.com/ghaoo/rboot)




