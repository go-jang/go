package collections_test

import (
	"testing"

	"github.com/go-jang/go/util/collections"
	"github.com/stretchr/testify/require"
)

func Test_Collections_SubtractSlice(t *testing.T) {
	t.Run("should produce sane results'", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{4, 3, 2}
		require.Equal(t, []int{1, 5}, collections.SubtractSlice(a, b))
	})
}

func Test_Collections_ReverseSlice(t *testing.T) {
	t.Run("should produce sane results'", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		require.Equal(t, []int{5, 4, 3, 2, 1}, collections.ReverseSlice(a))
	})
}

func TestCollectionsSort(t *testing.T) {
	t.Run("should produce sane results'", func(t *testing.T) {
		a := []int{1, 3, 5, 2, 4}
		require.Equal(t, []int{1, 2, 3, 4, 5}, collections.Sort(a))
	})
}

func TestCollectionsDistinct(t *testing.T) {
	t.Run("should produce sane results'", func(t *testing.T) {
		a := []int{1, 2, 2, 3, 4, 5, 2, 4, 3}
		require.Equal(t, []int{1, 2, 3, 4, 5}, collections.Distinct(a))
	})
}
