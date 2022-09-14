package application

import (
	"log"
)

//logf 打印日志(如果没有启用则打到控制台)
func (application *Application) logf(format string, v ...interface{}) {
	if application.logger == nil {
		log.Printf(format, v...)
	} else {
		application.logger.Printf(format, v...)
	}
}
