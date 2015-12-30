package reporter

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/carbin-gun/awesome-metrics/registry"
	"github.com/carbin-gun/awesome-metrics/mechanism"
)

var shortHostName string = ""

// OpenTSDBConfig provides a container with configuration parameters for
// the OpenTSDB exporter
type OpenTSDBConfig struct {
	Addr          *net.TCPAddr      // Network address to connect to
	Registry      registry.Registry // Registry to be exported
	FlushInterval time.Duration     // Flush interval
	DurationUnit  time.Duration     // Time conversion unit for durations
	Prefix        string            // Prefix to be prepended to metric names
}

// OpenTSDB is a blocking exporter function which reports metrics in r
// to a TSDB server located at addr, flushing them every d duration
// and prepending metric names with prefix.
func OpenTSDB(r registry.Registry, d time.Duration, prefix string, addr *net.TCPAddr) {
	OpenTSDBWithConfig(OpenTSDBConfig{
		Addr:          addr,
		Registry:      r,
		FlushInterval: d,
		DurationUnit:  time.Nanosecond,
		Prefix:        prefix,
	})
}

// OpenTSDBWithConfig is a blocking exporter function just like OpenTSDB,
// but it takes a OpenTSDBConfig instead.
func OpenTSDBWithConfig(c OpenTSDBConfig) {
	for _ = range time.Tick(c.FlushInterval) {
		if err := openTSDB(&c); nil != err {
			log.Println(err)
		}
	}
}

func getShortHostname() string {
	if shortHostName == "" {
		host, _ := os.Hostname()
		if index := strings.Index(host, "."); index > 0 {
			shortHostName = host[:index]
		} else {
			shortHostName = host
		}
	}
	return shortHostName
}

func openTSDB(c *OpenTSDBConfig) error {
	shortHostname := getShortHostname()
	now := time.Now().Unix()
	du := float64(c.DurationUnit)
	conn, err := net.DialTCP("tcp", nil, c.Addr)
	if nil != err {
		return err
	}
	defer conn.Close()
	w := bufio.NewWriter(conn)
	c.Registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case mechanism.Counter:
			fmt.Fprintf(w, "put %s.%s.count %d %d host=%s\n", c.Prefix, name, now, metric.Count(), shortHostname)
		case mechanism.Gauge:
			fmt.Fprintf(w, "put %s.%s.value %d %d host=%s\n", c.Prefix, name, now, metric.Value(), shortHostname)
		case mechanism.Gauge64:
			fmt.Fprintf(w, "put %s.%s.value %d %f host=%s\n", c.Prefix, name, now, metric.Value(), shortHostname)
		case mechanism.Histogram:
			h := metric.Snapshot()
			fmt.Fprintf(w, "put %s.%s.count %d %d host=%s\n", c.Prefix, name, now, metric.Count(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.min %d %d host=%s\n", c.Prefix, name, now, h.Min(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.max %d %d host=%s\n", c.Prefix, name, now, h.Max(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.mean %d %.2f host=%s\n", c.Prefix, name, now, h.Mean(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.std-dev %d %.2f host=%s\n", c.Prefix, name, now, h.StdDev(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.50-percentile %d %.2f host=%s\n", c.Prefix, name, now, h.Median(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.75-percentile %d %.2f host=%s\n", c.Prefix, name, now, h.Get75thPercentile(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.95-percentile %d %.2f host=%s\n", c.Prefix, name, now, h.Get95thPercentile(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.99-percentile %d %.2f host=%s\n", c.Prefix, name, now, h.Get99thPercentile(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.999-percentile %d %.2f host=%s\n", c.Prefix, name, now, h.Get999thPercentile()), shortHostname)
		case mechanism.Meter:
			fmt.Fprintf(w, "put %s.%s.count %d %d host=%s\n", c.Prefix, name, now, metric.Count(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.one-minute %d %.2f host=%s\n", c.Prefix, name, now, metric.Rate1(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.five-minute %d %.2f host=%s\n", c.Prefix, name, now, metric.Rate5(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.fifteen-minute %d %.2f host=%s\n", c.Prefix, name, now, metric.Rate15(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.mean %d %.2f host=%s\n", c.Prefix, name, now, metric.RateMean(), shortHostname)
		case mechanism.Timer:
			t := metric.Snapshot()
			fmt.Fprintf(w, "put %s.%s.count %d %d host=%s\n", c.Prefix, name, now, metric.Count(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.min %d %d host=%s\n", c.Prefix, name, now, t.Min()/int64(du), shortHostname)
			fmt.Fprintf(w, "put %s.%s.max %d %d host=%s\n", c.Prefix, name, now, t.Max()/int64(du), shortHostname)
			fmt.Fprintf(w, "put %s.%s.mean %d %.2f host=%s\n", c.Prefix, name, now, t.Mean()/du, shortHostname)
			fmt.Fprintf(w, "put %s.%s.std-dev %d %.2f host=%s\n", c.Prefix, name, now, t.StdDev()/du, shortHostname)
			fmt.Fprintf(w, "put %s.%s.50-percentile %d %.2f host=%s\n", c.Prefix, name, now, t.Median()/du, shortHostname)
			fmt.Fprintf(w, "put %s.%s.75-percentile %d %.2f host=%s\n", c.Prefix, name, now, t.Get75thPercentile()/du, shortHostname)
			fmt.Fprintf(w, "put %s.%s.95-percentile %d %.2f host=%s\n", c.Prefix, name, now, t.Get95thPercentile()/du, shortHostname)
			fmt.Fprintf(w, "put %s.%s.99-percentile %d %.2f host=%s\n", c.Prefix, name, now, t.Get99thPercentile()/du, shortHostname)
			fmt.Fprintf(w, "put %s.%s.999-percentile %d %.2f host=%s\n", c.Prefix, name, now, t.Get999thPercentile()/du, shortHostname)
			fmt.Fprintf(w, "put %s.%s.one-minute %d %.2f host=%s\n", c.Prefix, name, now, metric.Rate1(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.five-minute %d %.2f host=%s\n", c.Prefix, name, now, metric.Rate5(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.fifteen-minute %d %.2f host=%s\n", c.Prefix, name, now, metric.Rate15(), shortHostname)
			fmt.Fprintf(w, "put %s.%s.mean-rate %d %.2f host=%s\n", c.Prefix, name, now, metric.RateMean(), shortHostname)
		}
		w.Flush()
	})
	return nil
}
