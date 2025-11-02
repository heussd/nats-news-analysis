import json

from ai_search import search


if __name__ == "__main__":
    query = "space"
    print(json.dumps(search(query), indent=2))
