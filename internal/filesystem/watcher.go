// internal/filesystem/watcher.go
package filesystem

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// EventType represents the type of file event
type EventType string

const (
	Created  EventType = "CREATED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
	Renamed  EventType = "RENAMED"
)

// FileEvent represents a file system event
type FileEvent struct {
	Type      EventType
	Path      string
	Timestamp time.Time
}

// Watcher watches a directory for changes
type Watcher struct {
	watchDir   string
	eventChan  chan FileEvent
	fsWatcher  *fsnotify.Watcher
	stopChan   chan struct{}
	isWatching bool
}

// NewWatcher creates a new file system watcher
func NewWatcher(watchDir string, eventChan chan FileEvent) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Watcher{
		watchDir:   watchDir,
		eventChan:  eventChan,
		fsWatcher:  fsWatcher,
		stopChan:   make(chan struct{}),
		isWatching: false,
	}, nil
}

// Start begins watching the directory
func (w *Watcher) Start() error {
	if w.isWatching {
		return fmt.Errorf("watcher is already running")
	}

	err := w.fsWatcher.Add(w.watchDir)
	if err != nil {
		return fmt.Errorf("failed to add directory to watcher: %w", err)
	}

	w.isWatching = true

	go w.watchLoop()

	return nil
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	if !w.isWatching {
		return
	}

	w.isWatching = false
	close(w.stopChan)
	w.fsWatcher.Close()
}

// watchLoop processes file system events
func (w *Watcher) watchLoop() {
	for {
		select {
		case <-w.stopChan:
			return
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			eventType := w.getEventType(event)
			w.eventChan <- FileEvent{
				Type:      eventType,
				Path:      event.Name,
				Timestamp: time.Now(),
			}
		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Error: %v\n", err)
		}
	}
}

// getEventType determines the type of file event
func (w *Watcher) getEventType(event fsnotify.Event) EventType {
	if event.Has(fsnotify.Create) {
		return Created
	} else if event.Has(fsnotify.Write) {
		return Modified
	} else if event.Has(fsnotify.Remove) {
		return Deleted
	} else if event.Has(fsnotify.Rename) {
		return Renamed
	}
	return Modified // Default
}

// GetFileName extracts the filename from a path
func GetFileName(path string) string {
	return filepath.Base(path)
}
