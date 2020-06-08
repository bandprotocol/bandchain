from flask import Flask, jsonify, json
import os
import shlex
import subprocess
import base64
import werkzeug
from marshmallow import Schema, fields, ValidationError, validate

# Copy and paste this file on Google Cloud function
# Set environment flag of MAX_EXECUTABLE, MAX_CALLDATA, MAX_TIMEOUT, MAX_STDOUT, MAX_STDERR


def get_env(env, flag):
    if flag not in env:
        raise Exception(flag + " is missing")
    return int(env[flag])


def success(returncode, stdout, stderr, err):
    return (
        jsonify(
            {"returncode": returncode, "stdout": stdout, "stderr": stderr, "err": err}
        ),
        200,
    )


def bad_request(err):
    return jsonify({"error": err}), 400


def execute(request):
    """Responds to any HTTP request.
  Args:
      request (flask.Request): HTTP request object.
  Returns:
      The response text or any set of values that can be turned into a
      Response object using
      `make_response <http://flask.pocoo.org/docs/1.0/api/#flask.Flask.make_response>`.
  """
    env = os.environ.copy()

    MAX_EXECUTABLE = get_env(env, "MAX_EXECUTABLE")
    MAX_CALLDATA = get_env(env, "MAX_CALLDATA")
    MAX_TIMEOUT = get_env(env, "MAX_TIMEOUT")
    MAX_STDOUT = get_env(env, "MAX_STDOUT")
    MAX_STDERR = get_env(env, "MAX_STDERR")

    class Executable(fields.Field):
        def _deserialize(self, value, attr, data, **kwargs):
            try:
                return base64.b64decode(value).decode()
            except:
                raise ValidationError("Can't decoded executable")

    class RequestSchema(Schema):
        executable = Executable(
            required=True,
            validate=validate.Length(max=MAX_EXECUTABLE),
            error_messages={"required": "field is missing from JSON request"},
        )
        calldata = fields.Str(
            required=True,
            validate=validate.Length(max=MAX_CALLDATA),
            error_messages={"required": "field is missing from JSON request"},
        )
        timeout = fields.Int(
            required=True,
            validate=validate.Range(min=0, max=MAX_TIMEOUT),
            error_messages={"required": "field is missing from JSON request"},
        )

    try:
        request_json = request.get_json(force=True)
    except werkzeug.exceptions.BadRequest:
        return bad_request("invalid JSON request format")

    loaded_request = ""
    try:
        loaded_request = RequestSchema().load(request_json)
    except ValidationError as err:
        return bad_request(err.messages)

    path = "/tmp/execute.sh"
    with open(path, "w") as f:
        f.write(loaded_request.get("executable"))

    os.chmod(path, 0o775)

    try:

        timeout_millisec = loaded_request.get("timeout")
        timeout_sec = timeout_millisec / 1000

        proc = subprocess.Popen(
            [path] + shlex.split(loaded_request.get("calldata")),
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )

        proc.wait(timeout=timeout_sec)

        returncode = proc.returncode
        stdout = proc.stdout.read(MAX_STDOUT).decode()
        stderr = proc.stderr.read(MAX_STDERR).decode()

        return success(returncode, stdout, stderr, "")

    except OSError:
        return success(126, "", "", "Execution fail")

    except subprocess.TimeoutExpired:
        return success(111, "", "", "Execution time limit exceeded")
