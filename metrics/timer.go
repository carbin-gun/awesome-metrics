package metrics

import (
	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/carbin-gun/awesome-metrics/output"
)

type StandardTimer struct {
	histogram mechanism.Histogram
	meter     mechanism.Meter
}

//NewTimer return the default timer
func NewTimer() mechanism.Timer {
	return &StandardTimer{
		histogram: NewHistogram(NewExpDecayReservoir(DEFAULT_RESERVOIR_SIZE, DEFAULT_ALPHA)),
		meter:     NewMeter(),
	}
}

//CustomNewTimer with user specified histogram & meter
func CustomNewTimer(histogram mechanism.Histogram, meter mechanism.Meter) mechanism.Timer {
	return &StandardTimer{
		histogram: histogram,
		meter:     meter,
	}
}

func (timer *StandardTimer) Count() int64 {
	return 0

}
func (timer *StandardTimer) Rate1() float64 {
	return 0.0
}
func (timer *StandardTimer) Rate5() float64 {
	return 0.0

}
func (timer *StandardTimer) Rate15() float64 {
	return 0.0

}
func (timer *StandardTimer) RateMean() float64 {
	return 0.0

}

func (timer *StandardTimer) Snapshot() output.Snapshot {
	return nil
}
