package utils

import (
	"errors"
	"fmt"
	"sync"
)

const (
	_defaultCapacity = 10000
)

type Params[A any] func(q *queue[A])

func WithCapacity[A any](cap int) Params[A] {
	// force capacity to be >=1
	if cap < 1 {
		cap = 1
	}

	return func(q *queue[A]) {
		q.capacity = cap
		q.arr = make([]A, cap)
	}
}

type Queue[A any] interface {
	WaitAndPush(a A)
	WaitAndPop() A

	Push(A) error
	Pop() (A, error)

	Front() (A, error)
	Size() int
}

func NewQueue[A any](params ...Params[A]) Queue[A] {
	q := &queue[A]{
		arr:      make([]A, _defaultCapacity),
		head:     -1,
		tail:     0,
		capacity: _defaultCapacity,
		// default inits
		// mtx:      sync.RWMutex{},
	}
	q.cond = sync.NewCond(&q.mtx)

	for _, p := range params {
		p(q)
	}

	return q
}

type queue[A any] struct {
	arr                  []A
	head, tail, capacity int
	mtx                  sync.RWMutex
	cond                 *sync.Cond
}

// WaitAndPush to back of the queue
func (q *queue[A]) WaitAndPush(a A) {
	q.mtx.Lock()
	defer func() {
		q.mtx.Unlock()
		q.cond.Signal() // used to singal empty-q's pop
	}()
	// queue is full, wait for pushing
	fmt.Print("xxxsize is:", q.size())
	for q.size() == q.capacity {
		fmt.Print("size is:", q.size())
		q.cond.Wait()
	}

	q.forcepush(a)
}

// Push to back of the queue
func (q *queue[A]) Push(a A) error {
	q.mtx.Lock()
	defer func() {
		q.mtx.Unlock()
		q.cond.Signal() // used to singal empty-q's pop
	}()

	if q.size() == q.capacity {
		return errors.New("queue is full")
	}
	q.forcepush(a)
	return nil
}

// should be called inside lock, with pre conditions
func (q *queue[A]) forcepush(a A) {
	// initially empty
	if q.head == -1 {
		q.head = 0
	}

	q.arr[q.tail] = a
	q.tail = (q.tail + 1) % q.capacity
}

// WaitAndPop the front of the queue
func (q *queue[A]) WaitAndPop() A {
	q.mtx.Lock()
	defer func() {
		q.mtx.Unlock()
		q.cond.Signal() // used to singal full-q's push
	}()
	for q.size() == 0 {
		q.cond.Wait()
	}

	return q.forcepop()
}

// Pop the front of the queue
func (q *queue[A]) Pop() (A, error) {
	q.mtx.Lock()
	defer func() {
		q.mtx.Unlock()
		q.cond.Signal() // used to singal full-q's push
	}()

	if q.size() == 0 {
		var zero A
		return zero, errors.New("queue is empty")
	}

	return q.forcepop(), nil
}

// should be called inside lock, with pre conditions
func (q *queue[A]) forcepop() A {
	elem := q.arr[q.head]
	q.head = (q.head + 1) % q.capacity

	// queue is now empty
	if q.head == q.tail {
		q.head, q.tail = -1, 0
	}

	return elem
}

// Front for the queue
func (q *queue[A]) Front() (A, error) {
	q.mtx.RLock()
	defer q.mtx.RUnlock()

	// assume
	if q.size() == 0 {
		var zero A
		return zero, errors.New("queue is empty")
	}

	return q.arr[q.head], nil
}

// Size() for the queue
func (q *queue[A]) Size() int {
	q.mtx.RLock()
	defer q.mtx.RUnlock()

	return q.size()
}

// should be called inside lock
func (q *queue[A]) size() int {
	if q.head == -1 {
		return 0
	}
	if q.tail <= q.head {
		return q.capacity - q.head + q.tail
	}
	return q.tail - q.head

}
