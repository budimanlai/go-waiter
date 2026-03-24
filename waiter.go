package waiter

import (
	"context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("listener not found")

type Listener[T any] struct {
	store sync.Map // map[string]chan T
}

func NewWaiter[T any]() *Listener[T] {
	return &Listener[T]{}
}

// Register: buat slot untuk nunggu response
func (l *Listener[T]) Register(key string) {
	ch := make(chan T, 1) // buffered biar tidak blocking
	l.store.Store(key, ch)
}

// Listen: tunggu response (pakai context biar bisa timeout)
func (l *Listener[T]) Listen(ctx context.Context, key string) (T, error) {
	val, ok := l.store.Load(key)
	if !ok {
		var zero T
		return zero, ErrNotFound
	}

	ch := val.(chan T)

	select {
	case res := <-ch:
		return res, nil
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// Resolve: kirim response ke listener
func (l *Listener[T]) Resolve(key string, data T) {
	val, ok := l.store.Load(key)
	if !ok {
		return
	}

	ch := val.(chan T)

	// non-blocking send
	select {
	case ch <- data:
	default:
	}

	l.store.Delete(key)
}

// Remove: cleanup manual
func (l *Listener[T]) Remove(key string) {
	l.store.Delete(key)
}
