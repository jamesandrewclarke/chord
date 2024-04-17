#!/bin/bash

set -ex

docker build -t chord.azurecr.io/app/python:local_james -f python_library/Dockerfile --platform amd64 .
docker push chord.azurecr.io/app/python:local_james