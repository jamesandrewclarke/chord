import sys

import grpc
import dht.dht_pb2
import dht.dht_pb2_grpc

import hashlib

m = 64

DHT_PORT = 8081

def identifier_from_bytes(input: bytes) -> int:
    h = hashlib.sha256()
    h.update(input)
    
    bigid = int.from_bytes(h.digest())
    
    y = (1<<m-1)-1
    
    id = bigid % y

    return id

def set_key(addr: str, key: int, value: bytes):
    with grpc.insecure_channel(addr) as channel:
        stub = dht.dht_pb2_grpc.DHTStub(channel)
        req = dht.dht_pb2.SetKeyRequest(key=key, value=value)
        res = stub.SetKey(req)
        if res.forwardNode.address:
            forwardAddr = f"{res.forwardNode.address}:{DHT_PORT}"
            print("Forwarding...")
            return set_key(forwardAddr, key, value)
        else:
            return res, addr
            
        
def get_key(addr: str, key: int) -> str:
    with grpc.insecure_channel(addr) as channel:
        stub = dht.dht_pb2_grpc.DHTStub(channel)
        req = dht.dht_pb2.GetKeyRequest(key=key)
        res = stub.GetKey(req)
        return res

        
def main():
    if len(sys.argv) < 2:
        print("Usage: client.py <address> <key>")
        exit(1)
        
    ENTRY_ADDRESS = f"{sys.argv[1]}:{DHT_PORT}"
    test_input = " ".join(sys.argv[2:])

    key = test_input.encode("utf-8")
    id = identifier_from_bytes(key)

    res, addr = set_key(ENTRY_ADDRESS, id, key)
    print(get_key(addr, id))

if __name__ == "__main__":
    main()