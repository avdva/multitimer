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

// Timer sends payload to its chan after a delay.
// One can schedule several events, but only one real timer will be used.
// It is safe to use a Timer object concurrently.
// Timer may drop messages, if the reader does not read fast enough,
// so it's important to choose a correct chan capacity.
type Timer[T any] struct {
	C chan T

	m      sync.Mutex
	timer  *time.Timer
	timers *arrayHeap[item[T], func(a, b item[T]) bool]
}

// NewWithCapacity returns a timer for given capacity.
func NewWithCapacity[T any](cap int) *Timer[T] {
	return &Timer[T]{
		C:      make(chan T, cap),
		timers: newArrayHeap[item[T], func(a, b item[T]) bool](less[T]),
	}
}

// New returns a timer with capacity set to 1.
func New[T any]() *Timer[T] {
	return NewWithCapacity[T](1)
}

// Schedule schedules a timer.
// The payload will be sent to C after a delay.
func (mt *Timer[T]) Schedule(delay time.Duration, payload T) {
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
	mt.schedule(now)
}

// Stop cancels all timers.
func (mt *Timer[T]) Stop() {
	mt.m.Lock()
	defer mt.m.Unlock()
	mt.stopTimer()
	mt.timers.Clear()
}

func (mt *Timer[T]) schedule(now time.Time) {
	delay := mt.timers.Top().when.Sub(now)
	if mt.timer == nil {
		mt.timer = time.AfterFunc(delay, mt.fire)
	} else {
		mt.timer.Reset(delay)
	}
}

func (mt *Timer[T]) stopTimer() {
	if mt.timer != nil && !mt.timer.Stop() {
		select {
		case <-mt.timer.C:
		default:
		}
	}
}

func (mt *Timer[T]) fire() {
	now := time.Now()
	mt.m.Lock()
	defer mt.m.Unlock()
	toFire := mt.itemsToFire(now)
	if mt.timers.Len() == 0 {
		mt.stopTimer()
	} else {
		mt.schedule(now)
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
