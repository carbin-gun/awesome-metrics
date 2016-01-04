package metrics

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"math/rand"

	"github.com/carbin-gun/awesome-metrics/output"
	"github.com/carbin-gun/skiplist"
)

const (
	RescaleThreshold       = time.Hour
	DEFAULT_RESERVOIR_SIZE = 1028
	DEFAULT_ALPHA          = 0.015
)

type Reservoir interface {
	Size() int64
	Update(val int64)
	Snapshot() output.Snapshot
}

type ExpDecayReservoir struct {
	alpha         float64
	reservoirSize int64
	count         int64
	mutex         sync.RWMutex
	t0, t1        time.Time
	values        *WeightedSampleStorage
}

func NewExpDecayReservoir(reservoirSize int64, alpha float64) Reservoir {
	r := &ExpDecayReservoir{
		alpha:         alpha,
		reservoirSize: reservoirSize,
		t0:            time.Now(),
		values:        &WeightedSampleStorage{store: skiplist.NewList()},
	}
	r.t1 = r.t0.Add(RescaleThreshold)
	return r
}

func (r *ExpDecayReservoir) Size() int64 {
	count := r.count
	size := r.reservoirSize
	if count < size {
		return count
	} else {
		return size
	}
}

func (r *ExpDecayReservoir) Update(val int64) {
	r.UpdateBy(val, time.Now())
}
func (r *ExpDecayReservoir) UpdateBy(val int64, t time.Time) {
	r.rescaleIfNeeded(t)
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	newCount := atomic.AddInt64(&r.count, 1)
	weight := math.Exp(t.Sub(r.t0).Seconds() * r.alpha)
	sample := WeightedSample{weight: weight, value: val}
	priority := weight / rand.Float64()
	if newCount <= r.reservoirSize {
		r.values.Insert(priority, sample.MarshalBytes())
	} else {
		k, _ := r.values.First()
		if float64(k) < priority && r.values.Insert(priority, sample.MarshalBytes()) != nil {
			for r.values.Delete(k) {
				k, _ = r.values.First()
			}
		}
	}
}
func (r *ExpDecayReservoir) rescaleIfNeeded(t time.Time) {
	if t.After(r.t1) {
		r.mutex.Lock()
		defer r.mutex.Unlock()
		t0 := r.t0
		r.t0 = t
		scalingFactor := math.Exp(-r.alpha * (r.t0.Sub(t0).Seconds()))
		r.t1 = r.t0.Add(RescaleThreshold)
		iterator := r.values.Iterator()
		if iterator.Next() {
			key := iterator.Key()
			r.values.Delete(key)
			sample := UnMarshalFromBytes(iterator.Val())
			newVal := WeightedSample{weight: sample.weight * scalingFactor, value: sample.value}
			r.values.Insert(key*scalingFactor, newVal.MarshalBytes())
		}
	}
}

func (r *ExpDecayReservoir) Snapshot() output.Snapshot {

}

func (r *ExpDecayReservoir) duplicateVals() {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	iterator := r.values.Iterator()
	var samples []WeightedSample
	if iterator.Next() {
		value := UnMarshalFromBytes(iterator.Val())
		samples = append(samples, value)
	}
	var sumWeight int64
	values := make(int64, len(samples))
	for i := 0; i < len(samples); i++ {
		sumWeight += samples[i].weight
	}
	for i := 0; i < len(samples); i++ {
		values = append(values,samples[i].value)
	}



}
