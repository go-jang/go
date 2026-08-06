package objects_test

import (
	"errors"
	"testing"
	"time"

	"github.com/go-jang/go/util/objects"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToString(test *testing.T) {
	require.Equal(test, "nil", objects.ToString(nil))
	require.Equal(test, "abc", objects.ToString("abc"))
	require.Equal(test, "abc", objects.ToString([]byte("abc")))

	require.Equal(test, "true", objects.ToString(true))
	require.Equal(test, "false", objects.ToString(false))

	require.Equal(test, "1", objects.ToString(int(1)))
	require.Equal(test, "2", objects.ToString(int8(2)))
	require.Equal(test, "3", objects.ToString(int16(3)))
	require.Equal(test, "4", objects.ToString(int32(4)))
	require.Equal(test, "5", objects.ToString(int64(5)))

	require.Equal(test, "6", objects.ToString(uint(6)))
	require.Equal(test, "7", objects.ToString(uint8(7)))
	require.Equal(test, "8", objects.ToString(uint16(8)))
	require.Equal(test, "9", objects.ToString(uint32(9)))
	require.Equal(test, "10", objects.ToString(uint64(10)))

	require.Equal(test, "1.5", objects.ToString(float32(1.5)))
	require.Equal(test, "2.5", objects.ToString(float64(2.5)))

	require.Equal(test, "boom", objects.ToString(errors.New("boom")))

	require.Equal(test, "[1 2 3]", objects.ToString([]int{1, 2, 3}))
}

func Test_Equal(t *testing.T) {
	// Equal
	require.True(t, objects.Equal(255, 255))

	time1 := time.Now()
	time2 := time1
	require.True(t, objects.Equal(time1, time2))

	slice1 := make([]any, 0)
	slice2 := slice1
	require.True(t, objects.Equal(slice1, slice2))

	require.True(t, objects.Equal(slice1, make([]any, 0)))

	// Not eqal
	require.False(t, objects.Equal(255, 256))

	slice1 = make([]any, 0)
	slice2 = append(slice1, 1)
	require.False(t, objects.Equal(slice1, slice2))
}

func Test_HashCode(t *testing.T) {
	require.Equal(t, 8114084241814008568, objects.HashCode(int(123)))
	require.Equal(t, 4953291999513788008, objects.HashCode(int8(123)))
	require.Equal(t, -7208379536669766984, objects.HashCode(int16(123)))
	require.Equal(t, 5684921998043602808, objects.HashCode(int32(123)))
	require.Equal(t, 8114084241814008568, objects.HashCode(int64(123)))
	require.Equal(t, 8114084241814008568, objects.HashCode(uint(123)))
	require.Equal(t, 4953291999513788008, objects.HashCode(uint8(123)))
	require.Equal(t, -7208379536669766984, objects.HashCode(uint16(123)))
	require.Equal(t, 5684921998043602808, objects.HashCode(uint32(123)))
	require.Equal(t, 8114084241814008568, objects.HashCode(uint64(123)))
	require.Equal(t, 3326946953198266991, objects.HashCode(float32(123)))
	require.Equal(t, 2082925262096509245, objects.HashCode(float64(123)))
	require.Equal(t, 4591456448132287629, objects.HashCode("123"))
	type Port int8
	require.Equal(t, 4953291999513788008, objects.HashCode(Port(123)))
	require.Equal(t, 4953162257141659110, objects.HashCode(true))

	// Equal
	require.Equal(t, objects.HashCode(255), objects.HashCode(255))

	time1 := time.Now()
	time2 := time1
	require.Equal(t, objects.HashCode(time1), objects.HashCode(time2))

	slice1 := make([]any, 0)
	slice2 := slice1
	require.Equal(t, objects.HashCode(slice1), objects.HashCode(slice2))

	require.Equal(t, objects.HashCode(slice1), objects.HashCode(make([]any, 0)))

	// Not eqal
	require.NotEqual(t, objects.HashCode(255), objects.HashCode(256))

	slice1 = make([]any, 0)
	slice2 = append(slice1, 1)
	require.NotEqual(t, objects.HashCode(slice1), objects.HashCode(slice2))
}

func TestFirstNonNil(t *testing.T) {
	var v1 int
	v2 := 42
	assert.Equal(t, 0, *objects.FirstNonNil(nil, &v1, &v2))
}

func TestFirstNonZero(t *testing.T) {
	var v1 int
	v2 := 42
	assert.Equal(t, 0, *objects.FirstNonZero(nil, &v1, &v2))
	assert.Equal(t, "0", objects.FirstNonZero("", "0", "1"))
}
