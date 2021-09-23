package assistant

type Options func(*Assistant)

//WithConfig 自定义配置
func WithConfig(conf *Config) Options {
	return func(assistant *Assistant) { assistant.conf = conf }
}

//WithMySQLConfig 自定义MySQL配置(KV-参考调用示例)
//https://github.com/go-sql-driver/mysql
func WithMySQLConfig(conf map[string]string) Options {
	return func(assistant *Assistant) { assistant.mysqlConf = conf }
}
