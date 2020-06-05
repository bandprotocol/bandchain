from flask import Flask, abort, jsonify, json, request
import os
import shlex
import subprocess
import base64

app = Flask(__name__)

env = os.environ.copy()

MAX_EXECUTABLE = env['MAX_EXECUTABLE']
MAX_CALLDATA = env['MAX_CALLDATA']
MAX_TIMEOUT = env['MAX_TIMEOUT']
MAX_STDOUT = env['MAX_STDOUT']
MAX_STDERR = env['MAX_STDERR']

if not MAX_EXECUTABLE:
    raise Exception("MAX_EXECUTABLE is missing")
if not MAX_CALLDATA:
    raise Exception("MAX_CALLDATA is missing")
if not MAX_TIMEOUT:
    raise Exception("MAX_TIMEOUT is missing")
if not MAX_STDOUT:
    raise Exception("MAX_STDOUT is missing")
if not MAX_STDERR:
    raise Exception("MAX_STDERR is missing")

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
    try:
        request_json = request.get_json()
    except:
        return jsonify({
            "error": "invalid JSON request format",
        }), 400

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
        elif not request_json['timeout']:
            return jsonify({
                "error": "timeout field is empty",
            }), 400
        elif len(request_json['executable']) > int(MAX_EXECUTABLE):
            return jsonify({
                "error": "executable value exceed max size",
            }), 400
        elif len(request_json['calldata']) > int(MAX_CALLDATA):
            return jsonify({
                "error": "calldata value exceed max size",
            }), 400
        elif not isinstance(request_json['timeout'], int):
            return jsonify({
                "error": "timeout type is invalid",
            }), 400
        elif request_json['timeout'] > int(MAX_TIMEOUT):
            return jsonify({
                "error": "Runtime exceeded max size",
            }), 400
        elif request_json['timeout'] <= 0:
            return jsonify({
                "error": "Runtime must more than 0",
            }), 400

    executable = ""
    try:
        executable = base64.b64decode(request_json["executable"])
    except:
        return jsonify({
            "error": "Can't decoded executable",
        }), 400
    
    path = "/tmp/execute.sh"
    with open(path, "w") as f:
        f.write(executable.decode())
    os.chmod(path, 0o775)

    try:
        
        timeout_millisec = request_json['timeout']
        timeout_sec = timeout_millisec/1000

        result = subprocess.run(
            [path] + shlex.split(request_json["calldata"]), env=env, timeout=timeout_sec, capture_output=True
        )

        returncode = result.returncode
        stdout = result.stdout.decode()
        stderr = result.stderr.decode()
        print ("returncode", returncode)
        print ("stdout", stdout)
        print ("stderr", stderr)

        if len(stdout) > int(MAX_STDOUT):
            return jsonify({
                "returncode": 112,
                "stdout": "",
                "stderr": "",
                "err": "Stdout exceeded max size"
            }), 200
        
        if len(stderr) > int(MAX_STDERR):
            return jsonify({
                "returncode": 113,
                "stdout": "",
                "stderr": "",
                "err": "Stderr exceeded max size"
            }), 200
        

        return jsonify({
            "returncode": returncode,
            "stdout": stdout,
            "stderr": stderr,
            "err": ""
        }), 200
    
    except OSError:
        return jsonify({
            "returncode": 126,
            "stdout": "",
            "stderr": "",
            "err": "Execution fail"
        }), 200

    except subprocess.TimeoutExpired:
        return jsonify({
            "returncode": 111,
            "stdout": "",
            "stderr": "",
            "err": "Execution time limit exceeded"
        }), 200

