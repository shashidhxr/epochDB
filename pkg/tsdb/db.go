package tsdb

const numShards = 64

type DB interface {
	Append(labels []Label, t int64, v float64) error
	Query(labels []Label, start, end int64) ([]Point, error)
}

// sharded db engine
type MemDB struct {
	shards [numShards]*shard
}

// init db and shards
func NewMemDB() *MemDB {
	var db = &MemDB{}
	for i := 0; i < numShards; i++ {
		db.shards[i] = newShard()
	}
	return db
}

func (db *MemDB) getShard(hash uint64) *shard {
	return db.shards[hash % numShards]
}

func (db *MemDB) Append(labels []Label, t int64, v float64) error {
	var hash = hashLabels(labels)
	var shard = db.getShard(hash)

	var ser = shard.getSeries(hash, labels)
	return ser.appendPoint(t, v)
}

func (db *MemDB) Query(labels []Label, start, end int64) ([]Point, error) {
	var hash = hashLabels(labels)
	var shard = db.getShard(hash)

	// verify getSeries and getOrCreateSeries
	var ser = shard.getSeries(hash, labels)
	if ser == nil {
		return nil, nil
	}

	return ser.queryRange(start, end), nil
}