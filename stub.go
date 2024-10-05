// Package stub provides helper functions to replace global variables
// for testing, and restore them afterwards.
package stub

import (
	"fmt"
	"reflect"
)

// Value replaces the value of a pointer with a new value,
// and returns a function that restores the original value.
//
// Idiomatic usage will typically look like this:
//
//	func TestSomething(t *testing.T) {
//		defer stub.Value(&globalVar, newValue)()
//
//		// ...
//	}
//
// If the test has subtests that use t.Parallel,
// use t.Cleanup instead of defer:
//
//	func TestSomething(t *testing.T) {
//		t.Cleanup(stub.Value(&globalVar, newValue))
//
//		t.Run("subtest", func(t *testing.T) {
//			t.Parallel()
//
//			// ...
//		})
//	}
func Value[T any](ptr *T, value T) (restore func()) {
	original := *ptr
	*ptr = value
	return func() {
		*ptr = original
	}
}

// Func replaces the value of a function pointer
// with a function that returns the provided values.
// It returns a function that restores the original function.
// If the function has multiple return values, pass them all.
//
// Idiomatic usage will typically look like this:
//
//	func TestSomething(t *testing.T) {
//		defer stub.StubFunc(&globalFunc, 42)()
//
//		globalFunc() // returns 42
//		// ...
//	}
//
// If the test has subtests that use t.Parallel,
// use t.Cleanup instead of defer:
//
//	func TestSomething(t *testing.T) {
//		t.Cleanup(stub.StubFunc(&globalFunc, 42))
//
//		t.Run("subtest", func(t *testing.T) {
//			t.Parallel()
//
//			globalFunc() // returns 42
//			// ...
//		})
//	}
func Func(fnptr any, rets ...any) (restore func()) {
	fnptrv := reflect.ValueOf(fnptr)
	if fnptrv.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("want pointer, got %T", fnptr))
	}

	fnv := fnptrv.Elem()
	if fnv.Kind() != reflect.Func {
		panic(fmt.Sprintf("want pointer to function, got %T", fnptr))
	}

	fnt := fnv.Type()
	if fnt.NumOut() != len(rets) {
		panic(fmt.Sprintf("want %d return value(s), got %d", fnt.NumOut(), len(rets)))
	}

	vals := make([]reflect.Value, fnt.NumOut())
	for i, ret := range rets {
		v := reflect.ValueOf(ret)
		if !v.IsValid() {
			// nil is not a valid value for reflect.ValueOf,
			// but may be passed as a placeholder for zero values.
			v = reflect.New(fnt.Out(i)).Elem()
		}

		if !v.Type().AssignableTo(fnt.Out(i)) {
			panic(fmt.Sprintf("return type %v (%d) is not assignable to %v", v.Type(), i, fnt.Out(i)))
		}

		vals[i] = v
	}

	original := fnv.Interface()
	stub := reflect.MakeFunc(fnv.Type(), func([]reflect.Value) []reflect.Value {
		return vals
	}).Interface()

	fnv.Set(reflect.ValueOf(stub))
	return func() {
		fnv.Set(reflect.ValueOf(original))
	}
}
