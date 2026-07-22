package stream

// From creates a stream from a slice (no copy).
func From[T any](slice []T) *referencePipeline[T] {
	return referencePipelineFrom(slice)
}

// Of creates a stream from values.
func Of[T any](args ...T) *referencePipeline[T] {
	return referencePipelineOf(args...)
}

// Map transforms elements and returns a new Stream.
func Map[T any, R any](s *referencePipeline[T], f func(T) R) *referencePipeline[R] {
	out := make([]R, len(s.ToSlice()))
	for i, v := range s.ToSlice() {
		out[i] = f(v)
	}
	return From(out)
}

// FlatMap maps each element to a slice and flattens the result.
func FlatMap[T any, R any](s *referencePipeline[T], f func(T) []R) *referencePipeline[R] {
	total := 0
	tmp := make([][]R, len(s.ToSlice()))
	for i, v := range s.ToSlice() {
		part := f(v)
		tmp[i] = part
		total += len(part)
	}
	out := make([]R, 0, total)
	for _, part := range tmp {
		out = append(out, part...)
	}
	return From(out)
}

// Reduce folds the stream left-to-right starting with init.
func Reduce[T any, R any](s *referencePipeline[T], init R, f func(acc R, v T) R) R {
	acc := init
	for _, v := range s.ToSlice() {
		acc = f(acc, v)
	}
	return acc
}

// Distinct removes duplicates using a map with a key selector.
func Distinct[T any, K comparable](s *referencePipeline[T], key func(T) K) *referencePipeline[T] {
	seen := make(map[K]struct{}, len(s.ToSlice()))
	out := make([]T, 0, len(s.ToSlice()))
	for _, v := range s.ToSlice() {
		k := key(v)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			out = append(out, v)
		}
	}
	return From(out)
}

// GroupingBy returns elements grouped by a key function.
func GroupingBy[T any, K comparable](s *referencePipeline[T], keyFunc func(T) K) map[K][]T {
	return collect(s, &collector[T, map[K][]T]{
		supplier: func() map[K][]T { return make(map[K][]T) },
		accumulator: func(acc map[K][]T, elem T) {
			k := keyFunc(elem)
			acc[k] = append(acc[k], elem)
		},
		finisher: nil,
	})
}

// ToSet collects unique elements.
func ToSet[T comparable](s *referencePipeline[T]) []T {
	type setSliceAcc[T comparable] struct {
		seen map[T]struct{}
		out  []T
	}
	collected := collect(s, &collector[T, *setSliceAcc[T]]{
		supplier: func() *setSliceAcc[T] {
			return &setSliceAcc[T]{
				seen: make(map[T]struct{}, len(s.data)),
				out:  make([]T, 0, len(s.data)),
			}
		},
		accumulator: func(acc *setSliceAcc[T], elem T) {
			if _, ok := acc.seen[elem]; ok {
				return
			}
			acc.seen[elem] = struct{}{}
			acc.out = append(acc.out, elem)
		},
		finisher: nil,
	})
	return collected.out
}
