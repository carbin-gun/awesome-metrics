package reporter

import (
	"log"
	"time"

	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/carbin-gun/awesome-metrics/registry"
)

// Output each metric in the given registry periodically using the given
// logger.
func Log(r registry.Registry, d time.Duration, l *log.Logger) {
	for _ = range time.Tick(d) {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case mechanism.Counter:
				l.Printf("counter %s\n", name)
				l.Printf("  count:       %9d\n", metric.Count())
			case mechanism.Gauge:
				l.Printf("gauge %s\n", name)
				l.Printf("  value:       %9d\n", metric.Value())
			case mechanism.Gauge64:
				l.Printf("gauge %s\n", name)
				l.Printf("  value:       %f\n", metric.Value())
			case mechanism.Histogram:
				h := metric.Snapshot()
				l.Printf("histogram %s\n", name)
				l.Printf("  count:       %9d\n", metric.Count())
				l.Printf("  min:         %9d\n", h.Min())
				l.Printf("  max:         %9d\n", h.Max())
				l.Printf("  mean:        %12.2f\n", h.Mean())
				l.Printf("  stddev:      %12.2f\n", h.StdDev())
				l.Printf("  median:      %12.2f\n", h.Median())
				l.Printf("  75%%:         %12.2f\n", h.Get75thPercentile())
				l.Printf("  95%%:         %12.2f\n", h.Get95thPercentile())
				l.Printf("  99%%:         %12.2f\n", h.Get95thPercentile())
				l.Printf("  99.9%%:       %12.2f\n", h.Get999thPercentile())
			case mechanism.Meter:
				l.Printf("meter %s\n", name)
				l.Printf("  count:       %9d\n", metric.Count())
				l.Printf("  1-min rate:  %12.2f\n", metric.Rate1())
				l.Printf("  5-min rate:  %12.2f\n", metric.Rate5())
				l.Printf("  15-min rate: %12.2f\n", metric.Rate15())
				l.Printf("  mean rate:   %12.2f\n", metric.RateMean())
			case mechanism.Timer:
				t := metric.Snapshot()
				l.Printf("timer %s\n", name)
				l.Printf("  count:       %9d\n", metric.Count())
				l.Printf("  min:         %9d\n", t.Min())
				l.Printf("  max:         %9d\n", t.Max())
				l.Printf("  mean:        %12.2f\n", t.Mean())
				l.Printf("  stddev:      %12.2f\n", t.StdDev())
				l.Printf("  median:      %12.2f\n", t.Median())
				l.Printf("  75%%:         %12.2f\n", t.Get75thPercentile())
				l.Printf("  95%%:         %12.2f\n", t.Get95thPercentile())
				l.Printf("  99%%:         %12.2f\n", t.Get95thPercentile())
				l.Printf("  99.9%%:       %12.2f\n", t.Get999thPercentile())
				l.Printf("  1-min rate:  %12.2f\n", metric.Rate1())
				l.Printf("  5-min rate:  %12.2f\n", metric.Rate5())
				l.Printf("  15-min rate: %12.2f\n", metric.Rate15())
				l.Printf("  mean rate:   %12.2f\n", metric.RateMean())
			}
		})
	}
}
