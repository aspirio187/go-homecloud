// internal/sync/manager.go
package sync

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"homecloud/internal/filesystem"
	"homecloud/internal/models"
)

// SyncManager handles file synchronization
type SyncManager struct {
	watchDir   string
	eventChan  chan models.FileEvent
	watcher    *filesystem.Watcher
	fileInfos  map[string]*models.FileInfo
	mu         sync.RWMutex
	isRunning  bool
	statusChan chan *models.FileInfo
}

// NewSyncManager creates a new sync manager
func NewSyncManager(watchDir string) *SyncManager {
	return &SyncManager{
		watchDir:   watchDir,
		eventChan:  make(chan models.FileEvent),
		fileInfos:  make(map[string]*models.FileInfo),
		statusChan: make(chan *models.FileInfo, 100),
		isRunning:  false,
	}
}

// Start begins the sync manager
func (sm *SyncManager) Start() error {
	// Read initial files
	sm.initialRead()

	if sm.isRunning {
		return fmt.Errorf("sync manager is already running")
	}

	// Create and start a watcher
	watcher, err := filesystem.NewWatcher(sm.watchDir, sm.eventChan)
	if err != nil {
		return err
	}

	sm.watcher = watcher
	err = sm.watcher.Start()
	if err != nil {
		return err
	}

	sm.isRunning = true

	// Start processing events
	go sm.processEvents()

	return nil
}

// Stop stops the sync manager
func (sm *SyncManager) Stop() {
	if !sm.isRunning {
		return
	}

	sm.isRunning = false
	sm.watcher.Stop()
	close(sm.eventChan)
}

// GetFileInfos returns a copy of all file infos
func (sm *SyncManager) GetFileInfos() []*models.FileInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make([]*models.FileInfo, 0, len(sm.fileInfos))
	for _, info := range sm.fileInfos {
		result = append(result, info)
	}
	return result
}

// GetStatusChannel returns a channel for receiving status updates
func (sm *SyncManager) GetStatusChannel() <-chan *models.FileInfo {
	return sm.statusChan
}

// processEvents handles file events
func (sm *SyncManager) processEvents() {
	for event := range sm.eventChan {
		switch event.Type {
		case models.EventCreated, models.EventModified:
			sm.handleFileChange(event.Path, event.Timestamp)
		case models.EventDeleted:
			sm.handleFileDelete(event.Path)
		}
	}
}

// handleFileChange processes a file creation or modification
func (sm *SyncManager) handleFileChange(path string, timestamp time.Time) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// For now, just track the file status
	info := &models.FileInfo{
		Path:         path,
		Status:       models.StatusNotSynced,
		LastModified: timestamp,
		IsDownloaded: true, // It's a local file, so it's "downloaded"
		Version:      1,
	}

	sm.fileInfos[path] = info

	// In a real implementation, you would:
	// 1. Check if the file exists on the server
	// 2. Upload the file if needed
	// 3. Update the status accordingly

	// Mock a sync process
	go func() {
		// Update status to syncing
		sm.updateFileStatus(path, models.StatusSyncing)

		// Simulate sync time
		time.Sleep(2 * time.Second)

		// Update status to synced
		sm.updateFileStatus(path, models.StatusSynced)
	}()
}

// handleFileDelete processes a file deletion
func (sm *SyncManager) handleFileDelete(path string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Check if the file exists in our tracking map
	_, exists := sm.fileInfos[path]
	if exists {
		// Remove the file from tracking
		delete(sm.fileInfos, path)

		// Send a deleting notification through the status channel
		sm.statusChan <- &models.FileInfo{
			Path:   path,
			Status: models.StatusNotSynced,
		}
	}

	// In a real implementation, you would:
	// 1. Delete the file from the server
	// 2. Update any relevant tracking information
}

// updateFileStatus updates a file's status and notifies listeners
func (sm *SyncManager) updateFileStatus(path string, status models.SyncStatus) {
	sm.mu.Lock()
	info, exists := sm.fileInfos[path]
	if exists {
		info.Status = status
		// Clone the info to avoid race conditions
		updatedInfo := *info
		sm.statusChan <- &updatedInfo
	}
	sm.mu.Unlock()
}

func (sm *SyncManager) initialRead() {
	// This function is called when the app starts
	// It reads the initial files in the watch directory
	// and adds them to the sync manager
	// This is useful for syncing existing files on startup

	var initialContent []*models.FileInfo = make([]*models.FileInfo, 0)

	filepath.WalkDir(sm.watchDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == sm.watchDir {
			return nil
		}

		info, err := d.Info()

		if err != nil {
			fmt.Println("Error:", err)

			return nil
		}

		name := info.Name()
		isDir := d.IsDir()

		fileInfo := &models.FileInfo{
			Path:         path,
			Status:       models.StatusSynced,
			LastModified: info.ModTime(),
			Size:         info.Size(),
			IsDownloaded: true,
			IsDirectory:  isDir,
			Version:      1,
			LastSynced:   time.Now(),
		}

		if !isDir && path != sm.watchDir {

			filepathDir := filepath.Dir(path)
			// Get the file info directory and add this to it's content
			i := slices.IndexFunc(initialContent, func(e *models.FileInfo) bool {
				return e.Path == filepathDir
			})

			if i != -1 {
				if initialContent[i].FilesContent == nil {
					initialContent[i].FilesContent = make(map[string]*models.FileInfo)
				}

				// Add the file to the directory
				initialContent[i].FilesContent[name] = fileInfo

				return nil
			}
		}

		initialContent = append(initialContent, fileInfo)

		return nil
	})

	sm.fileInfos = make(map[string]*models.FileInfo)

	for _, info := range initialContent {
		sm.fileInfos[info.Path] = info
	}
}
