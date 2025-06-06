package wechatPong

import (
	"github.com/getsentry/sentry-go"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/trancecho/Wechat_OfficialAccount/server"
	"sync"
)

func init() {
	instance = &wechatPong{}
	server.RegisterModule(instance)
}

var instance *wechatPong

type wechatPong struct{}

func (m *wechatPong) GetModuleInfo() server.ModuleInfo {
	return server.ModuleInfo{
		ID:       server.NewModuleID("x1anyu", "wechatPong"),
		Instance: instance,
	}
}

func (w *wechatPong) Init() {}

func (w *wechatPong) PostInit() {}

func (w *wechatPong) Serve(s *server.Server) {
	s.MsgEngine.Group("^ping$", func(msg *server.Message) {
		msg.Reply = &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("pong\n" + "Current version: " + server.Version)}
	})
	MsgText("", 1).EventClick("")
	s.MsgEngine.Group("", func(msg *server.Message) {
		msg.Reply = &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText("Wechat ping"),
		}
	}).MsgText("^ping$", 1).MsgText("^666$", 1)
}

func (w *wechatPong) Start(server *server.Server) {
	defer sentry.Recover()
}

func (m *wechatPong) Stop(server *server.Server, wg *sync.WaitGroup) {
	defer wg.Done()
}
