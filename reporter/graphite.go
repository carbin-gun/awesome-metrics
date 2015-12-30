package reporter
import (
"net"
	"time"
	"bufio"
	"fmt"
	"log"
	"github.com/carbin-gun/awesome-metrics/output"
	"github.com/carbin-gun/awesome-metrics/mechanism"
	"github.com/carbin-gun/awesome-metrics/registry"
)


type GraphiteReporter struct{
	Addr          *net.TCPAddr     // TCP Address of server
	Registry      registry.Registry // data collector
	FlushInterval time.Duration    //data will flush from Registry to server address
	DurationUnit  time.Duration    // Time unit of flush interval
	Percentiles   []float64        // Percentiles to report from timers and histograms
}

//Report report data to server according to the FlushInterval
func (r *GraphiteReporter) Report() {
	for _ = range time.Tick(r.FlushInterval) {
		if err:=r.ReportOnce();nil!=err {
			log.Println("report error:",err)
		}
	}
}

//compute the report key universal prefix according to the registry if the registry with a prefix setting
func computeReportPrefix(registry registry.Registry) string{
	registryPrefix := registry.Prefix()
	var keyPrefix string
	if registryPrefix!="" {
		keyPrefix = fmt.Sprintf("%s.",registryPrefix) //there should be a point between the universal prefix and the real key
	}
	return keyPrefix
}

func outputPercentiles(w *bufio.Writer ,prefix string,name string,snapshot output.Snapshot,currentTime time.Time) {
	p75:=snapshot.Get75thPercentile()
	p95:=snapshot.Get95thPercentile()
	p98:=snapshot.Get98thPercentile()
	p99:=snapshot.Get99thPercentile()
	p999:=snapshot.Get999thPercentile()
	fmt.Fprintf(w, "%s%s.75-percentile %.2f %d\n", prefix, name, p75, currentTime)
	fmt.Fprintf(w, "%s%s.95-percentile %.2f %d\n", prefix, name, p95, currentTime)
	fmt.Fprintf(w, "%s%s.98-percentile %.2f %d\n", prefix, name, p98, currentTime)
	fmt.Fprintf(w, "%s%s.99-percentile %.2f %d\n", prefix, name, p99, currentTime)
	fmt.Fprintf(w, "%s%s.999-percentile %.2f %d\n", prefix, name, p999, currentTime)
}

//Report report data to server instantly
func (r *GraphiteReporter) ReportOnce() error{
	conn, err := net.DialTCP("tcp", nil, r.Addr)
	if nil != err {
		return err
	}
	defer conn.Close()
	w := bufio.NewWriter(conn)
	now:=time.Now()
	du := float64(r.DurationUnit)
	keyPrefix:= computeReportPrefix(r.Registry)
	r.Registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case mechanism.Counter:
			fmt.Fprintf(w, "%s%s.count %d %d\n", keyPrefix, name, metric.Count(), now)
		case mechanism.Gauge:
			fmt.Fprintf(w, "%s%s.value %d %d\n", keyPrefix, name, metric.Value(), now)
		case mechanism.Gauge64:
			fmt.Fprintf(w, "%s%s.value %f %d\n", keyPrefix, name, metric.Value(), now)
		case mechanism.Histogram:
			h := metric.Snapshot()
			fmt.Fprintf(w, "%s%s.count %d %d\n", keyPrefix, name, metric.Count(), now)
			fmt.Fprintf(w, "%s%s.min %d %d\n", keyPrefix, name, h.Min(), now)
			fmt.Fprintf(w, "%s%s.max %d %d\n", keyPrefix, name, h.Max(), now)
			fmt.Fprintf(w, "%s%s.mean %.2f %d\n", keyPrefix, name, h.Mean(), now)
			fmt.Fprintf(w, "%s%s.std-dev %.2f %d\n", keyPrefix, name, h.StdDev(), now)
			outputPercentiles(w,keyPrefix,name,h,now)
		case mechanism.Meter:
			fmt.Fprintf(w, "%s%s.count %d %d\n", keyPrefix, name, metric.Count(), now)
			fmt.Fprintf(w, "%s%s.one-minute %.2f %d\n", keyPrefix, name, metric.Rate1(), now)
			fmt.Fprintf(w, "%s%s.five-minute %.2f %d\n", keyPrefix, name, metric.Rate5(), now)
			fmt.Fprintf(w, "%s%s.fifteen-minute %.2f %d\n", keyPrefix, name, metric.Rate15(), now)
			fmt.Fprintf(w, "%s%s.mean %.2f %d\n", keyPrefix, name, metric.RateMean(), now)
		case mechanism.Timer:
			t := metric.Snapshot()
			fmt.Fprintf(w, "%s%s.count %d %d\n", keyPrefix, name, metric.Count(), now)
			fmt.Fprintf(w, "%s%s.min %d %d\n", keyPrefix, name, t.Min()/int64(du), now)
			fmt.Fprintf(w, "%s%s.max %d %d\n", keyPrefix, name, t.Max()/int64(du), now)
			fmt.Fprintf(w, "%s%s.mean %.2f %d\n", keyPrefix, name, t.Mean()/du, now)
			fmt.Fprintf(w, "%s%s.std-dev %.2f %d\n", keyPrefix, name, t.StdDev()/du, now)
			outputPercentiles(w,keyPrefix,name,h,now)
			fmt.Fprintf(w, "%s%s.one-minute %.2f %d\n", keyPrefix, name, metric.Rate1(), now)
			fmt.Fprintf(w, "%s%s.five-minute %.2f %d\n", keyPrefix, name, metric.Rate5(), now)
			fmt.Fprintf(w, "%s%s.fifteen-minute %.2f %d\n", keyPrefix, name, metric.Rate15(), now)
			fmt.Fprintf(w, "%s%s.mean-rate %.2f %d\n", keyPrefix, name, metric.RateMean(), now)
		}
		w.Flush()
	})
	return nil
}
}