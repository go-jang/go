package concurrent

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/go-errr/go/err"
)

type Future[T any] interface {
	Get() T
	Result() (T, error)
}

type futureImpl[T any] struct {
	once   sync.Once
	result chan resultWrapper[T]
	val    T
	err    any
}

type resultWrapper[T any] struct {
	val T
	err any
}

func (f *futureImpl[T]) wait() {
	f.once.Do(func() {
		r := <-f.result
		f.val, f.err = r.val, r.err
	})
}

func (f *futureImpl[T]) Get() T {
	f.wait()
	if f.err != nil {
		panic(NewExecutionExceptionFrom(fmt.Sprint(f.err), f.err))
	}
	return f.val
}

func (f *futureImpl[T]) Result() (T, error) {
	f.wait()
	if f.err == nil {
		return f.val, nil
	} else {
		return f.val, NewExecutionExceptionFrom(fmt.Sprint(f.err), f.err)
	}
}

var defaultExecutor *Executor
var defaultExecutorMu sync.Mutex

// Executor is a fixed-size worker pool for asynchronous task execution.
//
// A configured number of worker goroutines is started eagerly and
// tasks are distributed among them. Task submission uses synchronous
// handoff via an unbuffered channel. If all workers are busy, Submit()
// or Execute() blocks until a worker becomes available. Tasks are never
// executed by the caller goroutine.
//
// Submit() is intended for tasks whose completion is observed through
// a Future. A task panic is captured by the Future and is reported by
// Future.Get() or Future.Result(). If the returned Future is ignored,
// the failure is effectively ignored as well.
//
// Execute() is intended for fire-and-forget tasks. Because there is no
// Future through which to observe a failure, an uncaught task panic is
// handled by the default uncaught exception handler. If none is
// configured, it is printed to stderr.
//
// The executor provides natural backpressure and prevents unbounded
// task queuing. Cancellation, timeouts, and interruption are not
// supported.
//
// Close() stops the executor by closing the job channel. Workers exit
// after processing already accepted tasks. Submitting after Close()
// causes a panic.
type Executor struct {
	jobs chan func()
}

func DefaultExecutor() *Executor {
	if defaultExecutor == nil {
		Synchronized(&defaultExecutorMu, func() {
			if defaultExecutor == nil {
				defaultExecutor = NewExecutor(runtime.NumCPU())
			}
		})
	}
	return defaultExecutor
}

func NewExecutor(workers int) *Executor {
	e := &Executor{jobs: make(chan func())}
	for i := 0; i < workers; i++ {
		e.runWorker()
	}
	return e
}

func Submit[T any](executor *Executor, task func() T) Future[T] {
	f := &futureImpl[T]{result: make(chan resultWrapper[T], 1)}
	executor.jobs <- func() {
		defer err.Recover(func(r any) {
			f.result <- resultWrapper[T]{err: r}
		})
		res := task()
		f.result <- resultWrapper[T]{val: res, err: nil}
	}
	return f
}

func (e *Executor) Submit(task func()) Future[struct{}] {
	f := &futureImpl[struct{}]{result: make(chan resultWrapper[struct{}], 1)}
	e.jobs <- func() {
		defer err.Recover(func(r any) {
			f.result <- resultWrapper[struct{}]{err: r}
		})
		task()
		f.result <- resultWrapper[struct{}]{}
	}
	return f
}

func (e *Executor) Execute(task func()) {
	e.jobs <- func() {
		defer err.Recover()
		task()
	}
}

func (e *Executor) runWorker() {
	go func() {
		for job := range e.jobs {
			job()
		}
	}()
}

func (e *Executor) Close() {
	close(e.jobs)
}
