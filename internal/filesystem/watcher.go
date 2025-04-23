// internal/filesystem/watcher.go
package filesystem

import (
	"fmt"
	"path/filepath"
	"time"

	"homecloud/internal/models"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches a directory for changes
type Watcher struct {
	watchDir   string
	eventChan  chan models.FileEvent
	fsWatcher  *fsnotify.Watcher
	stopChan   chan struct{}
	isWatching bool
}

// NewWatcher creates a new file system watcher
func NewWatcher(watchDir string, eventChan chan models.FileEvent) (*Watcher, error) {
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
			w.eventChan <- models.FileEvent{
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
func (w *Watcher) getEventType(event fsnotify.Event) models.FileEventType {
	if event.Has(fsnotify.Create) {
		return models.EventCreated
	} else if event.Has(fsnotify.Write) {
		return models.EventModified
	} else if event.Has(fsnotify.Remove) {
		return models.EventDeleted
	} else if event.Has(fsnotify.Rename) {
		return models.EventRenamed
	}
	return models.EventModified // Default
}

// GetFileName extracts the filename from a path
func GetFileName(path string) string {
	return filepath.Base(path)
}
