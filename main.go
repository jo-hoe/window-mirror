package main

import (
	"log"
	"math"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

type WindowInfo struct {
	Handle       syscall.Handle
	Title        string
	OriginalRect win.RECT
	ScreenIndex  int
}

var (
	ole32                = syscall.MustLoadDLL("ole32.dll")
	procCoInitializeEx   = ole32.MustFindProc("CoInitializeEx")
	procCoUninitialize   = ole32.MustFindProc("CoUninitialize")
	
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

	// Detect all visible windows
	windows := detectVisibleWindows()

	// Mirror windows
	mirrorWindows(windows)
}

func detectVisibleWindows() []WindowInfo {
	var windowList []WindowInfo

	// Load user32.dll
	user32 := syscall.MustLoadDLL("user32.dll")
	enumWindows := user32.MustFindProc("EnumWindows")
	getWindowTextW := user32.MustFindProc("GetWindowTextW")
	getWindowTextLengthW := user32.MustFindProc("GetWindowTextLengthW")
	isWindowVisible := user32.MustFindProc("IsWindowVisible")
	isIconic := user32.MustFindProc("IsIconic")

	// Callback function to enumerate windows
	callback := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		// Check if window is visible and not minimized
		visibleRet, _, _ := isWindowVisible.Call(uintptr(hwnd))
		iconicRet, _, _ := isIconic.Call(uintptr(hwnd))
		
		if visibleRet == 0 || iconicRet != 0 {
			return 1 // continue enumeration
		}

		// Get window title length
		titleLen, _, _ := getWindowTextLengthW.Call(uintptr(hwnd))
		
		// Allocate buffer for title
		if titleLen > 0 {
			buffer := make([]uint16, titleLen+1)
			getWindowTextW.Call(
				uintptr(hwnd), 
				uintptr(unsafe.Pointer(&buffer[0])), 
				uintptr(len(buffer)),
			)

			// Get window rectangle
			var rect win.RECT
			win.GetWindowRect(win.HWND(hwnd), &rect)

			// Additional check to ensure window is on screen and not zero-sized
			if rect.Right > rect.Left && rect.Bottom > rect.Top {
				windowList = append(windowList, WindowInfo{
					Handle:       hwnd,
					Title:        syscall.UTF16ToString(buffer),
					OriginalRect: rect,
					ScreenIndex:  0, // Placeholder, implement screen detection
				})
			}
		}

		return 1
	})

	// Enumerate windows
	enumWindows.Call(callback, 0)

	return windowList
}

func getTotalScreenWidth() int {
	// Get number of monitors
	monitorCount := win.GetSystemMetrics(win.SM_CMONITORS)
	totalWidth := 0
	
	// Get width of each monitor
	for i := 0; i < int(monitorCount); i++ {
		width := win.GetSystemMetrics(win.SM_CXSCREEN)
		totalWidth += int(width)
	}
	
	return totalWidth
}

func mirrorWindows(windows []WindowInfo) {
	// Get total screen width
	totalScreenWidth := getTotalScreenWidth()

	// Load user32.dll
	user32 := syscall.MustLoadDLL("user32.dll")
	moveWindow := user32.MustFindProc("MoveWindow")

	for _, window := range windows {
		// Calculate new position
		newRect := calculateMirroredPosition(window, totalScreenWidth)

		// Move window
		moveWindow.Call(
			uintptr(window.Handle), 
			uintptr(newRect.Left), 
			uintptr(newRect.Top), 
			uintptr(newRect.Right-newRect.Left), 
			uintptr(newRect.Bottom-newRect.Top), 
			uintptr(1), // redraw flag
		)
	}
}

func calculateMirroredPosition(window WindowInfo, totalScreenWidth int) win.RECT {
	windowWidth := window.OriginalRect.Right - window.OriginalRect.Left
	windowHeight := window.OriginalRect.Bottom - window.OriginalRect.Top

	// Calculate window's current horizontal position relative to total screen width
	windowCenter := (window.OriginalRect.Left + window.OriginalRect.Right) / 2
	horizontalRatio := float64(windowCenter) / float64(totalScreenWidth)

	var newLeft int32
	if math.Abs(horizontalRatio-0.5) < 0.1 {
		// If window is near center, keep its position
		newLeft = window.OriginalRect.Left
	} else if horizontalRatio < 0.5 {
		// Left side windows move to right
		newLeft = int32(totalScreenWidth) - window.OriginalRect.Right
	} else {
		// Right side windows move to left
		newLeft = int32(totalScreenWidth) - window.OriginalRect.Left - windowWidth
	}

	return win.RECT{
		Left:   newLeft,
		Top:    window.OriginalRect.Top,
		Right:  newLeft + windowWidth,
		Bottom: window.OriginalRect.Top + windowHeight,
	}
}

func init() {
	// Error handling and logging setup
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}