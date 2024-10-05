package stub_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"go.abhg.dev/testing/stub"
)

func TestValue(t *testing.T) {
	v := 42

	restore := stub.Value(&v, 43)
	if v != 43 {
		t.Errorf("got %d, want 43", v)
	}

	restore()
	if v != 42 {
		t.Errorf("got %d, want 42", v)
	}
}

func TestFunc(t *testing.T) {
	fn := func() int { return 42 }

	restore := stub.Func(&fn, 43)
	if fn() != 43 {
		t.Errorf("got %d, want 43", fn())
	}

	restore()
	if fn() != 42 {
		t.Errorf("got %d, want 42", fn())
	}
}

func TestFuncNilIsZero(t *testing.T) {
	fn := func() int { return 42 }

	restore := stub.Func(&fn, nil)
	if fn() != 0 {
		t.Errorf("got %d, want 0", fn())
	}

	restore()
	if fn() != 42 {
		t.Errorf("got %d, want 42", fn())
	}
}

func TestFuncCast(t *testing.T) {
	fn := func() io.Reader { return strings.NewReader("hello") }

	restore := stub.Func(&fn, bytes.NewBufferString("world"))

	got1, err := io.ReadAll(fn())
	if err != nil {
		t.Fatal(err)
	}
	if string(got1) != "world" {
		t.Errorf("got %q, want 'world'", got1)
	}

	restore()

	got2, err := io.ReadAll(fn())
	if err != nil {
		t.Fatal(err)
	}
	if string(got2) != "hello" {
		t.Errorf("got %q, want 'hello'", got2)
	}
}

func TestFuncErrors(t *testing.T) {
	t.Run("NotAPointer", func(t *testing.T) {
		pval := expectPanic(t, func() {
			stub.Func(func() int { return 0 }, 43)
		})
		if want := "want pointer, got func() int"; pval != want {
			t.Errorf("got %q, want %q", pval, want)
		}
	})

	t.Run("NotAFunctionPointer", func(t *testing.T) {
		pval := expectPanic(t, func() {
			stub.Func(new(int), 43)
		})
		if want := "want pointer to function, got *int"; want != pval {
			t.Errorf("got %q, want %q", pval, want)
		}
	})

	t.Run("ReturnMismatch", func(t *testing.T) {
		fn := func() int { return 42 }
		pval := expectPanic(t, func() {
			stub.Func(&fn, 43, 44)
		})
		if want := "want 1 return value(s), got 2"; want != pval {
			t.Errorf("got %q, want %q", pval, want)
		}
	})

	t.Run("ReturnNotAssignable", func(t *testing.T) {
		fn := func() int { return 42 }
		pval := expectPanic(t, func() {
			stub.Func(&fn, "hello")
		})
		if want := "return type string (0) is not assignable to int"; want != pval {
			t.Errorf("got %q, want %q", pval, want)
		}
	})
}

func expectPanic(t testing.TB, fn func()) (pval any) {
	t.Helper()

	defer func() { pval = recover() }()
	fn()

	t.Fatalf("expected panic, got nil")
	return nil
}
