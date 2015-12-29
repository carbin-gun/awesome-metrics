package metrics

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

type WeightedSnapshot struct {
}

func (w *WeightedSnapshot) Value(percentile float64) float64 {
	return 0.0

}
func (w *WeightedSnapshot) Values() []float64 {
	return []float64{0}
}
func (w *WeightedSnapshot) Size() int64 {
	return 0
}
func (w *WeightedSnapshot) Max() int64 {
	return 0.0

}
func (w *WeightedSnapshot) StdDev() float64 {
	return 0.0

}
func (w *WeightedSnapshot) Mean() float64 {
	return 0.0

}

//Median is 50th
func (w *WeightedSnapshot) Median() float64 {
	return 0.0

}

func (w *WeightedSnapshot) Get75thPercentile() float64 {
	return 0.0

}
func (w *WeightedSnapshot) Get95thPercentile() float64 {
	return 0.0

}
func (w *WeightedSnapshot) Get98thPercentile() float64 {
	return 0.0

}
func (w *WeightedSnapshot) Get99thPercentile() float64 {
	return 0.0

}
func (w *WeightedSnapshot) Get999thPercentile() float64 {
	return 0.0
}
func (w *WeightedSnapshot) Percentiles() []float64 {
	median := w.Median()
	p75 := w.Get75thPercentile()
	p95 := w.Get95thPercentile()
	p99 := w.Get99thPercentile()
	p999 := w.Get999thPercentile()

	return []float64{median, p75, p95, p99, p999}
}
