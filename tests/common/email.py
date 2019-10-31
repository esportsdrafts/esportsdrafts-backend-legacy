"""
Email helper functions.

TODO: Make some kind of Email class that hold these functions and extracts
timestamp, content etc.

TODO: Support getting email from cloud DB in GCP.
"""

import codecs
import os.path
import re
from os import listdir
from os.path import isfile, join
from typing import List, Optional, Text, Tuple
from urllib.parse import parse_qs, urlparse

URL_REGEX = r'https:\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^;=%&:\/~+#-]*[\w@?^;=%&\/~+#-])?'  # noqa

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


def read_local_email(email_name: Text) -> Optional[Text]:
    """
    Read an email from local inbox and return its content as string.

    Arguments:
    email_name -- name of the email(Not full path!)

    Returns:
    Content of email as string

    """
    email_path = LOCAL_INBOX_PATH_OSX + '/' + email_name
    if not os.path.isfile(email_path):
        # Raise instead? Not found
        return None

    parsed = codecs.open(email_path, 'r')
    return parsed.read()


def get_verification_token(welcome_email: Text) -> Tuple[Optional[Text],
                                                         Optional[Text]]:
    """
    Extract email verification token from welcome email.

    Arguments:
    welcome_email -- body of welcome email as string

    Returns:
    Tuple containing user ID and token as string; if not found both values are
    None

    """
    matches = re.finditer(URL_REGEX, welcome_email, re.MULTILINE)
    user_id, token = None, None

    for match in matches:
        full_match = match.group(0)
        # This is dumb
        if 'token' in full_match and 'user' in full_match:
            parsed = urlparse(full_match)
            query = parse_qs(parsed.query)
            token = query['token'][0]
            user_id = query['user'][0]
            break

    return user_id, token
