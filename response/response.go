package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yShen868/go-common-util/log"
	"github.com/yShen868/go-common-util/response/errno"
	"go.uber.org/zap"
)

// Body 与 README 约定一致：code=0 成功
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// OK 成功 JSON 响应
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:    errno.CodeOK,
		Message: "success",
		Data:    data,
	})
}

// Error 将 error 转为统一结构：业务错误 *errno.AppError 原样；其它打日志并返回内部错误码
func Error(c *gin.Context, err error) {
	if err == nil {
		OK(c, nil)
		return
	}
	var ae *errno.AppError
	if errors.As(err, &ae) {
		c.JSON(http.StatusOK, Body{Code: ae.Code, Message: ae.Message, Data: nil})
		return
	}
	log.FromContext(c.Request.Context()).Error("unhandled error", zap.Error(err))
	c.JSON(http.StatusOK, Body{
		Code:    errno.CodeInternal,
		Message: errno.ErrInternal.Message,
		Data:    nil,
	})
}
