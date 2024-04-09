#!/bin/bash

set -ex

directory=${1-../protos}

# Produce generated Python code for gRPC
python -m grpc_tools.protoc -I$directory\
    --proto_path=$directory \
    --python_out=. \
    --pyi_out=. \
    --grpc_python_out=. \
    dht/dht.proto \

python -m grpc_tools.protoc -I$directory \
    --proto_path=$directory \
    --python_out=. \
    --pyi_out=. \
    --grpc_python_out=. \
    chord/chord.proto \