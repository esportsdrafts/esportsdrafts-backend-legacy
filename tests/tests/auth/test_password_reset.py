
import time

from tests.common.email import get_emails_from_local_inbox, read_local_email, get_verification_token
from tests.common.user import reset_password_request, verify_password_reset
from tests.common.utils import gen_random_chars


def test_reset_password(user, env):
    reset_password_request(user)
    sent_time = int(time.time())
    time_epsilon = 15

    # Encapsulate in an Email class
    if env == 'local':
        # Initial wait for email to get sent
        time.sleep(1)
        attempts = 0
        email_found = None

        while True:
            attempts += 1
            attempts <= 3, 'No reset password email found with correct content'

            emails = get_emails_from_local_inbox(
                user.username, 'reset_password')

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
                email_found = parsed
                break

            # Wait a bit for results to come in
            time.sleep(3)

        user_id, token = get_verification_token(email_found)
        new_password = gen_random_chars(20)
        verify_password_reset(user, token, new_password)

        user.password = new_password
        user.login()
    else:
        # TODO: Check cloud database for email
        # each 'test' email is copied to a DB in GCP
        pass
