package is

import (
	"bytes"
	"fmt"
	"reflect"
)

func objectTypeName(o any) string {
	return fmt.Sprintf("%T", o)
}

func objectTypeNames(o []any) string {
	if o == nil {
		return objectTypeName(o)
	}
	if len(o) == 1 {
		return objectTypeName(o[0])
	}
	var b bytes.Buffer
	b.WriteString(objectTypeName(o[0]))
	for _, e := range o[1:] {
		b.WriteString(",")
		b.WriteString(objectTypeName(e))
	}
	return b.String()
}

func isNil(o any) bool {
	if o == nil {
		return true
	}
	value := reflect.ValueOf(o)
	kind := value.Kind()
	if kind >= reflect.Chan &&
		kind <= reflect.Slice &&
		value.IsNil() {
		return true
	}
	return false
}

func isZero(o any) bool {
	if o == nil {
		return true
	}
	v := reflect.ValueOf(o)
	switch v.Kind() {
	case reflect.Ptr:
		return reflect.DeepEqual(o,
			reflect.New(v.Type().Elem()).Interface())
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		if v.Len() == 0 {
			return true
		}
		return false
	default:
		return reflect.DeepEqual(o,
			reflect.Zero(v.Type()).Interface())
	}
}

func isEqual(a any, b any) bool {
	if isNil(a) || isNil(b) {
		if isNil(a) && !isNil(b) {
			return false
		}
		if !isNil(a) && isNil(b) {
			return false
		}
		return a == b
	}

	// Call a.Equaler if it is implemented
	if e, ok := a.(Equaler); ok {
		return e.Equal(b)
	}

	if reflect.DeepEqual(a, b) {
		return true
	}

	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	// Convert types and compare
	if bValue.Type().ConvertibleTo(aValue.Type()) {
		return reflect.DeepEqual(a, bValue.Convert(aValue.Type()).Interface())
	}

	return false
}

// fail is a function variable that is called by test functions when they
// fail. It is overridden in test code for this package.
var fail = failDefault

// failDefault is the default failure function.
func failDefault(is *Is, format string, args ...any) {
	is.TB.Helper()

	failFmt := format
	if len(is.failFormat) != 0 {
		failFmt = fmt.Sprintf("%s%s%s", format, is.getMsgSep(), is.failFormat)
		args = append(args, is.failArgs...)
	}
	if is.strict {
		is.TB.Fatalf(failFmt, args...)
	} else {
		is.TB.Errorf(failFmt, args...)
	}
}
