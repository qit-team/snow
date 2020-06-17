package main

const (
	_tplConstantCommon = `package common
`

	_tplConstantErrorCode = `package errorcode

const (
	// 成功
	Success = 200

	// 参数错误
	ParamError = 400

	// 未经授权
	NotAuth = 401

	// 请求被禁止
	Forbidden = 403

	// 找不到页面
	NotFound = 404

	// 系统错误
	SystemError = 500
)

var MsgEN = map[int]string{
	Success:     "success",
	ParamError:  "param error",
	NotAuth: "not authorized",
	Forbidden:   "forbidden",
	NotFound:    "not found",
	SystemError: "system error",
}

func GetMsg(code int) string {
	if msg, ok := MsgEN[code]; ok {
		return msg
	}
	return ""
}
`

	_tplConstantLogType = `package logtype

const (
	Message = "message"
	GoPanic = "go.panic"
	HTTP    = "http"
)
`
)
