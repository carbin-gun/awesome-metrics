package metrics

import (
	"sync/atomic"

	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/carbin-gun/awesome-metrics/output"
)

type StandardHistogram struct {
	count     int64
	reservoir Reservoir
}

func NewHistogram(reservoir Reservoir) mechanism.Histogram {
	return &StandardHistogram{
		count:     0,
		reservoir: reservoir,
	}
}

//Counting interface
func (histogram *StandardHistogram) Count() int64 {
	return atomic.LoadInt64(&histogram.count)
}

//communication
func (histogram *StandardHistogram) Update(val int64) {
	atomic.AddInt64(&histogram.count, val)
	histogram.reservoir.Update(val)
}

//snapshot data about histogram
func (histogram *StandardHistogram) Snapshot() output.Snapshot {
	return histogram.reservoir.Snapshot()
}
