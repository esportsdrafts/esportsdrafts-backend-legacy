
import codecs
import time
from os import listdir
from os.path import isfile, join
from typing import List, Text

import requests
from tests.common.user import create_new_account
from tests.common.utils import gen_random_chars


def test_create_account(user):
    pass


def test_password_validation(user):
    try:
        # Test min length for password
        create_new_account(
            gen_random_chars(10),
            gen_random_chars(10) + '@test.nu',
            gen_random_chars(3))
        assert False

        # Test max length for password
        create_new_account(
            gen_random_chars(10),
            gen_random_chars(10) + '@test.nu',
            'a' * 500)
        assert False
    except requests.HTTPError:
        pass


def test_username_validation(user):
    try:
        create_new_account(
            user.username,
            gen_random_chars(10) + '@test.nu',
            gen_random_chars(12))
        assert False
    except requests.HTTPError:
        pass


def test_email_validation(user):
    try:
        create_new_account(
            gen_random_chars(20),
            user.email,
            gen_random_chars(12))
        assert False
    except requests.HTTPError:
        pass


def __get_local_email_inbox() -> List[Text]:
    inbox_path = '/Users/inbox'
    email_paths = [f for f in listdir(inbox_path)
                   if isfile(join(inbox_path, f))]
    return sorted(email_paths)


def test_verification_email_sent(user, env):
    # When doing local testing sent emails are stored on disk in /tmp/inbox\
    sent_time = int(time.time())
    time_epsilon = 30

    if env == 'local':
        # Initial wait for email to get sent
        time.sleep(1)

        for _ in range(3):
            emails = __get_local_email_inbox()
            user_emails = [
                e for e in emails if user.username in e and 'welcome' in e]

            # No emails with correct type and user
            if not user_emails:
                time.sleep(3)
                continue

            latest_email = user_emails[0]
            ts = int(latest_email.split('_')[0])

            # Too old email to be considered
            if ts > (sent_time + time_epsilon) or \
                    ts < (sent_time - time_epsilon):
                time.sleep(3)
                continue

            parsed = codecs.open('/Users/inbox/' + latest_email, 'r')
            parsed = parsed.read()

            if user.username in parsed:
                return

            # Wait a bit for results to come in
            time.sleep(3)

        assert False, "No welcome email found with correct content"
    else:
        # TODO: Check cloud database for email
        # each 'test' email is copied to a DB in GCP
        pass
