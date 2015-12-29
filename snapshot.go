package metrics

type Snapshot interface {
	GetValue(percentile float64) float64
	GetValues() []float64
	Size() int64
	GetMax() int64
	GetStdDev() float64
	GetMean() float64
	GetMin() int64
	GetMedian() float64 //Median is 50th
	Get75thPercentile() float64
	Get95thPercentile() float64
	Get98thPercentile() float64
	Get99thPercentile() float64
	Get999thPercentile() float64
}

type WeightedSnapshot struct {
}

func (w *WeightedSnapshot) GetValue(percentile float64) float64 {

}
func (w *WeightedSnapshot) GetValues() []float64 {

}
func (w *WeightedSnapshot) Size() int64 {

}
func (w *WeightedSnapshot) GetMax() int64 {

}
func (w *WeightedSnapshot) GetStdDev() float64 {

}
func (w *WeightedSnapshot) GetMean() float64 {

}

//Median is 50th
func (w *WeightedSnapshot) GetMedian() float64 {

}

func (w *WeightedSnapshot) Get75thPercentile() float64 {

}
func (w *WeightedSnapshot) Get95thPercentile() float64 {

}
func (w *WeightedSnapshot) Get98thPercentile() float64 {

}
func (w *WeightedSnapshot) Get99thPercentile() float64 {

}
func (w *WeightedSnapshot) Get999thPercentile() float64 {

}
