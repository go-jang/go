package stream

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-jang/go/util/optional"
)

// referencePipeline implements limitted stream API over a slice.
// Beware of unnececery calculations. Stage process all the elements before moving to the next processing stage unlike real stream.
type referencePipeline[T any] struct {
	data []T
}

// From creates a stream from a slice (no copy).
func referencePipelineFrom[T any](slice []T) *referencePipeline[T] {
	return &referencePipeline[T]{data: slice}
}

// Of creates a stream from values.
func referencePipelineOf[T any](args ...T) *referencePipeline[T] {
	return &referencePipeline[T]{data: args}
}

// ToSlice returns the underlying data (no copy).
func (s *referencePipeline[T]) ToSlice() []T {
	return s.data
}

// ForEach applies f to each element.
func (s *referencePipeline[T]) ForEach(f func(T)) {
	for _, v := range s.data {
		f(v)
	}
}

// Filter keeps elements where predicate returns true.
func (s *referencePipeline[T]) Filter(f func(T) bool) *referencePipeline[T] {
	out := make([]T, 0, len(s.data))
	for _, v := range s.data {
		if f(v) {
			out = append(out, v)
		}
	}
	return &referencePipeline[T]{data: out}
}

// FlatMap maps each element to a slice and flattens the result.
//
// See also: func Map[T any, R any](s *referencePipeline[T], f func(T) R) *referencePipeline[R]
func (s *referencePipeline[T]) Map(f func(T) T) *referencePipeline[T] {
	out := make([]T, len(s.ToSlice()))
	for i, v := range s.ToSlice() {
		out[i] = f(v)
	}
	return From(out)
}

// Peek applies action to each element and returns the same stream for chaining.
func (s *referencePipeline[T]) Peek(action func(T)) *referencePipeline[T] {
	for _, v := range s.data {
		action(v)
	}
	return s
}

// Limit returns the first n elements (or all if n > len).
func (s *referencePipeline[T]) Limit(n int) *referencePipeline[T] {
	if n >= len(s.data) {
		return s
	}
	if n <= 0 {
		return &referencePipeline[T]{data: nil}
	}
	return &referencePipeline[T]{data: s.data[:n]}
}

// Skip the first n elements.
func (s *referencePipeline[T]) Skip(n int) *referencePipeline[T] {
	if n <= 0 {
		return s
	}
	if n >= len(s.data) {
		return &referencePipeline[T]{data: nil}
	}
	return &referencePipeline[T]{data: s.data[n:]}
}

// Sorted returns a new ReferencePipeline with elements sorted according to the less function.
func (s *referencePipeline[T]) Sorted(less func(a, b T) bool) *referencePipeline[T] {
	// Make a copy of the slice so we don't mutate the original
	data := make([]T, len(s.data))
	copy(data, s.data)

	sort.Slice(data, func(i, j int) bool {
		return less(data[i], data[j])
	})

	return &referencePipeline[T]{data: data}
}

// FindFirst returns the first element as Optional[T], or empty if none.
func (s *referencePipeline[T]) FindFirst() *optional.Optional[T] {
	if len(s.data) == 0 {
		return optional.OfEmpty[T]()
	}
	return optional.OfValue(s.data[0])
}

// FindFirstPredicate returns the first element satisfying predicate.
func (s *referencePipeline[T]) FindFirstPredicate(f func(T) bool) *optional.Optional[T] {
	for _, v := range s.data {
		if f(v) {
			return optional.OfValue(v)
		}
	}
	return optional.OfEmpty[T]()
}

// AllMatch returns true if all elements of the stream match the predicate.
// Returns true for an empty stream.
func (s *referencePipeline[T]) AllMatch(pred func(T) bool) bool {
	for _, v := range s.data {
		if !pred(v) {
			return false
		}
	}
	return true
}

// AnyMatch returns true if any element of the stream matches the predicate.
// Returns false for an empty stream.
func (s *referencePipeline[T]) AnyMatch(pred func(T) bool) bool {
	for _, v := range s.data {
		if pred(v) {
			return true
		}
	}
	return false
}

// NoneMatch returns true if no elements of the stream match the predicate.
// Returns true for an empty stream.
func (s *referencePipeline[T]) NoneMatch(pred func(T) bool) bool {
	for _, v := range s.data {
		if pred(v) {
			return false
		}
	}
	return true
}

// Min returns the minimum element according to less function.
func (s *referencePipeline[T]) Min(less func(a, b T) bool) *optional.Optional[T] {
	if len(s.data) == 0 {
		return optional.OfEmpty[T]()
	}

	min := s.data[0]
	for _, v := range s.data[1:] {
		if less(v, min) {
			min = v
		}
	}
	return optional.OfValue(min)
}

// Max returns the maximum element according to less function.
func (s *referencePipeline[T]) Max(less func(a, b T) bool) *optional.Optional[T] {
	if len(s.data) == 0 {
		return optional.OfEmpty[T]()
	}

	max := s.data[0]
	for _, v := range s.data[1:] {
		if less(max, v) {
			max = v
		}
	}
	return optional.OfValue(max)
}

// Join concatenates elements into a string with optional prefix/suffix.
//
// Usage:
//
//	s.Join(",")                 -> "a,b,c"
//	s.Join(",", "[")            -> "[a,b,c"
//	s.Join(",", "[", "]")       -> "[a,b,c]"
func (p *referencePipeline[T]) Join(delim string, parts ...string) string {
	var prefix, suffix string
	switch len(parts) {
	case 1:
		prefix = parts[0]
	case 2:
		prefix = parts[0]
		suffix = parts[1]
	}

	var b strings.Builder
	b.WriteString(prefix)

	for i, v := range p.ToSlice() {
		if i > 0 {
			b.WriteString(delim)
		}
		b.WriteString(fmt.Sprint(v))
	}

	b.WriteString(suffix)
	return b.String()
}

// Count returns the number of elements in the stream.
func (s *referencePipeline[T]) Count() int {
	return len(s.data)
}
