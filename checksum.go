package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type worker struct {
	hasher hash.Hash
	buf    []byte
}

var workerPool = sync.Pool{
	New: func() interface{} {
		return &worker{
			hasher: sha256.New(),
			buf:    make([]byte, 64*1024),
		}
	},
}

func calculateSHA256(filePath string) (string, error) {
	w := workerPool.Get().(*worker)
	defer workerPool.Put(w)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	w.hasher.Reset()
	reader := bufio.NewReader(file)
	for {
		n, err := reader.Read(w.buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		w.hasher.Write(w.buf[:n])
	}

	return hex.EncodeToString(w.hasher.Sum(nil)), nil
}

func writeChecksums(results <-chan [2]string, checksumFile string, progress *ProgressBar) error {
	file, err := os.Create(checksumFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var parentDir string
	firstEntry := true

	for result := range results {
		checksum := result[0]
		path := result[1]

		// Extract parent directory from the path and write it only once at the start
		if firstEntry {
			parentDir = filepath.Dir(filepath.Clean(path))
			_, err := fmt.Fprintf(file, "dir  \"%s\"\n", parentDir)
			if err != nil {
				return err
			}
			firstEntry = false
		}

		// Use double quotes around the file path to handle spaces and special characters
		_, err := fmt.Fprintf(file, "%s  \"%s\"\n", checksum, path)
		if err != nil {
			return err
		}
		progress.IncrementFile()
	}
	progress.Complete()
	return nil
}

func GenerateChecksums(directory, checksumFile string) error {
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	files := make(chan string, numWorkers)
	results := make(chan [2]string, numWorkers)

	totalFiles, totalDirs := CountFilesAndDirs(directory)
	if totalFiles == 0 {
		return fmt.Errorf("no files found in directory")
	}

	progress := NewProgressBar(totalFiles, totalDirs)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go fileWorker(files, results, &wg)
	}

	go func() {
		err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error walking directory:", err)
				return nil
			}
			if info.Mode().IsRegular() {
				files <- path
			} else if info.IsDir() {
				progress.IncrementDir()
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error walking the directory:", err)
		}
		close(files)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	err := writeChecksums(results, checksumFile, progress)
	if err != nil {
		return err
	}

	fmt.Println("Checksums generated and saved in", checksumFile)
	return nil
}
