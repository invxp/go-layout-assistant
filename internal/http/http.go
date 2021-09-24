package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/invxp/go-layout-assistant/internal/util/convert"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type StatusCode uint

const (
	StatusCodeOk = StatusCode(100000 + iota)

	StatusCodeValidationError
	StatusCodeNotFound
	StatusCodeAuthFailed

	StatusCodeCreateServiceError = StatusCode(200000)

	StatusCodeMySQLError = StatusCode(300000)

	StatusCodeRedisError = StatusCode(400000)
)

const (
	APIPath = "/api/"
)

const (
	RouteCron   = "cron"
	RouteCreate = "create"
)

//Response HTTP统一回复结构
type Response struct {
	Code        StatusCode    `json:"code"`
	Description string        `json:"description"`
	MessageID   string        `json:"message_id"`
	Data        *ResponseData `json:"data,omitempty"`
}

//ResponseData 具体回复的内容
type ResponseData struct {
	Payload []byte `json:"payload,omitempty"`
}

//RequestGET GET请求(URLQuery)
type RequestGET struct {
	Key   string `json:"key" binding:"required" form:"key"`
	Value string `json:"value,omitempty" form:"value"`
}

//RequestPOST POST请求(Body)
type RequestPOST struct {
	Key   string `json:"key" binding:"required"`
	Value []byte `json:"value,omitempty"`
}

//RequestPOSTCreate POST-Create请求(Body)
type RequestPOSTCreate struct {
	ServiceName string `json:"service_name" binding:"required" validate:"lowercase"`
	ModName     string `json:"mod_name" binding:"required" validate:"lowercase"`
	InstallDir  string `json:"install_dir,omitempty" validate:"omitempty,lowercase"`
}

//Failed Client请求失败(通过panic终止代码)
func Failed(c *gin.Context, code StatusCode, description, messageID string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{code, description, messageID, nil})
	panic(nil)
}

//Fatal Server处理失败
func Fatal(c *gin.Context, code StatusCode, err error, messageID string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{code, fmt.Sprintf("%v", err), messageID, nil})
	panic(nil)
}

//Success 成功
func Success(c *gin.Context, messageID string, data *ResponseData) {
	c.JSON(http.StatusOK, Response{StatusCodeOk, "OK", messageID, data})
}

//Request HTTPClient请求助手
//如果POST, 数据为序列化后的JSONBytes
func Request(method, url string, body interface{}, customHeaders map[string]string) (*Response, error) {
	var request *http.Request
	var response *http.Response
	var respBytes []byte
	var err error

	request, err = http.NewRequest(method, url, bytes.NewReader(convert.MustMarshall(body)))
	if err != nil {
		return nil, errors.Wrap(err, "create request")
	}
	request.Header.Set("Content-Type", binding.MIMEJSON)

	for k, v := range customHeaders {
		request.Header.Set(k, v)
	}

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "client do request")
	}

	defer func() {
		_ = response.Body.Close()
	}()

	respBytes, err = ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, errors.Wrap(err, "io read response")
	}

	resp := &Response{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshall error: "+convert.ByteToString(respBytes))
	}

	if resp.Code != StatusCodeOk {
		err = fmt.Errorf("%v", resp.Description)
	}

	return resp, err
}
