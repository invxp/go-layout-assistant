package main

import (
	"flag"
	"github.com/invxp/go-layout-assistant/assistant"
	"github.com/invxp/go-layout-assistant/internal/util/config"
	"github.com/invxp/go-layout-assistant/internal/util/daemon"
	"log"
	"os"
	"os/signal"
)

const (
	version = "0.0.1-alpha"
)
const (
	flagConfig    = "c"      //-c %{path} 自定义配置文件路径
	flagDaemon    = "d"      //-d 是否后台启动
	flagDaemonize = "daemon" //-daemon 使用者无需关注, 内部后台启动参数
)

var (
	configFile   = flag.String(flagConfig, "config.toml", "set a config file")
	enableDaemon = flag.Bool(flagDaemon, false, "run program in daemonize")
	_            = flag.Int(flagDaemonize, 0, "daemonize pid flag")
)

//waitQuit 阻塞等待应用退出(ctrl+c / kill)
func waitQuit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	log.Printf("application exit")
}

//main 程序入口(尽量少做事)
//1. 解析参数
//2. 判断是否后台执行
//3. 加载配置
//4. 启动应用
//5. 定时任务
//6. 等待退出
func main() {
	flag.Parse()

	daemon.Daemonize(*enableDaemon, flagDaemon, flagDaemonize)

	conf := &assistant.Config{}
	config.MustLoad(*configFile, conf)

	log.Println("load config success:", conf)

	srv, err := assistant.New(
		assistant.WithConfig(conf),
		assistant.WithMySQLConfig(map[string]string{"timeout": "5s"}))

	if err != nil {
		log.Panic(err)
	}

	log.Println("start program version", version)

	go func() {
		log.Println("http server close status:", srv.Serv())
	}()

	waitQuit()
}
