
import uuid
from typing import Text

import pytest
from tests.common.user import User, create_new_account

env_urls = {
    'dev': 'dev.int.efantasy.com',
    'stage': 'stage.int.efantasy.com',
    'prod': 'efantasy.com',
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


@pytest.fixture()
def test_username() -> Text:
    return 'test_user_' + str(uuid.uuid4())


@pytest.fixture()
def test_email() -> Text:
    return 'test_user_' + str(uuid.uuid4()) + '@test.nu'


@pytest.fixture()
def test_password() -> Text:
    return str(uuid.uuid4())


@pytest.fixture()
def user(api_env_url: Text,
         test_username: Text,
         test_password: Text,
         test_email: Text) -> User:
    new_user = create_new_account(
        username=test_username,
        email=test_email,
        password=test_password,
        url=api_env_url
    )
    return new_user
