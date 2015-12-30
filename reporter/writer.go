package reporter

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/carbin-gun/awesome-metrics/registry"
)

// Write sorts writes each metric in the given registry periodically to the
// given io.Writer.
func Write(r registry.Registry, d time.Duration, w io.Writer) {
	for _ = range time.Tick(d) {
		WriteOnce(r, w)
	}
}

// WriteOnce sorts and writes metrics in the given registry to the given
// io.Writer.
func WriteOnce(r registry.Registry, w io.Writer) {
	var namedMetrics namedMetricSlice
	r.Each(func(name string, i interface{}) {
		namedMetrics = append(namedMetrics, namedMetric{name, i})
	})

	sort.Sort(namedMetrics)
	for _, namedMetric := range namedMetrics {
		switch metric := namedMetric.m.(type) {
		case mechanism.Counter:
			fmt.Fprintf(w, "counter %s\n", namedMetric.name)
			fmt.Fprintf(w, "  count:       %9d\n", metric.Count())
		case mechanism.Gauge:
			fmt.Fprintf(w, "gauge %s\n", namedMetric.name)
			fmt.Fprintf(w, "  value:       %9d\n", metric.Value())
		case mechanism.Gauge64:
			fmt.Fprintf(w, "gauge %s\n", namedMetric.name)
			fmt.Fprintf(w, "  value:       %f\n", metric.Value())

		case mechanism.Histogram:
			h := metric.Snapshot()
			fmt.Fprintf(w, "histogram %s\n", namedMetric.name)
			fmt.Fprintf(w, "  count:       %9d\n", metric.Count())
			fmt.Fprintf(w, "  min:         %9d\n", h.Min())
			fmt.Fprintf(w, "  max:         %9d\n", h.Max())
			fmt.Fprintf(w, "  mean:        %12.2f\n", h.Mean())
			fmt.Fprintf(w, "  stddev:      %12.2f\n", h.StdDev())
			fmt.Fprintf(w, "  median:      %12.2f\n", h.Median())
			fmt.Fprintf(w, "  75%%:         %12.2f\n", h.Get75thPercentile())
			fmt.Fprintf(w, "  95%%:         %12.2f\n", h.Get95thPercentile())
			fmt.Fprintf(w, "  99%%:         %12.2f\n", h.Get99thPercentile())
			fmt.Fprintf(w, "  99.9%%:       %12.2f\n", h.Get999thPercentile())
		case mechanism.Meter:
			fmt.Fprintf(w, "meter %s\n", namedMetric.name)
			fmt.Fprintf(w, "  count:       %9d\n", metric.Count())
			fmt.Fprintf(w, "  1-min rate:  %12.2f\n", metric.Rate1())
			fmt.Fprintf(w, "  5-min rate:  %12.2f\n", metric.Rate5())
			fmt.Fprintf(w, "  15-min rate: %12.2f\n", metric.Rate15())
			fmt.Fprintf(w, "  mean rate:   %12.2f\n", metric.RateMean())
		case mechanism.Timer:
			t := metric.Snapshot()
			fmt.Fprintf(w, "timer %s\n", namedMetric.name)
			fmt.Fprintf(w, "  count:       %9d\n", metric.Count())
			fmt.Fprintf(w, "  min:         %9d\n", t.Min())
			fmt.Fprintf(w, "  max:         %9d\n", t.Max())
			fmt.Fprintf(w, "  mean:        %12.2f\n", t.Mean())
			fmt.Fprintf(w, "  stddev:      %12.2f\n", t.StdDev())
			fmt.Fprintf(w, "  median:      %12.2f\n", t.Median())
			fmt.Fprintf(w, "  75%%:         %12.2f\n", t.Get75thPercentile())
			fmt.Fprintf(w, "  95%%:         %12.2f\n", t.Get95thPercentile())
			fmt.Fprintf(w, "  99%%:         %12.2f\n", t.Get99thPercentile())
			fmt.Fprintf(w, "  99.9%%:       %12.2f\n", t.Get999thPercentile())
			fmt.Fprintf(w, "  1-min rate:  %12.2f\n", metric.Rate1())
			fmt.Fprintf(w, "  5-min rate:  %12.2f\n", metric.Rate5())
			fmt.Fprintf(w, "  15-min rate: %12.2f\n", metric.Rate15())
			fmt.Fprintf(w, "  mean rate:   %12.2f\n", metric.RateMean())
		}
	}
}

type namedMetric struct {
	name string
	m    interface{}
}

// namedMetricSlice is a slice of namedMetrics that implements sort.Interface.
type namedMetricSlice []namedMetric

func (nms namedMetricSlice) Len() int { return len(nms) }

func (nms namedMetricSlice) Swap(i, j int) { nms[i], nms[j] = nms[j], nms[i] }

func (nms namedMetricSlice) Less(i, j int) bool {
	return nms[i].name < nms[j].name
}
