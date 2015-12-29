package metrics

import (
	"encoding/json"
	"io"
	"time"
)

// MarshalJSON returns a byte slice containing a JSON representation of all
// the metrics in the Registry.
func (r *StandardRegistry) MarshalJSON() ([]byte, error) {
	data := make(map[string]map[string]interface{})
	r.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case Counter:
			values["count"] = metric.Count()
		case Gauge:
			values["value"] = metric.Value()
		case GaugeFloat64:
			values["value"] = metric.Value()
		case Healthcheck:
			values["error"] = nil
			metric.Check()
			if err := metric.Error(); nil != err {
				values["error"] = metric.Error().Error()
			}
		case Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles()
			values["count"] = metric.Count()
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["stddev"] = h.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
		case Meter:
			values["count"] = metric.Count()
			values["1m.rate"] = metric.Rate1()
			values["5m.rate"] = metric.Rate5()
			values["15m.rate"] = metric.Rate15()
			values["mean.rate"] = metric.RateMean()
		case Timer:
			t := metric.Snapshot()
			ps := t.Percentiles()
			values["count"] = metric.Count()
			values["min"] = t.Min()
			values["max"] = t.Max()
			values["mean"] = t.Mean()
			values["stddev"] = t.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
			values["1m.rate"] = metric.Rate1()
			values["5m.rate"] = metric.Rate5()
			values["15m.rate"] = metric.Rate15()
			values["mean.rate"] = metric.RateMean()
		}
		data[name] = values
	})
	return json.Marshal(data)
}

// WriteJSON writes metrics from the given registry  periodically to the
// specified io.Writer as JSON.
func WriteJSON(r Registry, d time.Duration, w io.Writer) {
	for _ = range time.Tick(d) {
		WriteJSONOnce(r, w)
	}
}

// WriteJSONOnce writes metrics from the given registry to the specified
// io.Writer as JSON.
func WriteJSONOnce(r Registry, w io.Writer) {
	json.NewEncoder(w).Encode(r)
}
