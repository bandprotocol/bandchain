import json
import os
import shlex
import base64
import subprocess

HEADERS = {
    "content-type": "application/json",
    "x-lambda": "true",
    "access-control-allow-origin": "*",
    "access-control-allow-methods": "OPTIONS, POST",
}

runtime_version = "${RUNTIME_VERSION}"


def success(returncode, stdout, stderr, err):
    return {
        "statusCode": 200,
        "headers": HEADERS,
        "body": json.dumps(
            {
                "returncode": returncode,
                "stdout": stdout,
                "stderr": stderr,
                "error": err,
                "version": runtime_version,
            }
        ),
    }


def bad_request(err):
    return {
        "statusCode": 400,
        "headers": HEADERS,
        "body": json.dumps({"error": err}),
    }


def get_env(env, flag):
    if flag not in env:
        raise Exception(flag + " is missing")
    return int(env[flag])


def lambda_handler(event, context):
    try:
        body = json.loads(event["body"])
    except (json.decoder.JSONDecodeError, KeyError) as e:
        # Hack for preflight
        return {
            "statusCode": 200,
            "headers": HEADERS,
        }

    env = os.environ.copy()

    MAX_EXECUTABLE = get_env(env, "MAX_EXECUTABLE")
    MAX_DATA_SIZE = get_env(env, "MAX_DATA_SIZE")

    if "executable" not in body:
        return bad_request("Missing executable value")
    executable = base64.b64decode(body["executable"])
    if len(executable) > MAX_EXECUTABLE:
        return bad_request("Executable exceeds max size")
    if "calldata" not in body:
        return bad_request("Missing calldata value")
    if len(body["calldata"]) > MAX_DATA_SIZE:
        return bad_request("Calldata exceeds max size")
    if "timeout" not in body:
        return bad_request("Missing timeout value")
    try:
        timeout = int(body["timeout"])
    except ValueError:
        return bad_request("Timeout format invalid")

    path = "/tmp/execute.sh"
    with open(path, "w") as f:
        f.write(executable.decode())

    os.chmod(path, 0o775)
    try:
        env = os.environ.copy()
        env["PYTHONPATH"] = os.getcwd()
        for key, value in body.get("env", {}).items():
            env[key] = value

        proc = subprocess.Popen(
            [path] + shlex.split(body["calldata"]),
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
