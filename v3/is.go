package is

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
)

// Equaler is used to define equality for types.
//
// For example, this is useful if you have a struct that includes time.Time
// fields. You can implement this method and use time.Time.Equal() to do the
// comparison.
//
// Deprecated
type Equaler interface {
	Equal(in any) bool
}

// EqualityChecker is used to define equality for types during testing.
//
// For example, this is useful if you have a struct that includes time.Time
// fields. You can implement this method and use time.Time.Equal() to do the
// comparison.
type EqualityChecker interface {
	IsEqual(in any) bool
}

// Asserter provides methods that leverage the existing testing capabilities found
// in the Go test framework. The methods provided allow for a more natural,
// efficient and expressive approach to writing tests. The goal is to write
// fewer lines of code while improving communication of intent.
type Asserter interface {
	// tb returns the testing object with which this Asserter was originally
	// initialized.
	TB() testing.TB

	// Msg defines a message to print in the event of a failure. This allows you
	// to print out additional information about a failure if it happens.
	Msg(format string, args ...any) Asserter

	// AddMsg appends a message to print in the event of a failure. This allows
	// you to build a failure message in multiple steps. If no message was
	// previously set, simply sets the message.
	//
	// This method is most useful as a way of setting a default error message,
	// then adding additional information to the output for specific assertions.
	// For example:
	//
	// assert := is.New(t).Msg("User ID: %d",u.ID)
	// /*do things*/
	// assert.AddMsg("Raw Response: %s",body).Equal(res.StatusCode, http.StatusCreated)
	AddMsg(format string, args ...any) Asserter

	// Equal performs a deep compare of the provided objects and fails if they are
	// not equal.
	//
	// Equal does not respect type differences. If the types are different and
	// comparable (eg int32 and int64), they will be compared as though they are
	// the same type.
	Equal(actual any, expected any)

	// NotEqual performs a deep compare of the provided objects and fails if they are
	// equal.
	//
	// NotEqual does not respect type differences. If the types are different and
	// comparable (eg int32 and int64), they will be compared as though they are
	// the same type.
	NotEqual(a any, b any)

	// OneOf performs a deep compare of the provided object and an array of
	// comparison objects. It fails if the first object is not equal to one of the
	// comparison objects.
	//
	// OneOf does not respect type differences. If the types are different and
	// comparable (eg int32 and int64), they will be compared as though they are
	// the same type.
	OneOf(a any, b ...any)

	// NotOneOf performs a deep compare of the provided object and an array of
	// comparison objects. It fails if the first object is equal to one of the
	// comparison objects.
	//
	// NotOneOf does not respect type differences. If the types are different and
	// comparable (eg int32 and int64), they will be compared as though they are
	// the same type.
	NotOneOf(a any, b ...any)

	// Err checks the provided error object to determine if an error is present.
	Err(e error)

	// NotErr checks the provided error object to determine if an error is not
	// present.
	NotErr(e error)

	// Nil checks the provided object to determine if it is nil.
	Nil(o any)

	// NotNil checks the provided object to determine if it is not nil.
	NotNil(o any)

	// True checks the provided boolean to determine if it is true.
	True(b bool)

	// False checks the provided boolean to determine if is false.
	False(b bool)

	// Zero checks the provided object to determine if it is the zero value
	// for the type of that object. The zero value is the same as what the object
	// would contain when initialized but not assigned.
	//
	// This method, for example, would be used to determine if a string is empty,
	// an array is empty or a map is empty. It could also be used to determine if
	// a number is 0.
	//
	// In cases such as slice, map, array and chan, a nil value is treated the
	// same as an object with len == 0
	Zero(o any)

	// NotZero checks the provided object to determine if it is not the zero
	// value for the type of that object. The zero value is the same as what the
	// object would contain when initialized but not assigned.
	//
	// This method, for example, would be used to determine if a string is not
	// empty, an array is not empty or a map is not empty. It could also be used
	// to determine if a number is not 0.
	//
	// In cases such as slice, map, array and chan, a nil value is treated the
	// same as an object with len == 0
	NotZero(o any)

	// Len checks the provided object to determine if it is the same length as the
	// provided length argument.
	//
	// If the object is not one of type array, slice or map, it will fail.
	Len(o any, l int)

	// ShouldPanic expects the provided function to panic. If the function does
	// not panic, this assertion fails.
	ShouldPanic(f func())

	// EqualType checks the type of the two provided objects and
	// fails if they are not the same.
	EqualType(expected, actual any)

	// WaitForTrue waits until the provided func returns true. If the timeout is
	// reached before the function returns true, the test will fail.
	WaitForTrue(timeout time.Duration, f func() bool)

	// Lax accepts a function inside which a failed assertion will not halt
	// test execution. After the function returns, if any assertion had failed,
	// an additional message will be printed and test execution will be halted.
	//
	// This is useful for running assertions on, for example, many values in a struct
	// and having all the failed assertions print in one go, rather than having to run
	// the test multiple times, correcting a single failure per run.
	Lax(fn func(lax Asserter))
}

type asserter struct {
	tb         testing.TB
	strict     bool
	failFormat string
	failArgs   []any
	failed     bool
}

var _ Asserter = (*asserter)(nil)

// New returns a new Asserter containing the testing object provided.
func New(tb testing.TB) Asserter {
	if tb == nil {
		log.Fatalln("You must provide a testing object.")
	}
	return &asserter{tb: tb, strict: true}
}

