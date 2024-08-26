package main

import (
	"flag"
	"fmt"
)

func run(action, directory, checksumFile string) error {
	switch action {
	case "generate":
		if err := GenerateChecksums(directory, checksumFile); err != nil {
			return fmt.Errorf("error generating checksums: %w", err)
		}
	case "verify":
		if err := VerifyChecksums(checksumFile); err != nil {
			return fmt.Errorf("error verifying checksums: %w", err)
		}
	default:
		return fmt.Errorf("unknown action. Use 'generate' or 'verify'")
	}
	return nil
}

func main() {
	action := flag.String("action", "generate", "Action to perform: generate or verify checksums")
	directory := flag.String("directory", ".", "Directory to process")
	checksumFile := flag.String("checksumFile", "checksums.txt", "File to store or read checksums")
	flag.Parse()

	if err := run(*action, *directory, *checksumFile); err != nil {
		fmt.Println(err)
	}
}
