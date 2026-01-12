package ui

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ScannerService struct {
	entry            *widget.Entry
	onScan           func(string)
	mu               sync.RWMutex
	focusLockEnabled bool
	stopOnce         sync.Once
	stopCh           chan struct{}
}

func NewScannerService() *ScannerService {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("")
	entry.Hide()

	scanner := &ScannerService{
		entry:            entry,
		focusLockEnabled: true,
		stopCh:           make(chan struct{}),
	}

	entry.OnSubmitted = func(value string) {
		scanner.mu.RLock()
		handler := scanner.onScan
		scanner.mu.RUnlock()
		if handler != nil {
			handler(value)
		}
		entry.SetText("")
	}

	return scanner
}

func (s *ScannerService) Widget() *widget.Entry {
	return s.entry
}

func (s *ScannerService) OnScan(handler func(string)) {
	s.mu.Lock()
	s.onScan = handler
	s.mu.Unlock()
}

func (s *ScannerService) Start(window fyne.Window) {
	canvas := window.Canvas()
	canvas.Focus(s.entry)

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.mu.RLock()
				enabled := s.focusLockEnabled
				s.mu.RUnlock()
				if !enabled {
					continue
				}
				if canvas.Focused() != s.entry {
					canvas.Focus(s.entry)
				}
			case <-s.stopCh:
				return
			}
		}
	}()
}

func (s *ScannerService) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
	})
}

func (s *ScannerService) SetFocusLockEnabled(enabled bool) {
	s.mu.Lock()
	s.focusLockEnabled = enabled
	s.mu.Unlock()
}
