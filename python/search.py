import json

from ai_search import search


if __name__ == "__main__":
    query = "deutsch"
    print(json.dumps(search(query), indent=2))
