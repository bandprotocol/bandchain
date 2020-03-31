import json
import os
import shlex
import subprocess

HEADERS = {
    "content-type": "application/json",
    "x-lambda": "true",
    "access-control-allow-origin": "*",
    "access-control-allow-methods": "OPTIONS, POST",
}

def lambda_handler(event, context):
    try:
        body = json.loads(event["body"])
    except (json.decoder.JSONDecodeError, KeyError) as e:
        # Hack for preflight
        return {
            "statusCode": 200,
            "headers": HEADERS,
        }

    if "executable" not in body:
        return {
            "statusCode": 400,
            "headers": HEADERS,
            "body": json.dumps({"error": "Missing executable value",}),
        }
    if "calldata" not in body:
        return {
            "statusCode": 400,
            "headers": HEADERS,
            "body": json.dumps({"error": "Missing calldata value",}),
        }

    path = "/tmp/execute.sh"
    with open(path, "w") as f:
        f.write(body["executable"])

    os.chmod(path, 0o775)
    try:
        env = os.environ.copy()
        env["PATH"] = env["PATH"] + ":" + os.path.join(os.getcwd(), "exec", "usr", "bin")
        env["LD_LIBRARY_PATH"] = (
            env["LD_LIBRARY_PATH"]
            + ":"
            + os.path.join(os.getcwd(), "exec", "lib64")
            + ":"
            + os.path.join(os.getcwd(), "exec", "usr", "lib64")
        )

        result = subprocess.run(
            [path] + shlex.split(body["calldata"]), env=env, timeout=3, capture_output=True
        )

        return {
            "statusCode": 200,
            "headers": HEADERS,
            "body": json.dumps(
                {
                    "returncode": result.returncode,
                    "stdout": result.stdout.decode(),
                    "stderr": result.stderr.decode(),
                }
            ),
        }
    except subprocess.TimeoutExpired:
        return {
            "statusCode": 200,
            "headers": HEADERS,
            "body": json.dumps({"returncode": -1, "error": "Execution timeout",}),
        }

