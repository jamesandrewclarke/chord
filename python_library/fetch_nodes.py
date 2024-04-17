import sys

import numpy as np
import requests

if __name__ == "__main__":
    PROM_ENDPOINT = "localhost:9090"
    if len(sys.argv)> 1:
        PROM_ENDPOINT = sys.argv[1]

    QUERY_ENDPOINT = f"http://{PROM_ENDPOINT}/api/v1/query"

    res = requests.get(QUERY_ENDPOINT, params={
        'query': 'chord_successor',
    }) 

    data = res.json()
    result = data['data']['result'] 
    
    ids = []
    for data in result:
        metric = data['metric']
        id = metric['id']
        ids.append(int(id))
        
    n = len(ids)
    print(n)
    np.save(f'results/node_ids_{n}', ids)
