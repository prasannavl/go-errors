package httperror

import http "net/http"
import goerror "github.com/prasannavl/goerror"

type GoError interface {
	goerror.CodedError
	Stop() bool
}

func New(code int, message string, stop bool) GoError {
	return &GoErr{goerror.CodedErr{goerror.GoErr{&message, nil}, ErrorCode(code)}, stop}
}

func NewWithCause(code int, message string, cause error, stop bool) GoError {
	return &GoErr{goerror.CodedErr{goerror.GoErr{&message, cause}, ErrorCode(code)}, stop}
}

func From(err error, code int, stop bool) GoError {
	if err == nil {
		return nil
	}
	code = ErrorCode(code)
	if gerr, ok := err.(GoError); ok &&
		code == gerr.Code() && stop == gerr.Stop() {
		return gerr
	}
	return &GoErr{goerror.CodedErr{goerror.GoErr{nil, err}, code}, stop}
}

type GoErr struct {
	goerror.CodedErr
	ShouldStop bool
}

func (h *GoErr) Stop() bool {
	return h.ShouldStop
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
