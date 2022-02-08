package errs

import (
	"log"
	"net/http"
)

type Errors struct {
	code    int
	msg     string
	details []string
}

var errorCode = map[int]string{}

func NewErrors(code int, msg string) *Errors {

	if _, ok := errorCode[code]; ok {
		log.Fatalf("错误码：%d 已经存在，请更换一个", code)
	}

	errorCode[code] = msg

	return &Errors{
		code: code,
		msg:  msg,
	}
}

func (e *Errors) Code() int {
	return e.code
}

func (e *Errors) Msg() string {
	return e.msg
}

func (e *Errors) Error() string {
	return e.msg
}

func (e *Errors) Details() []string {
	return e.details
}

func (e *Errors) WithDetails(details ...string) *Errors {

	e.details = []string{}

	for _, detail := range details {
		e.details = append(e.details, detail)
	}

	return e
}

func (e *Errors) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ServerErrors.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	case NotFound.Code():
		return http.StatusNotFound
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenErrors.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code():
		return http.StatusTooManyRequests
	}

	return http.StatusInternalServerError
}
