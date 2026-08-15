package concurrent

import (
	"sync"

	"github.com/go-jang/go/lang"
)

var _ sync.Locker = (*FairLock)(nil)

// FairLock provides mutually exclusive locking with FIFO fairness.
// Goroutines acquire the lock in the order in which they enter the wait queue.
type FairLock struct {
	mutex      sync.Mutex
	condition  *sync.Cond
	nextTicket uint64
	serving    uint64
	locked     bool
}

func NewFairLock() *FairLock {
	result := &FairLock{}
	result.condition = sync.NewCond(&result.mutex)
	return result
}

func (this *FairLock) Lock() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	ticket := this.nextTicket
	this.nextTicket++

	for ticket != this.serving || this.locked {
		this.condition.Wait()
	}

	this.locked = true
}

func (this *FairLock) TryLock() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.locked || this.nextTicket != this.serving {
		return false
	}

	this.nextTicket++
	this.locked = true
	return true
}

func (this *FairLock) Unlock() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	lang.Assert(this.locked, "unlock of unlocked FairLock")

	this.serving++
	this.locked = false
	this.condition.Broadcast()
}

func (this *FairLock) IsLocked() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.locked
}

func (this *FairLock) HasQueuedGoroutines() bool {
	return this.QueueLength() > 0
}

func (this *FairLock) QueueLength() int {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	result := this.nextTicket - this.serving
	if this.locked {
		result--
	}
	return int(result)
}
