# go-layout-assistant
# 生成go标准应用框架库

### 快速开始
```
$ go get github.com/invxp/go-layout-assistant
$ cd go-layout-assistant
$ go build
$ ./go-layout-assistant
```
### 快速生成一个全新的Go应用(类似LNMP)
1. 自带HTTPServer(gin)
2. 能够打包成RPM(CI/CD)
3. 能够注册成系统服务(CentOS-7)

### 你能学到
1. Go应用工程框架模板
2. 良好的编码风格与习惯
3. Must/Options等常见Go设计模式

### 关于main函数(应该做最少的事情)
1. 解析启动参数
2. 判断是否后台执行
3. 加载应用配置
4. 启动各种服务
5. 启动定时任务
6. 等待应用退出

### 买一送十
#### 附赠各种工具类, 方便开发
* HTTP-Client
* 配置库(TOML)
* 高效类型转换(Convert)
* 定时任务(Cron)
* 加解密/HASH算法/压缩/解压
* 后台执行应用(Daemon)
* IO/路径工具(IO)
* 日志库(支持自动轮转)
* MySQL-Client(简单封装-需逐步完善)
* Redis-Client(简单封装-需逐步完善)