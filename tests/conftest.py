
from random import choice
from string import ascii_uppercase, ascii_lowercase
from typing import Text

import pytest
from tests.common.user import User, create_new_account

env_urls = {
    'dev': 'api.dev.int.efantasy.com',
    'stage': 'api.stage.int.efantasy.com',
    'prod': 'api.efantasy.com',
    'local': 'efantasy.localhost',
}


def pytest_addoption(parser):
    parser.addoption(
        '--env', action='store', default='local',
        choices=['local', 'dev', 'stage', 'prod'],
        help='Pick which environment to run integration tests against')


@pytest.fixture()
def env(pytestconfig):
    return pytestconfig.getoption('env')


@pytest.fixture()
def api_env_url(env: Text) -> Text:
    return 'api.' + env_urls[env]


@pytest.fixture()
def env_url(env: Text) -> Text:
    return env_urls[env]


def __gen_random_chars(n: int) -> Text:
    return ''.join(choice(ascii_uppercase + ascii_lowercase + '-_')
                   for i in range(n))


@pytest.fixture()
def test_username() -> Text:
    return 'test_user_' + __gen_random_chars(14)


@pytest.fixture()
def test_email() -> Text:
    return 'test_user_' + __gen_random_chars(14) + '@test.nu'


@pytest.fixture()
def test_password() -> Text:
    return __gen_random_chars(30)


@pytest.fixture()
def user(api_env_url: Text,
         test_username: Text,
         test_password: Text,
         test_email: Text) -> User:
    new_user = create_new_account(
        username=test_username,
        email=test_email,
        password=test_password,
        url=api_env_url,
    )
    new_user.login()
    return new_user
