import sys, secrets

import numpy as np
import time

import key_lib as dht

def main():
    if len(sys.argv) < 2:
        print("Usage: experiment.py <address> <output_name>") 
        exit(1)
        
    ENTRY_ADDRESS = f"{sys.argv[1]}:{dht.PORT}"
    TEST_NAME = sys.argv[2]

    keys = []
    
    N = 1000
    M = 255
    for i in range(N):
        random_stuff = secrets.token_bytes(M)
        dht.set_key(ENTRY_ADDRESS, str(random_stuff), random_stuff)
        keys.append(str(random_stuff))


    paths = np.empty((N,2))
    for i in range(len(keys)):
        start = time.time() 
        pl, _ = dht.get_key(ENTRY_ADDRESS, keys[i])
        latency = time.time() - start
        paths[i] = [pl,latency]
        
    np.save(TEST_NAME, paths)

if __name__ == "__main__":
    main()