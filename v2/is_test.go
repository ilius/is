package is

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var numberTypes = []reflect.Type{
	reflect.TypeOf(int(0)),
	reflect.TypeOf(int8(0)),
	reflect.TypeOf(int16(0)),
	reflect.TypeOf(int32(0)),
	reflect.TypeOf(int64(0)),
	reflect.TypeOf(uint(0)),
	reflect.TypeOf(uint8(0)),
	reflect.TypeOf(uint16(0)),
	reflect.TypeOf(uint32(0)),
	reflect.TypeOf(uint64(0)),
	reflect.TypeOf(float32(0)),
	reflect.TypeOf(float64(0)),
}

type testStruct struct {
	v int
}

var tests = []struct {
	a      interface{}
	b      interface{}
	c      interface{}
	d      interface{}
	e      interface{}
	cTypes []reflect.Type
}{
	{
		a:      0,
		b:      0,
		c:      1,
		d:      0,
		e:      1,
		cTypes: numberTypes,
	},
	{
		a: "test",
		b: "test",
		c: "testing",
		d: "",
		e: "testing",
	},
	{
		a: struct{}{},
		b: struct{}{},
		c: struct{ v int }{v: 1},
		d: testStruct{},
		e: testStruct{v: 1},
	},
	{
		a: &struct{}{},
		b: &struct{}{},
		c: &struct{ v int }{v: 1},
		d: &testStruct{},
		e: &testStruct{v: 1},
	},
	{
		a: []int64{0, 1},
		b: []int64{0, 1},
		c: []int64{0, 2},
		d: []int64{},
		e: []int64{0, 2},
	},
	{
		a: map[string]int64{"answer": 42},
		b: map[string]int64{"answer": 42},
		c: map[string]int64{"answer": 43},
		d: map[string]int64{},
		e: map[string]int64{"answer": 42},
	},
	{
		a: true,
		b: true,
		c: false,
		d: false,
		e: true,
	},
}

func TestNewPanic(t *testing.T) {
	is := New(t)
	is.ShouldPanic(func() {
		New(nil)
	})
}

func TestIs(t *testing.T) {
	is := New(t)
	is = is.New(t)

	for i, test := range tests {
		for _, cType := range test.cTypes {
			fail = func(is *Is, format string, args ...interface{}) {
				fmt.Print(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...))
				t.FailNow()
			}
			is.Equal(test.a, reflect.ValueOf(test.b).Convert(cType).Interface())
		}
		is.Equal(test.a, test.b)
	}

	for i, test := range tests {
		for _, cType := range test.cTypes {
			fail = func(is *Is, format string, args ...interface{}) {
				fmt.Print(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...))
				t.FailNow()
			}
			is.NotEqual(test.a, reflect.ValueOf(test.c).Convert(cType).Interface())
		}
		is.NotEqual(test.a, test.c)
	}

	for i, test := range tests {
		for _, cType := range test.cTypes {
			fail = func(is *Is, format string, args ...interface{}) {
				fmt.Print(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...))
				t.FailNow()
			}
			is.Zero(reflect.ValueOf(test.d).Convert(cType).Interface())
		}
		is.Zero(test.d)
	}

	for i, test := range tests {
		for _, cType := range test.cTypes {
			fail = func(is *Is, format string, args ...interface{}) {
				fmt.Print(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...))
				t.FailNow()
			}
			is.NotZero(reflect.ValueOf(test.e).Convert(cType).Interface())
		}
		is.NotZero(test.e)
	}

	fail = func(is *Is, format string, args ...interface{}) {
		fmt.Print(fmt.Sprintf(format, args...))
		t.FailNow()
	}
	is.Nil(nil)
	is.NotNil(&testStruct{v: 1})
	is.Err(errors.New("error"))
	is.ErrMsg(errors.New("another error"), "another error")
	is.NotErr(nil)
	is.True(true)
	is.False(false)
	is.Zero(nil)
	is.Nil((*testStruct)(nil))
	is.OneOf(1, 2, 3, 1)
	is.NotOneOf(1, 2, 3)
	is.EqualType(1, 2)

	lens := []interface{}{
		[]int{1, 2, 3},
		[3]int{1, 2, 3},
		map[int]int{1: 1, 2: 2, 3: 3},
	}
	for _, l := range lens {
		is.Len(l, 3)
	}

	fail = func(is *Is, format string, args ...interface{}) {}
	is.Equal((*testStruct)(nil), &testStruct{})
	is.Equal(&testStruct{}, (*testStruct)(nil))
	is.Equal((*testStruct)(nil), (*testStruct)(nil))

	fail = func(is *Is, format string, args ...interface{}) {
		fmt.Print(fmt.Sprintf(format, args...))
		t.FailNow()
	}
	is.ShouldPanic(func() {
		panic("The sky is falling!")
	})
}

func TestIsMsg(t *testing.T) {
	is := New(t)

	is = is.Msg("something %s", "else")
	if is.failFormat != "something %s" {
		t.Fatal("failFormat not set")
	}
	if is.failArgs[0].(string) != "else" {
		t.Fatal("failArgs not set")
	}

	is = is.AddMsg("another %s %s", "couple", "things")
	if is.failFormat != "something %s - another %s %s" {
		t.Fatal("AddMsg did not work")
	}
	is = is.PrependMsg("#%d message", 1)
	if is.failFormat != "#%d message - something %s - another %s %s" {
		t.Fatal("PrependMsg did not work")
	}
	if is.failArgs[0].(int) != 1 {
		t.Fatal("failArgs not set")
	}
	if is.failArgs[1].(string) != "else" {
		t.Fatal("failArgs not set")
	}
	if is.failArgs[2].(string) != "couple" {
		t.Fatal("failArgs not set")
	}
	if is.failArgs[3].(string) != "things" {
		t.Fatal("failArgs not set")
	}
}

