package httperror

import http "net/http"
import goerror "github.com/prasannavl/goerror"

func NewHttp(code int, message string, isFatal bool) goerror.CodedError {
	return goerror.NewCoded(ErrorCode(code), message, isFatal)
}

func NewHttpEx(code int, message string, cause error, isFatal bool) goerror.CodedError {
	return goerror.NewCodedEx(ErrorCode(code), message, cause, isFatal)
}

func HttpFrom(code int, err error, isFatal bool) goerror.CodedError {
	return goerror.CodedFrom(ErrorCode(code), err, isFatal)
}

func StatusCode(code int) int {
	if IsStatusCode(code) {
		return code
	}
	return http.StatusInternalServerError
}

func ErrorCode(code int) int {
	if IsErrorCode(code) {
		return code
	}
	return http.StatusInternalServerError
}

func IsStatusCode(code int) bool {
	return isInOpenRange(code, 99, 600)
}

func IsErrorCode(code int) bool {
	return isInOpenRange(code, 399, 600)
}

func IsClientErrorCode(code int) bool {
	return isInOpenRange(code, 399, 500)
}

func IsServerErrorCode(code int) bool {
	return isInOpenRange(code, 499, 600)
}

func isInOpenRange(subject int, min int, max int) bool {
	if subject < max && subject > min {
		return true
	}
	return false
}
