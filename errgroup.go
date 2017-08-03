package goerror

import (
	"bytes"
)

type ErrorGroup interface {
	error
	Errors() []error
}

const (
	ErrGroupDefaultPrefix    = "multiple errors:\r\n"
	ErrGroupDefaultSeparator = "\r\n"
)

func NewGroup(errors []error, prefix string, suffix string, sep string) ErrorGroup {
	errs := ValidErrors(errors)
	if errs == nil {
		return nil
	}
	return &ErrGroup{errs, &prefix, &sep, &suffix}
}

func GroupFrom(errors []error) ErrorGroup {
	prefix := ErrGroupDefaultPrefix
	sep := ErrGroupDefaultSeparator
	errs := ValidErrors(errors)
	if errs == nil {
		return nil
	}
	return &ErrGroup{errs, &prefix, &sep, nil}
}

type ErrGroup struct {
	Nodes        []error
	MsgPrefix    *string
	MsgSeparator *string
	MsgSuffix    *string
}

func (e *ErrGroup) Error() string {
	var buffer bytes.Buffer
	if e.MsgPrefix != nil {
		buffer.WriteString(*e.MsgPrefix)
	}
	l1 := len(e.Nodes) - 1
	// Reset l1 to -1 so that an explicit
	// check can be avoided during the loop
	if e.MsgSeparator == nil {
		l1 = -1
	}
	for i, err := range e.Nodes {
		buffer.WriteString(err.Error())
		if i < l1 {
			buffer.WriteString(*e.MsgSeparator)
		}
	}
	if e.MsgSuffix != nil {
		buffer.WriteString(*e.MsgSuffix)
	}
	return buffer.String()
}

func (e *ErrGroup) Errors() []error {
	return e.Nodes
}

func ValidErrors(errors []error) []error {
	clen := len(errors)
	if clen < 1 {
		return nil
	}
	errs := make([]error, 0, clen)
	for _, e := range errors {
		if e != nil {
			errs = append(errs, e)
		}
	}
	rlen := len(errs)
	if rlen < 1 {
		return nil
	}
	return errs
}
