# Produce generated Python code for gRPC
python -m grpc_tools.protoc -I../protos \
    --python_out=. \
    --pyi_out=. \
    --grpc_python_out=. \
    ../protos/dht/dht.proto

python -m grpc_tools.protoc -I../protos \
    --python_out=. \
    --pyi_out=. \
    --grpc_python_out=. \
    ../protos/chord/chord.proto