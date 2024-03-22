import sys

import grpc
import chord_pb2
import chord_pb2_grpc

def lookup(stub, id: int) -> int:
    req = chord_pb2.FindSuccessorRequest(id=int(id))
    res = stub.FindSuccessor(req)
    return res.identifier


def main():
    if len(sys.argv) < 3:
        print("Usage: client.py <address> <id>")
        exit(1)
        
    address = sys.argv[1]
    id = sys.argv[2]
    with grpc.insecure_channel(address) as channel:
        stub = chord_pb2_grpc.ChordStub(channel)
        
        first = lookup(stub, 0)
        count = 0
        nodes = set()
        nodes.add(first)
        id = first
        while True:
            result = lookup(stub, id+1)
            nodes.add(id)
            print(id)
            count += 1
            if result == first:
                break
            id = result 
            
        print(count)
        print(nodes)

if __name__ == "__main__":
    main()