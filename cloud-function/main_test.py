from main import app
from flask import json

def test_error_executable_empty():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'calldata': 'bitcoin', 'timeout': 123456}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))

  assert response.status_code == 400
  assert data['error'] == "executable field is missing from JSON request"

def test_error_calldata_empty():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'executable': '123', 'timeout': 123456}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))

  assert response.status_code == 400
  # assert data['error'] == "calldata field is missing from JSON request"