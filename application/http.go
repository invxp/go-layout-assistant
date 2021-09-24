package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/invxp/go-layout-assistant/internal/http"
	"github.com/invxp/go-layout-assistant/internal/util/convert"
	"github.com/invxp/go-layout-assistant/internal/util/system"
	"github.com/pkg/errors"
)

//httpAny Any路由演示(BindGET-URLQuery)
func (application *Application) httpAny(c *gin.Context) {
	get := &http.RequestGET{}
	err := c.ShouldBindQuery(get)
	if err != nil {
		http.Failed(c, http.StatusCodeValidationError, fmt.Sprintf("%v", errors.Wrap(err, "httpAny binding")), system.UniqueID())
	}

	application.logf("any router: %s", c.Param("any"))

	http.Success(c, system.UniqueID(), &http.ResponseData{Payload: convert.MustMarshall(get)})
}

//httpCron 添加一个Cron(演示用-BindPOST-JSON)
func (application *Application) httpCron(c *gin.Context) {
	post := &http.RequestPOST{}
	err := c.ShouldBindJSON(post)
	if err != nil {
		http.Failed(c, http.StatusCodeValidationError, fmt.Sprintf("%v", errors.Wrap(err, "httpCron binding")), system.UniqueID())
	}

	application.logf("post cron: %s.%s", post.Key, post.Value)

	application.AddCron(post.Key, func() { application.cronAlert() })
	application.StartCron()

	http.Success(c, system.UniqueID(), &http.ResponseData{Payload: convert.MustMarshall(post)})
}

//httpCreate 创建一个新的应用框架(BindPOST-JSON)
func (application *Application) httpCreate(c *gin.Context) {
	create := &http.RequestPOSTCreate{}
	err := c.ShouldBindJSON(create)
	if err != nil {
		http.Failed(c, http.StatusCodeValidationError, fmt.Sprintf("%v", errors.Wrap(err, "httpCreate binding")), system.UniqueID())
	}

	if err = validator.New().Struct(create); err != nil {
		http.Failed(c, http.StatusCodeValidationError, fmt.Sprintf("%v", errors.Wrap(err, "httpCreate validate")), system.UniqueID())
	}

	application.logf("create new application: %s, mod: %s, dir: %s", create.ServiceName, create.ModName, create.InstallDir)

	//random error
	if err = application.createService(create.ServiceName, create.ModName, create.InstallDir); err != nil {
		http.Fatal(c, http.StatusCodeCreateServiceError, err, system.UniqueID())
	}

	http.Success(c, system.UniqueID(), &http.ResponseData{Payload: convert.MustMarshall(create)})
}
