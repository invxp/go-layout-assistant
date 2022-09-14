package application

import (
	"github.com/invxp/go-layout-assistant/internal/util/log"
)

type Application struct {
	conf   *Config
	logger *log.Log
	close  chan struct{}
}

var defaultApplication = &Application{
	conf:   defaultConfig,
	logger: nil,
	close:  make(chan struct{}),
}

//New 新建一个应用, Options 动态传入配置(可选)
func New(opts ...Options) (application *Application, err error) {
	application = defaultApplication

	//遍历传入的Options方法
	for _, opt := range opts {
		opt(application)
	}

	//加载日志
	application.logger, err = log.New(
		application.conf.Log.FilePath,
		application.conf.Log.MaxAgeHours,
		application.conf.Log.MaxRotationMegabytes)

	return
}

func (application *Application) Logf(format string, v ...interface{}) {
	application.logf(format, v...)
}

//Close 退出应用
func (application *Application) Close() {
	close(application.close)
}
