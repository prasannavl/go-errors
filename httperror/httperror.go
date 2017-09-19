package httperror

import http "net/http"
import goerror "github.com/prasannavl/goerror"

type HttpError interface {
	goerror.CodedError
	Headers() http.Header
	End() bool
}

func New(code int, message string, end bool) HttpError {
	return &HttpErr{goerror.CodedErr{goerror.GoErr{&message, nil}, ErrorCode(code)}, end, nil}
}

func NewWithCause(code int, message string, cause error, end bool) HttpError {
	return &HttpErr{goerror.CodedErr{goerror.GoErr{&message, cause}, ErrorCode(code)}, end, nil}
}

func From(err error, code int, end bool) HttpError {
	if err == nil {
		return nil
	}
	code = ErrorCode(code)
	if gerr, ok := err.(HttpError); ok &&
		code == gerr.Code() && end == gerr.End() {
		return gerr
	}
	return &HttpErr{goerror.CodedErr{goerror.GoErr{nil, err}, code}, end, nil}
}

type HttpErr struct {
	goerror.CodedErr
	Stop       bool
	HeadersMap http.Header
}

func (h *HttpErr) Headers() http.Header {
	if h.HeadersMap == nil {
		h.HeadersMap = http.Header{}
	}
	return h.HeadersMap
}

func (h *HttpErr) End() bool {
	return h.Stop
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
