package server

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
	"time"
)

var ginLogger = logrus.WithField("server", "gin")
var wechatMsgLogger = logrus.WithField("server", "wechatMsg")

func ginRequestLog(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	endTime := time.Now()
	latencyTime := endTime.Sub(startTime)
	reqMethod := c.Request.Method
	reqURI := c.Request.RequestURI
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	ginLogger.Infof(
		"| %3d | %13v | %15s | %s | %s |",
		statusCode, latencyTime, reqMethod, reqURI, clientIP,
	)
}

func wechatMsgLog(m *Message) {
	startTime := time.Now()
	m.Next()
	endTime := time.Now()
	latencyTime := endTime.Sub(startTime)

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(m.Reply)
	reply := buffer.String()

	var msgType, key, id string
	if m.MsgType == message.MsgTypeEvent {
		msgType = string(m.Event)
		key = m.EventKey
	} else {
		msgType = string(m.MsgType)
		key = m.Content
	}
	if m.UnionID != "" {
		id = m.UnionID
	} else {
		id = string(m.FromUserName)
	}
	wechatMsgLogger.Infof(
		"| %10v | %v | %v | %s | %s |",
		latencyTime, msgType, key, id, reply)
}
