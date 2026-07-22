package util_test

import (
	"testing"

	"github.com/go-jang/go/util"
	"github.com/stretchr/testify/assert"
)

func TestNewHashSet_Empty(t *testing.T) {
	s := util.NewHashSet[string]()

	assert.NotNil(t, s)
	assert.Equal(t, 0, s.Size())
	assert.True(t, s.Empty())
	assert.False(t, s.Contains("a"))
	assert.Empty(t, s.ToSlice())
}

func TestHashSetOf(t *testing.T) {
	s := util.HashSetOf("a", "b", "a")

	assert.Equal(t, 2, s.Size())
	assert.True(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
	assert.False(t, s.Contains("c"))
}

func TestHashSet_Add(t *testing.T) {
	s := util.NewHashSet[string]()

	assert.True(t, s.Add("a"))
	assert.False(t, s.Add("a"))
	assert.True(t, s.Add("b"))

	assert.Equal(t, 2, s.Size())
	assert.True(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
}

func TestHashSet_Remove(t *testing.T) {
	s := util.HashSetOf("a", "b")

	assert.True(t, s.Remove("a"))
	assert.False(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
	assert.Equal(t, 1, s.Size())

	assert.False(t, s.Remove("a"))
	assert.False(t, s.Remove("c"))
	assert.Equal(t, 1, s.Size())
}

func TestHashSet_Clear(t *testing.T) {
	s := util.HashSetOf("a", "b")

	s.Clear()

	assert.Equal(t, 0, s.Size())
	assert.True(t, s.Empty())
	assert.False(t, s.Contains("a"))
	assert.False(t, s.Contains("b"))
	assert.Empty(t, s.ToSlice())
}

func TestHashSet_ContainsAll(t *testing.T) {
	s := util.HashSetOf("a", "b", "c")

	assert.True(t, s.ContainsAll(util.HashSetOf("a", "b")))
	assert.True(t, s.ContainsAll(util.NewHashSet[string]()))
	assert.False(t, s.ContainsAll(util.HashSetOf("a", "x")))
}

func TestHashSet_AddAll(t *testing.T) {
	s := util.HashSetOf("a")

	assert.False(t, s.AddAll(nil))
	assert.False(t, s.AddAll(util.NewHashSet[string]()))

	assert.True(t, s.AddAll(util.HashSetOf("a", "b", "c")))
	assert.Equal(t, 3, s.Size())
	assert.True(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
	assert.True(t, s.Contains("c"))

	assert.False(t, s.AddAll(util.HashSetOf("a", "b", "c")))
	assert.Equal(t, 3, s.Size())
}

func TestHashSet_RemoveAll(t *testing.T) {
	s := util.HashSetOf("a", "b", "c", "d")

	assert.False(t, s.RemoveAll(nil))
	assert.False(t, s.RemoveAll(util.NewHashSet[string]()))

	assert.True(t, s.RemoveAll(util.HashSetOf("b", "c", "x")))
	assert.Equal(t, 2, s.Size())
	assert.True(t, s.Contains("a"))
	assert.False(t, s.Contains("b"))
	assert.False(t, s.Contains("c"))
	assert.True(t, s.Contains("d"))

	assert.False(t, s.RemoveAll(util.HashSetOf("b", "c", "x")))
}

func TestHashSet_RetainAll(t *testing.T) {
	s := util.HashSetOf("a", "b", "c", "d")

	assert.True(t, s.RetainAll(util.HashSetOf("b", "d", "x")))
	assert.Equal(t, 2, s.Size())
	assert.False(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
	assert.False(t, s.Contains("c"))
	assert.True(t, s.Contains("d"))

	assert.False(t, s.RetainAll(util.HashSetOf("b", "d")))
	assert.Equal(t, 2, s.Size())
}

func TestHashSet_RetainAll_Nil(t *testing.T) {
	s := util.HashSetOf("a", "b")

	assert.True(t, s.RetainAll(nil))
	assert.True(t, s.Empty())
	assert.Equal(t, 0, s.Size())

	assert.False(t, s.RetainAll(nil))
}

func TestHashSet_ForEach(t *testing.T) {
	s := util.HashSetOf("a", "b", "c")
	visited := util.NewHashSet[string]()

	s.ForEach(func(v string) {
		visited.Add(v)
	})

	assert.Equal(t, 3, visited.Size())
	assert.True(t, visited.Contains("a"))
	assert.True(t, visited.Contains("b"))
	assert.True(t, visited.Contains("c"))
}

func TestHashSet_ToSlice(t *testing.T) {
	s := util.HashSetOf("a", "b", "c")

	values := s.ToSlice()
	actual := util.HashSetOf(values...)

	assert.Equal(t, 3, len(values))
	assert.Equal(t, 3, actual.Size())
	assert.True(t, actual.Contains("a"))
	assert.True(t, actual.Contains("b"))
	assert.True(t, actual.Contains("c"))
}
