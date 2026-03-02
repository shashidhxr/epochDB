package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/shashidhxr/epochDB/pkg/tsdb"
)

func main() {
	var db tsdb.DB = tsdb.NewMemDB()

	var labels1 = []tsdb.Label{
		{
			Name: "metric", Value: "cpu_usage",
		}, {
			Name: "host", Value: "server-na-1",
		},
	}
	var labels2 = []tsdb.Label{
		{
			Name: "metric", Value: "cpu_usage",
		}, {
			Name: "host", Value: "server-eu-1",
		},
	}

	fmt.Println("Starting concurrent ingestion")
	var baseTime = time.Now().UnixMilli()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			t := baseTime + int64(i * 10)
			_ = db.Append(labels1, t, float64(i))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			t := baseTime + int64(i * 10)
			_ = db.Append(labels2, t, float64(i * 2))
		}
	}()

	wg.Wait()
	fmt.Println("Ingestion complete")

	fmt.Println("Querying data")
	var startQuery = baseTime + 500
	var endQuery = baseTime + 650

	results, _ := db.Query(labels1, startQuery, endQuery)
	for _, p := range results {
		fmt.Printf("server-na -> T:%d, V:%.0f\n", p.T, p.V)
	}
}