package util

type List[E any] interface {
	Collection[E]

	AddFirst(e E)
	AddLast(e E)
	First() E
	Last() E
	At(index int) E
	Set(index int, element E) E
	AddAt(index int, element E)
	RemoveAt(index int) E
	IndexOf(element any) int
	LastIndexOf(element E) int
	SubList(fromIndex, toIndex int) List[E]
	Reversed() List[E]
	RemoveFirst() E
	RemoveLast() E
	Sort(func(i, j E) int)
}
