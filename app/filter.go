package app

import (
	"sync"
)

// ParallelFilter applies a single filter function (predicate) to filter items in a slice in parallel.
func ParallelWindowFilter(windows []WindowInfo, predicate func(WindowInfo) bool) []WindowInfo {
	if len(windows) == 0 {
		return []WindowInfo{}
	}

	var wg sync.WaitGroup
	ch := make(chan WindowInfo, len(windows)) // Buffered channel for results

	for _, windowInfo := range windows {
		wg.Add(1)
		// Process each item in parallel
		go func(p WindowInfo) {
			defer wg.Done()
			if predicate(p) {
				ch <- p // Send to channel if the predicate passes
			}
		}(windowInfo)
	}

	// Close the channel after all goroutines complete
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results from the channel
	var filtered []WindowInfo
	for p := range ch {
		filtered = append(filtered, p)
	}

	return filtered
}

// Check if a window is visible.
func IsWindowVisible(user32Api User32API, window WindowInfo) bool {
	return user32Api.IsWindowVisible(uintptr(window.Handle))
}

// Checks is a window is minimized.
func IsIconic(user32Api User32API, window WindowInfo) bool {
	return user32Api.IsIconic(uintptr(window.Handle))
}
