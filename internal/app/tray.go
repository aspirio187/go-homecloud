package app

import (
	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SetupSystemTray initializes and starts the system tray
func (a *App) SetupSystemTray() {
	go func() {
		systray.Run(func() {
			a.onSystrayReady()
		}, func() {
			a.onSystrayExit()
		})
	}()
}

func (a *App) onSystrayReady() {
	systray.SetIcon(a.iconData)
	systray.SetTitle("Home Cloud")
	systray.SetTooltip("Home Cloud - Your personal cloud")

	mOpen := systray.AddMenuItem("Open HomeCloud", "Open HomeCloud")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Exit", "Exit HomeCloud")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				runtime.WindowShow(a.ctx)
			case <-mQuit.ClickedCh:
				runtime.Quit(a.ctx)
				return
			}
		}
	}()
}

func (a *App) onSystrayExit() {
	runtime.WindowHide(a.ctx)
}