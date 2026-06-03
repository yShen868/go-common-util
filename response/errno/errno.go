package errno

import "fmt"

// AppError 业务可预期错误，对外仅暴露 Code 与 Message，不直接带底层细节
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

// New 构造业务错误
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Newf 使用格式化信息构造业务错误
func Newf(code int, format string, a ...any) *AppError {
	return &AppError{Code: code, Message: fmt.Sprintf(format, a...)}
}

// 常用预定义码（按项目约定扩展）
const (
	CodeOK           = 0
	CodeParam        = 1001
	CodeNotFound     = 1002
	CodeUnauthorized = 1003
	CodeForbidden    = 1004
	CodeInternal     = 5000
)

var (
	ErrParam        = New(CodeParam, "参数错误")
	ErrNotFound     = New(CodeNotFound, "资源不存在")
	ErrUnauthorized = New(CodeUnauthorized, "未授权")
	ErrForbidden    = New(CodeForbidden, "无权限")
	ErrInternal     = New(CodeInternal, "服务内部错误")
)
