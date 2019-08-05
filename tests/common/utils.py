
import requests


def raise_on_error(request):
    if request.status_code >= 400:
        json_res = request.json()
        raise requests.HTTPError(json_res)

    return None
