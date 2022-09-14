package application

//defaultConfig 默认配置
var defaultConfig = &Config{
	Log: struct {
		FilePath             string
		MaxAgeHours          uint
		MaxRotationMegabytes uint
	}{FilePath: "logs/application.log", MaxAgeHours: 3 * 24, MaxRotationMegabytes: 1 * 1024},
}

type Config struct {
	//Log 日志
	Log struct {
		FilePath             string //日志存储路径
		MaxAgeHours          uint   //日志轮转最大生命周期(小时)
		MaxRotationMegabytes uint   //日志最大文件大小(MB)
	}
}
