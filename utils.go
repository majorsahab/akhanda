package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CountFilesAndDirs(directory string) (int64, int64) {
	fileCount := int64(0)
	dirCount := int64(0)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			fileCount++
		} else if info.IsDir() {
			dirCount++
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error counting files and directories:", err)
	}
	return fileCount, dirCount
}

func FormatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func CountFilesAndDirsInCheckSumFile(checksumFile string) (files int64, dirs int64, err error) {
	file, err := os.Open(checksumFile)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var rootDir string
	var allPaths []string

	// Read the first line to get the root directory
	if scanner.Scan() {
		rootDirLine := scanner.Text()
		_, path := parseDirPath(rootDirLine)
		if path == "" {
			return 0, 0, fmt.Errorf("root directory line is empty")
		}
		rootDir = filepath.Clean(path)
	}

	// Process the rest of the lines
	for scanner.Scan() {
		line := scanner.Text()
		_, path := parseDirPath(line)
		if path == "" {
			continue
		}

		cleanPath := filepath.Clean(path)
		dirPath := filepath.Dir(cleanPath)

		// Add all directories leading up to the file
		for dirPath != rootDir && dirPath != "." {
			if !contains(allPaths, dirPath) {
				allPaths = append(allPaths, dirPath)
				dirs++
			}
			dirPath = filepath.Dir(dirPath)
		}

		// Mark the file path
		files++
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	// Handle the root directory separately if it was not already counted
	if !contains(allPaths, rootDir) && rootDir != "." {
		dirs++
	}

	return files, dirs, nil
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func parseDirPath(line string) (string, string) {
	parts := strings.SplitN(line, "  ", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], strings.Trim(parts[1], "\"")
}

// parseLine extracts checksum and file path from a line of the checksum file.
// It assumes that paths may be enclosed in double quotes.
func parseLine(line string) (string, string) {
	var checksum, path string
	_, err := fmt.Sscanf(line, "%s %q", &checksum, &path)
	if err != nil {
		// If parsing fails, return empty strings
		return "", ""
	}
	return checksum, path
}

func getDirectoryFromChecksumFile(checksumFile string) (string, error) {
	file, err := os.Open(checksumFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		_, path := parseLine(line)
		if path != "" {
			// Clean the path and get its directory
			return filepath.Dir(filepath.Clean(path)), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("no valid paths found in checksum file")
}
