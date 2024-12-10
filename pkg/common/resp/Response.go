package resp

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(code int, msg string, data interface{}) Response {
	return Response{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}
