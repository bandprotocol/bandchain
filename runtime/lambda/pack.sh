#!/bin/bash

yum install -y zip python3

cd /pack
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

cd venv/lib/python3.7/site-packages
zip -r9 ${OLDPWD}/function.zip .

cd ${OLDPWD}
zip -g function.zip run.py
deactivate
