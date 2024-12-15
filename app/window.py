import win32gui
import win32con

from app.appview import get_app_view_for_window


class Window:

    def __init__(self, window_handle: int, window_title: str):
        self._window_handle = window_handle
        self._window_title = window_title

    @property
    def window_handle(self) -> int:
        return self._window_handle

    @property
    def window_title(self) -> str:
        return self._window_title

    def __repr__(self):
        return f"Window(handle={self.window_handle}, title={self.window_title})"


def get_all_windows() -> list["Window"]:  # Use string literal here
    """
    Get list of window handles on the current virtual desktop.

    Returns:
        List[int]: List of window handles
    """
    visible_windows = []

    def enum_windows_callback(window_handle, _):
        visible_windows.append(window_handle)
        return True
    win32gui.EnumWindows(enum_windows_callback, None)

    windows = []
    for handle in visible_windows:
        title = get_window_title(handle)
        windows.append(Window(handle, title))
    return windows


def is_iconic(window_handle: int) -> bool:
    return win32gui.IsIconic(window_handle)


def is_visible(window_handle: int) -> bool:
    return win32gui.IsWindowVisible(window_handle)


def get_window_title(window_handle: int) -> str:
    return win32gui.GetWindowText(window_handle)


def is_window_pinned(window_handle: int) -> bool:
    app_view = get_app_view_for_window(window_handle)
    return app_view.is_pinned()

def pin_window(window_handle: int) -> None:
    app_view = get_app_view_for_window(window_handle)
    return app_view.pin()


def get_window_rectangle(window_handle: int) -> tuple[int, int, int, int]:
    """
    Get the window rectangle.

    Args:
        window_handle (int): Window handle

    Returns:
        tuple[int, int, int, int]: (left, top, right, bottom) window coordinates
    """
    return win32gui.GetWindowRect(window_handle)


def move_window(window_handle: int, left: int, top: int, width: int, height: int) -> None:
    # Move window
    win32gui.SetWindowPos(
        window_handle,
        win32con.HWND_TOP,
        left,
        top,
        width,
        height,
        win32con.SWP_SHOWWINDOW
    )
