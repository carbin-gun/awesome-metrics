package registry

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/fanliao/go-concurrentMap"
)

// DuplicateMetric is the error returned by Registry.Register when a metric
// already exists.  If you mean to Register that metric you must first
// Unregister the existing metric.
type MetricError string

func (err MetricError) Error() string {
	return fmt.Sprintf(" metric error: %s", string(err))
}

// A Registry holds references to a set of metrics 1by name and can iterate
// over them, calling callback functions provided by the user.
//
// This is an interface so as to encourage other structs to implement
// the Registry API as appropriate.
type Registry interface {

	// Call the given function for each registered metric.
	Each(func(string, interface{}))

	// Get the metric by the given name or nil if none is registered.
	Get(string) interface{}

	// Gets an existing metric or registers the given one.
	// The interface can be the metric to register if not found in registry,
	// or a function returning the metric for lazy instantiation.
	GetOrRegister(string, interface{}) interface{}

	// Register the given metric under the given name.
	Register(string, interface{}) error

	// Run all registered healthchecks.
	RunHealthChecks()

	// Unregister the metric with the given name.
	Unregister(string)

	// Unregister all metrics.  (Mostly for testing.)
	UnregisterAll()
	//get the universal prefix of all the metrics
	Prefix() string
	//MarshalJson output json
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

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
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
		return MetricError("register error for name:" + name)
	}
	switch i.(type) {
	case mechanism.Counter, mechanism.Gauge, mechanism.Gauge64, mechanism.Histogram, mechanism.Meter, mechanism.Timer:
		r.metrics.Put(name, i)
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
func (r *StandardRegistry) MarshalJSON() ([]byte, error) {
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
			values["median"] = h.Median()
			values["75%"] = h.Get75thPercentile()
			values["95%"] = h.Get95thPercentile()
			values["99%"] = h.Get99thPercentile()
			values["99.9%"] = h.Get999thPercentile()
			values["1m.rate"] = metric.Rate1()
			values["5m.rate"] = metric.Rate5()
			values["15m.rate"] = metric.Rate15()
			values["mean.rate"] = metric.RateMean()
		}
		data[name] = values
	})
	return json.Marshal(data)
}

var DefaultRegistry = NewRegistry()
