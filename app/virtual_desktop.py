from pyvda import get_apps_by_z_order


class VirtualDesktopManager:

    @staticmethod
    def is_window_on_current_desktop(window_handle:int) -> bool:
        """
        Identifies if the windows is on the currently active desktop

        Returns:
            bool: True if the window is on the currently active desktop, False otherwise
        """
        apps = get_apps_by_z_order()
        for app in apps:
            if app.hwnd == window_handle:
                return app.is_on_current_desktop()
            
        return False

