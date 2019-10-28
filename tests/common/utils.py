"""Utility HTTP and API classes and functions."""

from os import listdir
from os.path import isfile, join
from random import choice
from string import ascii_lowercase
from typing import List, Optional, Text

import requests

LOCAL_INBOX_PATH_OSX = '/Users/inbox'


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


def get_local_inbox_path(platform: Text) -> Text:
    """
    Return base path to local email inbox given a platform.

    Arguments:
    platform -- platform as lowercase string

    Returns:
    Path to inbox as string

    """
    platform = platform.lower()

    if platform == 'osx' or platform == 'mac':
        return LOCAL_INBOX_PATH_OSX
    if platform == 'linux':
        raise Exception('Linux local inbox not supported')

    raise Exception('Unknown platform')


def __get_local_email_inbox() -> List[Text]:
    email_paths = [f for f in listdir(LOCAL_INBOX_PATH_OSX)
                   if isfile(join(LOCAL_INBOX_PATH_OSX, f))]
    return sorted(email_paths)


def get_emails_from_local_inbox(
        username: Text, email_type: Optional[Text] = None) -> List[Text]:
    """
    Grabs email for a user and optionally type from 'local' email inbox.

    To get the full path to the email join with LOCAL_INBOX_PATH.

    Arguments:
    username -- username of a user as string
    email_type -- type of email to filter on, if none provided return all
        email types (default None)

    Returns:
    List of email file names in local storage inbox

    """
    emails = __get_local_email_inbox()

    user_emails = [e for e in emails if username in e]

    if email_type is not None:
        user_emails = [e for e in user_emails if email_type in e]

    return user_emails
