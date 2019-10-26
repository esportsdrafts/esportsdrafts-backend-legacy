
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


def test_verification_email_sent(user, env):
    if env == 'local':
        pass
    else:
        # TODO: Check cloud database for email
        pass
