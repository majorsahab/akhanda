package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		action       string
		directory    string
		checksumFile string
		expectError  bool
	}{
		{"generate", ".", "checksums.txt", false},
		{"verify", ".", "checksums.txt", false},
		{"invalid", ".", "checksums.txt", true},
	}

	for _, tt := range tests {
		err := run(tt.action, tt.directory, tt.checksumFile)
		if (err != nil) != tt.expectError {
			t.Errorf("run(%q, %q, %q) error = %v, expectError %v", tt.action, tt.directory, tt.checksumFile, err, tt.expectError)
		}
	}
}
