package main

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type ProgressBar struct {
	totalFiles      int64
	totalDirs       int64
	completedFiles  int64
	completedDirs   int64
	startTime       time.Time
	style           lipgloss.Style
	barStyle        lipgloss.Style
	completedStyle  lipgloss.Style
	incompleteStyle lipgloss.Style
	spinner         *Spinner
	stopSpinner     chan struct{}
}

func NewProgressBar(totalFiles, totalDirs int64) *ProgressBar {
	progress := &ProgressBar{
		totalFiles:      totalFiles,
		totalDirs:       totalDirs,
		startTime:       time.Now(),
		style:           lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true).Padding(0, 1),
		barStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("240")).Padding(0, 1),
		completedStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		incompleteStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		spinner:         NewSpinner(),
		stopSpinner:     make(chan struct{}),
	}
	go progress.runSpinner()
	return progress
}

func (p *ProgressBar) IncrementFile() {
	atomic.AddInt64(&p.completedFiles, 1)
	p.Display()
}

func (p *ProgressBar) IncrementDir() {
	atomic.AddInt64(&p.completedDirs, 1)
	p.Display()
}

func (p *ProgressBar) Complete() {
	atomic.StoreInt64(&p.completedFiles, p.totalFiles)
	atomic.StoreInt64(&p.completedDirs, p.totalDirs)
	p.Display()
	close(p.stopSpinner)
	fmt.Println()
}

func (p *ProgressBar) Display() {
	percent := float64(atomic.LoadInt64(&p.completedFiles)) / float64(p.totalFiles) * 100
	elapsed := time.Since(p.startTime)
	eta := time.Duration((elapsed.Seconds()/float64(atomic.LoadInt64(&p.completedFiles)))*float64(p.totalFiles-atomic.LoadInt64(&p.completedFiles))) * time.Second

	etaStr := FormatDuration(eta)

	barWidth := 40
	complete := int(float64(barWidth) * percent / 100)
	incomplete := barWidth - complete

	progressBar := fmt.Sprintf(
		"%s %s[%s%s] %.2f%% %s Files: %d/%d Dirs: %d/%d",
		p.spinner.Next(),
		p.style.Render("Progress: "),
		p.completedStyle.Render(strings.Repeat("█", complete)),
		p.incompleteStyle.Render(strings.Repeat("░", incomplete)),
		percent,
		etaStr,
		atomic.LoadInt64(&p.completedFiles),
		p.totalFiles,
		atomic.LoadInt64(&p.completedDirs),
		p.totalDirs,
	)
	fmt.Printf("\r%s", progressBar)
}

func (p *ProgressBar) runSpinner() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			p.Display()
		case <-p.stopSpinner:
			return
		}
	}
}
