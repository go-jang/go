package concurrent_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-errr/go/err"
	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/concurrent"
)

func TestExecutorWithPanicHandling(t *testing.T) {
	t.Skip("for manual run")
	exec := concurrent.NewExecutor(10, 0)

	fmt.Printf("Submit %v\n", time.Now())
	futures := []concurrent.Future[int]{}
	for i := 0; i < 50; i++ {
		futures = append(futures, concurrent.Submit(exec, func(_ context.Context) int {
			time.Sleep(time.Second)
			lang.Assert(i%2 != 0, "simulated panic %d", i)
			return i * i
		}))
	}
	exec.Shutdown()

	fmt.Printf("Consume %v\n", time.Now())
	for i, f := range futures {
		res, err := f.Result()
		// fmt.Printf("%d: %v, %v\n", i, res, err)
		if i%2 == 0 {
			lang.Assert(err != nil, "expected panic for index %d, got result %d", i, res)
		} else {
			lang.Assert(err == nil, "unexpected error for index %d: %v", i, err)
		}
	}
	fmt.Printf("Finish %v\n", time.Now())
}

func TestExecutorSubmit(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	future := exec.Submit(func(_ context.Context) {
		time.Sleep(10 * time.Millisecond)
	})
	future.Get()
	lang.Assert(future.IsDone(), "expected future to be done")
	exec.Shutdown()
}

func TestExecutorCancel(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	future := exec.Submit(func(ctx context.Context) {
		<-ctx.Done()
	})
	future.Cancel()
	defer err.Catch(func(e any) {
		lang.Assert(err.As1[*concurrent.CancellationException](e), "expected cancellation, got %T", e)
		lang.Assert(future.IsCancelled(), "expected future to be cancelled")
		lang.Assert(future.IsDone(), "expected future to be done")
		exec.Shutdown()
	})
	future.Get()
	lang.Assert(false, "expected cancellation")
}

func TestExecutorResultWithTimeout(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	future := concurrent.Submit(exec, func(_ context.Context) int {
		time.Sleep(100 * time.Millisecond)
		return 42
	})
	result, e := future.ResultWithTimeout(10 * time.Millisecond)
	lang.Assert(e != nil, "expected timeout, got result %d", result)
	lang.Assert(err.As1[*concurrent.TimeoutException](e), "expected TimeoutException, got %T", e)
	exec.Shutdown()
}

func TestExecutorGetWithTimeout(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	future := concurrent.Submit(exec, func(_ context.Context) int {
		time.Sleep(100 * time.Millisecond)
		return 42
	})
	defer err.Catch(func(e any) {
		lang.Assert(err.As1[*concurrent.TimeoutException](e), "expected TimeoutException, got %T", e)
		exec.Shutdown()
	})
	future.GetWithTimeout(10 * time.Millisecond)
	lang.Assert(false, "expected timeout")
}

func TestExecutorShutdown(t *testing.T) {
	exec := concurrent.NewExecutor(1, 2)
	futures := []concurrent.Future[int]{}
	for i := 0; i < 3; i++ {
		futures = append(futures, concurrent.Submit(exec, func(_ context.Context) int {
			time.Sleep(10 * time.Millisecond)
			return i
		}))
	}
	exec.Shutdown()
	for i, future := range futures {
		result := future.Get()
		lang.Assert(result == i, "expected %d, got %d", i, result)
	}
}

func TestExecutorSubmitAfterShutdown(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	exec.Shutdown()
	defer err.Catch(func(e any) {
		lang.Assert(err.As1[*concurrent.RejectedExecutionException](e), "expected RejectedExecutionException, got %T", e)
	})
	exec.Submit(func(_ context.Context) {})
	lang.Assert(false, "expected rejected execution")
}

func TestExecutorShutdownNow(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	future := exec.Submit(func(ctx context.Context) {
		<-ctx.Done()
	})
	exec.ShutdownNow()
	defer err.Catch(func(e any) {
		lang.Assert(err.Interrupted(e), "expected interruption, got %T", e)
		lang.Assert(!future.IsCancelled(), "expected future not to be cancelled")
		lang.Assert(future.IsDone(), "expected future to be done")
	})
	future.Get()
	lang.Assert(false, "expected interruption")
}

