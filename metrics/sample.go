package metrics

import "github.com/mtchavez/skiplist"

type WeightedSample struct {
	weight float64
	value  int64
}

type WeightedSampleStorage struct {
	store skiplist.SkipList
}

func (s *WeightedSampleStorage) Delete(key int) bool {
	return s.store.Delete(key)
}
func (s *WeightedSampleStorage) Insert(key int, val []byte) *skiplist.Node {
	return s.store.Insert(key, val)
}
func (s *WeightedSampleStorage) Iterator() skiplist.Iterator {
	return s.store.Iterator()
}
func (s *WeightedSampleStorage) Clear() {
	s.store = skiplist.NewList()
}
func (s *WeightedSampleStorage) First() (int, []byte) {
	iterator := s.store.Iterator()
	return iterator.Key(), iterator.Val()
}
