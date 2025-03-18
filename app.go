package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"homecloud/internal/sync"
)

// App struct
type App struct {
	ctx         context.Context
	syncManager *sync.SyncManager
	watchDir    string
}

// NewApp creates a new App application struct
func NewApp() *App {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home directory: %w", err))
	}
	watchDir := filepath.Join(homeDir, "homecloud")

	os.MkdirAll(watchDir, 0755)

	return &App{
		watchDir: watchDir,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.syncManager = sync.NewSyncManager(a.watchDir)
	err := a.syncManager.Start()
	if err != nil {
		fmt.Printf("failed to start sync manager: %v\n", err)
	}
}

func (a *App) GetWatchDir() string {
	return a.watchDir
}

func (a *App) setWatchDir(dir string) error {
	if a.syncManager != nil {
		a.syncManager.Stop()
	}

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	a.watchDir = dir

	a.syncManager = sync.NewSyncManager(a.watchDir)
	return a.syncManager.Start()
}

func (a *App) GetFiles() []sync.FileInfo {
	if a.syncManager == nil {
		return []sync.FileInfo{}
	}

	fileInfos := a.syncManager.GetFileInfos()

	result := make([]sync.FileInfo, len(fileInfos))
	for i, info := range fileInfos {
		result[i] = *info
	}

	return result
}

func (a *App) MinimizeToTray() {
	// This will be called from the frontend to minimize to tray
	// We'll implement this in tray.go
	// For now, just hide the window
	if a.ctx != nil {
		// Import "github.com/wailsapp/wails/v2/pkg/runtime" at the top
		// and uncomment this line when adding tray support
		// runtime.WindowHide(a.ctx)
	}
}
