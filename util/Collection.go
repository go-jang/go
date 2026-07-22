package util

type Collection[E any] interface {
	Iterable[E]

	Size() int
	Empty() bool
	Contains(o E) bool
	ToSlice() []E

	Add(e E) bool
	Remove(o E) bool
	Clear()

	ContainsAll(c Collection[E]) bool
	AddAll(c Collection[E]) bool
	RemoveAll(c Collection[E]) bool
	RetainAll(c Collection[E]) bool
}
