from main import app
from flask import json
from test.support import EnvironmentVarGuard
import os

def test_error_invalid_json_request():
  response = app.test_client().post(
      '/execute',
      data="{'executable': '123', 'calldata':}",
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 400
  assert data['error'] == "invalid JSON request format"

def test_error_missing_executable():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'calldata': 'bitcoin', 'timeout': 123456}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 400
  assert data['error'] == "executable field is missing from JSON request"

def test_error_executable_empty():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'executable': '', 'calldata': 'bitcoin', 'timeout': 123456}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 200
  assert data['returncode'] == 126
  assert data['stdout'] == ""
  assert data['stderr'] == ""
  assert data['err'] == "Execution fail"

def test_error_missing_calldata():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'executable': '123', 'timeout': 123456}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 400
  assert data['error'] == "calldata field is missing from JSON request"


def test_error_calldata_empty():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'executable': 'IyEvdXNyL2Jpbi9lbnYgcHl0aG9uMwpwcmludCgnaGVsbG8nKQ==', 'calldata': '', 'timeout': 123456}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 200
  assert data['returncode'] == 0
  assert data['stdout'] == "hello\n"
  assert data['stderr'] == ""
  assert data['err'] == ""

def test_error_missing_timeout():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'executable': '123', 'calldata': 'bitcoin'}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 400
  assert data['error'] == "timeout field is missing from JSON request"

def test_error_timeout_empty():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'executable': '123', 'calldata': 'bitcoin', 'timeout': ''}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 400
  assert data['error'] == "timeout field is empty"

def test_error_timeout_less_than_0():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'executable': '123', 'calldata': 'bitcoin', 'timeout': -5}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 400
  assert data['error'] == "Runtime must more than 0"

def test_error_timeout_more_than_max_timeout():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({'executable': '123', 'calldata': 'bitcoin', 'timeout': 1111111111111111111111111111111111111111}),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 400
  assert data['error'] == "Runtime exceeded max size"

def test_success_execution():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({
        "calldata": "123",
        "executable": "IyEvdXNyL2Jpbi9lbnYgcHl0aG9uMwpwcmludCgnaGVsbG8nKQ==", #!/usr/bin/env python3\nprint('hello')
        "timeout": 123456
      }),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 200
  assert data['returncode'] == 0
  assert data['stdout'] == "hello\n"
  assert data['stderr'] == ""
  assert data['err'] == ""

def test_error_execution_fail():
  response = app.test_client().post(
      '/execute',
      data=json.dumps({
        "calldata": "123",
        "executable": "IyEvdXNyL2Jpbi9lbnZlZWVlZWVlZWUKcHl0aG9uM1xucHJpbnQoJ2hlbGxvJyk=", #!/usr/bin/enveeeeeeeee python3\nprint('hello')
        "timeout": 123456
      }),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 200
  assert data['returncode'] == 126
  assert data['stdout'] == ""
  assert data['stderr'] == ""
  assert data['err'] == "Execution fail"

def test_error_execution_timeout():
  '''#!/usr/bin/env python3
      import time

      time.sleep(1)
  '''
  response = app.test_client().post(
      '/execute',
      data=json.dumps({
        "calldata": "123",
        "executable": "IyEvdXNyL2Jpbi9lbnYgcHl0aG9uMwppbXBvcnQgdGltZQoKdGltZS5zbGVlcCgxKQ==",
        "timeout": 100 #100 millisec
      }),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 200
  assert data['returncode'] == 111
  assert data['stdout'] == ""
  assert data['stderr'] == ""
  assert data['err'] == "Execution time limit exceeded"

def test_success_execution_timeout():
  '''#!/usr/bin/env python3
      import time

      time.sleep(1) # 1000 millisec
  '''
  response = app.test_client().post(
      '/execute',
      data=json.dumps({
        "calldata": "123",
        "executable": "IyEvdXNyL2Jpbi9lbnYgcHl0aG9uMwppbXBvcnQgdGltZQoKdGltZS5zbGVlcCgxKQpwcmludCgiaGVsbG8iKQ==",
        "timeout": 2000  # 2000 millisec
      }),
      content_type='application/json',
  )

  data = json.loads(response.get_data(as_text=True))
  assert response.status_code == 200
  assert data['returncode'] == 0
  assert data['stdout'] == "hello\n"
  assert data['stderr'] == ""
  assert data['err'] == ""


