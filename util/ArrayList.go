package util

import (
	"sort"

	"github.com/go-jang/go/util/objects"
)

type ArrayList[E any] struct {
	List[E]

	elementData []E
}

func NewArrayList[E any]() *ArrayList[E] {
	return &ArrayList[E]{
		elementData: make([]E, 0),
	}
}

func ArrayListOf[E any](vals ...E) *ArrayList[E] {
	l := NewArrayList[E]()
	l.elementData = make([]E, len(vals))
	copy(l.elementData, vals)
	return l
}

func (l *ArrayList[E]) At(index int) E {
	l.checkIndexOutOfBounds(index, l.Size()-1)
	return l.elementData[index]
}

func (l *ArrayList[E]) Set(index int, elem E) E {
	l.checkIndexOutOfBounds(index, l.Size()-1)
	old := l.elementData[index]
	l.elementData[index] = elem
	return old
}

func (l *ArrayList[E]) Size() int {
	return len(l.elementData)
}

func (l *ArrayList[E]) Empty() bool {
	return l.Size() == 0
}

func (l *ArrayList[E]) Contains(o E) bool {
	return l.IndexOf(o) >= 0
}

func (l *ArrayList[E]) ToSlice() []E {
	copySlice := make([]E, len(l.elementData))
	copy(copySlice, l.elementData)
	return copySlice
}

func (l *ArrayList[E]) Add(e E) bool {
	l.AddAt(l.Size(), e)
	return true
}

func (l *ArrayList[E]) AddAt(index int, e E) {
	l.checkIndexOutOfBounds(index, l.Size())
	l.elementData = append(l.elementData, e)
	if index == l.Size()-1 {
		return
	}
	n := len(l.elementData)
	copy(l.elementData[index+1:n], l.elementData[index:n-1])
	l.elementData[index] = e
}

func (l *ArrayList[E]) Remove(e E) bool {
	i := l.IndexOf(e)
	if i < 0 {
		return false
	}
	l.RemoveAt(i)
	return true
}

func (l *ArrayList[E]) RemoveAt(index int) E {
	l.checkIndexOutOfBounds(index, l.Size()-1)
	elem := l.elementData[index]
	copy(l.elementData[index:], l.elementData[index+1:])
	// clear last slot (avoid memory leaks for pointer types)
	var zero E
	l.elementData[len(l.elementData)-1] = zero
	// shrink slice length by 1
	l.elementData = l.elementData[:len(l.elementData)-1]
	return elem
}

func (l *ArrayList[E]) RemoveBy(predicate func(i int, value E) bool) bool {
	n := 0
	for i, v := range l.elementData {
		if !predicate(i, v) {
			l.elementData[n] = v
			n++
		}
	}

	removed := n < len(l.elementData)
	var zero E
	for i := n; i < len(l.elementData); i++ {
		l.elementData[i] = zero
	}

	l.elementData = l.elementData[:n]
	return removed
}

func (l *ArrayList[E]) IndexOf(o E) int {
	for i, elem := range l.elementData {
		if objects.Equal(o, elem) {
			return i
		}
	}
	return -1
}

func (l *ArrayList[E]) LastIndexOf(o E) int {
	for i := len(l.elementData) - 1; i >= 0; i-- {
		if objects.Equal(o, l.elementData[i]) {
			return i
		}
	}
	return -1
}

func (l *ArrayList[E]) AddFirst(element E) {
	l.AddAt(0, element)
}

func (l *ArrayList[E]) AddLast(element E) {
	l.Add(element)
}

func (l *ArrayList[E]) First() E {
	if l.Empty() {
		panic("No such element")
	}
	return l.At(0)
}

func (l *ArrayList[E]) Last() E {
	if l.Empty() {
		panic("No such element")
	}
	return l.At(l.Size() - 1)
}

func (l *ArrayList[E]) RemoveFirst() E {
	if l.Empty() {
		panic("No such element")
	}
	return l.RemoveAt(0)
}

func (l *ArrayList[E]) RemoveLast() E {
	if l.Empty() {
		panic("No such element")
	}
	return l.RemoveAt(l.Size() - 1)
}

func (l *ArrayList[E]) ContainsAll(c Collection[E]) bool {
	for _, o := range c.ToSlice() {
		if !l.Contains(o) {
			return false
		}
	}
	return true
}

func (l *ArrayList[E]) AddAll(c Collection[E]) bool {
	if c == nil || c.Size() == 0 {
		return false
	}
	for _, e := range c.ToSlice() {
		l.Add(e)
	}
	return true
}

func (l *ArrayList[E]) RemoveAll(c Collection[E]) bool {
	if c == nil || c.Size() == 0 {
		return false
	}

	modified := false
	// reuse underlying array
	dst := l.elementData[:0]
	for _, elem := range l.elementData {
		if c.Contains(elem) {
			modified = true
		} else {
			dst = append(dst, elem)
		}
	}

	// clear truncated tail to avoid keeping references
	for i := len(dst); i < len(l.elementData); i++ {
		var zero E
		l.elementData[i] = zero
	}
	l.elementData = dst

	return modified
}

func (l *ArrayList[E]) RetainAll(c Collection[E]) bool {
	if c == nil {
		// retain nothing
		if l.Empty() {
			return false
		}
		l.Clear()
		return true
	}

	modified := false
	dst := l.elementData[:0]

	for _, elem := range l.elementData {
		if c.Contains(elem) {
			// keep
			dst = append(dst, elem)
		} else {
			// drop
			modified = true
		}
	}

	// clear truncated tail
	for i := len(dst); i < len(l.elementData); i++ {
		var zero E
		l.elementData[i] = zero
	}
	l.elementData = dst

	return modified
}

func (l *ArrayList[E]) SubList(fromInclusive, toExclusive int) *ArrayList[E] {
	l.checkIndexOutOfBounds(fromInclusive, toExclusive)
	l.checkIndexOutOfBounds(toExclusive, l.Size())
	sub := NewArrayList[E]()
	n := toExclusive - fromInclusive
	if n == 0 {
		return sub
	}

	sub.elementData = make([]E, n)
	copy(sub.elementData, l.elementData[fromInclusive:toExclusive])
	return sub
}

func (l *ArrayList[E]) Reversed() *ArrayList[E] {
	res := NewArrayList[E]()
	n := l.Size()
	if n == 0 {
		return res
	}

	res.elementData = make([]E, n)
	for i := 0; i < n; i++ {
		res.elementData[i] = l.elementData[n-1-i]
	}
	return res
}

func (l *ArrayList[E]) Sort(cmp func(a, b E) int) {
	if l.Size() <= 1 {
		return
	}
	sort.Slice(l.elementData, func(i, j int) bool {
		return cmp(l.elementData[i], l.elementData[j]) < 0
	})
}

// Iterable
func (l *ArrayList[E]) ForEach(accept func(E)) {
	for _, v := range l.elementData {
		accept(v)
	}
}

func (l *ArrayList[E]) Clear() {
	// clear references for GC
	for i := range l.elementData {
		var zero E
		l.elementData[i] = zero
	}
	l.elementData = l.elementData[:0]
}

func (l *ArrayList[E]) TrimToSize() {
	if len(l.elementData) < cap(l.elementData) {
		newSlice := make([]E, len(l.elementData))
		copy(newSlice, l.elementData)
		l.elementData = newSlice
	}
}

func (l *ArrayList[E]) checkIndexOutOfBounds(i, lenInclusive int) {
	if i < 0 || i > lenInclusive {
		panic("index out of bounds")
	}
}
