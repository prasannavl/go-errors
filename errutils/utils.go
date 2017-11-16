package errutils

import "github.com/prasannavl/go-errors"

const DefaultIteratorLimit = 100

type ErrIterator struct {
	err   error
	limit int
	c     int
}

func (iter *ErrIterator) Next() error {
	if iter.c >= iter.limit {
		return nil
	}
	iter.c++
	eLast := iter.err
	if goerr, ok := iter.err.(errors.GoError); ok {
		iter.err = goerr.Cause()
		return eLast
	}
	iter.err = nil
	return eLast
}

func MakeIterator(err error) ErrIterator {
	return ErrIterator{err, DefaultIteratorLimit, 0}
}

func MakeIteratorLimited(err error, limit int) ErrIterator {
	return ErrIterator{err, limit, 0}
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

func MakeMsgIterator(err error) ErrMsgIterator {
	return ErrMsgIterator{MakeIterator(err)}
}

func MakeMsgIteratorLimited(err error, limit int) ErrMsgIterator {
	return ErrMsgIterator{MakeIteratorLimited(err, limit)}
}

func HasMessage(err error) bool {
	if err == nil {
		return false
	}
	if goerr, ok := err.(*errors.GoErr); ok {
		return goerr.Msg != nil
	}
	return true
}

func CollectMsgInto(err error, dest []string) []string {
	if err == nil {
		return nil
	}
	s := dest
	e := err
	for {
		if goerr, ok := e.(errors.GoError); ok {
			if HasMessage(goerr) {
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
