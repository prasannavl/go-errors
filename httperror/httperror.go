package httperror

import http "net/http"
import goerror "github.com/prasannavl/goerror"

type GoError interface {
	goerror.CodedError
	Stop() bool
}

func New(code int, message string, stop bool) GoError {
	return &Err{goerror.CodedErr{goerror.Err{&message, nil}, ErrorCode(code)}, stop}
}

func NewWithCause(code int, message string, cause error, stop bool) GoError {
	return &Err{goerror.CodedErr{goerror.Err{&message, cause}, ErrorCode(code)}, stop}
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
	return &Err{goerror.CodedErr{goerror.Err{nil, err}, code}, stop}
}

type Err struct {
	goerror.CodedErr
	ShouldStop bool
}

func (h *Err) Stop() bool {
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
