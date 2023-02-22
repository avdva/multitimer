package multitimer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMultiTimer(t *testing.T) {
	timer := NewWithCapacity[int](10)
	for i := 0; i < 10; i++ {
		timer.Add(100*time.Millisecond, i)
	}
	for i := 0; i < 10; i++ {
		val := <-timer.C
		assert.Equal(t, i, val)
	}
}
