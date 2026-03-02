package tsdb

import (
	"fmt"
	"sort"
	"sync"
)

// data points for unique set of labels
type series struct {
	mu sync.RWMutex
	labels []Label
	points []Point
}

func (s *series) appendPoint(t int64, v float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.points) > 0 {
		var lastTime = s.points[len(s.points)-1].T
		if t <= lastTime {
			return fmt.Errorf("out of order append")
		}
	}
	s.points = append(s.points, Point{T: t, V: v})
	return nil
}

func (s *series) queryRange(start, end int64) []Point {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var startIndex = sort.Search(len(s.points), func(i int) bool {
		return s.points[i].T >= start
	})

	var endIndex = sort.Search(len(s.points), func(i int) bool {
		return s.points[i].T > end
 	})

	if startIndex >= endIndex {
		return nil
	}

	// avoid external slice corruption
	var result = make([]Point, endIndex - startIndex)
	copy(result, s.points[startIndex:endIndex])

	return result
}