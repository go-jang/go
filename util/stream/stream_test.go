package stream_test

import (
	"math"
	"testing"

	"github.com/go-jang/go/util/stream"
	"github.com/stretchr/testify/require"
)

func TestStreamFrom(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		str1 := stream.From([]int{1, 2, 3})
		str2 := stream.Of(1, 2, 3)
		require.Equal(t, str1.ToSlice(), str2.ToSlice())
	})
}

func TestStreamFilter(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		evens := s.Filter(func(v int) bool { return v%2 == 0 })
		require.Equal(t, []int{2, 4}, evens.ToSlice())
	})
}

func TestStreamMapTT(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		squares := s.Map(func(v int) int { return v * v })
		require.Equal(t, []int{1, 4, 9, 16, 25}, squares.ToSlice())
	})
}

func TestStreamMapTR(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		sqrts := stream.Map(s, func(v int) float64 { return math.Sqrt(float64(v)) })
		require.Equal(t, []float64{1, 1.4142135623730951, 1.7320508075688772, 2, 2.23606797749979}, sqrts.ToSlice())
	})
}

func TestStreamLimit(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		require.Equal(t, []int{1, 2}, s.Limit(2).ToSlice())
	})
}

func TestStreamSkip(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		require.Equal(t, []int{4, 5}, s.Skip(3).ToSlice())
	})
}

func TestStreamReduce(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		sum := stream.Reduce(s, 0, func(acc, v int) int { return acc + v })
		require.Equal(t, 15, sum)
	})
}

func TestStreamFlatMap(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		pairs := stream.FlatMap(s, func(i int) []int { return []int{i, -i} })
		require.Equal(t, []int{1, -1, 2, -2, 3, -3, 4, -4, 5, -5}, pairs.ToSlice())
	})
}

func TestStreamDistinct(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		type user struct {
			Id   int
			Name string
		}
		users := stream.Of(
			user{1, "Ann"}, user{2, "Bob"}, user{1, "Ann2"},
		)
		unique := stream.Distinct(users, func(u user) int { return u.Id })
		require.Equal(t, []user{{1, "Ann"}, {2, "Bob"}}, unique.ToSlice())
	})
}

func TestStreamForEach(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		sum := 0
		s.ForEach(func(v int) { sum += v })
		require.Equal(t, 15, sum)
	})
}

func TestStreamSorted(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		require.Equal(t, []int{1, 2, 3, 4, 5}, s.ToSlice())
	})
}

func TestStreamFindFirst(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		require.Equal(t, true, s.FindFirst().Present())
		require.Equal(t, 1, s.FindFirst().Value())
	})
}

func TestStreamFindFirstPredicate(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		firstEven := s.FindFirstPredicate(func(v int) bool { return v%2 == 0 })
		require.Equal(t, true, firstEven.Present())
		require.Equal(t, 2, firstEven.Value())
	})
}

func TestStreamAllMatch(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s1 := stream.Of(2, 4, 6)
		s2 := stream.Of(2, 3, 6)
		s3 := stream.Of[int]() // empty stream

		require.Equal(t, true, s1.AllMatch(func(v int) bool { return v%2 == 0 }))
		require.Equal(t, false, s2.AllMatch(func(v int) bool { return v%2 == 0 }))
		require.Equal(t, true, s3.AllMatch(func(v int) bool { return v%2 == 0 }))
	})
}

func TestStreamAnyMatch(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s1 := stream.Of(1, 3, 5)
		s2 := stream.Of(1, 2, 5)
		s3 := stream.Of[int]() // empty stream

		require.Equal(t, false, s1.AnyMatch(func(v int) bool { return v%2 == 0 }))
		require.Equal(t, true, s2.AnyMatch(func(v int) bool { return v%2 == 0 }))
		require.Equal(t, false, s3.AnyMatch(func(v int) bool { return v%2 == 0 }))
	})
}

func TestStreamNoneMatch(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s1 := stream.Of(1, 3, 5)
		s2 := stream.Of(1, 2, 5)
		s3 := stream.Of[int]() // empty stream

		require.Equal(t, true, s1.NoneMatch(func(v int) bool { return v%2 == 0 }))
		require.Equal(t, false, s2.NoneMatch(func(v int) bool { return v%2 == 0 }))
		require.Equal(t, true, s3.NoneMatch(func(v int) bool { return v%2 == 0 }))
	})
}

func TestStreamMinMax(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(3, 1, 4, 1, 5, 9, 2)

		min := s.Min(func(a, b int) bool { return a < b })
		require.Equal(t, 1, min.Value())

		max := s.Max(func(a, b int) bool { return a < b })
		require.Equal(t, 9, max.Value())
	})
}

func TestStreamPeek(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		sum := 0
		s.Peek(func(v int) {
			sum += v
		}).ForEach(func(v int) {
			// Could do something else, but Peek already ran
		})
		require.Equal(t, 15, sum)
	})
}
func TestStreamJoin(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		require.Equal(t, "[1, 2, 3, 4, 5]", s.Join(", ", "[", "]"))
		require.Equal(t, "1,2,3,4,5", s.Join(","))
	})
}

func TestStreamCount(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 3, 4, 5)
		require.Equal(t, 5, s.Count())
	})
}

func TestStreamGroupingBy(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of("apple", "banana", "apricot", "blueberry", "avocado")
		grouped := stream.GroupingBy(s, func(s string) string {
			return string(s[0]) // group by first letter
		})
		require.Equal(t, map[string][]string{"a": {"apple", "apricot", "avocado"}, "b": {"banana", "blueberry"}}, grouped)
	})
}

func TestStreamToSet(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		s := stream.Of(1, 2, 2, 3, 3, 3)
		set := stream.ToSet(s)
		require.Equal(t, []int{1, 2, 3}, set)
	})
}
