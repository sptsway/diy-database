package utils

import (
	"sync"
	"testing"
	"time"
)

func TestPushPop(t *testing.T) {
	q := NewQueue[int](WithCapacity[int](3))

	if err := q.Push(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := q.Push(2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, err := q.Pop()
	if err != nil || v != 1 {
		t.Fatalf("expected 1, nil; got %v, %v", v, err)
	}

	v, err = q.Pop()
	if err != nil || v != 2 {
		t.Fatalf("expected 2, nil; got %v, %v", v, err)
	}
}

func TestPopEmpty(t *testing.T) {
	q := NewQueue[int](WithCapacity[int](3))

	_, err := q.Pop()
	if err == nil {
		t.Fatal("expected error popping empty queue")
	}
}

func TestPushFull(t *testing.T) {
	q := NewQueue[int](WithCapacity[int](2))

	if err := q.Push(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := q.Push(2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := q.Push(3); err == nil {
		t.Fatal("expected error pushing to full queue")
	}
}

func TestFrontAndSize(t *testing.T) {
	q := NewQueue[int](WithCapacity[int](3))

	if q.Size() != 0 {
		t.Fatalf("expected size 0, got %d", q.Size())
	}

	_ = q.Push(10)
	_ = q.Push(20)

	if q.Size() != 2 {
		t.Fatalf("expected size 2, got %d", q.Size())
	}

	v, err := q.Front()
	if err != nil || v != 10 {
		t.Fatalf("expected front 10, nil; got %v, %v", v, err)
	}
}

func TestWrapAround(t *testing.T) {
	q := NewQueue[int](WithCapacity[int](2))

	_ = q.Push(1)
	_ = q.Push(2)
	_, _ = q.Pop() // head moves forward
	_ = q.Push(3)  // tail wraps around

	v, err := q.Pop()
	if err != nil || v != 2 {
		t.Fatalf("expected 2, nil; got %v, %v", v, err)
	}
	v, err = q.Pop()
	if err != nil || v != 3 {
		t.Fatalf("expected 3, nil; got %v, %v", v, err)
	}
}

func TestWaitAndPop_BlocksUntilPush(t *testing.T) {
	q := NewQueue[int](WithCapacity[int](3))

	done := make(chan int)
	go func() {
		done <- q.WaitAndPop()
	}()

	// give the goroutine time to block on empty queue
	time.Sleep(50 * time.Millisecond)
	q.WaitAndPush(42)

	select {
	case v := <-done:
		if v != 42 {
			t.Fatalf("expected 42, got %d", v)
		}
	case <-time.After(time.Second):
		t.Fatal("WaitAndPop did not unblock after push")
	}
}

func TestWaitAndPush_BlocksUntilPop(t *testing.T) {
	q := NewQueue[int](WithCapacity[int](1))
	_ = q.Push(1) // fill it

	done := make(chan struct{})
	go func() {
		q.WaitAndPush(2) // should block until a slot frees up
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	if _, err := q.Pop(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("WaitAndPush did not unblock after pop")
	}
}

func TestConcurrentPushPop(t *testing.T) {
	q := NewQueue[int](WithCapacity[int](100))
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			q.WaitAndPush(v)
		}(i)
	}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q.WaitAndPop()
		}()
	}

	wg.Wait()
	if q.Size() != 0 {
		t.Fatalf("expected size 0 after equal push/pop, got %d", q.Size())
	}
}
