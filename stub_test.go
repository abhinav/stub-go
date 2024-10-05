package stub_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.abhg.dev/testing/stub"
)

func TestValue(t *testing.T) {
	v := 42
	restore := stub.Value(&v, 43)
	assert.Equal(t, 43, v)
	restore()
	assert.Equal(t, 42, v)
}

func TestFunc(t *testing.T) {
	fn := func() int { return 42 }

	restore := stub.Func(&fn, 43)
	assert.Equal(t, 43, fn())
	restore()
	assert.Equal(t, 42, fn())
}

func TestFuncNilIsZero(t *testing.T) {
	fn := func() int { return 42 }

	restore := stub.Func(&fn, nil)
	assert.Equal(t, 0, fn())
	restore()
	assert.Equal(t, 42, fn())
}

func TestFuncCast(t *testing.T) {
	fn := func() io.Reader { return strings.NewReader("hello") }

	restore := stub.Func(&fn, bytes.NewBufferString("world"))

	got1, err := io.ReadAll(fn())
	require.NoError(t, err)
	assert.Equal(t, "world", string(got1))

	restore()

	got2, err := io.ReadAll(fn())
	require.NoError(t, err)
	assert.Equal(t, "hello", string(got2))
}

func TestFuncErrors(t *testing.T) {
	t.Run("NotAPointer", func(t *testing.T) {
		assert.PanicsWithValue(t, "want pointer, got func() int", func() {
			stub.Func(func() int { return 0 }, 43)
		})
	})

	t.Run("NotAFunctionPointer", func(t *testing.T) {
		assert.PanicsWithValue(t, "want pointer to function, got *int", func() {
			stub.Func(new(int), 43)
		})
	})

	t.Run("ReturnMismatch", func(t *testing.T) {
		fn := func() int { return 42 }
		assert.PanicsWithValue(t, "want 1 return value(s), got 2", func() {
			stub.Func(&fn, 43, 44)
		})
	})

	t.Run("ReturnNotAssignable", func(t *testing.T) {
		fn := func() int { return 42 }
		assert.PanicsWithValue(t, "return type string (0) is not assignable to int", func() {
			stub.Func(&fn, "hello")
		})
	})
}
