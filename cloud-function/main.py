from flask import Flask, abort, jsonify, json, request
import os
import shlex
import subprocess

app = Flask(__name__)

env = os.environ.copy()

MAX_EXECUTABLE = env['MAX_EXECUTABLE']
MAX_CALLDATA = env['MAX_CALLDATA']
MAX_TIMEOUT = env['MAX_TIMEOUT']
if not MAX_EXECUTABLE:
    raise Exception("MAX_EXECUTABLE is missing")
if not MAX_CALLDATA:
    raise Exception("MAX_CALLDATA is missing")
if not MAX_TIMEOUT:
    raise Exception("MAX_TIMEOUT is missing")

@app.route('/execute', methods=['POST'])
def execute():
    """Responds to any HTTP request.
    Args:
        request (flask.Request): HTTP request object.
    Returns:
        The response text or any set of values that can be turned into a
        Response object using
        `make_response <http://flask.pocoo.org/docs/1.0/api/#flask.Flask.make_response>`.
    """
    
    request_json = request.get_json()
    if request_json:
        if not 'executable' in request_json:
            return jsonify({
                "error": "executable field is missing from JSON request",
            }), 400
        elif not 'calldata' in request_json:
            return jsonify({
                "error": "calldata field is missing from JSON request",
            }), 400
        elif not 'timeout' in request_json:
            return jsonify({
                "error": "timeout field is missing from JSON request",
            }), 400
        elif len(request_json['executable']) > int(MAX_EXECUTABLE):
            return jsonify({
                "error": "executable value exceed max size",
            }), 400
        elif len(request_json['calldata']) > int(MAX_CALLDATA):
            return jsonify({
                "error": "calldata value exceed max size",
            }), 400
        elif not isinstance(request_json['timeout'], int) or request_json['timeout'] < 0:
            return jsonify({
                "error": "timeout type is invalid",
            }), 400
        elif request_json['timeout'] > int(MAX_TIMEOUT):
            return jsonify({
                "error": "Runtime exceeded max size",
            }), 400
    
    path = "/tmp/execute.sh"
    with open(path, "w") as f:
        f.write(request_json["executable"])
    os.chmod(path, 0o775)

    try:
        env["PATH"] = env["PATH"] + ":" + os.path.join(os.getcwd(), "exec", "usr", "bin")
        print("PATH", env["PATH"])

        result = subprocess.run(
            [path] + shlex.split(request_json["calldata"]), env=env, timeout=request_json['timeout'], capture_output=True
        )

        return jsonify({
            "returncode": result.returncode,
            "stdout": result.stdout.decode(),
            "stderr": result.stderr.decode(),
            "err": ""
        }), 200
    
    except OSError:
        return jsonify({
            "returncode": 126,
            "stdout": "",
            "stderr": "",
            "err": "Data source could not be invoked"
        }), 200

    except subprocess.TimeoutExpired:
        return jsonify({
            "returncode": 111,
            "stdout": "",
            "stderr": "",
            "err": "Execution time limit exceeded"
        }), 200

