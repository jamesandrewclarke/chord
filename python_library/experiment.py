import sys, secrets

import key_lib as dht

def main():
    if len(sys.argv) < 1:
        print("Usage: client.py <address>") 
        exit(1)
        
    ENTRY_ADDRESS = f"{sys.argv[1]}:{dht.PORT}"

    keys = []
    for i in range(2000):
        random_stuff = secrets.token_bytes(255)
        dht.set_key(ENTRY_ADDRESS, str(random_stuff), random_stuff)
        keys.append(str(random_stuff))


    paths = []
    for i in range(len(keys)):
        pl, _ = dht.get_key(ENTRY_ADDRESS, keys[i])
        paths.append(pl)

    print(paths)
    print(f"Average: {sum(paths)/len(paths)}")

if __name__ == "__main__":
    main()