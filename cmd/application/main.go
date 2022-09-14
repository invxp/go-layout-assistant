package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/invxp/go-layout-assistant/application"
	"github.com/invxp/go-layout-assistant/internal/util/config"
)

const (
	version = "0.0.1-alpha"
)
const (
	flagConfig = "c" //-c %{path} 自定义配置文件路径
)

var (
	configFile = flag.String(flagConfig, "config.toml", "set a config file")
)

//waitQuit 阻塞等待应用退出(ctrl+c / kill)
func waitQuit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	log.Printf("application exit")
}

//main函数(做该做的事情)
//1. 解析启动参数
//2. 加载应用配置
//3. 等待应用退出
func main() {
	flag.Parse()

	conf := &application.Config{}
	if err := config.Load(*configFile, conf); err != nil {
		log.Panic(err)
	}

	log.Println("load config success:", conf)

	srv, err := application.New(
		application.WithConfig(conf),
	)

	if err != nil {
		log.Panic(err)
	}

	//写入你的逻辑
	srv.Logf("application %s running...", version)
	srv.Close()
	//下面是等待应用退出

	waitQuit()
}
