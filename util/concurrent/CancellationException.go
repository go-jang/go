package concurrent

import (
	"fmt"

	"github.com/go-errr/go/err"
)

type CancellationException struct {
	err.IllegalStateException
}

func NewCancellationException(message string) *CancellationException {
	return &CancellationException{
		IllegalStateException: *err.NewIllegalStateExceptionWith(message, nil, err.StackTrace(1)),
	}
}

func NewCancellationExceptionFrom(message string, cause any) *CancellationException {
	return &CancellationException{
		IllegalStateException: *err.NewIllegalStateExceptionWith(message, cause, err.StackTrace(1)),
	}
}

func NewCancellationExceptionWith(message string, cause any, stackTrace []uintptr) *CancellationException {
	return &CancellationException{
		IllegalStateException: *err.NewIllegalStateExceptionWith(message, cause, stackTrace),
	}
}

func (this *CancellationException) Format(s fmt.State, verb rune) {
	this.DefaultFormat(s, verb, this)
}
