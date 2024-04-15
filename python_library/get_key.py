import sys 

import key_lib as dht

def main():
    if len(sys.argv) < 1:
        print("Usage: client.py <address> <key>") 
        exit(1)
        
    ENTRY_ADDRESS = f"{sys.argv[1]}:{dht.PORT}"
    key = sys.argv[2]

    try:
        pl, b = dht.get_key(ENTRY_ADDRESS, key)
        print(pl, file=sys.stderr, flush=True)
        # sys.stdout.buffer.write(b)
    except Exception as e:
        sys.stderr.write(f"Error getting key: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()