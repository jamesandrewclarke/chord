import sys

import grpc
import chord.chord_pb2
import chord.chord_pb2_grpc

def lookup(stub, id: int) -> int:
    req = chord.chord_pb2.FindSuccessorRequest(id=int(id))
    res = stub.FindSuccessor(req)
    return res.identifier


def main():
    if len(sys.argv) < 2:
        print("Usage: client.py <address>") 
        exit(1)
        
    address = sys.argv[1]
    with grpc.insecure_channel(address) as channel:
        stub = chord.chord_pb2_grpc.ChordStub(channel)
        
        print("Starting lookup...")
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
            
        print()
        print(f"Total nodes: {count}")

if __name__ == "__main__":
    main()