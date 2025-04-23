package models

import (
	"time"
)

// SyncStatus represents the sync status of a file
type SyncStatus string

const (
	StatusNotSynced SyncStatus = "NOT_SYNCED"
	StatusSyncing   SyncStatus = "SYNCING"
	StatusSynced    SyncStatus = "SYNCED"
	StatusError     SyncStatus = "ERROR"
)

// FileInfo represents a file's sync information
type FileInfo struct {
	Path         string                `json:"path"`
	Status       SyncStatus            `json:"status"`
	LastModified time.Time             `json:"lastModified"`
	Size         int64                 `json:"size"`
	IsDownloaded bool                  `json:"isDownloaded"`
	IsDirectory  bool                  `json:"isDirectory"`
	Version      int                   `json:"version"`
	Checksum     string                `json:"checksum,omitempty"`
	LastSynced   time.Time             `json:"lastSynced"`
	FilesContent map[string]*FileInfo  `json:"filesContent,omitempty"`
}

// FileEvent represents a file system event
type FileEvent struct {
	Type      FileEventType `json:"type"`
	Path      string        `json:"path"`
	Timestamp time.Time     `json:"timestamp"`
}

// FileEventType represents the type of file event
type FileEventType string

const (
	EventCreated  FileEventType = "CREATED"
	EventModified FileEventType = "MODIFIED"
	EventDeleted  FileEventType = "DELETED"
	EventRenamed  FileEventType = "RENAMED"
)