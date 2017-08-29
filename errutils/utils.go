package errutils

import goerror "github.com/prasannavl/goerror"

type ErrIterator struct {
	err error
}

func (iter *ErrIterator) Next() error {
	eLast := iter.err
	if goerr, ok := iter.err.(goerror.Error); ok {
		iter.err = goerr.Cause()
		return eLast
	}
	iter.err = nil
	return eLast
}

func MakeIterator(err error) ErrIterator {
	return ErrIterator{err}
}

type ErrMsgIterator struct {
	errIter ErrIterator
}

func (iter *ErrMsgIterator) Next() *string {
	for {
		e := iter.errIter.Next()
		if e == nil {
			return nil
		}
		if HasMessage(e) {
			m := e.Error()
			return &m
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

func MakeMsgIterator(err error) ErrMsgIterator {
	return ErrMsgIterator{MakeIterator(err)}
}

func CollectMsgInto(err error, dest []string) []string {
	if err == nil {
		return nil
	}
	s := dest
	e := err
	for {
		if goerr, ok := e.(goerror.Error); ok {
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

func CollectMsg(err error) []string {
	return CollectMsgInto(err, nil)
}

func CollectAllMsgInto(errs []error, dest []string) []string {
	if len(errs) < 1 {
		return nil
	}
	s := dest
	for _, err := range errs {
		m := CollectMsgInto(err, dest)
		if m != nil {
			s = append(s, m...)
		}
	}
	return s
}

func CollectAllMsg(errs []error) []string {
	return CollectAllMsgInto(errs, nil)
}
