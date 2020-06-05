from flask import Flask, jsonify, json
import os
import shlex
import subprocess
import base64

# Copy and paste this file on Google Cloud function
# Set environment flag of MAX_EXECUTABLE, MAX_CALLDATA, MAX_TIMEOUT, MAX_STDOUT, MAX_STDERR

def get_env(env, flag):
  if not env[flag]:
    raise Exception(flag + " is missing")
  return int(env[flag])

def check_field_in_json(request_json, field):
  if not field in request_json:
    return (
        jsonify({"error": field + " field is missing from JSON request",}),
        400,
    )
  return ()

def check_max_exceed(request_json, field, max):
  if len(request_json[field]) > max:
    return jsonify({"error": field + " value exceeded max size",}), 400
  return ()

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

    try:
        request_json = request.get_json()
    except:
        return jsonify({"error": "invalid JSON request format",}), 400

    request_json = request.get_json()
    if request_json:

        validate_executable = check_field_in_json(request_json, "executable")
        if validate_executable: return validate_executable

        validate_calldata = check_field_in_json(request_json, "calldata")
        if validate_calldata: return validate_calldata

        validate_timeout = check_field_in_json(request_json, "timeout")
        if validate_timeout: return validate_timeout

        valadate_max_executable = check_max_exceed(request_json, "executable", MAX_EXECUTABLE)
        if valadate_max_executable: return valadate_max_executable

        valadate_max_calldata = check_max_exceed(request_json, "calldata", MAX_CALLDATA)
        if valadate_max_calldata: return valadate_max_calldata
        
        if not isinstance(request_json["timeout"], int):
            return jsonify({"error": "timeout type is invalid",}), 400
        elif request_json["timeout"] <= 0:
            return jsonify({"error": "Runtime must more than 0",}), 400
        elif request_json["timeout"] > MAX_TIMEOUT:
            return jsonify({"error": "Runtime exceeded max size",}), 400

    executable = ""
    try:
        executable = base64.b64decode(request_json["executable"])
    except:
        return jsonify({"error": "Can't decoded executable",}), 400

    path = "/tmp/execute.sh"
    with open(path, "w") as f:
        f.write(executable.decode())

    os.chmod(path, 0o775)

    try:

        timeout_millisec = request_json["timeout"]
        timeout_sec = timeout_millisec / 1000

        proc = subprocess.Popen(
            [path] + shlex.split(request_json["calldata"]),
            env=env,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
  
        proc.wait(timeout=timeout_sec)

        returncode = proc.returncode
        stdout = proc.stdout.read(MAX_STDOUT).decode()
        stderr = proc.stderr.read(MAX_STDERR).decode()

        return (
            jsonify(
                {
                    "returncode": returncode,
                    "stdout": stdout,
                    "stderr": stderr,
                    "err": "",
                }
            ),
            200,
        )

    except OSError:
        return (
            jsonify(
                {"returncode": 126, "stdout": "", "stderr": "", "err": "Execution fail"}
            ),
            200,
        )

    except subprocess.TimeoutExpired:
        return (
            jsonify(
                {
                    "returncode": 111,
                    "stdout": "",
                    "stderr": "",
                    "err": "Execution time limit exceeded",
                }
            ),
            200,
        )
