package assistant

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/invxp/go-layout-assistant/internal/http"
	"github.com/invxp/go-layout-assistant/internal/util/convert"
	"github.com/invxp/go-layout-assistant/internal/util/system"
	"github.com/pkg/errors"
)

//httpAny Any路由演示(BindGET-URLQuery)
func (assistant *Assistant) httpAny(c *gin.Context) {
	get := &http.RequestGET{}
	err := c.ShouldBindQuery(get)
	if err != nil {
		http.Failed(c, http.StatusCodeValidationError, fmt.Sprintf("%v", errors.Wrap(err, "httpAny binding")), system.UniqueID())
	}

	assistant.logf("any router: %s", c.Param("any"))

	http.Success(c, system.UniqueID(), &http.ResponseData{Payload: convert.MustMarshall(get)})
}

//httpCron 添加一个Cron(演示用-BindPOST-JSON)
func (assistant *Assistant) httpCron(c *gin.Context) {
	post := &http.RequestPOST{}
	err := c.ShouldBindJSON(post)
	if err != nil {
		http.Failed(c, http.StatusCodeValidationError, fmt.Sprintf("%v", errors.Wrap(err, "httpCron binding")), system.UniqueID())
	}

	assistant.logf("post cron: %s.%s", post.Key, post.Value)

	assistant.AddCron(post.Key, func() { assistant.cronAlert() })
	assistant.StartCron()

	http.Success(c, system.UniqueID(), &http.ResponseData{Payload: convert.MustMarshall(post)})
}
