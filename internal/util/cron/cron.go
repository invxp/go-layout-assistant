package cron

import (
	"github.com/robfig/cron"
	"log"
)

/*
工具包
定时任务
*/

type Cron struct {
	c *cron.Cron
}

//New 新建一个Cron
func New() *Cron {
	return &Cron{cron.New()}
}

//MustAdd 新增一个定时任务
func (c *Cron) MustAdd(spec string, cmd func()) {
	if err := c.c.AddFunc(spec, cmd); err != nil {
		log.Panic(err)
	}
}

//Start 开始所有任务(可重入)
func (c *Cron) Start() {
	c.c.Start()
}

//Stop 停止所有任务(可重入)
func (c *Cron) Stop() {
	c.c.Stop()
}
