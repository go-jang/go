package util_test

import (
	"testing"

	"github.com/go-jang/go/util"
	"github.com/stretchr/testify/require"
)

func TestArrayList_SizeAndEmpty(t *testing.T) {
	l := util.NewArrayList[int]()
	require.Equal(t, 0, l.Size())
	require.True(t, l.Empty())

	l.Add(1)
	require.Equal(t, 1, l.Size())
	require.False(t, l.Empty())
}

func TestArrayList_At(t *testing.T) {
	l := util.ArrayListOf(10, 20, 30)
	require.Equal(t, 10, l.At(0))
	require.Equal(t, 20, l.At(1))
	require.Equal(t, 30, l.At(2))

	require.Panics(t, func() { l.At(-1) })
	require.Panics(t, func() { l.At(3) })
}

func TestArrayList_Set(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3)
	old := l.Set(1, 42)
	require.Equal(t, 2, old)
	require.Equal(t, []int{1, 42, 3}, l.ToSlice())

	require.Panics(t, func() { l.Set(-1, 0) })
	require.Panics(t, func() { l.Set(3, 0) })
}

func TestArrayList_Add(t *testing.T) {
	l := util.NewArrayList[int]()
	ok := l.Add(1)
	require.True(t, ok)
	require.Equal(t, []int{1}, l.ToSlice())

	l.Add(2)
	l.Add(3)
	require.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestArrayList_AddAt(t *testing.T) {
	l := util.ArrayListOf(1, 3, 4)

	l.AddAt(1, 2)
	require.Equal(t, []int{1, 2, 3, 4}, l.ToSlice())

	l.AddAt(l.Size(), 5)
	require.Equal(t, []int{1, 2, 3, 4, 5}, l.ToSlice())

	l.AddAt(0, 0)
	require.Equal(t, []int{0, 1, 2, 3, 4, 5}, l.ToSlice())

	require.Panics(t, func() { l.AddAt(-1, 100) })
	require.Panics(t, func() { l.AddAt(l.Size()+1, 100) })
}

func TestArrayList_Remove(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3, 2)

	ok := l.Remove(2)
	require.True(t, ok)
	require.Equal(t, []int{1, 3, 2}, l.ToSlice())

	ok = l.Remove(99)
	require.False(t, ok)
	require.Equal(t, []int{1, 3, 2}, l.ToSlice())
}

func TestArrayList_RemoveAt(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3, 4)

	e := l.RemoveAt(1)
	require.Equal(t, 2, e)
	require.Equal(t, []int{1, 3, 4}, l.ToSlice())

	e = l.RemoveAt(0)
	require.Equal(t, 1, e)
	require.Equal(t, []int{3, 4}, l.ToSlice())

	e = l.RemoveAt(l.Size() - 1)
	require.Equal(t, 4, e)
	require.Equal(t, []int{3}, l.ToSlice())

	require.Panics(t, func() { l.RemoveAt(-1) })
	require.Panics(t, func() { l.RemoveAt(1) })
}

func TestArrayList_RemoveBy_RemoveEvenValues(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3, 4, 5, 6)

	removed := l.RemoveBy(func(i int, v int) bool { return v%2 == 0 })

	require.True(t, removed)
	require.Equal(t, []int{1, 3, 5}, l.ToSlice())
}

func TestArrayList_ContainsAndIndexOf(t *testing.T) {
	l := util.ArrayListOf(10, 20, 30, 20)

	require.True(t, l.Contains(10))
	require.True(t, l.Contains(20))
	require.False(t, l.Contains(99))

	require.Equal(t, 0, l.IndexOf(10))
	require.Equal(t, 1, l.IndexOf(20))
	require.Equal(t, -1, l.IndexOf(99))

	require.Equal(t, 3, l.LastIndexOf(20))
	require.Equal(t, 0, l.LastIndexOf(10))
	require.Equal(t, -1, l.LastIndexOf(99))
}

