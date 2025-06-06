package pong

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/trancecho/Wechat_OfficialAccount/server"
	"sync"
)

func init() {
	instance = &pong{}
	server.RegisterModule(instance)
}

var instance *pong

type pong struct{}

func (m *pong) GetModuleInfo() server.ModuleInfo {
	return server.ModuleInfo{
		ID:       server.NewModuleID("x1anyu", "pong"),
		Instance: instance,
	}
}

func (p *pong) Init() {}

func (p *pong) PostInit() {}

func (p *pong) Serve(server *server.Server) {
	server.HttpEngine.Get("/ping", handlePingPong)
}

func (p *pong) Start(server *server.Server) {}

func (p *pong) Stop(server *server.Server, wg *sync.WaitGroup) {
	defer wg.Done()
}

func handlePingPong(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg":        "pong",
		"User-Agent": c.GetHeader("User-Agent"),
	})
}
