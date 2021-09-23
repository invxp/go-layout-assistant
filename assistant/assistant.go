package assistant

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	internal "github.com/invxp/go-layout-assistant/internal/http"
	"github.com/invxp/go-layout-assistant/internal/util/cron"
	"github.com/invxp/go-layout-assistant/internal/util/log"
	"github.com/invxp/go-layout-assistant/internal/util/mysql"
	"github.com/invxp/go-layout-assistant/internal/util/redis"
	"net/http"
	"time"
)

//Assistant - 面向函数编程, 任何对象皆通过New创建
//conf   	- 配置文件
//logger 	- 日志
//redis  	- 缓存/KV
//mysql  	- 数据库/SQL
//mysqlConf - 自定义数据库配置(kv)
//server    - HTTP服务器
//router    - HTTP路由(gin)
//cron      - 定时任务(cron)
//close  	- 关闭应用的channel
type Assistant struct {
	conf      *Config
	logger    *log.Log
	redis     *redis.Redis
	mysql     *mysql.MySQL
	mysqlConf map[string]string
	server    *http.Server
	router    *gin.Engine
	cron      *cron.Cron
	close     chan struct{}
}

//New 新建一个应用, Options 动态传入配置(可选)
func New(opts ...Options) (*Assistant, error) {
	assistant := &Assistant{close: make(chan struct{}), cron: cron.New()}

	//加载默认Config
	assistant.conf = defaultConfig

	//遍历传入的Options方法
	for _, opt := range opts {
		opt(assistant)
	}

	//加载日志
	if err := assistant.loadLogger(); err != nil {
		return nil, err
	}

	assistant.logf("load log success: %v", assistant.logger)

	//加载MySQL
	if err := assistant.loadMySQL(assistant.mysqlConf); err != nil {
		return nil, err
	}

	assistant.logf("load mysql success: %v", assistant.mysql)

	//加载Redis
	if err := assistant.loadRedis(); err != nil {
		return nil, err
	}

	assistant.logf("load redis success: %v", assistant.redis)

	return assistant, nil
}

//AddCron 添加一个定时任务
//此处与Linux不同, 格式为: 秒/分/时/月/日/年
func (assistant *Assistant) AddCron(spec string, f func()) {
	assistant.cron.MustAdd(spec, f)
}

//StartCron 开始定时任务
func (assistant *Assistant) StartCron() {
	assistant.cron.Start()
}

//Serv 开启HTTPServer
func (assistant *Assistant) Serv() error {
	if !assistant.conf.HTTP.Enable {
		return fmt.Errorf("http service was not enabled")
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	assistant.router = gin.Default()
	assistant.server = &http.Server{
		Addr:    assistant.conf.HTTP.Address,
		Handler: assistant.router}

	//加载日志中间件(可选)
	assistant.router.Use(assistant.httpStatics())
	//加载校验中间件(可选)
	assistant.router.Use(assistant.httpAuth())

	//默认AnyRouter
	assistant.router.Any(":any", assistant.httpAny)

	//APIRouter
	api := assistant.router.Group(internal.APIPath)
	{
		api.Any(internal.RouteCron, assistant.httpCron)
	}

	assistant.logf("start http service on: %s", assistant.conf.HTTP.Address)

	return assistant.server.ListenAndServe()
}

//Close 优雅的退出应用
func (assistant *Assistant) Close(timeout time.Duration) {
	assistant.cron.Stop()

	close(assistant.close)

	if assistant.server != nil {
		assistant.logf("shutting down http server...")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := assistant.server.Shutdown(ctx); err != nil {
			assistant.logf("http server forced to shutdown: %v", err)
		} else {
			assistant.logf("http server close success")
		}
	}

	if assistant.redis != nil {
		assistant.logf("close redis client: %v", assistant.redis.Close())
	}

	if assistant.mysql != nil {
		assistant.logf("close mysql client: %v", assistant.mysql.Close())
	}
}
