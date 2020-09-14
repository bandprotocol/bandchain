import os
import jwt
import json


def response(is_authorized):
    return {
        "isAuthorized": is_authorized,
        "context": {},
    }


def lambda_handler(event, context):
    if "headers" not in event:
        return response(False)
    headers = event["headers"]

    if "authorization" not in headers:
        return response(False)
    token = headers["authorization"]

    env = os.environ.copy()
    if "JWT_SECRET_KEY" not in env:
        return response(False)

    SINGING_KEY = env["JWT_SECRET_KEY"]
    try:
        jwt.decode(token, SINGING_KEY, algorithms=['HS256'])
    except:
        return response(False)

    return response(True)
