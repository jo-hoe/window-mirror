package main

import (
	"log"
	"syscall"

	"github.com/jo-hoe/window-mirror/app"
	"github.com/lxn/win"
)

var (
	ole32              = syscall.MustLoadDLL("ole32.dll")
	procCoInitializeEx = ole32.MustFindProc("CoInitializeEx")
	procCoUninitialize = ole32.MustFindProc("CoUninitialize")

	// COM initialization flags
	COINIT_APARTMENTTHREADED = 0x2
	COINIT_DISABLE_OLE1DDE   = 0x4
)

func main() {
	// Initialize COM library
	hr, _, _ := procCoInitializeEx.Call(
		0,
		uintptr(COINIT_APARTMENTTHREADED|COINIT_DISABLE_OLE1DDE),
	)

	// 0 indicates S_OK (success)
	if hr != 0 && hr != 0x00000001 { // S_OK or S_FALSE
		log.Fatalf("CoInitializeEx failed with code: %x", hr)
	}
	defer procCoUninitialize.Call()

	// Detect all visible windows with their monitor information
	user32Api := app.GetUser32Instance()
	windows := user32Api.GetAllWindows()

	// remove windows without title
	windows = app.ParallelWindowFilter(windows, func(window app.WindowInfo) bool {
		return !(window.Title == "")
	})
	// remove minimized windows
	windows = app.ParallelWindowFilter(windows, func(window app.WindowInfo) bool {
		return !app.IsIconic(user32Api, window)
	})
	// remove invisible windows
	windows = app.ParallelWindowFilter(windows, func(window app.WindowInfo) bool {
		return app.IsWindowVisible(user32Api, window)
	})

	// Mirror windows on their respective monitors
	mirrorWindows(user32Api, windows)
}

func mirrorWindows(user32Api app.User32API, windows []app.WindowInfo) {
	for _, window := range windows {
		// Calculate new position within the monitor's work area
		newRect := calculateMirroredPosition(window)

		// Move window
		user32Api.MoveWindows(window,
			newRect.Left,
			newRect.Top,
			newRect.Right-newRect.Left,
			newRect.Bottom-newRect.Top)
	}
}

func calculateMirroredPosition(window app.WindowInfo) win.RECT {
	// Get monitor work area
	monitorRect := window.MonitorInfo.RcWork

	windowWidth := window.Rectangle.Right - window.Rectangle.Left
	windowHeight := window.Rectangle.Bottom - window.Rectangle.Top

	// Calculate window's current horizontal position relative to monitor width
	windowCenter := (window.Rectangle.Left + window.Rectangle.Right) / 2
	monitorWidth := monitorRect.Right - monitorRect.Left

	// Normalize window position within the monitor
	localWindowCenter := windowCenter - monitorRect.Left
	horizontalRatio := float64(localWindowCenter) / float64(monitorWidth)

	var newLeft int32
	if horizontalRatio < 0.5 {
		// Left side windows move to right
		newLeft = monitorRect.Right - (window.Rectangle.Right - monitorRect.Left)
	} else {
		// Right side windows move to left
		newLeft = monitorRect.Left + (monitorRect.Right - window.Rectangle.Left)
		newLeft -= windowWidth
	}

	return win.RECT{
		Left:   newLeft,
		Top:    window.Rectangle.Top,
		Right:  newLeft + windowWidth,
		Bottom: window.Rectangle.Top + windowHeight,
	}
}

func init() {
	// Error handling and logging setup
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
