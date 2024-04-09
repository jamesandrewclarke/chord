import sys
import grpc
import dht.dht_pb2
import dht.dht_pb2_grpc

PORT = 8081

def set_key(addr: str, key: str, value: bytes):
    with grpc.insecure_channel(addr) as channel:
        stub = dht.dht_pb2_grpc.DHTStub(channel)
        req = dht.dht_pb2.SetKeyRequest(key=key, value=value)
        res = stub.SetKey(req)
        if res.forwardNode.address:
            forwardAddr = f"{res.forwardNode.address}:{PORT}"
            #sys.stderr.write(f"forwarding to {res.forwardNode.address}")
            return set_key(forwardAddr, key, value)
        else:
            return res, addr
            
        
def get_key(addr: str, key: str) -> bytes:
    with grpc.insecure_channel(addr) as channel:
        stub = dht.dht_pb2_grpc.DHTStub(channel)
        req = dht.dht_pb2.GetKeyRequest(key=key)
        res = stub.GetKey(req)
        if res.forwardNode.address:
            forwardAddr = f"{res.forwardNode.address}:{PORT}"
            #sys.stderr.write(f"forwarding to {res.forwardNode.address}")
            return get_key(forwardAddr, key)

        return res.value