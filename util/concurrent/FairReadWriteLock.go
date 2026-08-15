package concurrent

import (
	"sync"

	"github.com/go-jang/go/lang"
)

var _ sync.Locker = (*FairReadWriteLock)(nil)
var _ sync.Locker = (*fairReadLocker)(nil)

// FairReadWriteLock provides read/write locking with FIFO fairness between reader and writer groups.
// Consecutive readers may acquire the lock concurrently, while queued writers cannot be bypassed by later readers.
//
// Given arrival order: R1 R2 R3 W1 R4 R5 W2 R6
//
// guarantee: [R1 R2 R3] → W1 → [R4 R5] → W2 → [R6]
type FairReadWriteLock struct {
	mutex      sync.Mutex
	condition  *sync.Cond
	nextTicket uint64

	readerHead *fairReadWriteLockWaiter
	readerTail *fairReadWriteLockWaiter
	writerHead *fairReadWriteLockWaiter
	writerTail *fairReadWriteLockWaiter

	readers     int
	writeLocked bool
	queueLength int
}

func NewFairReadWriteLock() *FairReadWriteLock {
	result := &FairReadWriteLock{}
	result.condition = sync.NewCond(&result.mutex)
	return result
}

func (this *FairReadWriteLock) Lock() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	waiter := this.enqueueWriter()

	for !this.canWrite(waiter) {
		this.condition.Wait()
	}

	this.dequeueWriter()
	this.writeLocked = true
}

func (this *FairReadWriteLock) Unlock() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	lang.Assert(this.writeLocked, "unlock of unlocked FairReadWriteLock")

	this.writeLocked = false
	this.condition.Broadcast()
}

func (this *FairReadWriteLock) RLock() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	waiter := this.enqueueReader()

	for !this.canRead(waiter) {
		this.condition.Wait()
	}

	this.dequeueReader()
	this.readers++
	this.condition.Broadcast()
}

func (this *FairReadWriteLock) RUnlock() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	lang.Assert(this.readers > 0, "runlock of unlocked FairReadWriteLock")

	this.readers--
	if this.readers == 0 {
		this.condition.Broadcast()
	}
}

func (this *FairReadWriteLock) TryLock() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.writeLocked || this.readers > 0 || this.queueLength > 0 {
		return false
	}

	this.writeLocked = true
	return true
}

func (this *FairReadWriteLock) TryRLock() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.writeLocked || this.queueLength > 0 {
		return false
	}

	this.readers++
	return true
}

func (this *FairReadWriteLock) IsWriteLocked() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.writeLocked
}

func (this *FairReadWriteLock) ReadLockCount() int {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.readers
}

func (this *FairReadWriteLock) HasQueuedGoroutines() bool {
	return this.QueueLength() > 0
}

func (this *FairReadWriteLock) QueueLength() int {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.queueLength
}

func (this *FairReadWriteLock) RLocker() sync.Locker {
	return &fairReadLocker{lock: this}
}

func (this *FairReadWriteLock) canRead(waiter *fairReadWriteLockWaiter) bool {
	if this.writeLocked || this.readerHead != waiter {
		return false
	}

	return this.writerHead == nil || ticketBefore(waiter.ticket, this.writerHead.ticket)
}

func (this *FairReadWriteLock) canWrite(waiter *fairReadWriteLockWaiter) bool {
	if this.writeLocked || this.readers > 0 || this.writerHead != waiter {
		return false
	}

	return this.readerHead == nil || ticketBefore(waiter.ticket, this.readerHead.ticket)
}

func ticketBefore(left, right uint64) bool {
	return int64(left-right) < 0
}

func (this *FairReadWriteLock) enqueueReader() *fairReadWriteLockWaiter {
	waiter := this.newWaiter()

	if this.readerTail == nil {
		this.readerHead = waiter
		this.readerTail = waiter
	} else {
		this.readerTail.next = waiter
		this.readerTail = waiter
	}

	this.queueLength++
	return waiter
}

func (this *FairReadWriteLock) enqueueWriter() *fairReadWriteLockWaiter {
	waiter := this.newWaiter()

	if this.writerTail == nil {
		this.writerHead = waiter
		this.writerTail = waiter
	} else {
		this.writerTail.next = waiter
		this.writerTail = waiter
	}

	this.queueLength++
	return waiter
}

func (this *FairReadWriteLock) newWaiter() *fairReadWriteLockWaiter {
	result := &fairReadWriteLockWaiter{
		ticket: this.nextTicket,
	}
	this.nextTicket++
	return result
}

func (this *FairReadWriteLock) dequeueReader() {
	lang.Assert(this.readerHead != nil, "empty FairReadWriteLock reader queue")

	this.readerHead = this.readerHead.next
	this.queueLength--

	if this.readerHead == nil {
		this.readerTail = nil
	}
}

func (this *FairReadWriteLock) dequeueWriter() {
	lang.Assert(this.writerHead != nil, "empty FairReadWriteLock writer queue")

	this.writerHead = this.writerHead.next
	this.queueLength--

	if this.writerHead == nil {
		this.writerTail = nil
	}
}

type fairReadWriteLockWaiter struct {
	ticket uint64
	next   *fairReadWriteLockWaiter
}

type fairReadLocker struct {
	lock *FairReadWriteLock
}

func (this *fairReadLocker) Lock() {
	this.lock.RLock()
}

func (this *fairReadLocker) Unlock() {
	this.lock.RUnlock()
}
