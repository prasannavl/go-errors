package utils_test

import (
	"errors"
	"github.com/prasannavl/goerror"
	"github.com/prasannavl/goerror/utils"
	"testing"
)

func TestIterator(t *testing.T) {
	generalMsg := "some error"
	innerMsg := "inner error"
	wrapperMsg := "some wrapper error"

	errs := []goerror.Error{
		goerror.New(generalMsg),
		goerror.From(errors.New(innerMsg)),
		goerror.NewEx(wrapperMsg, nil),
		goerror.NewEx(wrapperMsg, errors.New(innerMsg)),
	}

	res := make([]string, 0, len(errs))

	for _, err := range errs {
		iter := utils.MakeIterator(err)
		for {
			e := iter()
			if e == nil {
				break
			}
			if utils.HasMessage(e) {
				res = append(res, e.Error())
			}
		}
	}

	expected := []string{
		generalMsg,
		innerMsg,
		wrapperMsg,
		wrapperMsg,
		innerMsg,
	}

	for i, exp := range expected {
		if res[i] != exp {
			t.Fatalf("result[%v]:%q doesn't match %q\nres: %#v\nexp: %#v", i, res[i], exp, res, expected)
		}
	}
}
