from app.window import get_all_windows
from app.window_mirror import WindowMirror

def main():
    windows = get_all_windows()
    windows = WindowMirror.filter_windows(windows)

    for window in windows:
        WindowMirror.mirror_window(window.window_handle)

if __name__ == '__main__':
    main()
