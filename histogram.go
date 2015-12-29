package metrics

// Histograms calculate distribution statistics from a series of int64 values.
type Histogram struct {
	count     int64
	reservoir Reservoir
	Snapshot  func() Snapshot
	Update    func(int64)
	Count     func() int64
}

func NewHistogram(r Reservoir) *Histogram {
	return &Histogram{
		count:     0,
		reservoir: r,
		Snapshot:  StandardHistogramSnapshot,
		Update:    StandardHistogramUpdate,
		Count:     StandardHistogramCount,
	}
}

func StandardHistogramSnapshot() Snapshot {
	return nil
}
func StandardHistogramUpdate(int64) {
}
func StandardHistogramCount() int64 {
	return 0

}
