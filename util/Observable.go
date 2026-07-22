package util

import (
	"sync"

	"github.com/go-jang/go/util/objects"
)

type Observable[T any] struct {
	observers   *ArrayList[func(observable *Observable[T], arg T)]
	observersMu sync.RWMutex
}

func NewObservable[T any]() *Observable[T] {
	return &Observable[T]{
		observers: NewArrayList[func(observable *Observable[T], arg T)](),
	}
}

func (this *Observable[T]) Empty() bool {
	this.observersMu.RLock()
	defer this.observersMu.RUnlock()

	return this.observers.Empty()
}

func (this *Observable[T]) AddObserver(obs func(observable *Observable[T], arg T)) {
	this.observersMu.Lock()
	defer this.observersMu.Unlock()

	this.observers.Add(obs)
}

func (this *Observable[T]) RemoveObserver(obs func(observable *Observable[T], arg T)) {
	this.observersMu.Lock()
	defer this.observersMu.Unlock()

	this.observers.RemoveBy(func(idx int, value func(observable *Observable[T], arg T)) bool {
		return objects.FuncPtr(value) == objects.FuncPtr(obs)
	})
}

func (this *Observable[T]) NotifyObservers(data T) {
	if this.observers.Empty() {
		return
	}
	this.observersMu.RLock()
	observers := this.observers.ToSlice()
	this.observersMu.RUnlock()

	for _, observer := range observers {
		observer(this, data)
	}
}
