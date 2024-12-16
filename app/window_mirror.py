from app.monitor import get_window_monitor_info
from app.virtual_desktop import is_window_on_current_desktop
from app.window import Window, get_window_rectangle, is_iconic, is_visible, move_window


class WindowMirror:

    @staticmethod
    def filter_windows(windows: Window):
        result = []

        for window in windows:
            if window.window_title == "":
                continue
            if is_iconic(window.window_handle):
                continue
            if not is_visible(window.window_handle):
                continue
            if not is_window_on_current_desktop(window.window_handle):
                continue

            result.append(window)

        return result

    @staticmethod
    def mirror_window(window_handle):
        """
        Mirror a single window horizontally.

        Args:
            window_handle (int): Window handle
        """
        # Get current window rectangle
        left, top, right, bottom = get_window_rectangle(window_handle)

        # Get monitor info
        monitor_left, _, monitor_right, _ = \
            get_window_monitor_info(window_handle)

        # Calculate new position
        width = right - left
        height = bottom - top
        new_left = monitor_right - (left - monitor_left) - width
        new_top = top

        move_window(window_handle, new_left, new_top, width, height)
