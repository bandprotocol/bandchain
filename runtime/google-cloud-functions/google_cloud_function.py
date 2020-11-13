from flask import jsonify
import os
import shlex
import subprocess
import base64

# Set environment flag of MAX_EXECUTABLE, MAX_DATA_SIZE


runtime_version = "${RUNTIME_VERSION}"


def get_env(env, flag):
    if flag not in env:
        raise Exception(flag + " is missing")
    return int(env[flag])


def success(returncode, stdout, stderr, err):
    return (
        jsonify(
            {
                "returncode": returncode,
                "stdout": stdout,
                "stderr": stderr,
                "err": err,
                "version": runtime_version,
            }
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
    MAX_DATA_SIZE = get_env(env, "MAX_DATA_SIZE")

    request_json = request.get_json(force=True)
    if "executable" not in request_json:
        return bad_request("Missing executable value")
    if len(request_json["executable"]) > MAX_EXECUTABLE:
        return bad_request("Executable exceeds max size")
    if "calldata" not in request_json:
        return bad_request("Missing calldata value")
    if len(request_json["calldata"]) > MAX_DATA_SIZE:
        return bad_request("Calldata exceeds max size")
    if "timeout" not in request_json:
        return bad_request("Missing timeout value")
    try:
        timeout = int(request_json["timeout"])
    except ValueError:
        return bad_request("Timeout format invalid")

    path = "/tmp/execute.sh"
    with open(path, "w") as f:
        f.write(base64.b64decode(request_json["executable"]).decode())

    os.chmod(path, 0o775)
    try:
        env = os.environ.copy()
        for key, value in request_json.get("env", {}).items():
            env[key] = value

        proc = subprocess.Popen(
            [path] + shlex.split(request_json["calldata"]),
            env=env,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )

        proc.wait(timeout=(timeout / 1000))
        returncode = proc.returncode
        stdout = proc.stdout.read(MAX_DATA_SIZE).decode()
        stderr = proc.stderr.read(MAX_DATA_SIZE).decode()
        return success(returncode, stdout, stderr, "")
    except OSError:
        return success(126, "", "", "Execution fail")
    except subprocess.TimeoutExpired:
        return success(111, "", "", "Execution time limit exceeded")
