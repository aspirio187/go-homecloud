package main

import (
	"embed"

	"homecloud/internal/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed assets/icons/icon_win.ico
var iconData []byte

func main() {
	// Create a new app instance
	application := app.NewApp(iconData)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Home Cloud",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:         application.Startup,
		HideWindowOnClose: true,
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
