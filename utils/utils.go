package utils

import goerror "github.com/prasannavl/goerror"

func MakeIterator(err error) func() error {
	e := err
	return func() error {
		for {
			eLast := e
			if goerr, ok := e.(goerror.Error); ok {
				e = goerr.Cause()
				return eLast
			}
			e = nil
			return eLast
		}
	}
}

func HasMessage(err error) bool {
	if err == nil {
		return false
	}
	if goerr, ok := err.(goerror.Error); ok {
		return goerr.IsSource()
	}
	return true
}

func MakeMsgIterator(err error) func() *string {
	iter := MakeIterator(err)
	return func() *string {
		for {
			e := iter()
			if e == nil {
				return nil
			}
			if HasMessage(e) {
				errStr := e.Error()
				return &errStr
			}
		}
	}
}

func CollectMsg(err error) []string {
	if err == nil {
		return nil
	}
	s := make([]string, 0, 2)
	e := err
	for {
		if goerr, ok := err.(goerror.Error); ok {
			if goerr.IsSource() {
				s = append(s, goerr.Error())
			}
			e = goerr.Cause()
			if e == nil {
				break
			}
			continue
		}
		s = append(s, e.Error())
		break
	}
	return s
}
