package metrics

import "github.com/carbin-gun/awesome-metrics/output"

type StandardHistogram struct {
	count     int64
	reservoir Reservoir
}

func NewHistogram(reservoir Reservoir) {
	return &StandardHistogram{
		count:     0,
		reservoir: reservoir,
	}
}

//Counting interface
func (histogram *StandardHistogram) Count() int64 {
	return 0
}

//communication
func (histogram *StandardHistogram) Update(int64) {

}

//snapshot data about histogram
func (histogram *StandardHistogram) Snapshot() output.Snapshot {
	return nil
}
