package metrics

import "sync"

//implements Gauge interface
type StandardGauge64 struct {
	value float64
	mutex sync.RWMutex
}

// Update updates the gauge's value.
func (g *StandardGauge64) Update(v float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value = v
}

// Value returns the gauge's current value.
func (g *StandardGauge64) Value() float64 {
	g.mutex.RLock()
	g.mutex.RUnlock()
	return g.value
}
