
import glob
import os
import platform as pl
from typing import Text

import pytest
from tests.common.email import get_local_inbox_path
from tests.common.user import User, create_new_account
from tests.common.utils import gen_random_chars

env_urls = {
    'local': 'esportsdrafts.localhost',
    'staging': 'esportsdrafts.dev',
    'production': 'esportsdrafts.app',
}


def pytest_addoption(parser):
    parser.addoption(
        '--env', action='store', default='local',
        choices=['local', 'staging', 'production'],
        help='Pick which environment to run integration tests against')


@pytest.fixture(scope='session')
def env(pytestconfig):
    return pytestconfig.getoption('env')


@pytest.fixture(scope='session')
def api_env_url(env: Text) -> Text:
    return 'api.' + env_urls[env]


@pytest.fixture(scope='session')
def env_url(env: Text) -> Text:
    return env_urls[env]


@pytest.fixture(scope='session')
def test_username() -> Text:
    return 'test_user_' + gen_random_chars(14)


@pytest.fixture(scope='session')
def test_email() -> Text:
    return 'test_user_' + gen_random_chars(14) + '@test.nu'


@pytest.fixture(scope='session')
def test_password() -> Text:
    return gen_random_chars(30)


@pytest.fixture(scope='session')
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


@pytest.fixture(scope='session')
def platform():
    return pl.system().lower()


@pytest.fixture(scope='session', autouse=True)
def cleanup(request, platform):
    """Cleanup local emails after each test session."""
    def remove_local_emails():
        files = glob.glob(get_local_inbox_path(platform) + '/*')
        for f in files:
            os.remove(f)

    request.addfinalizer(remove_local_emails)
