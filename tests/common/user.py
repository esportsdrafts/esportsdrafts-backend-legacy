"""User class and functions to CRUD users."""

import time
from typing import List, Optional, Text  # noqa

import jwt
import requests
from tests.common.utils import raise_on_error

requests.packages.urllib3.disable_warnings()


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
        self.roles = []  # type: List[Text]
        self.user_id = None  # type: Optional[Text]

        if url.endswith('/'):
            self.url = url[:-1]

        self.__auth_token = None  # type: Optional[Text]
        self.__user_roles = []  # type: List[Text]

    def login(self):
        """Authenticate the user using email + password."""
        payload = {
            'username': self.username,
            'password': self.password,
            'claim': 'username+password',
        }
        res = requests.post(self.url + '/v1/auth/auth',
                            json=payload,
                            verify=not self.url.endswith('.localhost'))
        raise_on_error(res)

        res_json = res.json()
        self.__auth_token = res_json['access_token']
        self.__auth_expires_in = res_json['expires_in']

        # Grab claims without verifying validity
        claims = jwt.decode(self.__auth_token, verify=False)
        self.user_id = claims.get('user_id')
        self.roles = claims.get('roles', [])

    def logout(self):
        """Clear authentication for the user."""
        payload = {}
        res = requests.post(self.url + '/v1/auth/auth',
                            json=payload,
                            verify=not self.url.endswith('.localhost'))
        raise_on_error(res)

    @property
    def is_authenticated(self) -> bool:
        """Indicate if the user is currently authenticated; otherwise False."""
        if self.__auth_token is None or \
                int(time.time()) > self.__auth_expires_in:
            self.__auth_token = None
            self.__auth_expires_in = None
            self.roles = []
            return False

        res = requests.get(self.url + '/v1/user/me',
                           verify=not self.url.endswith('.localhost'))
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

    res = requests.post(url + '/v1/auth/register', json=payload,
                        verify=not url.endswith('.localhost'))
    raise_on_error(res)

    return User(username, email, password, url)


def verify_email(user: User, token: Text) -> None:
    payload = {
        'username': user.username,
        'token': token,
    }
    res = requests.post(user.url + '/v1/auth/verifyemail', json=payload,
                        verify=not user.url.endswith('.localhost'))
    raise_on_error(res)


def check_username_available(username: Text, env: Text) -> bool:
    if username is None:
        return False
    res = requests.get(f'{env}/v1/auth/check?username={username}',
                       verify=not env.endswith('.localhost'))
    return res.status_code == 200


def reset_password_request(user: User):
    payload = {
        'username': user.username,
        'email': user.email,
    }
    res = requests.post(
        user.url + '/v1/auth/passwordreset/request', json=payload,
        verify=not user.url.endswith('.localhost'))
    raise_on_error(res)


def verify_password_reset(user: User, token: Text, new_password: Text):
    payload = {
        'username': user.username,
        'token': token,
        'password': new_password,
    }
    res = requests.post(
        user.url + '/v1/auth/passwordreset/verify', json=payload,
        verify=not user.url.endswith('.localhost'))
    raise_on_error(res)
