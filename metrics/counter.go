package metrics

import (
	"sync/atomic"

	"github.com/carbin-gun/awesome-metrics/mechanism"
)

type StandardCounter struct {
	count int64
}

func NewCounter() mechanism.Counter {
	return &StandardCounter{count: 0}
}

//implement the Counter interface
func (c *StandardCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

func (c *StandardCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

// Inc increments the counter by the given amount.
func (c *StandardCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}
