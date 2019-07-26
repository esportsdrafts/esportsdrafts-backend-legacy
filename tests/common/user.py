"""User class and functions to CRUD users."""

from typing import Text, List

import requests


class User():
    """eFantasy user account."""

    def __init__(self,
                 username: Text,
                 email: Text,
                 password: Text,
                 url: Text = 'https://api.efantasy.com'):
        """
        Initialize a user.

        Note that the user wont be authenticated until login() is called.

        Arguments:
        username -- String username
        email -- Email as string
        password -- User password as string

        """
        self.username = username
        self.email = email
        self.password = password
        self.url = url

        if url.endswith('/'):
            self.url = url[:-1]

        self.__signature_token = None
        self.__payload_token = None
        self.__user_roles = []  # type: List[Text]

    def login(self):
        """Authenticate the user using email + password."""
        payload = {
            'username': self.username,
            'password': self.password,
            'claim': 'username+password',
        }
        res = requests.post(self.url + '/v1/auth/auth', json=payload)
        res.raise_for_status()

    def logout(self):
        """Clear authentication for the user."""
        payload = {}
        requests.post(self.url + '/v1/auth/logout', json=payload)

    @property
    def is_authenticated(self) -> bool:
        """Indicate if the user is currently authenticated; otherwise False."""
        if self.__signature_token is None or self.__payload_token is None:
            return False
        res = requests.get(self.url + '/v1/user/me')
        return res.status_code == 200


def create_new_account(username: Text,
                       email: Text,
                       password: Text,
                       url: Text = 'https://api.efantasy.com') -> User:
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

    res = requests.post(url + '/v1/auth/register', payload)
    res.raise_for_status()

    return User(username, email, password, url)