func TestArrayList_AddFirstAndAddLast(t *testing.T) {
	l := util.NewArrayList[int]()
	l.AddLast(2)
	l.AddFirst(1)
	l.AddLast(3)

	require.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestArrayList_FirstAndLast(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3)
	require.Equal(t, 1, l.First())
	require.Equal(t, 3, l.Last())

	empty := util.NewArrayList[int]()
	require.Panics(t, func() { empty.First() })
	require.Panics(t, func() { empty.Last() })
}

func TestArrayList_RemoveFirstAndRemoveLast(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3)

	f := l.RemoveFirst()
	require.Equal(t, 1, f)
	require.Equal(t, []int{2, 3}, l.ToSlice())

	l2 := util.ArrayListOf(10, 20, 30)
	lr := l2.RemoveLast()
	require.Equal(t, 30, lr)
	require.Equal(t, []int{10, 20}, l2.ToSlice())

	empty := util.NewArrayList[int]()
	require.Panics(t, func() { empty.RemoveFirst() })
	require.Panics(t, func() { empty.RemoveLast() })
}

func TestArrayList_ToSlice(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3)
	s := l.ToSlice()

	require.Equal(t, []int{1, 2, 3}, s)

	s[0] = 99
	require.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestArrayList_ContainsAll(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3, 4)
	c1 := util.ArrayListOf(1, 2)
	c2 := util.ArrayListOf(2, 5)

	require.True(t, l.ContainsAll(c1))
	require.False(t, l.ContainsAll(c2))
}

func TestArrayList_AddAll(t *testing.T) {
	l := util.ArrayListOf(1, 2)
	c := util.ArrayListOf(3, 4)

	ok := l.AddAll(c)
	require.True(t, ok)
	require.Equal(t, []int{1, 2, 3, 4}, l.ToSlice())

	empty := util.NewArrayList[int]()
	ok = l.AddAll(empty)
	require.False(t, ok)
}

func TestArrayList_RemoveAll(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3, 2, 4)
	c := util.ArrayListOf(2, 4)

	ok := l.RemoveAll(c)
	require.True(t, ok)
	require.Equal(t, []int{1, 3}, l.ToSlice())

	l2 := util.ArrayListOf(1, 2, 3)
	empty := util.NewArrayList[int]()
	ok = l2.RemoveAll(empty)
	require.False(t, ok)
	require.Equal(t, []int{1, 2, 3}, l2.ToSlice())
}

func TestArrayList_RetainAll(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3, 4, 5)
	c := util.ArrayListOf(2, 4, 6)

	ok := l.RetainAll(c)
	require.True(t, ok)
	require.Equal(t, []int{2, 4}, l.ToSlice())
}

func TestArrayList_SubList(t *testing.T) {
	l := util.ArrayListOf(0, 1, 2, 3, 4)

	sub := l.SubList(1, 4)
	require.Equal(t, []int{1, 2, 3}, sub.ToSlice())

	empty := l.SubList(2, 2)
	require.True(t, empty.Empty())

	require.Panics(t, func() { l.SubList(-1, 2) })
	require.Panics(t, func() { l.SubList(0, 10) })
}

func TestArrayList_Reversed(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3, 4)
	r := l.Reversed()

	require.Equal(t, []int{4, 3, 2, 1}, r.ToSlice())
	require.Equal(t, []int{1, 2, 3, 4}, l.ToSlice())

	empty := util.NewArrayList[int]()
	r2 := empty.Reversed()
	require.True(t, r2.Empty())
}

func TestArrayList_Sort(t *testing.T) {
	l := util.ArrayListOf(5, 3, 4, 1, 2)

	l.Sort(func(a, b int) int {
		return a - b
	})
	require.Equal(t, []int{1, 2, 3, 4, 5}, l.ToSlice())

	l2 := util.ArrayListOf(10)
	l2.Sort(func(a, b int) int { return a - b })
	require.Equal(t, []int{10}, l2.ToSlice())
}

func TestArrayList_Clear(t *testing.T) {
	l := util.ArrayListOf(1, 2, 3)
	require.False(t, l.Empty())

	l.Clear()
	require.True(t, l.Empty())
	require.Equal(t, 0, l.Size())
	require.Equal(t, 0, len(l.ToSlice()))
}
