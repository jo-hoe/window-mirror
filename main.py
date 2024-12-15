from app.window import Window, get_all_windows
from app.window_mirror import WindowMirror


if __name__ == '__main__':
    windows = get_all_windows()
    windows = WindowMirror.filter_windows(windows)

    for window in windows:
        WindowMirror.mirror_window(window.window_handle)
