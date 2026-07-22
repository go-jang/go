package util

type Iterable[E any] interface {
	ForEach(accept func(E))
}
