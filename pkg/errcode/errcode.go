package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	code    int
	msg     string
	details []string
}

// key-->code value-->msg
var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在", code))
	}

	// HttpError 中存储了 msg，codes 中也存储了 msg
	codes[code] = msg

	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.GetCode(), e.GetMsg())
}

func (e *Error) GetCode() int {
	return e.code
}

func (e *Error) GetMsg() string {
	return e.msg
}

func (e *Error) GetDetails() []string {
	return e.details
}

// 增加说明
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e // 获取切片结构体
	newError.details = []string{}
	for _, d := range details {
		newError.details = append(newError.details, d)
	}

	return &newError
}

// 将内部定义的错误码转换成 HTTP 协议的状态码
func (e *Error) StatusCode() int {
	switch e.GetCode() {
	case Success.GetCode():
		return http.StatusOK
	case ServerError.GetCode():
		return http.StatusInternalServerError
	case InvalidParams.GetCode():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.GetCode():
		fallthrough
	case UnauthorizedTokenError.GetCode():
		fallthrough
	case UnauthorizedTokenGenerate.GetCode():
		fallthrough
	case UnauthorizedTokenTimeout.GetCode():
		return http.StatusUnauthorized
	case TooManyRequests.GetCode():
		return http.StatusTooManyRequests
	}

	return http.StatusInternalServerError
}
