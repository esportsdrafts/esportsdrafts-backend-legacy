"""Email helper functions."""

from os import listdir
from os.path import isfile, join
from typing import List, Optional, Text


LOCAL_INBOX_PATH_OSX = '/Users/inbox'


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
