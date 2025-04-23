package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"homecloud/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// MetadataStore handles storage of file metadata
type MetadataStore struct {
	db *sql.DB
}

// NewMetadataStore creates a new metadata storage instance
func NewMetadataStore(appDataPath string) (*MetadataStore, error) {
	// Ensure directory exists
	dbDir := filepath.Join(appDataPath, "db")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "homecloud.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tables if they don't exist
	if err := initDatabase(db); err != nil {
		db.Close()
		return nil, err
	}

	return &MetadataStore{db: db}, nil
}

// Close closes the database connection
func (m *MetadataStore) Close() error {
	return m.db.Close()
}

// SaveFileInfo saves or updates file information
func (m *MetadataStore) SaveFileInfo(info *models.FileInfo) error {
	// Check if file exists
	var exists bool
	err := m.db.QueryRow("SELECT EXISTS(SELECT 1 FROM files WHERE path = ?)", info.Path).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	if exists {
		// Update existing file
		_, err = m.db.Exec(
			"UPDATE files SET status = ?, last_modified = ?, size = ?, is_downloaded = ?, is_directory = ?, version = ?, checksum = ?, last_synced = ? WHERE path = ?",
			string(info.Status),
			info.LastModified.Unix(),
			info.Size,
			info.IsDownloaded,
			info.IsDirectory,
			info.Version,
			info.Checksum,
			info.LastSynced.Unix(),
			info.Path,
		)
	} else {
		// Insert new file
		_, err = m.db.Exec(
			"INSERT INTO files (path, status, last_modified, size, is_downloaded, is_directory, version, checksum, last_synced) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			info.Path,
			string(info.Status),
			info.LastModified.Unix(),
			info.Size,
			info.IsDownloaded,
			info.IsDirectory,
			info.Version,
			info.Checksum,
			info.LastSynced.Unix(),
		)
	}

	return err
}

// GetFileInfo retrieves file information by path
func (m *MetadataStore) GetFileInfo(path string) (*models.FileInfo, error) {
	var info models.FileInfo
	var lastModified, lastSynced int64
	var statusStr string

	err := m.db.QueryRow(
		"SELECT path, status, last_modified, size, is_downloaded, is_directory, version, checksum, last_synced FROM files WHERE path = ?",
		path,
	).Scan(
		&info.Path,
		&statusStr,
		&lastModified,
		&info.Size,
		&info.IsDownloaded,
		&info.IsDirectory,
		&info.Version,
		&info.Checksum,
		&lastSynced,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	info.Status = models.SyncStatus(statusStr)
	info.LastModified = time.Unix(lastModified, 0)
	info.LastSynced = time.Unix(lastSynced, 0)

	return &info, nil
}

// ListAllFiles returns all file metadata
func (m *MetadataStore) ListAllFiles() ([]*models.FileInfo, error) {
	rows, err := m.db.Query("SELECT path, status, last_modified, size, is_downloaded, is_directory, version, checksum, last_synced FROM files")
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []*models.FileInfo
	for rows.Next() {
		var info models.FileInfo
		var lastModified, lastSynced int64
		var statusStr string

		err := rows.Scan(
			&info.Path,
			&statusStr,
			&lastModified,
			&info.Size,
			&info.IsDownloaded,
			&info.IsDirectory,
			&info.Version,
			&info.Checksum,
			&lastSynced,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}

		info.Status = models.SyncStatus(statusStr)
		info.LastModified = time.Unix(lastModified, 0)
		info.LastSynced = time.Unix(lastSynced, 0)

		files = append(files, &info)
	}

	return files, nil
}

// DeleteFileInfo removes a file from the database
func (m *MetadataStore) DeleteFileInfo(path string) error {
	_, err := m.db.Exec("DELETE FROM files WHERE path = ?", path)
	return err
}

// GetSyncQueue returns files that need to be synchronized
func (m *MetadataStore) GetSyncQueue() ([]*models.FileInfo, error) {
	rows, err := m.db.Query("SELECT path, status, last_modified, size, is_downloaded, is_directory, version, checksum, last_synced FROM files WHERE status != ?", string(models.StatusSynced))
	if err != nil {
		return nil, fmt.Errorf("failed to query sync queue: %w", err)
	}
	defer rows.Close()

	var files []*models.FileInfo
	for rows.Next() {
		var info models.FileInfo
		var lastModified, lastSynced int64
		var statusStr string

		err := rows.Scan(
			&info.Path,
			&statusStr,
			&lastModified,
			&info.Size,
			&info.IsDownloaded,
			&info.IsDirectory,
			&info.Version,
			&info.Checksum,
			&lastSynced,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}

		info.Status = models.SyncStatus(statusStr)
		info.LastModified = time.Unix(lastModified, 0)
		info.LastSynced = time.Unix(lastSynced, 0)

		files = append(files, &info)
	}

	return files, nil
}

// Initialize database schema
func initDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			path TEXT PRIMARY KEY,
			status TEXT NOT NULL,
			last_modified INTEGER NOT NULL,
			size INTEGER NOT NULL,
			is_downloaded BOOLEAN NOT NULL,
			is_directory BOOLEAN NOT NULL,
			version INTEGER NOT NULL,
			checksum TEXT,
			last_synced INTEGER NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS idx_files_status ON files(status);
	`)
	
	return err
}