package tsdb

type Point struct {
	T int64
	V float64
}

type Label struct {
	Name  string
	Value string
}