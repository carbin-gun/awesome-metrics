package metrics

import (
	"sync/atomic"
	"time"
)

var (
	TickInterval int64 = 5 * time.Second.Nanoseconds()
)

// Meters count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
type Meter interface {
	Count() int64
	Mark(int64)
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
}

// GetOrRegisterMeter returns an existing Meter or constructs and registers a
// new StandardMeter.
func GetOrRegisterMeter(name string, r Registry) Meter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewMeter).(Meter)
}

// NewMeter constructs a new StandardMeter and launches a goroutine.
func NewMeter() Meter {
	if UseNilMetrics {
		return NilMeter{}
	}
	m := newStandardMeter()
	return m
}

// NewMeter constructs and registers a new StandardMeter and launches a
// goroutine.
func NewRegisteredMeter(name string, r Registry) Meter {
	c := NewMeter()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NilMeter is a no-op Meter.
type NilMeter struct{}

// Count is a no-op.
func (NilMeter) Count() int64 { return 0 }

// Mark is a no-op.
func (NilMeter) Mark(n int64) {}

// Rate1 is a no-op.
func (NilMeter) Rate1() float64 { return 0.0 }

// Rate5 is a no-op.
func (NilMeter) Rate5() float64 { return 0.0 }

// Rate15is a no-op.
func (NilMeter) Rate15() float64 { return 0.0 }

// RateMean is a no-op.
func (NilMeter) RateMean() float64 { return 0.0 }

// StandardMeter is the standard implementation of a Meter.
type StandardMeter struct {
	a1, a5, a15 EWMA
	startTime   time.Time //start time ,not updated when set
	count       int64
	lastTick    int64 //last time tick,update every time when tick happens
}

func newStandardMeter() *StandardMeter {
	return &StandardMeter{
		a1:        NewEWMA1(),
		a5:        NewEWMA5(),
		a15:       NewEWMA15(),
		startTime: time.Now(),
		lastTick:  int64(time.Now().Nanosecond()),
	}
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

// Mark records the occurance of n events.
func (m *StandardMeter) Mark(n int64) {
	atomic.AddInt64(&m.count, n)
	m.a1.Update(n)
	m.a5.Update(n)
	m.a15.Update(n)
}

// Count returns the number of events recorded.
func (m *StandardMeter) Count() int64 {
	count := atomic.LoadInt64(&m.count)
	return count
}

// Rate1 returns the one-minute moving average rate of events per second.
func (m *StandardMeter) Rate1() float64 {
	m.tickIfNecessary()
	return m.a1.Rate()
}

// Rate5 returns the five-minute moving average rate of events per second.
func (m *StandardMeter) Rate5() float64 {
	m.tickIfNecessary()
	return m.a1.Rate()
}

// Rate15 returns the fifteen-minute moving average rate of events per second.
func (m *StandardMeter) Rate15() float64 {
	m.tickIfNecessary()
	return m.a15.Rate()
}

// RateMean returns the meter's mean rate of events per second.
func (m *StandardMeter) RateMean() float64 {
	if m.Count() == 0 {
		return float64(0.0)
	} else {
		elapsed := time.Now().Sub(m.startTime).Nanoseconds()
		return float64(m.Count()) / float64(elapsed)
	}

}

func (m *StandardMeter) tick() {
	m.a1.Tick()
	m.a5.Tick()
	m.a15.Tick()
}
