package concurrent_test

import (
	"sync"
	"testing"
	"time"

	"github.com/go-jang/go/util/concurrent"
	"github.com/stretchr/testify/require"
)

func TestFairLock(test *testing.T) {
	lock := concurrent.NewFairLock()
	lock.Lock()

	result := make([]int, 0, 1000)
	var waitGroup sync.WaitGroup

	for i := 1; i <= 1000; i++ {
		waitGroup.Add(1)

		go func(value int) {
			defer waitGroup.Done()

			lock.Lock()
			defer lock.Unlock()

			result = append(result, value)
		}(i)

		time.Sleep(time.Millisecond)
	}

	require.True(test, lock.IsLocked())
	require.Equal(test, 1000, lock.QueueLength())
	require.True(test, lock.HasQueuedGoroutines())

	lock.Unlock()
	waitGroup.Wait()

	require.Len(test, result, 1000)
	for i := 1; i <= 1000; i++ {
		require.Equal(test, i, result[i-1])
	}

	require.False(test, lock.IsLocked())
	require.Equal(test, 0, lock.QueueLength())
	require.False(test, lock.HasQueuedGoroutines())
}

func TestFairLockTryLock(test *testing.T) {
	lock := concurrent.NewFairLock()

	require.True(test, lock.TryLock())
	require.True(test, lock.IsLocked())
	require.False(test, lock.TryLock())

	lock.Unlock()

	require.True(test, lock.TryLock())
	lock.Unlock()
}
