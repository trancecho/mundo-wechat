package templateMessage

import (
	"encoding/json"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
	"github.com/trancecho/mundo-wechat/server"
	"github.com/trancecho/mundo-wechat/utils"
	"strings"
	"sync"
)

const maxSenderNumber = 50
const globalMaxRetryTime = int64(3)

func init() {
	instance = &Module{}
	server.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

var logger *logrus.Entry

var instance *Module

type Module struct {
	template               *message.Template
	MessageQueue           chan *TemplateMessage
	MessageSenderWaitGroup sync.WaitGroup
	senderLimit            chan struct{} // 最大协程限制
}

type TemplateMessage struct {
	Message      *message.TemplateMessage
	Resend       bool
	RetriedTime  int64
	MaxRetryTime int64
}

func (m *Module) GetModuleInfo() server.ModuleInfo {
	return server.ModuleInfo{
		ID:       server.NewModuleID("x1anyu", "templateMessage"),
		Instance: instance,
	}
}

func (m *Module) Init() {
	logger.Info("Init template message sender...")

	m.MessageQueue = make(chan *TemplateMessage, maxSenderNumber)
	m.senderLimit = make(chan struct{}, maxSenderNumber)
}

func (m *Module) PostInit() {}

func (m *Module) Serve(s *server.Server) {
	m.template = message.NewTemplate(s.WechatEngine.GetContext())
	go m.registerMessageSender(m.MessageQueue)
}

func (m *Module) Start(s *server.Server) {}

func (m *Module) Stop(s *server.Server, wg *sync.WaitGroup) {
	close(m.senderLimit)
	close(m.MessageQueue)
	m.MessageSenderWaitGroup.Wait()
	wg.Done()
}

func (m *Module) PushMessage(superMessage *TemplateMessage) {
	m.MessageQueue <- superMessage
}

func (m *Module) registerMessageSender(msgChannel chan *TemplateMessage) {
	for msg := range msgChannel {
		m.senderLimit <- struct{}{}
		m.MessageSenderWaitGroup.Add(1)
		go func(t *TemplateMessage) {
			defer m.MessageSenderWaitGroup.Done()
			m.sendMessage(t)
			<-m.senderLimit
		}(msg)
	}
}

func (m *Module) sendMessage(templateMessage *TemplateMessage) {
	_, err := m.template.Send(templateMessage.Message)

	msgMarshal, _ := json.Marshal(templateMessage)
	if err != nil {
		if strings.Contains(err.Error(), "43004") {
			logger.Warn("Send templateMsg failed, errcode 43004")
			return
		}
		logger.Warnf("Send templateMsg failed, msg: %v , errMsg: %v", string(msgMarshal), err.Error())
		if templateMessage.Resend {
			if templateMessage.MaxRetryTime == -1 {
				m.MessageQueue <- templateMessage
			} else if templateMessage.MaxRetryTime != 0 && templateMessage.RetriedTime < templateMessage.MaxRetryTime {
				templateMessage.RetriedTime += 1
				m.MessageQueue <- templateMessage
			} else if templateMessage.MaxRetryTime == 0 && templateMessage.RetriedTime < globalMaxRetryTime {
				templateMessage.RetriedTime += 1
				m.MessageQueue <- templateMessage
			} else {
				logger.Warnf("Send template message failed, msg: %v , stop retry.", string(msgMarshal)) // 这部分log需要优化
			}
		}
		return
	}
	logger.Info("Send templateMsg sucess:" + string(msgMarshal))
}
