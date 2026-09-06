package concurrent

import (
	"fmt"

	"github.com/go-errr/go/err"
)

type TimeoutException struct {
	err.AbstractException
}

func NewTimeoutException(message string) *TimeoutException {
	return &TimeoutException{
		AbstractException: *err.NewAbstractException(message, nil, err.StackTrace(1)),
	}
}

func NewTimeoutExceptionFrom(message string, cause any) *TimeoutException {
	return &TimeoutException{
		AbstractException: *err.NewAbstractException(message, cause, err.StackTrace(1)),
	}
}

func NewTimeoutExceptionWith(message string, cause any, stackTrace []uintptr) *TimeoutException {
	return &TimeoutException{
		AbstractException: *err.NewAbstractException(message, cause, stackTrace),
	}
}

func (this *TimeoutException) Format(s fmt.State, verb rune) {
	this.DefaultFormat(s, verb, this)
}
