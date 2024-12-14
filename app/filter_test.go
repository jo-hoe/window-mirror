package app

import (
	"reflect"
	"syscall"
	"testing"

	"github.com/lxn/win"
)

func TestParallelWindowFilter(t *testing.T) {
	type args struct {
		windows   []WindowInfo
		predicate func(WindowInfo) bool
	}
	tests := []struct {
		name string
		args args
		want []WindowInfo
	}{
		{
			name: "Filter Out Items",
			args: args{
				windows: []WindowInfo{
					{Title: "test1"},
					{Title: "test2"},
					{Title: "test3"},
				},
				predicate: func(w WindowInfo) bool {
					return w.Title == "test2"
				},
			},
			want: []WindowInfo{
				{Title: "test2"},
			},
		},
		{
			name: "Empty Window List",
			args: args{
				windows:   []WindowInfo{},
				predicate: func(w WindowInfo) bool { return true },
			},
			want: []WindowInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParallelWindowFilter(tt.args.windows, tt.args.predicate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParallelWindowFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// MockUser32API is a mock implementation of the User32API interface.
type MockUser32API struct {
	visibleReturn bool
	iconicReturn  bool
}

func (m *MockUser32API) IsWindowVisible(handle uintptr) bool {
	return m.visibleReturn
}

func (m *MockUser32API) IsIconic(handle uintptr) bool {
	return m.iconicReturn
}

func (m *MockUser32API) GetAllWindows() []WindowInfo {
	return []WindowInfo{}
}

func (m *MockUser32API) GetWindowTitle(windowHandle syscall.Handle) string {
	return ""
}

func (m *MockUser32API) GetWindowRectangle(windowHandle syscall.Handle) win.RECT {
	return win.RECT{}
}

func TestIsWindowVisible(t *testing.T) {
	// Mock setup
	mockAPI := &MockUser32API{visibleReturn: true}
	user32Instance = mockAPI // Replace the real instance with the mock

	window := WindowInfo{Handle: 123}
	if !IsWindowVisible(user32Instance, window) {
		t.Errorf("Expected window to be visible")
	}
}

func TestIsIconic(t *testing.T) {
	// Mock setup
	mockAPI := &MockUser32API{iconicReturn: true}
	user32Instance = mockAPI // Replace the real instance with the mock

	window := WindowInfo{Handle: 456}

	if !IsIconic(user32Instance, window) {
		t.Errorf("Expected window to be iconic")
	}
}
