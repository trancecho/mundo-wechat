package main

import (
	"github.com/trancecho/mundo-wechat/config"
	"github.com/trancecho/mundo-wechat/server"
	"github.com/trancecho/mundo-wechat/utils"
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
