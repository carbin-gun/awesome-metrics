package metrics

import (
	"encoding/json"

	"github.com/carbin-gun/skiplist"
)

type WeightedSample struct {
	weight float64
	value  int64
}

func (s WeightedSample) MarshalBytes() []byte {
	v, _ := json.Marshal(s)
	return v
}
func UnMarshalFromBytes(bytes []byte) WeightedSample {
	var val WeightedSample
	json.Unmarshal(bytes, &val)
	return val
}

type WeightedSampleStorage struct {
	store skiplist.SkipList
}

func (s *WeightedSampleStorage) Delete(key float64) bool {
	return s.store.Delete(key)
}
func (s *WeightedSampleStorage) Insert(key float64, val []byte) *skiplist.Node {
	return s.store.Insert(key, val)
}
func (s *WeightedSampleStorage) Iterator() skiplist.Iterator {
	return s.store.Iterator()
}
func (s *WeightedSampleStorage) Clear() {
	s.store = skiplist.NewList()
}
func (s *WeightedSampleStorage) First() (float64, []byte) {
	iterator := s.store.Iterator()
	return iterator.Key(), iterator.Val()
}
