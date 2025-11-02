import json
from typing import List

from numpyencoder import NumpyEncoder
import config
import requests


def add(searchDocs: List[dict]):
    data = {"value": [searchDoc for searchDoc in searchDocs]}

    res = requests.post(
        url=f"{config.AI_SEARCH_ENDPOINT}/indexes/index/docs/index?api-version={config.AI_SEARCH_API_VERSION}",
        headers={
            "Content-Type": "application/json",
            "api-key": config.AI_SEARCH_API_KEY,
        },
        data=json.dumps(data, cls=NumpyEncoder),
    )

    res.raise_for_status()
    if res.status_code == 200:
        print(f"{len(searchDocs)} documents added to the search index successfully.")


def search(query: str):
    data = json.dumps(
        {
            "search": query,
            "select": "title, content, url, author, language, date, excerpt, baseUrl",
            "facets": ["date", "language", "author", "baseUrl"],
            "vectorQueries": [
                {
                    "kind": "vector",
                    "vector": list(config.model.query_embed(query=query))[0],
                    "k": 50,
                    "fields": "vector",
                },
            ],
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
