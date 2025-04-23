package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"homecloud/internal/config"
	"homecloud/internal/models"
	"homecloud/internal/server"
	"homecloud/internal/storage"
	"homecloud/internal/sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	syncManager   *sync.SyncManager
	configManager *config.Config
	serverClient  *server.Client
	metadataStore *storage.MetadataStore
	appDataPath   string
	iconData      []byte
	isConnected   bool
}

// NewApp creates a new App application struct
func NewApp(iconData []byte) *App {
	// Get application data directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home directory: %w", err))
	}

	appDataPath := filepath.Join(homeDir, ".homecloud")
	os.MkdirAll(appDataPath, 0755)

	// Create default watch directory
	watchDir := filepath.Join(homeDir, "homecloud")
	os.MkdirAll(watchDir, 0755)

	// Load configuration
	configPath := filepath.Join(appDataPath, "config.json")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// If loading fails, use default config
		cfg = config.DefaultConfig()
	}

	// Create a server client
	client := server.NewClient(cfg.ServerURL)

	return &App{
		appDataPath:  appDataPath,
		configManager: cfg,
		serverClient: client,
		iconData:    iconData,
		isConnected: false,
	}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize metadata store
	metadataStore, err := storage.NewMetadataStore(a.appDataPath)
	if err != nil {
		fmt.Printf("failed to create metadata store: %v\n", err)
	} else {
		a.metadataStore = metadataStore
	}

	// Initialize sync manager with default watch directory
	a.syncManager = sync.NewSyncManager(a.configManager.WatchDirs[0])
	err = a.syncManager.Start()
	if err != nil {
		fmt.Printf("failed to start sync manager: %v\n", err)
	}

	// Try to load auth info and authenticate with the server
	authInfo, err := server.LoadAuth(a.appDataPath)
	if err == nil && authInfo != nil {
		// Auto-authenticate with saved token
		a.autoConnect(authInfo)
	}

	// Start the system tray
	a.SetupSystemTray()
}

// autoConnect attempts to connect to the server using saved auth info
func (a *App) autoConnect(authInfo *server.AuthInfo) {
	// We would implement the logic to automatically connect to the server
	// using saved authentication information
	// For now, just update the connection status
	a.isConnected = true
}

// Connect authenticates with the server
func (a *App) Connect(username, password string) error {
	err := a.serverClient.Authenticate(username, password)
	if err != nil {
		return err
	}

	// Save auth info for future auto-connection
	authInfo := &server.AuthInfo{
		Token:     "dummy-token", // In real implementation, get from server client
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Username:  username,
	}

	err = server.SaveAuth(a.appDataPath, authInfo)
	if err != nil {
		return err
	}

	a.isConnected = true
	return nil
}

// GetWatchDir returns the current watch directory
func (a *App) GetWatchDir() string {
	if len(a.configManager.WatchDirs) == 0 {
		return ""
	}
	return a.configManager.WatchDirs[0]
}

// SetWatchDir changes the watch directory and restarts the sync manager
func (a *App) SetWatchDir(dir string) error {
	if a.syncManager != nil {
		a.syncManager.Stop()
	}

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Update the first watch directory in config
	if len(a.configManager.WatchDirs) == 0 {
		a.configManager.WatchDirs = []string{dir}
	} else {
		a.configManager.WatchDirs[0] = dir
	}

	// Save the updated configuration
	configPath := filepath.Join(a.appDataPath, "config.json")
	if err := config.SaveConfig(configPath, a.configManager); err != nil {
		return err
	}

	a.syncManager = sync.NewSyncManager(dir)
	return a.syncManager.Start()
}

// GetFiles returns the list of files being tracked
func (a *App) GetFiles() []models.FileInfo {
	if a.syncManager == nil {
		return []models.FileInfo{}
	}

	fileInfos := a.syncManager.GetFileInfos()

	result := make([]models.FileInfo, len(fileInfos))
	for i, info := range fileInfos {
		result[i] = *info
	}

	return result
}

// MinimizeToTray minimizes the application to system tray
func (a *App) MinimizeToTray() {
	// This will be called from the frontend to minimize to tray
	if a.ctx != nil {
		runtime.WindowHide(a.ctx)
	}
}

// IsConnected returns the server connection status
func (a *App) IsConnected() bool {
	return a.isConnected
}

// Shutdown is called when the application is closing
func (a *App) Shutdown() {
	if a.metadataStore != nil {
		a.metadataStore.Close()
	}

	if a.syncManager != nil {
		a.syncManager.Stop()
	}
}