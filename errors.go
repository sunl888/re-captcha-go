package re_captcha_go

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var reCaptchaError = map[string]string{
	"missing-input-secret":   "The secret parameter is missing",
	"missing-input-response": "The response parameter is missing",
	"bad-request":            "The request is invalid or malformed",
	"invalid-input-secret":   "The secret parameter is invalid or malformed",
	"invalid-input-response": "The response parameter is invalid or malformed",
}

type CustomError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e *CustomError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func New(statusCode int, message string, err ...error) error {
	return &CustomError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func ReCaptchaVerifyError(errMsg []string) error {
	msg, ok := reCaptchaError[errMsg[0]]
	if !ok {
		msg = "uncaught error"
	}
	return &CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    msg,
	}
}

func BadRequest(message string, err ...error) error {
	return &CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func BodyReadError(err error) error {
	return &CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    fmt.Sprintf("body 读取失败: %+v", err),
	}
}

func JsonUnmarshalError(err error) error {
	return &CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    fmt.Sprintf("json 解析失败: %+v", err),
	}
}

func JsonMarshalError(err error) error {
	return &CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    fmt.Sprintf("转换到 json 失败: %+v", err),
	}
}
