package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func VerifyChecksums(checksumFile string) error {
	totalFiles, totalDirs, err := CountFilesAndDirsInCheckSumFile(checksumFile)
	if err != nil {
		return err
	}

	directory, err := getDirectoryFromChecksumFile(checksumFile)
	if err != nil {
		return err
	}

	progress := NewProgressBar(totalFiles, totalDirs)

	file, err := os.Open(checksumFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	lines := make(chan string, numWorkers)
	results := make(chan string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				var storedChecksum, relativePath string
				_, err := fmt.Sscanf(line, "%s %q", &storedChecksum, &relativePath)
				if err != nil {
					results <- fmt.Sprintf("Error parsing line: %s", line)
					continue
				}

				// Remove quotes from the file path if necessary
				filePath := strings.Trim(relativePath, "\"")
				filePath = filepath.Clean(filePath)
				if !filepath.IsAbs(filePath) {
					filePath = filepath.Join(directory, filePath)
				}

				// Skip directories
				info, err := os.Stat(filePath)
				if err != nil {
					results <- fmt.Sprintf("Error accessing file: %s", filePath)
					continue
				}
				if info.IsDir() {
					progress.IncrementDir()
					continue
				}

				// Calculate checksum for regular files only
				currentChecksum, err := calculateSHA256(filePath)
				if err != nil {
					results <- fmt.Sprintf("Error calculating checksum for %s: %v", filePath, err)
					continue
				}

				if currentChecksum != storedChecksum {
					results <- fmt.Sprintf("File %s has been modified.", filePath)
				}

				progress.IncrementFile()
			}
		}()
	}

	go func() {
		defer close(lines)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines <- line
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading checksum file:", err)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

	progress.Complete()
	return nil
}
