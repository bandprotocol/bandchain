from flask import Flask, request, abort, jsonify
app = Flask(__name__)

@app.route('/hello', methods=['POST'])
def hello_world():
    """Responds to any HTTP request.
    Args:
        request (flask.Request): HTTP request object.
    Returns:
        The response text or any set of values that can be turned into a
        Response object using
        `make_response <http://flask.pocoo.org/docs/1.0/api/#flask.Flask.make_response>`.
    """
    request_json = request.get_json()
    if request_json and not 'executable' in request_json:
      return jsonify({
        "error": "executable field is missing from JSON request",
      }), 400
    elif request_json and not 'calldata' in request_json:
      return jsonify({
        "error": "calldata field is missing from JSON request",
      }), 400
    elif request_json and not 'timeout' in request_json:
      return jsonify({
        "error": "timeout field is missing from JSON request",
      }), 400
    else:
        return f'Hello World2!'
