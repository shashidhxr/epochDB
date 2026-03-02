package tsdb

import (
	"hash/fnv"
	"sort"
)

func hashLabels(labels []Label) uint64 {
	var sorted = make([]Label, len(labels))
	copy(sorted, labels)

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Name == sorted[j].Name {
			return sorted[i].Value < sorted[j].Value
		}
		return sorted[i].Name < sorted[j].Name
	})

	var h = fnv.New64a()
	for _, l := range sorted {
		h.Write([]byte(l.Name))
		h.Write([]byte{'\xff'})
		h.Write([]byte(l.Value))
		h.Write([]byte{'\xff'})
	}
	return h.Sum64()
}