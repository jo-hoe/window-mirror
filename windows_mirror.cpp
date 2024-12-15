#include <windows.h>
#include <dwmapi.h>
#include <vector>
#include <iostream>
#include <algorithm> 
#include <ShObjIdl.h>  // For IVirtualDesktopManager

#pragma comment(lib, "Dwmapi.lib")
#pragma comment(lib, "Ole32.lib")  // For COM initialization

// Helper struct to represent monitor boundaries
struct MonitorInfo {
    RECT rect;
};
// Explicit definition of CLSID_VirtualDesktopManager
const CLSID CLSID_VirtualDesktopManager = {0xaa509086, 0x5ca9, 0x4c25, {0x8f, 0x95, 0x58, 0x9d, 0x3c, 0x07, 0x98, 0x62}};
// Global COM Pointer for Virtual Desktop Manager
IVirtualDesktopManager* virtualDesktopManager = nullptr;

// Function to enumerate monitors and store their boundaries
BOOL CALLBACK MonitorEnumProc(HMONITOR hMonitor, HDC hdcMonitor, LPRECT lprcMonitor, LPARAM dwData) {
    std::vector<MonitorInfo>* monitors = reinterpret_cast<std::vector<MonitorInfo>*>(dwData);
    MonitorInfo info;
    info.rect = *lprcMonitor;
    monitors->push_back(info);
    return TRUE;
}

// Function to calculate the bounds of the entire virtual desktop
RECT GetVirtualDesktopBounds(const std::vector<MonitorInfo>& monitors) {
    RECT desktopBounds = {LONG_MAX, LONG_MAX, LONG_MIN, LONG_MIN};
    for (const auto& monitor : monitors) {
        desktopBounds.left = std::min(desktopBounds.left, monitor.rect.left);
        desktopBounds.top = std::min(desktopBounds.top, monitor.rect.top);
        desktopBounds.right = std::max(desktopBounds.right, monitor.rect.right);
        desktopBounds.bottom = std::max(desktopBounds.bottom, monitor.rect.bottom);
    }
    return desktopBounds;
}

// Function to mirror a window across the virtual desktop
void MirrorWindow(HWND hwnd, const RECT& desktopBounds) {
    RECT windowRect;
    if (!GetWindowRect(hwnd, &windowRect)) return;

    // Calculate the window's mirrored position
    int windowWidth = windowRect.right - windowRect.left;
    int mirroredLeft = desktopBounds.right - (windowRect.left - desktopBounds.left) - windowWidth;
    int mirroredTop = windowRect.top;  // Vertical position remains unchanged

    // Move the window to its mirrored position
    SetWindowPos(hwnd, nullptr, mirroredLeft, mirroredTop, 0, 0, SWP_NOZORDER | SWP_NOSIZE | SWP_NOACTIVATE);
}

// Function to check if a window is on the active virtual desktop
bool IsWindowOnActiveVirtualDesktop(HWND hwnd) {
    if (!virtualDesktopManager) return true;  // If COM interface fails, assume the window is on the active desktop
    BOOL onCurrentDesktop = FALSE;
    if (SUCCEEDED(virtualDesktopManager->IsWindowOnCurrentVirtualDesktop(hwnd, &onCurrentDesktop))) {
        return onCurrentDesktop == TRUE;
    }
    return true;
}

// Callback to process each window
BOOL CALLBACK EnumWindowsProc(HWND hwnd, LPARAM lParam) {
    // Check if the window is visible
    if (!IsWindowVisible(hwnd)) return TRUE;

    // Get the extended window styles and skip tooltips, etc.
    LONG style = GetWindowLong(hwnd, GWL_EXSTYLE);
    if (style & WS_EX_TOOLWINDOW) return TRUE;  // Skip tool windows

    // Check if the window is on the active virtual desktop
    if (virtualDesktopManager) {
        BOOL onCurrentDesktop = FALSE;
        if (FAILED(virtualDesktopManager->IsWindowOnCurrentVirtualDesktop(hwnd, &onCurrentDesktop))) {
            std::cerr << "Failed to determine if window is on current virtual desktop.\n";
            return TRUE;  // If check fails, skip this window
        }
        if (!onCurrentDesktop) return TRUE;  // Skip windows not on the active virtual desktop
    }

    // Get desktop bounds
    RECT* desktopBounds = reinterpret_cast<RECT*>(lParam);

    // Mirror the window across the virtual desktop
    MirrorWindow(hwnd, *desktopBounds);

    return TRUE;
}

int main() {
    // Initialize COM for Virtual Desktop Manager
    if (FAILED(CoInitialize(nullptr))) {
        std::cerr << "Failed to initialize COM.\n";
        return -1;
    }

    if (FAILED(CoCreateInstance(CLSID_VirtualDesktopManager, nullptr, CLSCTX_INPROC_SERVER, IID_PPV_ARGS(&virtualDesktopManager)))) {
        std::cerr << "Failed to initialize IVirtualDesktopManager. Continuing without virtual desktop support.\n";
    }

    // Enumerate all monitors
    std::vector<MonitorInfo> monitors;
    EnumDisplayMonitors(nullptr, nullptr, MonitorEnumProc, reinterpret_cast<LPARAM>(&monitors));

    // Calculate the virtual desktop bounds
    RECT desktopBounds = GetVirtualDesktopBounds(monitors);

    // Enumerate all windows and process them
    EnumWindows(EnumWindowsProc, reinterpret_cast<LPARAM>(&desktopBounds));

    // Cleanup COM
    if (virtualDesktopManager) virtualDesktopManager->Release();
    CoUninitialize();

    return 0;
}
