package main

import (
	"sync"
	"testing"
)

func TestFileWorker(t *testing.T) {
	files := make(chan string, 2)
	results := make(chan [2]string, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go fileWorker(files, results, &wg)

	// Send test data
	files <- "file1.txt"
	files <- "file2.txt"
	close(files)

	// Wait for the worker to finish
	wg.Wait()
	close(results)

	// Check results
	expectedResults := map[string]string{
		"file1.txt": "mockChecksum",
		"file2.txt": "mockChecksum",
	}

	for result := range results {
		checksum, filePath := result[0], result[1]
		if expectedChecksum, ok := expectedResults[filePath]; !ok || expectedChecksum != checksum {
			t.Errorf("unexpected result: got %v, want %v", result, [2]string{expectedChecksum, filePath})
		}
	}
}
