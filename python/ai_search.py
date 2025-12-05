import json
from typing import List

from numpyencoder import NumpyEncoder
import config
import requests


def add(searchDocs: List[dict]) -> bool:
    if len(searchDocs) == 0:
        print("No documents to add.")
        return

    data = {"value": [searchDoc for searchDoc in searchDocs]}

    res = requests.post(
        url=f"{config.AI_SEARCH_ENDPOINT}/indexes/index/docs/index?api-version={config.AI_SEARCH_API_VERSION}",
        headers={
            "Content-Type": "application/json",
            "api-key": config.AI_SEARCH_API_KEY,
        },
        data=json.dumps(data, cls=NumpyEncoder),
    )

    if res.status_code == 200:
        print(f"{len(searchDocs)} documents added to the search index successfully.")
        return True
    else:
        print(
            f"Failed to add documents to the search index. Status code: {res.status_code}, Status text: {res.text}"
        )
        return False


def search(query: str, top: int = 10, baseUrl: str = None) -> dict:
    data = json.dumps(
        {
            "search": query,
            "top": top,
            "select": "title, excerpt, date, url",
            "vectorQueries": [
                {
                    "kind": "vector",
                    "vector": list(config.model.query_embed(query=query))[0],
                    "k": 50,
                    "fields": "vector",
                },
            ],
            "filter": f"search.ismatch('*{baseUrl}', 'baseUrl')" if baseUrl else None,
            "orderby": "date desc" if baseUrl else None,
        },
        cls=NumpyEncoder,
    )

    res = requests.post(
        url=f"{config.AI_SEARCH_ENDPOINT}/indexes/index/docs/search?api-version={config.AI_SEARCH_API_VERSION}",
        headers={
            "Content-Type": "application/json",
            "api-key": config.AI_SEARCH_API_KEY,
        },
        data=data,
    )

    res.raise_for_status()
    return res.json()
