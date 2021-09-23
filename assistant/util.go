package assistant

import (
	"github.com/gin-gonic/gin"
	internal "github.com/invxp/go-layout-assistant/internal/http"
	"github.com/invxp/go-layout-assistant/internal/util/system"
	"log"
	"net/http"
	"time"
)

//logf 打印日志(如果没有启用则打到控制台)
func (assistant *Assistant) logf(format string, v ...interface{}) {
	if assistant.logger == nil {
		log.Printf(format, v...)
	} else {
		assistant.logger.Printf(format, v...)
	}
}

//httpStatics HTTP统计执行时间(中间件)
func (assistant *Assistant) httpStatics() gin.HandlerFunc {
	return func(c *gin.Context) {
		lastTime := time.Now()
		c.Next()
		assistant.logf("HTTP %d %s@%s: %s - LATENCY: %v, HEADERS: %v", c.Writer.Status(), c.ClientIP(), c.Request.Method, c.FullPath(), time.Since(lastTime), c.Request.Header)
	}
}

//httpAuth HTTP校验是否非法(中间件）
func (assistant *Assistant) httpAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Auth") == "FALSE" {
			c.AbortWithStatusJSON(http.StatusBadRequest, internal.Response{Code: internal.StatusCodeAuthFailed, Description: "Auth Failed", MessageID: system.UniqueID()})
		}
		c.Next()
	}
}
