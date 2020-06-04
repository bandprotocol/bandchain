from flask import Flask, request, abort, jsonify
import os

app = Flask(__name__)

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

    MAX_EXECUTABLE = os.getenv('MAX_EXECUTABLE')
    MAX_CALLDATA = os.getenv('MAX_CALLDATA')
    MAX_TIMEOUT = os.getenv('MAX_TIMEOUT')
    if not MAX_EXECUTABLE:
        exit(101)
    if not MAX_CALLDATA:
        exit(102)
    if not MAX_TIMEOUT:
        exit(103)

    print ("MAX_EXECUTABLE", MAX_EXECUTABLE)
    print ("MAX_CALLDATA", MAX_CALLDATA)
    print ("MAX_TIMEOUT", MAX_TIMEOUT)
    
    request_json = request.get_json()
    if request_json:
        if not 'executable' in request_json:
            return jsonify({
                "error": "executable field is missing from JSON request",
            }), 400
        elif len(request_json['executable']) > int(MAX_EXECUTABLE):
            return jsonify({
                "error": "executable value exceed max size",
            }), 400
    
        if not 'calldata' in request_json:
            return jsonify({
                "error": "calldata field is missing from JSON request",
            }), 400
        elif len(request_json['calldata']) > int(MAX_CALLDATA):
            return jsonify({
                "error": "calldata value exceed max size",
            }), 400
        
        if not 'timeout' in request_json:
            return jsonify({
                "error": "timeout field is missing from JSON request",
            }), 400
        elif not isinstance(request_json['timeout'], int) or request_json['timeout'] < 0:
            return jsonify({
                "error": "timeout type is invalid",
            }), 400
        elif request_json['timeout'] > int(MAX_TIMEOUT):
            return jsonify({
                "error": "Runtime exceeded max size",
            }), 400
        
    else:
        return f'Hello World2!'
