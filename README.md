# is [![GoDoc](https://godoc.org/github.com/tylerb/is?status.png)](http://godoc.org/github.com/tylerb/is) [![Build Status](https://circleci.com/gh/tylerb/is.svg?style=shield&circle-token=94428439ffc6eda6471dc218471dab20985f444c)](https://circleci.com/gh/tylerb/is)

Is provides a quick, clean and simple framework for writing Go tests.

## Installation

To install, simply execute:

```
go get github.com/tylerb/is
```

## Vendoring

Vendoring is recommended, as this library can change from time to time. The last change was updating it to use the new `Helper()` method in the 1.9 testing framework.

Check out the official Go dependency manager, [dep](https://github.com/golang/dep). Alternatively, I also like [glide](https://github.com/Masterminds/glide).

## Usage

Using `Is` is simple:

```go
func TestSomething(t *testing.T) {
	is := is.New(t)

	expected := 10
	result, _ := awesomeFunction()
	is.Equal(expected,result)
}
```

If you'd like a bit more information when a test fails, you may use the `Msg()` method:

```go
func TestSomething(t *testing.T) {
	is := is.New(t)

	expected := 10
	result, details := awesomeFunction()
	is.Msg("result details: %s", details).Equal(expected,result)
}
```

By default, Is fails and stops the test immediately. If you prefer to run multiple assertions to see them all fail at once, use the `Lax` method:

```go
func TestSomething(t *testing.T) {
	is := is.New(t).Lax()

	is.Equal(1,someFunc()) // if this fails, a message is printed and the test continues
	is.Equal(2,someOtherFunc()) // if this fails, a message is printed and the test continues
```

If you are using a relaxed instance of Is, you can switch it back to strict mode with `Strict`. This is useful when an assertion *must* be correct, or subsequent calls will panic:

```go
func TestSomething(t *testing.T) {
	is := is.New(t).Lax()

	results := someFunc()
	is.Strict().Equal(len(results),3) // if this fails, a message is printed and testing stops
	is.Equal(results[0],1) // if this fails, a message is printed and testing continues
	is.Equal(results[1],2)
	is.Equal(results[2],3)
```

Strict mode, in this case, applies only to the line on which it is invoked, as we don't overwrite our copy of the `is` variable.

## Contributing

All pull requests should:

- Pass [golangci-lint run](https://github.com/golangci/golangci-lint) with no warnings.
- Be `go fmt` formatted.
