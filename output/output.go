package output

//Counting
type Counting interface {
	Count() int64
}

//Marker interface
type Metric interface {
}

//Metered Counting+rate
type Metered interface {
	Count() int64
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
}
type Gauged interface {
	Value() int64
}

type GaugedFloat64 interface {
	Value() float64
}

type Histogram interface {
	Count() int64
}

type Snapshot interface {
	Value(percentile float64) float64
	Values() []float64
	Size() int64
	Max() int64
	StdDev() float64
	Mean() float64
	Min() int64
	Median() float64 //Median is 50th
	Get75thPercentile() float64
	Get95thPercentile() float64
	Get98thPercentile() float64
	Get99thPercentile() float64
	Get999thPercentile() float64
	Percentiles() []float64
}
