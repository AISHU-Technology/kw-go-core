package assert

import (
	"fmt"
	"github.com/json-iterator/go"
	"reflect"
	"regexp"
	"strings"
)

// T is the minimum interface of *testing.T.
type T interface {
	Error(args ...interface{})
}

func Fail(t T, str string, msg ...string) {
	args := append([]string{str}, msg...)
	t.Error(strings.Join(args, "; "))
}

// True assertion failed when got is false.
func True(t T, got bool, msg ...string) {
	if !got {
		Fail(t, "got false but expect true", msg...)
	}
}

// False assertion failed when got is true.
func False(t T, got bool, msg ...string) {
	if got {
		Fail(t, "got true but expect false", msg...)
	}
}

// isNil reports v is nil, but will not panic.
func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice,
		reflect.UnsafePointer:
		return v.IsNil()
	}
	return !v.IsValid()
}

// Nil assertion failed when got is not nil.
func Nil(t T, got interface{}, msg ...string) {

	// Why can't we use got==nil to judge？Because if
	// a := (*int)(nil)        // %T == *int
	// b := (interface{})(nil) // %T == <nil>
	// then a==b is false, because they are different types.
	if !isNil(reflect.ValueOf(got)) {
		str := fmt.Sprintf("got (%T) %v but expect nil", got, got)
		Fail(t, str, msg...)
	}
}

// NotNil assertion failed when got is nil.
func NotNil(t T, got interface{}, msg ...string) {

	if isNil(reflect.ValueOf(got)) {
		Fail(t, "got nil but expect not nil", msg...)
	}
}

// Equal assertion failed when got and expect are not `deeply equal`.
func Equal(t T, got interface{}, expect interface{}, msg ...string) {

	if !reflect.DeepEqual(got, expect) {
		str := fmt.Sprintf("got (%t) %v but expect (%t) %v", got, got, expect, expect)
		Fail(t, str, msg...)
	}
}

// NotEqual assertion failed when got and expect are `deeply equal`.
func NotEqual(t T, got interface{}, expect interface{}, msg ...string) {

	if reflect.DeepEqual(got, expect) {
		str := fmt.Sprintf("got (%T) %v but expect not (%T) %v", got, got, expect, expect)
		Fail(t, str, msg...)
	}
}

// JsonEqual assertion failed when got and expect are not `json equal`.
func JsonEqual(t T, got string, expect string, msg ...string) {

	var gotJson interface{}
	if err := jsoniter.Unmarshal([]byte(got), &gotJson); err != nil {
		Fail(t, err.Error(), msg...)
		return
	}
	var expectJson interface{}
	if err := jsoniter.Unmarshal([]byte(expect), &expectJson); err != nil {
		Fail(t, err.Error(), msg...)
		return
	}
	if !reflect.DeepEqual(gotJson, expectJson) {
		str := fmt.Sprintf("got (%T) %v but expect (%T) %v", got, got, expect, expect)
		Fail(t, str, msg...)
	}
}

// Same assertion failed when got and expect are not same.
func Same(t T, got interface{}, expect interface{}, msg ...string) {

	if got != expect {
		str := fmt.Sprintf("got (%T) %v but expect (%T) %v", got, got, expect, expect)
		Fail(t, str, msg...)
	}
}

// NotSame assertion failed when got and expect are same.
func NotSame(t T, got interface{}, expect interface{}, msg ...string) {

	if got == expect {
		str := fmt.Sprintf("expect not (%T) %v", expect, expect)
		Fail(t, str, msg...)
	}
}

// Panic assertion failed when fn doesn't panic or not match expr expression.
func Panic(t T, fn func(), expr string, msg ...string) {

	str := recovery(fn)
	if str == "<<SUCCESS>>" {
		Fail(t, "did not panic", msg...)
	} else {
		matches(t, str, expr, msg...)
	}
}

func recovery(fn func()) (str string) {
	defer func() {
		if r := recover(); r != nil {
			str = fmt.Sprint(r)
		}
	}()
	fn()
	return "<<SUCCESS>>"
}

// Matches assertion failed when got doesn't match expr expression.
func Matches(t T, got string, expr string, msg ...string) {

	matches(t, got, expr, msg...)
}

// Error assertion failed when got `error` doesn't match expr expression.
func Error(t T, got error, expr string, msg ...string) {

	if got == nil {
		Fail(t, "expect not nil error", msg...)
		return
	}
	matches(t, got.Error(), expr, msg...)
}

func matches(t T, got string, expr string, msg ...string) {

	if ok, err := regexp.MatchString(expr, got); err != nil {
		Fail(t, "invalid pattern", msg...)
	} else if !ok {
		str := fmt.Sprintf("got %q which does not match %q", got, expr)
		Fail(t, str, msg...)
	}
}

// TypeOf assertion failed when got and expect are not same type.
func TypeOf(t T, got interface{}, expect interface{}, msg ...string) {

	e1 := reflect.TypeOf(got)
	e2 := reflect.TypeOf(expect)
	if e2.Kind() == reflect.Ptr && e2.Elem().Kind() == reflect.Interface {
		e2 = e2.Elem()
	}

	if !e1.AssignableTo(e2) {
		str := fmt.Sprintf("got type (%s) but expect type (%s)", e1, e2)
		Fail(t, str, msg...)
	}
}

