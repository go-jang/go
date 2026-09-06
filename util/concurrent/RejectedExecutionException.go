package concurrent

import (
	"fmt"

	"github.com/go-errr/go/err"
)

type RejectedExecutionException struct {
	err.RuntimeException
}

func NewRejectedExecutionException(message string) *RejectedExecutionException {
	return &RejectedExecutionException{
		RuntimeException: *err.NewRuntimeExceptionWith(message, nil, err.StackTrace(1)),
	}
}

func NewRejectedExecutionExceptionFrom(message string, cause any) *RejectedExecutionException {
	return &RejectedExecutionException{
		RuntimeException: *err.NewRuntimeExceptionWith(message, cause, err.StackTrace(1)),
	}
}

func NewRejectedExecutionExceptionWith(message string, cause any, stackTrace []uintptr) *RejectedExecutionException {
	return &RejectedExecutionException{
		RuntimeException: *err.NewRuntimeExceptionWith(message, cause, stackTrace),
	}
}

func (this *RejectedExecutionException) Format(s fmt.State, verb rune) {
	this.DefaultFormat(s, verb, this)
}
