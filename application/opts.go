package application

type Options func(*Application)

//WithConfig 自定义配置
func WithConfig(conf *Config) Options {
	return func(application *Application) { application.conf = conf }
}
