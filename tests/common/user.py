"""User class and functions to CRUD users."""

import time
from typing import List, Text  # noqa

import requests
from tests.common.utils import raise_on_error


class User():
    """esportsdrafts user account."""

    def __init__(self,
                 username: Text,
                 email: Text,
                 password: Text,
                 url: Text = 'https://api.esportsdrafts.com'):
        """
        Initialize a user.

        Note that the user wont be authenticated until login() is called.

        Arguments:
        username -- String username
        email -- Email as string
        password -- User password as string

        """
        self.username = username.lower()
        self.email = email
        self.password = password
        self.url = url

        if url.endswith('/'):
            self.url = url[:-1]

        self.__auth_token = None
        self.__user_roles = []  # type: List[Text]

    def login(self):
        """Authenticate the user using email + password."""
        payload = {
            'username': self.username,
            'password': self.password,
            'claim': 'username+password',
        }
        res = requests.post(self.url + '/v1/auth/auth',
                            json=payload, verify=False)
        raise_on_error(res)

        res_json = res.json()
        self.__auth_token = res_json["access_token"]
        self.__auth_expires_in = res_json["expires_in"]

    def logout(self):
        """Clear authentication for the user."""
        payload = {}
        res = requests.post(self.url + '/v1/auth/auth',
                            json=payload, verify=False)
        raise_on_error(res)

    @property
    def is_authenticated(self) -> bool:
        """Indicate if the user is currently authenticated; otherwise False."""
        if self.__auth_token is None or \
                int(time.time()) > self.__auth_expires_in:
            self.__auth_token = None
            self.__auth_expires_in = None
            return False

        res = requests.get(self.url + '/v1/user/me', verify=False)
        raise_on_error(res)

        return True

    def __repr__(self):  # noqa
        return self.__str__()

    def __str__(self):  # noqa
        return f'User(username={self.username}, email={self.email})'

# TODO: Move all below to api_client class


def create_new_account(
        username: Text, email: Text, password: Text,
        url: Text = 'https://api.esportsdrafts.localhost') -> User:
    """
    Create a new user account and return a User object.

    Arguments:
    username -- The username as string
    email -- Email as string
    password -- Password as string
    url -- URL of API to call

    Returns:
    A fully contructed User object. NOTE: Not authenticated.

    """
    # Always using HTTPS, even for local development
    if not url.startswith('https://'):
        url_split = url.split('//')
        url = 'https://' + url_split[-1]

    if url.endswith('/'):
        url = url[:-1]

    payload = {
        'username': username,
        'email': email,
        'password': password,
    }

    res = requests.post(url + '/v1/auth/register', json=payload, verify=False)
    raise_on_error(res)

    return User(username, email, password, url)


def verify_email():
    pass


def reset_password():
    pass


def verify_password_reset():
    pass
