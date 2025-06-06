package newUser

import (
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/trancecho/mundo-wechat/server"
	"sync"
)

func init() {
	instance = &newUser{}
	server.RegisterModule(instance)
}

var instance *newUser

type newUser struct {
	WelcomeMessage string
}

func (newUser) GetModuleInfo() server.ModuleInfo {
	return server.ModuleInfo{
		ID:       server.NewModuleID("x1anyu", "newUser"),
		Instance: instance,
	}
}

func (newUser) Init() {}

func (newUser) PostInit() {}

func (newUser) Serve(s *server.Server) {
	s.MsgEngine.EventSubscribe(0, func(msg *server.Message) {
		text := message.NewText("你好！欢迎来到未央学社！")
		msg.Reply = &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})
}

func (newUser) Start(server *server.Server) {}

func (newUser) Stop(server *server.Server, wg *sync.WaitGroup) {
	defer wg.Done()
}
