package multitimer

import (
	"container/heap"
	"sync"
	"time"
)

type item[T any] struct {
	when    time.Time
	payload T
}

func less[T any](a, b item[T]) bool {
	return a.when.Before(b.when)
}

type Timer[T any] struct {
	C chan T

	m      sync.RWMutex
	timer  *time.Timer
	timers *arrayHeap[item[T], func(a, b item[T]) bool]
}

func NewWithCapacity[T any](cap int) *Timer[T] {
	return &Timer[T]{
		C:      make(chan T, cap),
		timers: newArrayHeap[item[T], func(a, b item[T]) bool](less[T]),
	}
}

func New[T any]() *Timer[T] {
	return NewWithCapacity[T](1)
}

func (mt *Timer[T]) Add(delay time.Duration, payload T) {
	if delay < 0 {
		panic("negative delay")
	}
	now := time.Now()
	when := now.Add(delay)
	mt.m.Lock()
	defer mt.m.Unlock()
	heap.Push(mt.timers, item[T]{
		when:    when,
		payload: payload,
	})
	mt.schedule(delay)
}

func (mt *Timer[T]) Stop() {
	mt.m.Lock()
	defer mt.m.Unlock()
	if mt.timer != nil && !mt.timer.Stop() {
		<-mt.timer.C
	}
	mt.timers.Clear()
}

func (mt *Timer[T]) schedule(delay time.Duration) {
	if mt.timer == nil {
		mt.timer = time.AfterFunc(delay, mt.fire)
	} else {
		mt.timer.Reset(delay)
	}
}

func (mt *Timer[T]) fire() {
	now := time.Now()
	mt.m.Lock()
	defer mt.m.Unlock()
	toFire := mt.itemsToFire(now)
	if mt.timers.Len() == 0 {
		if !mt.timer.Stop() {
			select {
			case <-mt.timer.C:
			default:
			}
		}
	} else {
		mt.schedule(mt.timers.Top().when.Sub(now))
	}
	for _, item := range toFire {
		select {
		case mt.C <- item.payload:
		default:
		}
	}
}

func (mt *Timer[T]) itemsToFire(now time.Time) []item[T] {
	var result []item[T]
	for mt.timers.Len() > 0 {
		top := mt.timers.Top()
		if top.when.After(now) {
			break
		}
		heap.Pop(mt.timers)
		result = append(result, top)
	}
	return result
}