func TestIsAddMsg(t *testing.T) {
	is := New(t)
	is = is.AddMsg("something %s %s", "new", "here")
	if is.failFormat != "something %s %s" {
		t.Fatal("AddMsg: bad failFormat")
	}
	if is.failArgs[0].(string) != "new" {
		t.Fatal("AddMsg: bad failArgs[0]")
	}
	if is.failArgs[1].(string) != "here" {
		t.Fatal("AddMsg: bad failArgs[1]")
	}
}

func TestIsPrependMsg(t *testing.T) {
	is := New(t)
	is = is.PrependMsg("something %s %s", "new", "here")
	if is.failFormat != "something %s %s" {
		t.Fatal("PrependMsg: bad failFormat")
	}
	if is.failArgs[0].(string) != "new" {
		t.Fatal("PrependMsg: bad failArgs[0]")
	}
	if is.failArgs[1].(string) != "here" {
		t.Fatal("PrependMsg: bad failArgs[1]")
	}
}

func TestIsMsgSep(t *testing.T) {
	is := New(t)
	is = is.MsgSep(", ")
	is = is.AddMsg("msg one")
	is = is.AddMsg("msg two")
	if is.failFormat != "msg one, msg two" {
		t.Fatal("bad failFormat")
	}
}

func TestIsLax(t *testing.T) {
	is := New(t)

	hit := 0

	fail = func(is *Is, format string, args ...interface{}) {
		if is.strict {
			t.FailNow()
		}
		hit++
	}

	is.Lax().Equal(1, 2)

	fail = failDefault

	is.Strict().Equal(hit, 1)
}

func TestIsOneOf(t *testing.T) {
	is := New(t)

	hit := 0
	fail = func(is *Is, format string, args ...interface{}) {
		hit++
	}
	is.OneOf(2, 1, 2, 3)
	is.OneOf(4, 1, 2, 3)
	is.NotOneOf(2, 1, 2, 3)
	is.NotOneOf(4, 1, 2, 3)

	fail = failDefault
	is.Strict().Equal(hit, 2)
}

func TestContains(t *testing.T) {
	is := New(t)

	hit := 0
	fail = func(is *Is, format string, args ...interface{}) {
		hit++
	}

	is.True(is.Contains("hello", "ell"))
	is.True(is.Contains([]string{"hello", "world"}, "hello"))
	is.False(is.Contains("hello", "elf"))
	is.False(is.Contains([]string{"hello", "world"}, "test"))

	fail = failDefault
	is.Strict().Equal(hit, 2)
}

func TestIsFailures(t *testing.T) {
	is := New(t)

	hit := 0
	fail = func(is *Is, format string, args ...interface{}) {
		hit++
	}

	is.NotEqual(1, 1)
	is.Err(nil)
	is.ErrMsg(errors.New("error 1"), "error 2")
	is.NotErr(errors.New("error"))
	is.Nil(&hit)
	is.NotNil(nil)
	is.True(false)
	is.False(true)
	is.Zero(1)
	is.NotZero(0)
	is.Len([]int{}, 1)
	is.Len(nil, 1)
	is.ShouldPanic(func() {})

	fail = failDefault
	is.Strict().Equal(hit, 13)
}

func TestWaitForTrue(t *testing.T) {
	is := New(t)

	hit := 0
	fail = func(is *Is, format string, args ...interface{}) {
		hit++
	}

	is.WaitForTrue(200*time.Millisecond, func() bool {
		return false
	})
	is.Strict().Equal(hit, 1)

	is.WaitForTrue(200*time.Millisecond, func() bool {
		return true
	})
	is.Strict().Equal(hit, 1)
}

type equaler struct {
	equal  bool
	called bool
}

func (e *equaler) Equal(in interface{}) bool {
	e.called = true
	v, ok := in.(*equaler)
	if !ok {
		return false
	}
	return e.equal == v.equal
}

func TestEqualer(t *testing.T) {
	is := New(t)

	hit := 0
	fail = func(is *Is, format string, args ...interface{}) {
		hit++
	}

	a := &equaler{equal: true}
	b := &equaler{}

	is.Equal(a, b)
	if !a.called {
		t.Fatalf("a.Equal should have been called")
	}

	is.Equal(b, a)
	if !b.called {
		t.Fatalf("b.Equal should have been called")
	}

	if hit != 2 {
		t.Fatalf("fail func should have been called 2 times, but was called %d times", hit)
	}

	a.called = false
	b.called = false
	b.equal = true
	hit = 0

	is.NotEqual(a, b)
	if !a.called {
		t.Fatalf("a.Equal should have been called")
	}

	is.NotEqual(b, a)
	if !b.called {
		t.Fatalf("b.Equal should have been called")
	}

	if hit != 2 {
		t.Fatalf("fail func should have been called 2 times, but was called %d times", hit)
	}
}
