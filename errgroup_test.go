package goerror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/prasannavl/goerror"
)

func TestMultiError(t *testing.T) {
	errMsgs := []string{
		"error 1",
		"error 2",
		"error 3",
	}

	errs := make([]error, 0, len(errMsgs))
	for _, e := range errMsgs {
		errs = append(errs, errors.New(e))
	}

	err := goerror.GroupFrom(errs)
	fmt.Println(err)

}

func TestMultiErrorNil(t *testing.T) {
	err := goerror.GroupFrom(nil)
	fmt.Println(err)
}

func TestMultiErrorNilMsg(t *testing.T) {
	errMsgs := []string{
		"error 1",
		"error 2",
	}

	errs := make([]error, 0, len(errMsgs))
	for _, e := range errMsgs {
		errs = append(errs, errors.New(e))
	}

	err := goerror.ErrGroup{errs, nil, nil, nil}
	fmt.Println(err.Error())
}
