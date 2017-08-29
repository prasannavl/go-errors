package errutils_test

import (
	"errors"
	"testing"

	"github.com/prasannavl/goerror"
	"github.com/prasannavl/goerror/errutils"
)

func TestIterator(t *testing.T) {
	generalMsg := "some error"
	innerMsg := "inner error"
	wrapperMsg := "some wrapper error"

	errs := []goerror.Error{
		goerror.New(generalMsg),
		goerror.From(errors.New(innerMsg)),
		goerror.New(wrapperMsg),
		goerror.NewWithCause(wrapperMsg, errors.New(innerMsg)),
	}

	res := make([]string, 0, len(errs))

	for _, err := range errs {
		iter := errutils.MakeIterator(err)
		for {
			e := iter()
			if e == nil {
				break
			}
			if errutils.HasMessage(e) {
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
			break
		}
	}
}

func TestCollectMsg(t *testing.T) {
	generalMsg := "some error"
	innerMsg := "inner error"
	wrapperMsg := "some wrapper error"

	errs := []error{
		goerror.New(generalMsg),
		goerror.From(errors.New(innerMsg)),
		goerror.New(wrapperMsg),
		goerror.NewWithCause(wrapperMsg, errors.New(innerMsg)),
	}

	res := errutils.CollectMsg(errs[0], nil)
	if res[0] != generalMsg {
		t.Fatalf("CollectMsg failed, expected %q got: %q", generalMsg, res)
	}

	res = errutils.CollectMsgAll(errs, nil)

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
			break
		}
	}
}
