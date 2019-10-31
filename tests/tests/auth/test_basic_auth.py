
import time

import requests
from tests.common.email import get_emails_from_local_inbox, read_local_email
from tests.common.user import create_new_account
from tests.common.utils import gen_random_chars


def test_create_account(user):
    pass


def __check_fails(fn):
    try:
        fn()
        assert False
    except requests.HTTPError:
        pass


def test_username_validation():
    # Check too long username
    __check_fails(lambda: create_new_account(
        gen_random_chars(100),
        gen_random_chars(10) + '@test.nu',
        gen_random_chars(30)
    ))

    # Test max length for password
    __check_fails(lambda: create_new_account(
        '[][]-()2312dsaz...>>>//',
        gen_random_chars(10) + '@test.nu',
        gen_random_chars(30)
    ))

    # Test profanity filter doing basic validation
    __check_fails(lambda: create_new_account(
        'penis',
        gen_random_chars(10) + '@test.nu',
        gen_random_chars(30)
    ))


def test_password_validation():
    # Test min length for password
    __check_fails(lambda: create_new_account(
        gen_random_chars(10),
        gen_random_chars(10) + '@test.nu',
        gen_random_chars(3)))
    # Test max length for password
    __check_fails(lambda: create_new_account(
        gen_random_chars(10),
        gen_random_chars(10) + '@test.nu',
        'a' * 500))


def test_username_duplication_validation(user):
    __check_fails(lambda: create_new_account(
        user.username,
        gen_random_chars(10) + '@test.nu',
        gen_random_chars(12)))


def test_email_validation(user):
    __check_fails(lambda: create_new_account(
        gen_random_chars(20),
        user.email,
        gen_random_chars(12)))


def test_verification_email_sent(user, env):
    # When doing local testing sent emails are stored on disk in /tmp/inbox\
    sent_time = int(time.time())
    time_epsilon = 15

    # Encapsulate in an Email class
    if env == 'local':
        # Initial wait for email to get sent
        time.sleep(1)

        for _ in range(3):
            emails = get_emails_from_local_inbox(user.username, 'welcome')

            # No emails with correct type and user
            if not emails:
                time.sleep(3)
                continue

            latest_email = emails[0]
            ts = int(latest_email.split('_')[0])

            # Too old email to be considered
            if ts > (sent_time + time_epsilon) or \
                    ts < (sent_time - time_epsilon):
                time.sleep(3)
                continue

            parsed = read_local_email(latest_email)

            if user.username in parsed:
                return

            # Wait a bit for results to come in
            time.sleep(3)

        assert False, 'No welcome email found with correct content'
    else:
        # TODO: Check cloud database for email
        # each 'test' email is copied to a DB in GCP
        pass
