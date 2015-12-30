package reporter

import (
	"fmt"
	"log/syslog"
	"time"

	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/carbin-gun/awesome-metrics/registry"
)

// Output each metric in the given registry to syslog periodically using
// the given syslogger.
func Syslog(r registry.Registry, d time.Duration, w *syslog.Writer) {
	for _ = range time.Tick(d) {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case mechanism.Counter:
				w.Info(fmt.Sprintf("counter %s: count: %d", name, metric.Count()))
			case mechanism.Gauge:
				w.Info(fmt.Sprintf("gauge %s: value: %d", name, metric.Value()))
			case mechanism.Gauge64:
				w.Info(fmt.Sprintf("gauge %s: value: %f", name, metric.Value()))
			case mechanism.Histogram:
				h := metric.Snapshot()
				w.Info(fmt.Sprintf(
					"histogram %s: count: %d min: %d max: %d mean: %.2f stddev: %.2f median: %.2f 75%%: %.2f 95%%: %.2f 99%%: %.2f 99.9%%: %.2f",
					name,
					metric.Count(),
					h.Min(),
					h.Max(),
					h.Mean(),
					h.StdDev(),
					h.Median(),
					h.Get75thPercentile(),
					h.Get95thPercentile(),
					h.Get99thPercentile(),
					h.Get999thPercentile(),
				))
			case mechanism.Meter:
				w.Info(fmt.Sprintf(
					"meter %s: count: %d 1-min: %.2f 5-min: %.2f 15-min: %.2f mean: %.2f",
					name,
					metric.Count(),
					metric.Rate1(),
					metric.Rate5(),
					metric.Rate15(),
					metric.RateMean(),
				))
			case mechanism.Timer:
				t := metric.Snapshot()
				w.Info(fmt.Sprintf(
					"timer %s: count: %d min: %d max: %d mean: %.2f stddev: %.2f median: %.2f 75%%: %.2f 95%%: %.2f 99%%: %.2f 99.9%%: %.2f 1-min: %.2f 5-min: %.2f 15-min: %.2f mean-rate: %.2f",
					name,
					metric.Count(),
					t.Min(),
					t.Max(),
					t.Mean(),
					t.StdDev(),
					t.Median(),
					t.Get75thPercentile(),
					t.Get95thPercentile(),
					t.Get99thPercentile(),
					t.Get999thPercentile(),
					metric.Rate1(),
					metric.Rate5(),
					metric.Rate15(),
					metric.RateMean(),
				))
			}
		})
	}
}
