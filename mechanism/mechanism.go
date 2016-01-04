package mechanism

import (
	"time"

	"github.com/carbin-gun/awesome-metrics/output"
)

type Counter interface {
	//Counting interface
	Count() int64
	//communication
	Dec(i int64)
	Inc(i int64)
}

type Meter interface {
	//Metered interface
	Count() int64
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	//communication
	Mark()
}
type Histogram interface {
	//Counting interface
	Count() int64
	//communication
	Update(int64)
	//snapshot data about histogram
	Snapshot() output.Snapshot
}
type Timer interface {
	//counting
	Count() int64
	//from meter interface
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	//snapshot data about histogram
	Snapshot() output.Snapshot

	//communication
	Time(func())
	Update(duration time.Duration)
}
type Gauge interface {
	//Gauged interface
	Value() int64
	//communication
	Update(int64)
}

type Gauge64 interface {
	//Gauged64 interface
	Value() float64
	//communication
	Update(float64)
}
type EWMA interface {
	Rate() float64
	Update(int64)
	Tick()
}
