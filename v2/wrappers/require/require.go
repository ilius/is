package require

import (
	"reflect"
	"testing"

	"github.com/ilius/is/v2"
)

type PanicTestFunc func()

type TestingT = testing.TB

func Contains(t TestingT, s any, contains any, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.Contains(s, contains)
}

func Equal(t TestingT, expected any, actual any, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.Equal(actual, expected)
}

func EqualError(t TestingT, theError error, errString string, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.ErrMsg(theError, errString)
}

func EqualValuesf(t TestingT, expected any, actual any, msg string, args ...any) {
	is := is.New(t)
	is.AddMsg(msg, args...)
	is.Equal(actual, expected)
	is.EqualType(expected, actual)
}

func Error(t TestingT, err error, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.Err(err)
}

func False(t TestingT, value bool, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.False(value)
}

func IsType(t TestingT, expectedType any, object any, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.IsType(expectedType.(reflect.Type), object)
}

func Len(t TestingT, object any, length int, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.Len(object, length)
}

func Nil(t TestingT, object any, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.Nil(object)
}

func NoError(t TestingT, err error, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.NotErr(err)
}

func NotNil(t TestingT, object any, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.NotNil(object)
}

func Panics(t TestingT, f PanicTestFunc, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.ShouldPanic(f)
}

func True(t TestingT, value bool, msgAndArgs ...any) {
	is := is.New(t)
	if len(msgAndArgs) > 0 {
		format := msgAndArgs[0].(string)
		is.AddMsg(format, msgAndArgs[1:]...)
	}
	is.True(value)
}
