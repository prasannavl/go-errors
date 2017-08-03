# goerror

A super tiny package for error encapsulation in idiomatic Go.

**Godoc:** https://godoc.org/github.com/prasannavl/goerror  
PS: `Read the source, Luke` - it's tiny.

## Install

`go get github.com/prasannavl/goerror`

## Sub-packages

- `httperror` - Provides `HttpError` interface and `HttpErr` type
- `utils` - Provides a set of helpers to ease iteration, and message aggregation

## Types

Provides 3 interfaces
- `Error`
- `CodedError`
- `ErrorGroup`

Provides 3 concrete types
- `Err`
- `CodedErr`
- `ErrGroup`

### Error

`Error` provides a straightforward way to wrap the `error` interface as `Cause`. The `utils` subpackage provides a way to iterate over chained errors and to collect error messages easily.

### CodedError

`CodedError` adds just one item to `Error` - `Code`. It's useful for building any errors with an `int` code. The `httperror` subpackage, for example simply builds over this with an added code validation.

### ErrorGroup

`ErrGroup` provides a clean way to aggregate errors as `[]error`. But unlike other libraries out there - doesn't provides any way to add, combine, or remove errors. It prefers to sticks to Go's idiomatic way of keeping things minimalistic. You can use a Go's own `make`, `append` and other slice function to do the same. Once done, you simply use `GroupFrom` to add your slice to create a group. It's all just slices. The only thing `ErrGroup` does is provides you an `error` interface that by default prints a nice message, and a way to retrive all the original errors.

## Notes

`HttpError` provides one additional method `Stop` that's useful to signify any middleware chain to stop processing.

Combining `Error` and `ErrorGroup` should be sufficient to handle most complex error wrapping, merging scenarios in Go without the use of other packages that add too many whistles which are, in my opinion completely unnecessary - and frankly not idiomatic Go.
