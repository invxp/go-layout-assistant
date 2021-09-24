package application

//defaultConfig 默认配置
var defaultConfig = &Config{
	HTTP: struct {
		Enable  bool
		Address string
	}{Enable: false, Address: ":80"},

	Log: struct {
		Enable               bool
		Path                 string
		MaxAgeHours          uint
		MaxRotationMegabytes uint
	}{Enable: false, Path: "logs", MaxAgeHours: 3 * 24, MaxRotationMegabytes: 1 * 1024},

	MySQL: struct {
		Enable                bool
		Host                  string
		Username              string
		Password              string
		Database              string
		MaxConnIdleTimeMinute uint
		MaxConnLifeTimeMinute uint
		MaxOpenConnections    uint
		MaxIdleConnections    uint
	}{Enable: false, Host: "127.0.0.1:3306", Username: "user", Password: "pwd", Database: "test", MaxConnIdleTimeMinute: 10, MaxConnLifeTimeMinute: 10, MaxOpenConnections: 1000, MaxIdleConnections: 100},

	Redis: struct {
		Enable                bool
		Host                  string
		Password              string
		Database              uint
		MaxIdle               uint
		MaxActive             uint
		MaxConnTimeoutSecond  uint
		MaxConnIdleTimeMinute uint
	}{Enable: false, Host: "127.0.0.1:6379", Password: "pwd", Database: 1, MaxIdle: 100, MaxActive: 1000, MaxConnTimeoutSecond: 5, MaxConnIdleTimeMinute: 10},
}

type Config struct {
	//HTTP HTTPServer(Service)
	HTTP struct {
		Enable  bool   //是否开启
		Address string //监听地址
	}

	//Log 日志
	Log struct {
		Enable               bool   //是否开启
		Path                 string //日志存储路径
		MaxAgeHours          uint   //日志轮转最大生命周期(小时)
		MaxRotationMegabytes uint   //日志最大文件大小(MB)
	}

	//MySQL 数据库(Client)
	MySQL struct {
		Enable                bool   //是否开启
		Host                  string //主机地址
		Username              string //用户名
		Password              string //密码
		Database              string //数据库名称
		MaxConnIdleTimeMinute uint   //连接池空闲时间(分钟)
		MaxConnLifeTimeMinute uint   //连接池生命周期(分钟)
		MaxOpenConnections    uint   //最大多少个连接池
		MaxIdleConnections    uint   //最大多少个空闲连接
	}

	//Redis 缓存(Client)
	Redis struct {
		Enable                bool   //是否开启
		Host                  string //主机地址
		Password              string //密码
		Database              uint   //数据库序列号
		MaxIdle               uint   //最大多少个空闲连接
		MaxActive             uint   //最大多少个连接池
		MaxConnTimeoutSecond  uint   //多少秒连接失败算超时
		MaxConnIdleTimeMinute uint   //连接池空闲时间(分钟)
	}
}
