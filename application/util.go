package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	internal "github.com/invxp/go-layout-assistant/internal/http"
	"github.com/invxp/go-layout-assistant/internal/util/system"
	"log"
	"net/http"
	"time"
)

//logf 打印日志(如果没有启用则打到控制台)
func (application *Application) logf(format string, v ...interface{}) {
	if application.logger == nil {
		log.Printf(format, v...)
	} else {
		application.logger.Printf(format, v...)
	}
}

//httpStatics HTTP统计执行时间(中间件)
func (application *Application) httpStatics() gin.HandlerFunc {
	return func(c *gin.Context) {
		lastTime := time.Now()
		c.Next()
		application.logf("HTTP %d %s@%s: %s - LATENCY: %v, HEADERS: %v", c.Writer.Status(), c.ClientIP(), c.Request.Method, c.FullPath(), time.Since(lastTime), c.Request.Header)
	}
}

//httpAuth HTTP校验是否非法(中间件）
func (application *Application) httpAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Auth") == "FALSE" {
			c.AbortWithStatusJSON(http.StatusBadRequest, internal.Response{Code: internal.StatusCodeAuthFailed, Description: "Auth Failed", MessageID: system.UniqueID()})
		}
		c.Next()
	}
}

//crateService 随机出错(只是个测试)
func (application *Application) createService(serviceName, modPath, installDir string) error {
	if time.Now().Unix()%2 == 0 {
		return fmt.Errorf("create %s.%s-%s error", modPath, serviceName, installDir)
	}
	return nil
}
