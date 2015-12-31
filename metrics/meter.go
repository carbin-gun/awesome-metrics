package metrics

import (
	"sync/atomic"
	"time"

	"github.com/carbin-gun/awesome-metrics/mechanism"
)

const TickInterval int64 = 5e9 //5s

type StandardMeter struct {
	a1, a5, a15 mechanism.EWMA
	startTime   time.Time //start time ,not updated when set
	count       int64
	lastTick    int64 //last time tick,update every time when tick happens
}

func NewMeter() mechanism.Meter {
	return &StandardMeter{
		a1:        NewEWMA1(),
		a5:        NewEWMA5(),
		a15:       NewEWMA15(),
		startTime: time.Now(),
		lastTick:  int64(time.Now().Nanosecond()),
	}
}

//Metered interface
func (meter *StandardMeter) Count() int64 {
	return 0
}

func (meter *StandardMeter) Rate1() float64 {
	return 0.0
}
func (meter *StandardMeter) Rate5() float64 {
	return 0.0
}
func (meter *StandardMeter) Rate15() float64 {
	return 0.0
}
func (meter *StandardMeter) RateMean() float64 {
	return 0.0
}

//communication
func (meter *StandardMeter) Mark() {
	meter.mark(1)
}
func (meter *StandardMeter) mark(val int64) {
	meter.tickIfNecessary()
}

func (m *StandardMeter) tickIfNecessary() {
	old := m.lastTick
	current := int64(time.Now().Nanosecond())
	age := current - old
	if age > TickInterval {
		newStick := current - age%TickInterval
		if atomic.CompareAndSwapInt64(&m.lastTick, old, newStick) {
			requiredTicks := age / TickInterval
			var i int64
			for ; i < requiredTicks; i++ {
				m.tick()
			}
		}
	}
}
func (m *StandardMeter) tick() {
	m.a1.Tick()
	m.a5.Tick()
	m.a15.Tick()
}
