#!/bin/bash

# clean old build
rm -r build dist pyband.egg-info

# build new package
python3 setup.py sdist bdist_wheel

# publish to testpypi
python3 -m twine upload --repository testpypi dist/*
