package main

import (
	"log"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

var g singleflight.Group

type DataRow int64

func main() {
	loadFn := func(k string) (DataRow, error) {
		v, err, _ := g.Do(k, func() (interface{}, error) {
			return loadDataFromDB(), nil
		})
		if err != nil {
			return 0, err
		}
		return v.(DataRow), nil
	}

	var wg sync.WaitGroup
	size := 100
	wg.Add(size)

	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			data, err := loadFn("mock-key")
			if err != nil {
				log.Print(err)
				return
			}
			// Result: only query once
			log.Println(data)
		}()
	}
	wg.Wait()
}

// loadDataFromDB
func loadDataFromDB() DataRow {
	log.Println("Loading data from DB")
	return DataRow(time.Now().UnixNano())
}
