import sys

import key_lib as dht

def main():
    if len(sys.argv) < 2:
        print("Usage: client.py <address> <key> <value>")
        exit(1)
        
    ENTRY_ADDRESS = f"{sys.argv[1]}:{dht.PORT}"

    key = sys.argv[2]
    value = " ".join(sys.argv[3:]).encode("utf-8")
    if not value:
        value = sys.stdin.buffer.read()

    try:
        dht.set_key(ENTRY_ADDRESS, key, value)
    except Exception as e:
        sys.stderr.write(f"Error setting key: {e}")
        sys.exit(1)

    print(f"Key {key} set successfully")

if __name__ == "__main__":
    main()