func TestExecutorShutdownNowDoesNotExecuteQueuedTasks(t *testing.T) {
	exec := concurrent.NewExecutor(1, 1)
	started := make(chan struct{})
	release := make(chan struct{})
	exec.Submit(func(_ context.Context) {
		close(started)
		<-release
	})
	<-started
	exec.Submit(func(_ context.Context) {
		lang.Assert(false, "queued task must not execute")
	})
	exec.ShutdownNow()
	close(release)
}

func TestExecutorExecute(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	done := make(chan struct{})
	exec.Execute(func(_ context.Context) {
		close(done)
	})
	<-done
	exec.Shutdown()
}

func TestExecutorExecuteAfterShutdown(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	exec.Shutdown()
	defer err.Catch(func(e any) {
		lang.Assert(err.As1[*concurrent.RejectedExecutionException](e), "expected RejectedExecutionException, got %T", e)
	})
	exec.Execute(func(_ context.Context) {})
	lang.Assert(false, "expected rejected execution")
}

func TestExecutorQueue(t *testing.T) {
	exec := concurrent.NewExecutor(1, 1)
	release := make(chan struct{})
	exec.Submit(func(_ context.Context) {
		<-release
	})
	exec.Submit(func(_ context.Context) {})
	submitted := make(chan struct{})
	go func() {
		exec.Submit(func(_ context.Context) {})
		close(submitted)
	}()
	select {
	case <-submitted:
		lang.Assert(false, "expected submission to block")
	case <-time.After(10 * time.Millisecond):
	}
	close(release)
	select {
	case <-submitted:
	case <-time.After(time.Second):
		lang.Assert(false, "expected submission after queue capacity became available")
	}
	exec.Shutdown()
}

func TestExecutorShutdownRejectsBlockedSubmit(t *testing.T) {
	exec := concurrent.NewExecutor(1, 1)
	release := make(chan struct{})
	exec.Submit(func(_ context.Context) {
		<-release
	})
	exec.Submit(func(_ context.Context) {})
	rejected := make(chan bool)
	go func() {
		defer err.Catch(func(e any) {
			rejected <- err.As1[*concurrent.RejectedExecutionException](e)
		})
		exec.Submit(func(_ context.Context) {})
	}()
	time.Sleep(10 * time.Millisecond)
	exec.Shutdown()
	lang.Assert(<-rejected, "expected RejectedExecutionException")
	close(release)
}

func TestExecutorTimeoutDoesNotCancel(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	future := concurrent.Submit(exec, func(_ context.Context) int {
		time.Sleep(50 * time.Millisecond)
		return 42
	})
	_, e := future.ResultWithTimeout(10 * time.Millisecond)
	lang.Assert(err.As1[*concurrent.TimeoutException](e), "expected TimeoutException, got %T", e)
	lang.Assert(!future.IsCancelled(), "timeout must not cancel future")
	result := future.Get()
	lang.Assert(result == 42, "expected 42, got %d", result)
	exec.Shutdown()
}

func TestExecutorNonErrorPanic(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	future := concurrent.Submit(exec, func(_ context.Context) int {
		panic("simulated panic")
	})
	_, e := future.Result()
	lang.Assert(err.As1[*concurrent.ExecutionException](e), "expected ExecutionException, got %T", e)
	exec.Shutdown()
}

func TestExecutorResultAfterCancel(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	future := exec.Submit(func(ctx context.Context) {
		<-ctx.Done()
	})
	future.Cancel()
	_, e := future.Result()
	lang.Assert(err.As1[*concurrent.CancellationException](e), "expected cancellation, got %T", e)
	lang.Assert(future.IsCancelled(), "expected future to be cancelled")
	lang.Assert(future.IsDone(), "expected future to be done")
	exec.Shutdown()
}

func TestExecutorTaskInterruption(t *testing.T) {
	exec := concurrent.NewExecutor(1, 0)
	future := concurrent.Submit(exec, func(_ context.Context) int {
		panic(err.NewInterruptedException("Task interrupted"))
	})
	_, e := future.Result()
	lang.Assert(err.As1[*concurrent.ExecutionException](e), "expected execution exception, got %T", e)
	lang.Assert(err.As1[*err.InterruptedException](e), "expected interrupted cause")
	exec.Shutdown()
}
