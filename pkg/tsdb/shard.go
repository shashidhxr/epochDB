package tsdb

import "sync"

type shard struct {
	mu sync.RWMutex
	series map[uint64]*series
}

func newShard() *shard {
	return &shard{
		series: make(map[uint64]*series),
	}
}

func (s *shard) getSeries(hash uint64, labels []Label) *series {
	// fast path
	s.mu.RLock()
	ser, exists := s.series[hash]
	s.mu.RUnlock()

	if exists {
		return ser
	}

	// slow path
	s.mu.Lock()
	defer s.mu.Unlock()

	// check if other routine created while waiting for the Rlock
	ser, exists = s.series[hash]
	if exists {
		return ser
	}

	ser = &series {
		labels: labels,
		points: make([]Point, 0, 128),
	}
	s.series[hash] = ser
	return ser
}