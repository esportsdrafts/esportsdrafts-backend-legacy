"""Utility HTTP and API classes and functions."""

from typing import Text

from random import choice
from string import ascii_lowercase, ascii_uppercase

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


def gen_random_chars(n: int = 10) -> Text:
    """
    Generate a string of random lower ascii characters.

    Arguments:
    n -- number of characters to generate (default 10)

    Returns:
    A string with random characters

    """
    if n < 1:
        raise Exception('Number of random chars to generate has to be > 0')

    return ''.join(choice(ascii_lowercase + '-_')
                   for i in range(n))
