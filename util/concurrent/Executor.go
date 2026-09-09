package concurrent

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-errr/go/err"
	"github.com/go-jang/go/lang"
)

type Future[T any] interface {
	Get() T
	Result() (T, error)
	GetWithTimeout(timeout time.Duration) T
	ResultWithTimeout(timeout time.Duration) (T, error)
	Cancel() bool
	IsCancelled() bool
	IsDone() bool
}

type futureImpl[T any] struct {
	once      sync.Once
	result    chan resultWrapper[T]
	ctx       context.Context
	cancel    context.CancelFunc
	done      atomic.Bool
	cancelled atomic.Bool
	val       T
	err       error
}

type resultWrapper[T any] struct {
	val T
	err any
}

func (f *futureImpl[T]) wait() {
	f.once.Do(func() {
		<-f.ctx.Done()
		if f.cancelled.Load() {
			f.err = NewCancellationExceptionFrom("Task canceled", f.ctx.Err())
			return
		}
		if f.done.Load() {
			r := <-f.result
			f.val = r.val
			if r.err != nil {
				f.err = NewExecutionExceptionFrom(fmt.Sprint(r.err), r.err)
			}
			return
		}
		f.done.Store(true)
		f.err = err.NewInterruptedExceptionFrom("Task interrupted", f.ctx.Err())
	})
}

func (f *futureImpl[T]) Get() T {
	f.wait()
	if f.err != nil {
		panic(f.err)
	}
	return f.val
}

func (f *futureImpl[T]) Result() (T, error) {
	f.wait()
	return f.val, f.err
}

func (f *futureImpl[T]) GetWithTimeout(timeout time.Duration) T {
	val, e := f.ResultWithTimeout(timeout)
	if e != nil {
		panic(e)
	}
	return val
}

func (f *futureImpl[T]) ResultWithTimeout(timeout time.Duration) (T, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-f.ctx.Done():
		return f.Result()
	case <-timer.C:
		var zero T
		return zero, NewTimeoutExceptionFrom("Timeout", context.DeadlineExceeded)
	}
}

func (f *futureImpl[T]) Cancel() bool {
	if !f.done.CompareAndSwap(false, true) {
		return false
	}
	f.cancelled.Store(true)
	f.cancel()
	return true
}

func (f *futureImpl[T]) IsCancelled() bool {
	return f.cancelled.Load()
}

func (f *futureImpl[T]) IsDone() bool {
	return f.done.Load()
}

var defaultExecutor *Executor
var defaultExecutorOnce sync.Once

// Executor is a fixed-size worker pool for asynchronous task execution.
//
// A configured number of worker goroutines is started eagerly and tasks are
// distributed among them. queueSize controls the number of tasks that may wait
// for execution. A queue size of zero provides synchronous handoff. If the
// queue is full, submission blocks until a worker or queue slot becomes
// available.
//
// Submit() is intended for tasks whose completion is observed through a Future.
// A task panic is captured by the Future and reported as an ExecutionException
// by Future.Get() or Future.Result(). If the returned Future is ignored, the
// failure is effectively ignored as well.
//
// Execute() is intended for fire-and-forget tasks. Because there is no Future
// through which to observe a failure, an uncaught task panic is handled by the
// default uncaught exception handler. If none is configured, it is printed to
// stderr.
//
// Tasks receive a context used for cooperative interruption. Future.Cancel()
// cancels the Future and interrupts its task through the context. The task must
// observe the context; running tasks cannot be forcibly stopped.
//
// Shutdown() stops accepting new tasks, completes already accepted tasks, and
// terminates the workers afterwards. ShutdownNow() additionally interrupts
// running tasks and prevents queued tasks from being executed. Futures of tasks
// interrupted before completion report InterruptedException. Submission after
// shutdown, including submissions already blocked waiting for queue capacity,
// panics with RejectedExecutionException.
type Executor struct {
	tasks        chan func()
	ctx          context.Context
	cancel       context.CancelFunc
	shutdown     chan struct{}
	shutdownOnce sync.Once
}

func DefaultExecutor() *Executor {
	defaultExecutorOnce.Do(func() {
		defaultExecutor = NewExecutor(runtime.NumCPU(), 0)
	})
	return defaultExecutor
}

func NewExecutor(poolSize, queueSize int) *Executor {
	lang.Assert(poolSize > 0, "poolSize must be greater than zero")
	lang.Assert(queueSize >= 0, "queueSize must not be negative")

	ctx, cancel := context.WithCancel(context.Background())
	executor := &Executor{
		tasks:    make(chan func(), queueSize),
		ctx:      ctx,
		cancel:   cancel,
		shutdown: make(chan struct{}),
	}
	for i := 0; i < poolSize; i++ {
		executor.runWorker()
	}
	return executor
}

func Submit[T any](executor *Executor, task func(context.Context) T) Future[T] {
	ctx, cancel := context.WithCancel(executor.ctx)
	future := &futureImpl[T]{
		result: make(chan resultWrapper[T], 1),
		ctx:    ctx,
		cancel: cancel,
	}
	wrappedTask := func() {
		defer future.cancel()
		defer future.done.CompareAndSwap(false, true)
		defer err.Recover(func(e any) {
			future.result <- resultWrapper[T]{err: e}
		})
		future.result <- resultWrapper[T]{val: task(ctx)}
	}
	select {
	case <-executor.shutdown:
		cancel()
		panic(NewRejectedExecutionException("Task rejected: executor is shut down"))
	default:
	}

	select {
	case executor.tasks <- wrappedTask:
		return future
	case <-executor.shutdown:
		cancel()
		panic(NewRejectedExecutionException("Task rejected: executor is shut down"))
	}
}

func (this *Executor) Submit(task func(context.Context)) Future[struct{}] {
	return Submit(this, func(ctx context.Context) struct{} {
		task(ctx)
		return struct{}{}
	})
}

func (this *Executor) Execute(task func(context.Context)) {
	select {
	case <-this.shutdown:
		panic(NewRejectedExecutionException("Task rejected: executor is shut down"))
	default:
	}

	wrappedTask := func() {
		defer err.Recover()
		task(this.ctx)
	}
	select {
	case this.tasks <- wrappedTask:
	case <-this.shutdown:
		panic(NewRejectedExecutionException("Task rejected: executor is shut down"))
	}
}

func (this *Executor) runWorker() {
	go func() {
		for {
			select {
			case task := <-this.tasks:
				if this.ctx.Err() != nil {
					return
				}
				task()
			case <-this.shutdown:
				if this.ctx.Err() != nil {
					return
				}
				for {
					select {
					case task := <-this.tasks:
						if this.ctx.Err() != nil {
							return
						}
						task()
					default:
						return
					}
				}
			}
		}
	}()
}

func (this *Executor) Shutdown() {
	this.shutdownOnce.Do(func() {
		close(this.shutdown)
	})
}

func (this *Executor) ShutdownNow() {
	this.cancel()
	this.Shutdown()
}
