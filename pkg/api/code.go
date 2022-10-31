package api

const (
	// 通用code
	CodeSuccess    = iota
	CodeException  // 暴露给前端的异常错误码
	CodeParamError // 参数有误
	CodeNotExists  // 不存在
)
