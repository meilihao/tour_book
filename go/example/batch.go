package main

import (
	"fmt"
)

const batchSize = 3


func processBatch(list []int) {
	fmt.Println(list)
}

func process(data []int) {
	for start, end := 0, 0; start <= len(data)-1; start = end {
		end = start + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[start:end]
		processBatch(batch)
	}
	fmt.Println("done processing all data")
}

func main() {
	data := make([]int, 0)
	for i := 1; i < 8; i++ {
		data = append(data, i)
	}
	process(data)
}