package stream

type collector[T any, R any] struct {
	supplier    func() R
	accumulator func(R, T)
	finisher    func(R) R
}

// collect consumes a stream into the Collector's result type
func collect[T any, R any](s *referencePipeline[T], c *collector[T, R]) R {
	result := c.supplier()
	for _, v := range s.data {
		c.accumulator(result, v)
	}
	if c.finisher != nil {
		return c.finisher(result)
	}
	return result
}
