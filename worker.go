package main

import (
	"sync"
)

func fileWorker(files <-chan string, results chan<- [2]string, wg *sync.WaitGroup) {
	defer wg.Done()
	for filePath := range files {
		checksum, err := calculateSHA256(filePath)
		if err == nil {
			results <- [2]string{checksum, filePath}
		}
	}
}
