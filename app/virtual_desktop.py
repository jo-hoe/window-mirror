from pyvda import get_apps_by_z_order

from app.appview import get_app_view_for_window


@staticmethod
def is_window_on_current_desktop(window_handle: int) -> bool:
    """
    Identifies if the windows is on the currently active desktop

    Returns:
        bool: True if the window is on the currently active desktop, False otherwise
    """

    apps_view = get_app_view_for_window(window_handle)

    if apps_view is None:
        return False

    return apps_view.is_on_current_desktop()
