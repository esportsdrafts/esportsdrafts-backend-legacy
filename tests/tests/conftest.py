
import pytest


def pytest_addoption(parser):
    parser.addoption(
        "--env", action="store", default="local",
        choices=['local', 'dev', 'stage', 'prod'],
        help='Pick which environment to run integration tests against')


@pytest.fixture()
def env(pytestconfig):
    return pytestconfig.getoption("env")
