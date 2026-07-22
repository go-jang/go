package util

/*
Set is a collection that contains no duplicate elements.

Add returns false if the element is already present.
Element equality is defined by the concrete implementation.
*/
type Set[E any] interface {
	Collection[E]
}
