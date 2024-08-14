package main

import (
	"sync"
)

type Spinner struct {
	frames []string
	index  int
	mu     sync.Mutex
}

func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{
			"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷",
		},
		index: 0,
	}
}

func (s *Spinner) Next() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	frame := s.frames[s.index]
	s.index = (s.index + 1) % len(s.frames)
	return frame
}
