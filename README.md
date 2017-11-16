# go-errors

A super tiny package for error encapsulation in idiomatic Go. A better and lighter alternative to [`pkg/errors`](https://github.com/pkg/errors).  
(See below for comparison to `pkg/errors`)

> **Note:** This package was previously called `goerror`. Now renamed to `errors` to be a drop-in replacement for std `errors` package.

**Documentation:** `Read the source, Luke` - it's tiny.

## Get

`go get -u github.com/prasannavl/go-errors`

## Sub-packages

- `httperror` - Provides `HttpError` interface and `HttpErr` type
- `errutils` - Provides a set of helpers to ease iteration, and message aggregation

## Types

Provides 3 interfaces
- `GoError`
- `CodedError`
- `ErrorGroup`

Provides 3 concrete types
- `GoErr`
- `CodedErr`
- `ErrGroup`

### GoError

`GoError` provides a straightforward way to wrap the `error` interface as `Cause`. The `utils` subpackage provides a way to iterate over chained errors and to collect error messages easily.

### CodedError

`CodedError` adds just one item to `Error` - `Code`. It's useful for building any errors with an `int` code. The `httperror` subpackage, for example simply builds over this with an added code validation.

### ErrorGroup

`ErrGroup` provides a clean way to aggregate errors as `[]error`. But unlike other libraries out there - doesn't provides any way to add, combine, or remove errors. It prefers to sticks to Go's idiomatic way of keeping things minimalistic. You can use a Go's own `make`, `append` and other slice function to do the same. Once done, you simply use `GroupFrom` to add your slice to create a group. It's all just slices. The only thing `ErrGroup` does is provides you an `error` interface that by default prints a nice message, and a way to retrive all the original errors.

## Example

```go

func TestErrors() {
    // New errors

    err := errors.New("some error")
    // Error with codes
    codedErr := errors.NewCoded(42, "compute overwhelmed. crashing")
    // HttpErrors are essentially, coded errors,
    // with http code validation
    httpErr := httperror.New(400, "naughty you! that's not valid", true)

    // Wrap existing existing errors
    wrapped1 := errors.From(errors.New("some other error"))
    wrapped2 := errors.From(wrapped1)
    // With additional message
    wrapper3 := errors.NewWithCause("wrapped as different msg", wrapped2)

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
    errGroup := errors.GroupFrom(errs)

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

## pkg/errors

Some of the features provided by this package are very similar to what the package `pkg/errors` offers. So, why do another error package?

- `pkg/errors` has a very nice and simple way to chain errors with `Cause`. However, other than that I tend to disagree with the pattern the pkg/errors package tends to build on. **In a language with exceptions, stack traces are critical - and so even though it's very expensive to take them - we try to get it where-ever we can. But Go takes a drastically different approach towards error handling. It's very simple.** You just pass a simple structure around. And as such, is very lightweight. And it's very clean because, in order to identify where an error originated from, you can create errors that are unique with codes (which is the one of the purpose of error codes in the first place, beyond internal handling of error such as comparisons). But `pkg/errors` destroys this "lightness" by calling `runtime.Callers` everytime you create an error. `It is very expensive to create a new error or wrap one in pkg/errors`. In Go, stack traces need to be only taken during `panic`s, or when you explicitly need them. I find taking stack traces along with errors to be harmful, and not idiomatic go.

    This package solves the problems that `pkg/errors` try to solve, and a lot more - by ErrorGroups, CodedError and utils, all while still retaining more simplicity than `pkg/errors`.

Other than this, `goerror` also provides a more complete error encapsulation solution:

- It provides `CodedError` pattern by out of the box - which is what a lot of usage requires rather than stack traces.
- Defines a neat way to group errors with `GroupError`.
- It also provides `httperror` that's built over `CodedError`, out of the box but as separate package.

## How is this more Go idiomatic?

As mentioned, stack-traces absolutely essential when you work with language with execptions. Because they can throw from anywhere, and you can be `catch`ing it elsewhere - you're bound to lose track of where you are very quickly. However the approach Go takes is simple. You check your errors at every surface of an API. So, simple error type with kind can work very nicely. And that doesn't mean you can't or shoudln't take stack traces. You can always just wrap them in higher order functions. I've found it to be a better practice *to take traces at component or contextual boundaries, when needed.*

**Simple composable errors with codes, should be the basic building block of error handling - while a package like `pkg/errors` that's intended to be a fundamental building block, unfortunately does expensive caller stacks, and is structured far more complex internally** 

## I like stack-traces! 

> Then take them! 

It's just not intended to go into the building block. They should be above it. So you have exactly that choice. There **are times when stack traces are critical, and there are times when they are just a waste of system resources** - but the truth is, error handling helpers like these can never know that regardless of any number of debates or discussions - and so this package does the absolute minimum to act as a building block, yet provides nice composition. Rest can be built over it.

## How is this simpler?

Don't let the number of types fool you. Conceptually, it's all just in one file `errors.go`, which is about 40 lines long. (Yup, that's it). All it really does is just wrap things up neatly in a type, and allows you to compose them. If you're a reasonably proficient in Go, you should be able to read and understand the whole code of the library in just under 5 minutes - I encourage you to - *you should know exactly what happens* you instantiate an error type.

## Drop-in replacement for std `errors`?

Yup. It is!

## Notes

`HttpError` provides one additional method `End` that's useful to signify any middleware chain to stop processing.

Combining `GoError` and `ErrorGroup` should be sufficient to handle most complex error wrapping, merging scenarios in Go without the use of other packages that add too many whistles which are, in my opinion completely unnecessary - and frankly not idiomatic Go.
