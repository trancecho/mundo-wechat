package server

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"sort"
	"sync"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/sirupsen/logrus"
	"github.com/trancecho/mundo-wechat/config"
)

var Version = "debug"
var Instance *Server
var logger = logrus.WithField("server", "internal")

type Server struct {
	HttpEngine   *gin.Engine
	WechatEngine *officialaccount.OfficialAccount
	MsgEngine    *MsgEngine
}

// Init 快速初始化
func Init() {
	logger.Info("wechat_mp_server version: ", Version)

	// 初始化网络服务
	logger.Info("start init gin...")
	gin.SetMode(gin.ReleaseMode)
	httpEngine := gin.New()
	httpEngine.Use(ginRequestLog, sentrygin.New(sentrygin.Options{}))

	// 初始化微信相关
	logger.Info("start init wechat...")
	wc := wechat.NewWechat()
	memoryCache := cache.NewMemory()

	cfg := &offConfig.Config{
		AppID:          config.GlobalConfig.GetString("wechat.appID"),
		AppSecret:      config.GlobalConfig.GetString("wechat.appSecret"),
		Token:          config.GlobalConfig.GetString("wechat.token"),
		EncodingAESKey: config.GlobalConfig.GetString("wechat.encodingAESKey"),
		Cache:          memoryCache,
	}
	wcOfficialAccount := wc.GetOfficialAccount(cfg)

	Instance = &Server{
		HttpEngine:   httpEngine,
		WechatEngine: wcOfficialAccount,
		MsgEngine:    NewMsgEngine(),
	}
	Instance.MsgEngine.Use(wechatMsgLog) // 注册log中间件
}

// StartService 启动服务
func StartService() {
	logger.Infof("initializing modules ...")
	for _, mi := range Modules {
		mi.Instance.Init()
	}
	for _, mi := range Modules {
		mi.Instance.PostInit()
	}
	logger.Info("all modules initialized")

	logger.Info("register modules serve functions ...")
	// 微信接入验证接口（GET）
	Instance.HttpEngine.GET("/serve", WXCheckSignature)

	// 微信接入验证：GET 用于验证，POST 用于事件消息
	Instance.HttpEngine.POST("/serve", Instance.MsgEngine.Serve)

	for _, mi := range Modules {
		mi.Instance.Serve(Instance)
	}
	logger.Info("all modules serve functions registered")

	logger.Info("starting modules tasks ...")
	for _, mi := range Modules {
		go mi.Instance.Start(Instance)
	}

	logger.Info("tasks running")
}

// Run 正式开启服务
func Run() {
	go func() {
		logger.Info("http engine starting...")
		if err := Instance.HttpEngine.Run("127.0.0.1:" + config.GlobalConfig.GetString("httpEngine.port")); err != nil {
			logger.Fatal(err)
		} else {
			logger.Info("http engine running...")
		}
	}()
}

// Stop 停止所有服务
func Stop() {
	logger.Warn("stopping ...")
	wg := sync.WaitGroup{}
	for _, mi := range Modules {
		wg.Add(1)
		mi.Instance.Stop(Instance, &wg)
	}
	wg.Wait()
	logger.Info("stopped")
	Modules = make(map[string]ModuleInfo)
}

// ========== 微信接入验证 ==========
func WXCheckSignature(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp") // 修正变量名
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	token := "123456" // 建议从 config 中读取

	log.Printf("收到验证请求: signature=%s, timestamp=%s, nonce=%s, token=%s",
		signature, timestamp, nonce, token)

	if !CheckSignature(signature, timestamp, nonce, token) {
		log.Println("签名验证失败")
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid signature"})
		return
	}

	log.Println("签名验证通过，返回 echostr:", echostr)
	c.String(200, echostr) // 必须原样返回 echostr 字符串
}

func CheckSignature(signature, timestamp, nonce, token string) bool {
	strs := []string{token, timestamp, nonce}
	sort.Strings(strs)

	h := sha1.New()
	h.Write([]byte(strs[0] + strs[1] + strs[2]))

	return hex.EncodeToString(h.Sum(nil)) == signature
}
