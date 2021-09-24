package application

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

//Application   - 面向函数编程, 任何对象皆通过New创建
//conf   	- 配置文件
//logger 	- 日志
//redis  	- 缓存/KV
//mysql  	- 数据库/SQL
//mysqlConf - 自定义数据库配置(kv)
//server    - HTTP服务器
//router    - HTTP路由(gin)
//cron      - 定时任务(cron)
//close  	- 关闭应用的channel
type Application struct {
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
func New(opts ...Options) (*Application, error) {
	application := &Application{close: make(chan struct{}), cron: cron.New()}

	//加载默认Config
	application.conf = defaultConfig

	//遍历传入的Options方法
	for _, opt := range opts {
		opt(application)
	}

	//加载日志
	if err := application.loadLogger(); err != nil {
		return nil, err
	}

	application.logf("load log success: %v", application.logger)

	//加载MySQL
	if err := application.loadMySQL(application.mysqlConf); err != nil {
		return nil, err
	}

	application.logf("load mysql success: %v", application.mysql)

	//加载Redis
	if err := application.loadRedis(); err != nil {
		return nil, err
	}

	application.logf("load redis success: %v", application.redis)

	return application, nil
}

//AddCron 添加一个定时任务
//此处与Linux不同, 格式为: 秒/分/时/月/日/年
func (application *Application) AddCron(spec string, f func()) {
	application.cron.MustAdd(spec, f)
}

//StartCron 开始定时任务
func (application *Application) StartCron() {
	application.cron.Start()
}

//Serv 开启HTTPServer
func (application *Application) Serv() error {
	if !application.conf.HTTP.Enable {
		return fmt.Errorf("http service was not enabled")
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	application.router = gin.Default()
	application.server = &http.Server{
		Addr:    application.conf.HTTP.Address,
		Handler: application.router}

	//加载日志中间件(可选)
	application.router.Use(application.httpStatics())
	//加载校验中间件(可选)
	application.router.Use(application.httpAuth())

	//默认AnyRouter
	application.router.Any(":any", application.httpAny)

	//APIRouter
	api := application.router.Group(internal.APIPath)
	{
		api.Any(internal.RouteCron, application.httpCron)
		api.POST(internal.RouteCreate, application.httpCreate)
	}

	application.logf("start http service on: %s", application.conf.HTTP.Address)

	return application.server.ListenAndServe()
}

//Close 优雅的退出应用
func (application *Application) Close(timeout time.Duration) {
	application.cron.Stop()

	close(application.close)

	if application.server != nil {
		application.logf("shutting down http server...")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := application.server.Shutdown(ctx); err != nil {
			application.logf("http server forced to shutdown: %v", err)
		} else {
			application.logf("http server close success")
		}
	}

	if application.redis != nil {
		application.logf("close redis client: %v", application.redis.Close())
	}

	if application.mysql != nil {
		application.logf("close mysql client: %v", application.mysql.Close())
	}
}
