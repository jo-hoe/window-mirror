from app.monitor import MonitorManager
from app.virtual_desktop import VirtualDesktopManager
from app.window import Window


class WindowMirror:

    @staticmethod
    def filter_windows(windows: Window):
        result = []

        for window in windows:
            if window.window_title == "":
                continue
            if Window.is_iconic(window.window_handle):
                continue
            if not Window.is_visible(window.window_handle):
                continue
            if not VirtualDesktopManager.is_window_on_current_desktop(window.window_handle):
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
        left, top, right, bottom = Window.get_window_rectangle(window_handle)

        # Get monitor info
        monitor_left, monitor_top, monitor_right, monitor_bottom = \
            MonitorManager.get_window_monitor_info(window_handle)

        # Calculate new position
        width = right - left
        height = bottom - top
        new_left = monitor_right - (left - monitor_left) - width
        new_top = top

        Window.move_window(window_handle, new_left, new_top, width, height)
