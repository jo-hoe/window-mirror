from app.window import Window
from app.window_mirror import WindowMirror


if __name__ == '__main__':
    windows = Window.get_all_windows()
    windows = WindowMirror.filter_windows(windows)

    for window in windows:
        WindowMirror.mirror_window(window)
