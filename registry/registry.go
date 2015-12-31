package registry

import (
	"encoding/json"
	"reflect"

	"errors"

	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/fanliao/go-concurrentMap"
)

type Registry interface {
	Each(func(string, interface{}))
	Get(string) interface{}
	GetOrRegister(string, interface{}) interface{}
	Register(string, interface{}) error
	Unregister(string)
	UnregisterAll()
	Prefix() string
	MarshalJson() ([]byte, error)
}

type StandardRegistry struct {
	universalPrefix string
	metrics         *concurrent.ConcurrentMap
}

//Registry creation with specifying the universal prefix of all the metrics-keys
func NewPrefixRegistry(prefix string) Registry {
	return &StandardRegistry{universalPrefix: prefix, metrics: concurrent.NewConcurrentMap()}
}

// Create a new registry.
func NewRegistry() Registry {
	return &StandardRegistry{metrics: concurrent.NewConcurrentMap()}
}

// Call the given function for each registered metric.
func (r *StandardRegistry) Each(f func(string, interface{})) {
	for name, i := range r.registered() {
		f(name, i)
	}
}

// Get the metric by the given name or nil if none is registered.
func (r *StandardRegistry) Get(name string) interface{} {
	val, _ := r.metrics.Get(name)
	return val
}

// put if absent.
func (r *StandardRegistry) GetOrRegister(name string, i interface{}) interface{} {
	val := r.Get(name)
	if val != nil {
		return val
	}
	if v := reflect.ValueOf(i); v.Kind() == reflect.Func {
		i = v.Call(nil)[0].Interface()
	}
	r.register(name, i)
	return i
}

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func (r *StandardRegistry) Register(name string, i interface{}) error {
	return r.register(name, i)
}

// Unregister the metric with the given name.
func (r *StandardRegistry) Unregister(name string) {
	r.metrics.Remove(name)
}

// Unregister all metrics.  (Mostly for testing.)
func (r *StandardRegistry) UnregisterAll() {
	all := r.metrics.ToSlice()
	for name, _ := range all {
		r.metrics.Remove(name)
	}
}

//get the universal prefix of all the metrics
func (r *StandardRegistry) Prefix() string {
	return r.universalPrefix
}

func (r *StandardRegistry) register(name string, i interface{}) error {
	if val, err := r.metrics.Get(name); err != nil || val != nil {
		return errors.New("register error for name:" + name)
	}
	switch i.(type) {
	case mechanism.Counter, mechanism.Gauge, mechanism.Gauge64, mechanism.Histogram, mechanism.Meter, mechanism.Timer:
		r.metrics.PutIfAbsent(name, i)
	}
	return nil
}

func (r *StandardRegistry) registered() map[string]interface{} {
	metrics := make(map[string]interface{}, r.metrics.Size())
	for _, entry := range r.metrics.ToSlice() {
		keyString := entry.Key().(string)
		metrics[keyString] = entry.Value()
	}
	return metrics
}
func (r *StandardRegistry) MarshalJson() ([]byte, error) {
	data := make(map[string]map[string]interface{})
	r.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case mechanism.Counter:
			values["count"] = metric.Count()
		case mechanism.Gauge:
			values["value"] = metric.Value()
		case mechanism.Gauge64:
			values["value"] = metric.Value()
		case mechanism.Histogram:
			h := metric.Snapshot()
			values["count"] = metric.Count()
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["stddev"] = h.StdDev()
			values["median"] = h.Median()
			values["75%"] = h.Get75thPercentile()
			values["95%"] = h.Get95thPercentile()
			values["99%"] = h.Get99thPercentile()
			values["99.9%"] = h.Get999thPercentile()
		case mechanism.Meter:
			values["count"] = metric.Count()
			values["1m.rate"] = metric.Rate1()
			values["5m.rate"] = metric.Rate5()
			values["15m.rate"] = metric.Rate15()
			values["mean.rate"] = metric.RateMean()
		case mechanism.Timer:
			t := metric.Snapshot()
			values["count"] = metric.Count()
			values["min"] = t.Min()
			values["max"] = t.Max()
			values["mean"] = t.Mean()
			values["stddev"] = t.StdDev()
			values["median"] = t.Median()
			values["75%"] = t.Get75thPercentile()
			values["95%"] = t.Get95thPercentile()
			values["99%"] = t.Get99thPercentile()
			values["99.9%"] = t.Get999thPercentile()
			values["1m.rate"] = metric.Rate1()
			values["5m.rate"] = metric.Rate5()
			values["15m.rate"] = metric.Rate15()
			values["mean.rate"] = metric.RateMean()
		}
		data[name] = values
	})
	return json.Marshal(data)
}
