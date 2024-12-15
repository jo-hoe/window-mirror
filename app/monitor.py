import win32con
import win32api


class MonitorManager:

    @staticmethod
    def get_window_monitor_info(window_handle: int) -> tuple[int, int, int, int]:
        """
        Get the monitor information for a given window.

        Args:
            window_handle (int): Window handle

        Returns:
            Tuple[int, int, int, int]: (left, top, right, bottom) monitor coordinates
        """
        monitor = win32api.MonitorFromWindow(
            window_handle, win32con.MONITOR_DEFAULTTONEAREST)
        monitor_info = win32api.GetMonitorInfo(monitor)
        return monitor_info['Monitor']
