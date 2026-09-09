# go-jang

High-level language and utility APIs for Go, inspired by useful Java abstractions and adapted to Go.

`go-jang` does not try to turn Go into Java. It takes selected high-level concepts that make application code clearer and implements them using native Go mechanisms such as generics, goroutines, channels, contexts, and panic/error handling.

The library contains a number of small utilities and higher-level APIs. A few examples are shown below; browse the packages for the rest.

## Language helpers

The `lang` package contains small helpers for common language-level patterns.

### Ternary operator

Go intentionally has no ternary operator. For simple value selection, `lang.If` provides a concise equivalent:

```go
name := lang.If(user != nil, user.Name, "anonymous")
```

Equivalent to:

```text
user != nil ? user.Name : "anonymous"
```

It is intended for expressions where a full `if` statement would add more structure than meaning.

### Assertions

`lang.Assert` provides concise invariant checking:

```go
lang.Assert(poolSize > 0, "poolSize must be greater than zero")
lang.Assert(queueSize >= 0, "queueSize must not be negative")
```

Assertions are useful for programming errors, violated assumptions, and invalid internal state where continuing execution would be incorrect.

## Concurrent execution

One of the more substantial APIs in `go-jang` is `util/concurrent`.

It provides a fixed-size thread-pool-style executor implemented with Go goroutines, channels, and contexts.

The goal is similar to Java's `ExecutorService`: application code submits tasks to a controlled pool instead of creating an unrestricted number of goroutines itself.

A complete example:

```go
package main

import (
	"context"
	"fmt"

	"github.com/go-jang/go/util/concurrent"
)

func main() {
	executor := concurrent.NewExecutor(4, 16)
	defer executor.Shutdown()

	future := concurrent.Submit(executor, func(ctx context.Context) int {
		return 42
	})

	fmt.Println(future.Get())
}
```

`NewExecutor` accepts:

```go
concurrent.NewExecutor(poolSize, queueSize)
```

`poolSize` limits the number of tasks executing concurrently.

`queueSize` limits how many submitted tasks may wait for execution. A queue size of zero provides synchronous handoff. When the queue is full, submission waits until capacity becomes available.

This provides natural back-pressure rather than silently creating an unbounded number of goroutines or queued tasks.

### Processing a batch of tasks

A common use case is submitting a batch of independent tasks to a bounded executor and then collecting their results:

```go
executor := concurrent.NewExecutor(runtime.NumCPU(), 0)
defer executor.Shutdown()

futures := make([]concurrent.Future[int], 0, len(values))
for _, value := range values {
	futures = append(futures, concurrent.Submit(executor, func(ctx context.Context) int {
		return process(ctx, value)
	}))
}

results := make([]int, 0, len(futures))
for _, future := range futures {
	results = append(results, future.Get())
}
```

Only `poolSize` tasks execute concurrently. Additional tasks wait in the bounded queue, providing back-pressure instead of creating one goroutine per item.

Tasks are submitted first and their results are collected afterwards. `Future.Get()` keeps the processing code concise while propagating task failures as `ExecutionException`.

### Tasks without a result

For tasks that only need completion tracking:

```go
future := executor.Submit(func(ctx context.Context) {
	process(ctx)
})

future.Get()
```

### Fire-and-forget tasks

When the caller does not need a `Future`:

```go
executor.Execute(func(ctx context.Context) {
	process(ctx)
})
```

Execute is safe for fire-and-forget tasks that may fail. An uncaught panic is handled by the default uncaught exception handler: it is printed to stderr by default, or handled by a custom handler configured through go-errr.

With Submit, failures are captured by the returned Future instead. The caller should eventually call Future.Get() to observe them; ignoring the Future also ignores any task failure it contains.

### Future

`Future` provides both exception-style and explicit Go-style result handling:

```go
value := future.Get()
```

or:

```go
value, err := future.Result()
```

Timed variants are also available:

```go
value := future.GetWithTimeout(time.Second)

value, err := future.ResultWithTimeout(time.Second)
```

The executor distinguishes different terminal conditions rather than reducing everything to `context.Canceled`:

```text
task failure              -> ExecutionException
Future.Cancel()           -> CancellationException
executor interruption     -> InterruptedException
wait timeout              -> TimeoutException
submission after shutdown -> RejectedExecutionException
```

This preserves useful high-level semantics while still using standard Go primitives underneath.

### Cancellation and interruption

Tasks receive a `context.Context`:

```go
future := executor.Submit(func(ctx context.Context) {
	for _, path := range paths {
		err.Assert(ctx.Err(), "Interrupted")
		processFile(path)
	}
})
```

Context cancellation plays the role of cooperative interruption.

Calling:

```go
cancelled := future.Cancel()
```

attempts to cancel the Future and signals interruption to its task through the context. It returns `true` if the It returns true if the Future was successfully cancelled, or false if it had already completed or been cancelled.

A successfully cancelled Future is immediately considered both done and cancelled. A running task, however, may continue until it observes the cancelled context. Running goroutines are never forcibly terminated, so task code must observe the context at appropriate interruption points when cancellation matters.

The cancelled Future is reported by `Get()` as:

```text
CancellationException: Task canceled
Caused by: context canceled
```

Future cancellation takes precedence over any error produced by the task while reacting to that cancellation.

### Shutdown

Graceful shutdown:

```go
executor.Shutdown()
```

stops accepting new tasks and allows already accepted tasks to finish.

Immediate shutdown:

```go
executor.ShutdownNow()
```

additionally interrupts executor tasks through their contexts and prevents queued work from continuing where possible.

Submissions made after shutdown are rejected with `RejectedExecutionException`.

### Why use an Executor instead of starting goroutines directly?

Starting a goroutine is intentionally cheap and simple:

```go
go process()
```

For many cases, that is exactly the right solution.

An executor becomes useful when the application needs explicit control over execution:

* limit concurrent work;
* apply back-pressure through a bounded queue;
* wait for task results;
* capture task failures;
* cancel individual tasks;
* interrupt a whole group of tasks;
* wait with a timeout;
* distinguish cancellation, interruption, timeout, rejection, and execution failure;
* shut an execution facility down in a controlled way.

The abstraction is therefore not a replacement for goroutines. It is a higher-level execution tool for places where unrestricted goroutine creation is not enough.

## More

`go-jang` contains additional language, optional-value, stream-style, collection, synchronization, and utility APIs.

The project deliberately keeps these abstractions small and focused. Explore the individual packages and their Go documentation for the complete API.

## Design principle

Java is used as a source of mature API ideas, not as an implementation blueprint. The result is intended to feel like a Go library with a richer application-level toolbox, not Java rewritten in Go.
