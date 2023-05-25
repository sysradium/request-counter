// Package counter provides ...
package counter

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"
)

type ListCounter struct {
	m        sync.RWMutex
	requests *list.List
	window   time.Duration
	ctx      context.Context
}

func (l *ListCounter) Get() int {
	l.m.Lock()
	defer l.m.Unlock()

	count := l.requests.Len()

	expirationTime := time.Now().Add(-l.window)
	for c := l.requests.Front(); c != nil; c = c.Next() {
		if c.Value.(time.Time).After(expirationTime) {
			break
		}
		count--
	}

	return count
}

func (l *ListCounter) Increment() {
	l.m.Lock()
	defer l.m.Unlock()

	l.requests.PushBack(time.Now())
}

func (l *ListCounter) Vacuum() {
	timer := time.NewTicker(10 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			l.cleanup()
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *ListCounter) Flush() {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			l.flush()
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *ListCounter) flush() {
	l.m.Lock()
	defer l.m.Unlock()
}

func (l *ListCounter) cleanup() {
	l.m.Lock()
	defer l.m.Unlock()

	expirationTime := time.Now().Add(-l.window)

	count := 0
	for l.requests.Len() > 0 {
		head := l.requests.Front()
		if !head.Value.(time.Time).Before(expirationTime) {
			break
		}
		count++
		l.requests.Remove(head)
	}
	fmt.Printf("cleaned %d\n", count)
}

func NewCounter(
	ctx context.Context,
	window time.Duration,
) *ListCounter {
	return &ListCounter{
		requests: list.New(),
		window:   window,
		ctx:      ctx,
	}
}
