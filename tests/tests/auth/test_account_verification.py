
from tests.common.email import (get_emails_from_local_inbox,
                                get_verification_token, read_local_email)
from tests.common.user import verify_email


def test_verify_email(user):
    email = get_emails_from_local_inbox(user.username, 'welcome')[0]
    email_content = read_local_email(email)
    user_id, token = get_verification_token(email_content)

    assert user_id == user.username

    verify_email(user, token)


def test_verify_verified(user):
    email = get_emails_from_local_inbox(user.username, 'welcome')[0]
    email_content = read_local_email(email)
    user_id, token = get_verification_token(email_content)

    assert user_id == user.username

    verify_email(user, token)
    verify_email(user, token)


def test_verify_invalid_token(user):
    try:
        verify_email(user, 'random_token_that_is_wrong')
        assert False
    except Exception:
        pass
