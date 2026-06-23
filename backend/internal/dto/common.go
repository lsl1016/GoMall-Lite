package dto

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(data interface{}) Response {
	return Response{Code: 200, Message: "success", Data: data}
}

func MessageOK(message string, data interface{}) Response {
	return Response{Code: 200, Message: message, Data: data}
}

func Fail(code int, message string) Response {
	return Response{Code: code, Message: message, Data: nil}
}
