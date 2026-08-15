package concurrent_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-jang/go/util/concurrent"
	"github.com/stretchr/testify/require"
)

func TestFairReadWriteLock(test *testing.T) {
	lock := concurrent.NewFairReadWriteLock()
	lock.Lock()

	var stage atomic.Int32
	var invalid atomic.Int32
	var waitGroup sync.WaitGroup

	readers := func(expected int32, count int) {
		for range count {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				lock.RLock()
				defer lock.RUnlock()

				if stage.Load() != expected {
					invalid.Add(1)
				}
			}()
		}
	}

	writer := func(expected, next int32) {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			lock.Lock()
			defer lock.Unlock()

			if stage.Load() != expected {
				invalid.Add(1)
			}
			stage.Store(next)
		}()
	}

	readers(0, 100)
	require.Eventually(test, func() bool { return lock.QueueLength() == 100 }, time.Second, time.Millisecond)

	writer(0, 1)
	require.Eventually(test, func() bool { return lock.QueueLength() == 101 }, time.Second, time.Millisecond)

	readers(1, 100)
	require.Eventually(test, func() bool { return lock.QueueLength() == 201 }, time.Second, time.Millisecond)

	writer(1, 2)
	require.Eventually(test, func() bool { return lock.QueueLength() == 202 }, time.Second, time.Millisecond)

	readers(2, 100)
	require.Eventually(test, func() bool { return lock.QueueLength() == 302 }, time.Second, time.Millisecond)

	lock.Unlock()
	waitGroup.Wait()

	require.Zero(test, invalid.Load())
	require.Equal(test, int32(2), stage.Load())
	require.Equal(test, 0, lock.QueueLength())
}
