package util

import "github.com/go-jang/go/util/objects"

/*
HashSet is a hash table based set implementation.

It stores unique elements and uses objects.HashCode and objects.Equal
to resolve equality and bucket placement.

The outer hash table growth is delegated to Go map runtime.
*/
type HashSet[E any] struct {
	table map[int][]E
	size  int
}

func NewHashSet[E any]() *HashSet[E] {
	return &HashSet[E]{
		table: make(map[int][]E),
	}
}

func HashSetOf[E any](vals ...E) *HashSet[E] {
	s := NewHashSet[E]()
	for _, v := range vals {
		s.Add(v)
	}
	return s
}

func (s *HashSet[E]) Size() int {
	return s.size
}

func (s *HashSet[E]) Empty() bool {
	return s.size == 0
}

func (s *HashSet[E]) Contains(o E) bool {
	hash := objects.HashCode(o)
	bucket := s.table[hash]
	for _, elem := range bucket {
		if objects.Equal(elem, o) {
			return true
		}
	}
	return false
}

func (s *HashSet[E]) ToSlice() []E {
	res := make([]E, 0, s.size)
	for _, bucket := range s.table {
		res = append(res, bucket...)
	}
	return res
}

func (s *HashSet[E]) Add(e E) bool {
	hash := objects.HashCode(e)
	bucket := s.table[hash]

	for _, elem := range bucket {
		if objects.Equal(elem, e) {
			return false
		}
	}

	s.table[hash] = append(bucket, e)
	s.size++
	return true
}

func (s *HashSet[E]) Remove(o E) bool {
	hash := objects.HashCode(o)
	bucket := s.table[hash]
	if len(bucket) == 0 {
		return false
	}

	for i, elem := range bucket {
		if objects.Equal(elem, o) {
			copy(bucket[i:], bucket[i+1:])

			var zero E
			bucket[len(bucket)-1] = zero
			bucket = bucket[:len(bucket)-1]

			if len(bucket) == 0 {
				delete(s.table, hash)
			} else {
				s.table[hash] = bucket
			}

			s.size--
			return true
		}
	}

	return false
}

func (s *HashSet[E]) Clear() {
	var zero E
	for hash, bucket := range s.table {
		for i := range bucket {
			bucket[i] = zero
		}
		delete(s.table, hash)
	}
	s.size = 0
}

func (s *HashSet[E]) ContainsAll(c Collection[E]) bool {
	for _, o := range c.ToSlice() {
		if !s.Contains(o) {
			return false
		}
	}
	return true
}

func (s *HashSet[E]) AddAll(c Collection[E]) bool {
	if c == nil || c.Size() == 0 {
		return false
	}

	modified := false
	for _, e := range c.ToSlice() {
		if s.Add(e) {
			modified = true
		}
	}
	return modified
}

func (s *HashSet[E]) RemoveAll(c Collection[E]) bool {
	if c == nil || c.Size() == 0 {
		return false
	}

	modified := false
	for _, e := range c.ToSlice() {
		if s.Remove(e) {
			modified = true
		}
	}
	return modified
}

func (s *HashSet[E]) RetainAll(c Collection[E]) bool {
	if c == nil {
		if s.Empty() {
			return false
		}
		s.Clear()
		return true
	}

	modified := false
	newTable := make(map[int][]E, len(s.table))
	newSize := 0

	for _, bucket := range s.table {
		for _, elem := range bucket {
			if c.Contains(elem) {
				hash := objects.HashCode(elem)
				newTable[hash] = append(newTable[hash], elem)
				newSize++
			} else {
				modified = true
			}
		}
	}

	if modified {
		s.table = newTable
		s.size = newSize
	}
	return modified
}

// Iterable
func (s *HashSet[E]) ForEach(accept func(E)) {
	for _, bucket := range s.table {
		for _, elem := range bucket {
			accept(elem)
		}
	}
}
