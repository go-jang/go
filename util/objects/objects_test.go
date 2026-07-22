package objects_test

import (
	"testing"
	"time"

	"github.com/go-jang/go/util/objects"

	"github.com/stretchr/testify/require"
)

func TestToString(t *testing.T) {
	t.Run("should produce sane results'", func(t *testing.T) {
		require.Equal(t, "123", objects.ToString(int(123)))
		require.Equal(t, "123", objects.ToString(int8(123)))
		require.Equal(t, "123", objects.ToString(int16(123)))
		require.Equal(t, "123", objects.ToString(int32(123)))
		require.Equal(t, "123", objects.ToString(int64(123)))
		require.Equal(t, "123", objects.ToString(uint(123)))
		require.Equal(t, "123", objects.ToString(uint8(123)))
		require.Equal(t, "123", objects.ToString(uint16(123)))
		require.Equal(t, "123", objects.ToString(uint32(123)))
		require.Equal(t, "123", objects.ToString(uint64(123)))
		require.Equal(t, "123", objects.ToString(float32(123)))
		require.Equal(t, "123", objects.ToString(float64(123)))
		require.Equal(t, "123", objects.ToString("123"))
		type Port int8
		require.Equal(t, "123", objects.ToString(Port(123)))
		require.Equal(t, "true", objects.ToString(true))
	})
}

func Test_Equal(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
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
	})
}

func Test_HashCode(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
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
	})
}
