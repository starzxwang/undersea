package api

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) *Resp {
	return &Resp{
		Data: data,
	}
}

func Failed(code int, message string, data interface{}) *Resp {
	return &Resp{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
