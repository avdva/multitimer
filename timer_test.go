package multitimer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMultiTimer(t *testing.T) {
	timer := NewWithCapacity[int](10)
	for i := 0; i < 10; i++ {
		timer.Schedule(100*time.Millisecond, i)
	}
	for i := 0; i < 10; i++ {
		val := <-timer.C
		assert.Equal(t, i, val)
	}
}

func TestMultiTimer_Stop(t *testing.T) {
	timer := NewWithCapacity[int](10)
	for i := 0; i < 10; i++ {
		timer.Schedule(100*time.Millisecond*time.Duration(i+1), i)
	}
	for i := 0; i < 5; i++ {
		val := <-timer.C
		assert.Equal(t, i, val)
	}
	timer.Stop()
	select {
	case <-time.After(time.Second):
	case val := <-timer.C:
		t.Fatalf("received unexpected value %d", val)
	}
}
