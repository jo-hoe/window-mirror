package app

import (
	"sync"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

type WindowInfo struct {
	Handle      syscall.Handle
	Title       string
	Rectangle   win.RECT
	MonitorInfo win.MONITORINFO
}

// User32API defines the interface for interacting with user32.dll.
type User32API interface {
	IsWindowVisible(handle uintptr) bool
	IsIconic(handle uintptr) bool
	GetAllWindows() []WindowInfo
	GetWindowTitle(windowHandle syscall.Handle) string
	GetWindowRectangle(windowHandle syscall.Handle) win.RECT
}

// User32DLL implements the User32API interface using user32.dll.
type User32DLL struct {
	isWindowVisible      *syscall.Proc
	isIconic             *syscall.Proc
	enumWindows          *syscall.Proc
	getWindowTitle       *syscall.Proc
	getWindowTitleLength *syscall.Proc
}

var (
	user32Instance User32API
	once           sync.Once
)

// getUser32Instance initializes and returns a singleton instance of User32API.
func GetUser32Instance() User32API {
	once.Do(func() {
		dll := syscall.MustLoadDLL("user32.dll")
		user32Instance = &User32DLL{
			isWindowVisible:      dll.MustFindProc("IsWindowVisible"),
			isIconic:             dll.MustFindProc("IsIconic"),
			enumWindows:          dll.MustFindProc("EnumWindows"),
			getWindowTitle:       dll.MustFindProc("GetWindowTextW"),
			getWindowTitleLength: dll.MustFindProc("GetWindowTextLengthW"),
		}
	})
	return user32Instance
}

// IsWindowVisible checks if a window is visible using the real implementation.
func (u *User32DLL) IsWindowVisible(handle uintptr) bool {
	visibleRet, _, _ := u.isWindowVisible.Call(handle)
	return visibleRet == 1
}

// IsIconic checks if a window is minimized (iconic) using the real implementation.
func (u *User32DLL) IsIconic(handle uintptr) bool {
	iconicRet, _, _ := u.isIconic.Call(handle)
	return iconicRet != 0
}

func (u *User32DLL) GetWindowTitle(windowHandle syscall.Handle) string {
	titleLength, _, _ := u.getWindowTitleLength.Call(uintptr(windowHandle))
	if titleLength < 1 {
		return ""
	}

	buffer := make([]uint16, titleLength)

	u.getWindowTitle.Call(
		uintptr(windowHandle),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
	)

	return syscall.UTF16ToString(buffer)
}

func (u *User32DLL) GetWindowRectangle(windowHandle syscall.Handle) win.RECT {
	var rect win.RECT
	win.GetWindowRect(win.HWND(windowHandle), &rect)
	return rect
}

func (u *User32DLL) GetAllWindows() []WindowInfo {
	var windowList []WindowInfo

	callback := syscall.NewCallback(func(windowHandle syscall.Handle, _ uintptr) uintptr {
		windowList = append(windowList, WindowInfo{
			Handle:      windowHandle,
			Title:       u.GetWindowTitle(windowHandle),
			Rectangle:   u.GetWindowRectangle(windowHandle),
			MonitorInfo: getMonitorInfo(windowHandle),
		})

		return 1
	})

	// Enumerate windows
	u.enumWindows.Call(callback, 0)

	return windowList
}

func getMonitorInfo(windowHandle syscall.Handle) win.MONITORINFO {
	// Find the monitor that contains the window
	hMonitor := win.MonitorFromWindow(win.HWND(windowHandle), win.MONITOR_DEFAULTTONEAREST)

	// Create a MONITORINFO struct
	var monitorInfo win.MONITORINFO
	monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))

	// Retrieve monitor information
	if win.GetMonitorInfo(hMonitor, &monitorInfo) {
		return monitorInfo
	}

	// Return an empty/zero MONITORINFO if retrieval fails
	return win.MONITORINFO{}
}
