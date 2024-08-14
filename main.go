package main

import (
	"flag"
	"fmt"
)

func main() {
	action := flag.String("action", "generate", "Action to perform: generate or verify checksums")
	directory := flag.String("directory", ".", "Directory to process")
	checksumFile := flag.String("checksumFile", "checksums.txt", "File to store or read checksums")
	flag.Parse()

	switch *action {
	case "generate":
		if err := GenerateChecksums(*directory, *checksumFile); err != nil {
			fmt.Println("Error generating checksums:", err)
		}
	case "verify":
		if err := VerifyChecksums(*checksumFile); err != nil {
			fmt.Println("Error verifying checksums:", err)
		}
	default:
		fmt.Println("Unknown action. Use 'generate' or 'verify'.")
	}
}
