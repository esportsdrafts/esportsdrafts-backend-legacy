"""Utility HTTP and API classes and functions."""

import requests


def raise_on_error(request: requests.Response) -> None:
    """
    Check if the response contains a 4xx or 5xx return code.

    If it does contain one of those error codes, raise an exception and let the
    user deal with it. Otherwise returns nothing.

    Raises:
    HTTPError -- If request contains a 4xx or 5xx HTTP code

    """
    if request.status_code >= 400:
        json_res = request.json()
        raise requests.HTTPError(json_res)

    return None