func (self *asserter) TB() testing.TB {
	return self.tb
}

// Msg defines a message to print in the event of a failure. This allows you
// to print out additional information about a failure if it happens.
func (self *asserter) Msg(format string, args ...any) Asserter {
	return &asserter{
		tb:         self.tb,
		strict:     self.strict,
		failFormat: format,
		failArgs:   args,
	}
}

func (self *asserter) AddMsg(format string, args ...any) Asserter {
	if self.failFormat == "" {
		return self.Msg(format, args...)
	}
	return &asserter{
		tb:         self.tb,
		strict:     self.strict,
		failFormat: fmt.Sprintf("%s - %s", self.failFormat, format),
		failArgs:   append(self.failArgs, args...),
	}
}

func (self *asserter) Equal(actual any, expected any) {
	self.tb.Helper()
	if !isEqual(actual, expected) {
		fail(self, "actual value '%v' (%s) should be equal to expected value '%v' (%s)%s",
			actual, objectTypeName(actual),
			expected, objectTypeName(expected),
			diff(actual, expected),
		)
	}
}

func (self *asserter) NotEqual(actual any, expected any) {
	self.tb.Helper()
	if isEqual(actual, expected) {
		fail(self, "actual value '%v' (%s) should not be equal to expected value '%v' (%s)",
			actual, objectTypeName(actual),
			expected, objectTypeName(expected))
	}
}

func (self *asserter) OneOf(a any, b ...any) {
	self.tb.Helper()
	result := false
	for _, o := range b {
		result = isEqual(a, o)
		if result {
			break
		}
	}
	if !result {
		fail(self, "expected object '%s' to be equal to one of '%s', but got: %v and %v",
			objectTypeName(a),
			objectTypeNames(b), a, b)
	}
}

func (self *asserter) NotOneOf(a any, b ...any) {
	self.tb.Helper()
	result := false
	for _, o := range b {
		result = isEqual(a, o)
		if result {
			break
		}
	}
	if result {
		fail(self, "expected object '%s' not to be equal to one of '%s', but got: %v and %v",
			objectTypeName(a),
			objectTypeNames(b), a, b)
	}
}

func (self *asserter) Err(err error) {
	self.tb.Helper()
	if isNil(err) {
		fail(self, "expected error")
	}
}

func (self *asserter) NotErr(err error) {
	self.tb.Helper()
	if !isNil(err) {
		fail(self, "expected no error, but got: %v", err)
	}
}

func (self *asserter) Nil(o any) {
	self.tb.Helper()
	if !isNil(o) {
		fail(self, "expected object '%s' to be nil, but got: %v", objectTypeName(o), o)
	}
}

func (self *asserter) NotNil(o any) {
	self.tb.Helper()
	if isNil(o) {
		fail(self, "expected object '%s' not to be nil", objectTypeName(o))
	}
}

func (self *asserter) True(b bool) {
	self.tb.Helper()
	if !b {
		fail(self, "expected boolean to be true")
	}
}

func (self *asserter) False(b bool) {
	self.tb.Helper()
	if b {
		fail(self, "expected boolean to be false")
	}
}

func (self *asserter) Zero(o any) {
	self.tb.Helper()
	if !isZero(o) {
		fail(self, "expected object '%s' to be zero value, but it was: %v", objectTypeName(o), o)
	}
}

func (self *asserter) NotZero(o any) {
	self.tb.Helper()
	if isZero(o) {
		fail(self, "expected object '%s' not to be zero value", objectTypeName(o))
	}
}

func (self *asserter) Len(obj any, length int) {
	self.tb.Helper()
	t := reflect.TypeOf(obj)
	if obj == nil ||
		(t.Kind() != reflect.Array &&
			t.Kind() != reflect.Slice &&
			t.Kind() != reflect.Map) {
		fail(
			self,
			"expected object '%s' to be of length '%d', but the object is not one of array, slice or map",
			objectTypeName(obj),
			length,
		)
		return
	}

	rLen := reflect.ValueOf(obj).Len()
	if rLen != length {
		fail(self, "expected object '%s' to be of length '%d' but it was: %d", objectTypeName(obj), length, rLen)
	}
}

func (self *asserter) ShouldPanic(fn func()) {
	self.tb.Helper()
	defer func() {
		r := recover()
		if r == nil {
			fail(self, "expected function to panic")
		}
	}()
	fn()
}

func (self *asserter) EqualType(expected, actual any) {
	self.tb.Helper()
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		fail(
			self,
			"expected objects '%s' to be of the same type as object '%s'",
			objectTypeName(expected),
			objectTypeName(actual),
		)
	}
}

func (self *asserter) WaitForTrue(timeout time.Duration, f func() bool) {
	self.tb.Helper()
	after := time.After(timeout)
	for {
		select {
		case <-after:
			fail(self, "function did not return true within the timeout of %v", timeout)
			return
		default:
			if f() {
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (self *asserter) Lax(fn func(lax Asserter)) {
	lax := &asserter{
		tb:         self.tb,
		strict:     false,
		failFormat: self.failFormat,
		failArgs:   self.failArgs,
		failed:     false,
	}

	fn(lax)

	if lax.failed {
		fail(self, "at least one assertion in the Lax function failed")
	}
}
