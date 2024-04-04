import sys

import grpc
import chord_pb2
import chord_pb2_grpc

import hashlib

m = 64

def lookup(stub, id: int) -> int:
    req = chord_pb2.FindSuccessorRequest(id=int(id))
    res = stub.FindSuccessor(req)
    return res

def identifier_from_bytes(input: bytes) -> int:
    h = hashlib.sha256()
    h.update(input)
    
    bigid = int.from_bytes(h.digest())
    
    y = (1<<m-1)-1
    
    id = bigid % y

    return id

def set_key(stub, key: int, value: bytes):
    req = chord_pb2.SetKeyRequest(key=key, value=value)
    res = stub.SetKey(req)
    
    return res


def get_key(stub, key: int) -> str:
    req = chord_pb2.LookupRequest(key=key)
    res = stub.Lookup(req) 
    
    return res.value.decode('utf-8')

def main():
    if len(sys.argv) < 2:
        print("Usage: client.py <address> <key>")
        exit(1)
        
    address = sys.argv[1]
    test_input = " ".join(sys.argv[2:])

    key = test_input.encode("utf-8")
    id = identifier_from_bytes(key)

    with grpc.insecure_channel(address) as channel:
        stub = chord_pb2_grpc.ChordStub(channel)
        
        node = lookup(stub, id)
        print("on node: ", node.identifier)
        
    with grpc.insecure_channel(node.address) as channel:
        stub = chord_pb2_grpc.ChordStub(channel)
        set_key(stub, id, key)
        print(get_key(stub, id))


if __name__ == "__main__":
    main()