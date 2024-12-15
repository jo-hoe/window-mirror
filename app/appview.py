from functools import cache
from pyvda import get_apps_by_z_order, AppView


@cache
@staticmethod
def get_app_view_for_window(window_handle: int) -> AppView:
    """
    Identifies the app view object for a given window

    Returns:
        AppView: AppView object for the given window, None otherwise
    """
    apps = get_apps_by_z_order()
    for app in apps:
        if app.hwnd == window_handle:
            return app

    return None