// Implements assertion failed when got doesn't implement expect.
func Implements(t T, got interface{}, expect interface{}, msg ...string) {

	e1 := reflect.TypeOf(got)
	e2 := reflect.TypeOf(expect)
	if e2.Kind() == reflect.Ptr {
		if e2.Elem().Kind() == reflect.Interface {
			e2 = e2.Elem()
		} else {
			Fail(t, "expect should be interface", msg...)
			return
		}
	}

	if !e1.Implements(e2) {
		str := fmt.Sprintf("got type (%s) but expect type (%s)", e1, e2)
		Fail(t, str, msg...)
	}
}

// InSlice assertion failed when got is not in expect array & slice.
func InSlice(t T, got interface{}, expect interface{}, msg ...string) {

	v := reflect.ValueOf(expect)
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		str := fmt.Sprintf("unsupported expect value (%s) %s", expect, expect)
		Fail(t, str, msg...)
		return
	}

	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(got, v.Index(i).Interface()) {
			return
		}
	}

	str := fmt.Sprintf("got (%T) %v is not in (%T) %v", got, got, expect, expect)
	Fail(t, str, msg...)
}

// NotInSlice assertion failed when got is in expect array & slice.
func NotInSlice(t T, got interface{}, expect interface{}, msg ...string) {

	v := reflect.ValueOf(expect)
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		str := fmt.Sprintf("unsupported expect value (%v) %v", expect, expect)
		Fail(t, str, msg...)
		return
	}

	e := reflect.TypeOf(got)
	if e != v.Type().Elem() {
		str := fmt.Sprintf("got type (%s) doesn't match expect type (%s)", e, v.Type())
		Fail(t, str, msg...)
		return
	}

	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(got, v.Index(i).Interface()) {
			str := fmt.Sprintf("got (%T) %v is in (%T) %v", got, got, expect, expect)
			Fail(t, str, msg...)
			return
		}
	}
}

// SubInSlice assertion failed when got is not sub in expect array & slice.
func SubInSlice(t T, got interface{}, expect interface{}, msg ...string) {

	v1 := reflect.ValueOf(got)
	if v1.Kind() != reflect.Array && v1.Kind() != reflect.Slice {
		str := fmt.Sprintf("unsupported got value (%s) %s", got, got)
		Fail(t, str, msg...)
		return
	}

	v2 := reflect.ValueOf(expect)
	if v2.Kind() != reflect.Array && v2.Kind() != reflect.Slice {
		str := fmt.Sprintf("unsupported expect value (%s) %s", expect, expect)
		Fail(t, str, msg...)
		return
	}

	for i := 0; i < v1.Len(); i++ {
		for j := 0; j < v2.Len(); j++ {
			if reflect.DeepEqual(v1.Index(i).Interface(), v2.Index(j).Interface()) {
				return
			}
		}
	}

	str := fmt.Sprintf("got (%T) %v is not sub in (%T) %v", got, got, expect, expect)
	Fail(t, str, msg...)
}

// InMapKeys assertion failed when got is not in keys of expect map.
func InMapKeys(t T, got interface{}, expect interface{}, msg ...string) {

	switch v := reflect.ValueOf(expect); v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if reflect.DeepEqual(got, key.Interface()) {
				return
			}
		}
	default:
		str := fmt.Sprintf("unsupported expect value (%T) %v", expect, expect)
		Fail(t, str, msg...)
		return
	}

	str := fmt.Sprintf("got (%T) %v is not in keys of (%T) %v", got, got, expect, expect)
	Fail(t, str, msg...)
}

// InMapValues assertion failed when got is not in values of expect map.
func InMapValues(t T, got interface{}, expect interface{}, msg ...string) {

	switch v := reflect.ValueOf(expect); v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if reflect.DeepEqual(got, v.MapIndex(key).Interface()) {
				return
			}
		}
	default:
		str := fmt.Sprintf("unsupported expect value (%T) %v", expect, expect)
		Fail(t, str, msg...)
		return
	}
	str := fmt.Sprintf("got (%T) %v is not in values of (%T) %v", got, got, expect, expect)
	Fail(t, str, msg...)
}

// HasPrefix assertion failed when v doesn't have prefix `prefix`.
func HasPrefix(t T, value, prefix string, msg ...string) bool {

	if !strings.HasPrefix(value, prefix) {
		Fail(t, fmt.Sprintf("'%s' doesn't have prefix '%s'", value, prefix), msg...)
		return false
	}
	return true
}

// EqualFold assertion failed when v doesn't equal to `s` under Unicode case-folding.
func EqualFold(t T, value, s string, msg ...string) bool {

	if !strings.EqualFold(value, s) {
		Fail(t, fmt.Sprintf("'%s' doesn't equal fold to '%s'", value, s), msg...)
		return false
	}
	return true
}

// HasSuffix assertion failed when v doesn't have suffix `suffix`.
func HasSuffix(t T, value, suffix string, msg ...string) bool {

	if !strings.HasSuffix(value, suffix) {
		Fail(t, fmt.Sprintf("'%s' doesn't have suffix '%s'", value, suffix), msg...)
		return false
	}
	return true
}

// Contains assertion failed when v doesn't contain substring `substr`.
func Contains(t T, value, substr string, msg ...string) bool {

	if !strings.Contains(value, substr) {
		Fail(t, fmt.Sprintf("'%s' doesn't contain substr '%s'", value, substr), msg...)
		return false
	}
	return true
}
