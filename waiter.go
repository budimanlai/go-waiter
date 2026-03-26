package waiter

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
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
		l.Remove(key)
		return res, nil
	case <-ctx.Done():
		l.Remove(key)
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
}

// Remove: cleanup manual
func (l *Listener[T]) Remove(key string) {
	l.store.Delete(key)
}

// RequestWithKey: helper buat register, send, dan listen sekaligus
func (l *Listener[T]) RequestWithKey(ctx context.Context, key string, send func(key string) error) (T, error) {
	l.Register(key)

	err := send(key)

	if err != nil {
		l.Remove(key)
		var zero T
		return zero, err
	}

	return l.Listen(ctx, key)
}

// Request: helper buat register, send, dan listen sekaligus
func (l *Listener[T]) Request(ctx context.Context, send func(key string) error) (T, error) {
	key := uuid.New().String()
	return l.RequestWithKey(ctx, key, send)
}

// RequestWithKeyAndTimeout: helper buat register, send, dan listen sekaligus dengan timeout
func (l *Listener[T]) RequestWithKeyAndTimeout(key string, timeout time.Duration, send func(key string) error) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return l.RequestWithKey(ctx, key, send)
}

// RequestWithTimeout: helper buat register, send, dan listen sekaligus dengan timeout
func (l *Listener[T]) RequestWithTimeout(timeout time.Duration, send func(key string) error) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return l.RequestWithKey(ctx, uuid.New().String(), send)
}
