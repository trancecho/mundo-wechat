package main

import (
	"github.com/trancecho/Wechat_OfficialAccount/config"
	"github.com/trancecho/Wechat_OfficialAccount/module/newUser"         // 新关注的欢迎模块
	"github.com/trancecho/Wechat_OfficialAccount/module/pong"            // gin的ping-pong模块
	"github.com/trancecho/Wechat_OfficialAccount/module/templateMessage" // 集中管理发送模板消息的模组
	"github.com/trancecho/Wechat_OfficialAccount/module/wechatPong"      // 微信ping-pong模块
	"github.com/trancecho/Wechat_OfficialAccount/server"
	"github.com/trancecho/Wechat_OfficialAccount/utils"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	utils.WriteLogToFS()
	config.Init()
}

func main() {
	server.Init()

	server.StartService()

	server.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	server.Stop()
}
