# goerror

A super tiny package for error encapsulation in idiomatic Go. A better and lighter alternative to [`pkg/errors`](https://github.com/pkg/errors).  
(See below for comparison to `pkg/errors`)

**Godoc:** https://godoc.org/github.com/prasannavl/goerror  
PS: `Read the source, Luke` - it's tiny.

## Get

`go get -u github.com/prasannavl/goerror`

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

## Example

```go

func TestErrors() {
    // New errors

    err := goerror.New("some error")
    // Error with codes
    codedErr := goerror.NewCoded(42, "compute overwhelmed. crashing")
    // HttpErrors are essentially, coded errors,
    // with http code validation
    httpErr := httperror.New(400, "naughty you! that's not valid", true)

    // Wrap existing existing errors
    wrapped1 := goerror.From(errors.New("some other error"))
    wrapped2 := goerror.From(wrapped1)
    wrapper3 := goerror.From("wrapped as different msg", wrapped2)

    // Intuitively messages:
    fmt.Println(wrapped1)
    // Output 
    // some other error

    // Same output
    fmt.Println(wrapped2)
    // Output
    // some other error

    fmt.Println(wrapped3)
    // Output
    // wrapped as different msg

    // You can get the real err by accessing `Cause`. Or more intuitively,
    // Collect from messages

    msgs := errutils.CollectMsg(wrapped3, nil)
    // msg has [ 'wrapped as different msg', 'some other error' ]

    // Error group
    errs := []error { err, codedErr, httpErr, wrapped1, wrapped2}
    errGroup := goerror.GroupFrom(errs)

    // Let's make things interesting
    fmt.Println(errGroup)
    // Output
    // multiple errors:
    // ...
    // prints all of the given errors
    // ...

    // Or take control yourself
    allErrs = errGroup.Errors()
    msgs = errutils.CollectMsg(allErrs, nil)
}
```

A real example in middleware using [`prasannavl/mchain`](https://www.github.com/prasannavl/mchain)

```go
func RequestIDMustInitHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		c := FromRequest(r)
		if _, ok := r.Header[RequestIDHeaderKey]; ok {
			msg := fmt.Sprintf("error: illegal header (%s)", RequestIDHeaderKey)
			return httperror.New(400, msg, true)
		}
		var uid uuid.UUID
		mustNewUUID(&uid)
		c.RequestID = uid
		return next.ServeHTTP(w, r)
	}
	return mchain.HandlerFunc(f)
}
```

## Notes

`HttpError` provides one additional method `Stop` that's useful to signify any middleware chain to stop processing.

Combining `Error` and `ErrorGroup` should be sufficient to handle most complex error wrapping, merging scenarios in Go without the use of other packages that add too many whistles which are, in my opinion completely unnecessary - and frankly not idiomatic Go.

## pkg/errors

Some of the features provided by this package are very similar to what the package `pkg/errors` offers. So, why do another error package?

- `pkg/errors` has a very nice and simple way to chain errors with `Cause`. However, other than that I tend to disagree with the pattern the pkg/errors package tends to build on. In a language with exceptions, stack traces are critical - and so even though it's very expensive to take them - we try to get it where-ever we can. But Go takes a drastically different approach towards error handling. It's very simple. You just pass a simple structure around. And as such, is very lightweight. And it's very clean because, in order to identify where an error originated from, you can create errors that are unique with codes (which is the one of the purpose of error codes in the first place, beyond internal handling of error such as comparisons). But `pkg/errors` destroys this "lightness" by calling `runtime.Callers` everytime you create an error. `It is very expensive to create a new error or wrap one in pkg/errors`. In Go, stack traces need to be only taken during `panic`s, or when you explicitly need them. I find taking stack traces along with errors to be harmful, and not idiomatic go.

    This package solves the problems that `pkg/errors` try to solve, and a lot more - by ErrorGroups, CodedError and utils, all while still retaining more simplicity than `pkg/errors`.

Other than this, `goerror` also provides a more complete error encapsulation solution:

- It provides `CodedError` pattern by out of the box - which is what a lot of usage requires rather than stack traces.
- Defines a neat way to group errors with `GroupError`.
- It also provides `httperror` that's built over `CodedError`, out of the box but as separate package.
