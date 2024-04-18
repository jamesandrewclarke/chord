import sys

import numpy as np
import requests

if __name__ == "__main__":
    PROM_ENDPOINT = "localhost:9090"
    if len(sys.argv)> 1:
        PROM_ENDPOINT = sys.argv[1]

    QUERY_ENDPOINT = f"http://{PROM_ENDPOINT}/api/v1/query"

    res = requests.get(QUERY_ENDPOINT, params={
        'query': 'dht_keys_total',
    }) 

    data = res.json()
    result = data['data']['result'] 

    
    key_totals = [
        [int(value['value'][1]), int(value['metric']['id'])] for value in result
    ]

    n = len(key_totals)

    print(f"Got {n} totals")
    arr = np.array(key_totals)

    keys = arr[:,0]
    print(f"1% = {np.percentile(keys, 1)}")
    print(f"50% = {np.percentile(keys, 50)}")
    print(f"99% = {np.percentile(keys, 99)}")

    print(f"Sum = {np.sum(keys)}")

    np.save(f"results/dht_key_totals_{n}", key_totals)
