package reporter
import (
"net"
"github.com/carbin-gun/awesome-metrics"
	"time"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"log"
)


type Reporter struct{
	Addr          *net.TCPAddr     // TCP Address of server
	Registry      metrics.Registry // data collector
	FlushInterval time.Duration    //data will flush from Registry to server address
	DurationUnit  time.Duration    // Time unit of flush interval
	Percentiles   []float64        // Percentiles to report from timers and histograms
}

//Report report data to server according to the FlushInterval
func (r *Reporter) Report() {
	for _ = range time.Tick(r.FlushInterval) {
		if err:=r.ReportOnce();nil!=err {
			log.Println("report error:",err)
		}
	}
}

//compute the report key universal prefix according to the registry if the registry with a prefix setting
func computeReportPrefix(registry metrics.Registry) string{
	registryPrefix := registry.Prefix()
	var keyPrefix string
	if registryPrefix!="" {
		keyPrefix = fmt.Sprintf("%s.",registryPrefix) //there should be a point between the universal prefix and the real key
	}
	return keyPrefix
}
//Report report data to server instantly
func (r *Reporter) ReportOnce() error{
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
		case metrics.Counter:
			fmt.Fprintf(w, "%s%s.count %d %d\n", keyPrefix, name, metric.Count(), now)
		case metrics.Gauge:
			fmt.Fprintf(w, "%s%s.value %d %d\n", keyPrefix, name, metric.Value(), now)
		case metrics.GaugeFloat64:
			fmt.Fprintf(w, "%s%s.value %f %d\n", keyPrefix, name, metric.Value(), now)
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles(r.Percentiles)
			fmt.Fprintf(w, "%s%s.count %d %d\n", keyPrefix, name, h.Count(), now)
			fmt.Fprintf(w, "%s%s.min %d %d\n", keyPrefix, name, h.Min(), now)
			fmt.Fprintf(w, "%s%s.max %d %d\n", keyPrefix, name, h.Max(), now)
			fmt.Fprintf(w, "%s%s.mean %.2f %d\n", keyPrefix, name, h.Mean(), now)
			fmt.Fprintf(w, "%s%s.std-dev %.2f %d\n", keyPrefix, name, h.StdDev(), now)
			for psIdx, psKey := range r.Percentiles {
				key := strings.Replace(strconv.FormatFloat(psKey*100.0, 'f', -1, 64), ".", "", 1)
				fmt.Fprintf(w, "%s%s.%s-percentile %.2f %d\n", keyPrefix, name, key, ps[psIdx], now)
			}
		case metrics.Meter:
			m := metric.Snapshot()
			fmt.Fprintf(w, "%s%s.count %d %d\n", keyPrefix, name, m.Count(), now)
			fmt.Fprintf(w, "%s%s.one-minute %.2f %d\n", keyPrefix, name, m.Rate1(), now)
			fmt.Fprintf(w, "%s%s.five-minute %.2f %d\n", keyPrefix, name, m.Rate5(), now)
			fmt.Fprintf(w, "%s%s.fifteen-minute %.2f %d\n", keyPrefix, name, m.Rate15(), now)
			fmt.Fprintf(w, "%s%s.mean %.2f %d\n", keyPrefix, name, m.RateMean(), now)
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles(r.Percentiles)
			fmt.Fprintf(w, "%s%s.count %d %d\n", keyPrefix, name, t.Count(), now)
			fmt.Fprintf(w, "%s%s.min %d %d\n", keyPrefix, name, t.Min()/int64(du), now)
			fmt.Fprintf(w, "%s%s.max %d %d\n", keyPrefix, name, t.Max()/int64(du), now)
			fmt.Fprintf(w, "%s%s.mean %.2f %d\n", keyPrefix, name, t.Mean()/du, now)
			fmt.Fprintf(w, "%s%s.std-dev %.2f %d\n", keyPrefix, name, t.StdDev()/du, now)
			for psIdx, psKey := range r.Percentiles {
				key := strings.Replace(strconv.FormatFloat(psKey*100.0, 'f', -1, 64), ".", "", 1)
				fmt.Fprintf(w, "%s%s.%s-percentile %.2f %d\n", keyPrefix, name, key, ps[psIdx]/du, now)
			}
			fmt.Fprintf(w, "%s%s.one-minute %.2f %d\n", keyPrefix, name, t.Rate1(), now)
			fmt.Fprintf(w, "%s%s.five-minute %.2f %d\n", keyPrefix, name, t.Rate5(), now)
			fmt.Fprintf(w, "%s%s.fifteen-minute %.2f %d\n", keyPrefix, name, t.Rate15(), now)
			fmt.Fprintf(w, "%s%s.mean-rate %.2f %d\n", keyPrefix, name, t.RateMean(), now)
		}
		w.Flush()
	})
	return nil
}
}