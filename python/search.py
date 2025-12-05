import json

from ai_search import search


if __name__ == "__main__":
    query = "llm"
    print(json.dumps(search(query, top=200, baseUrl="github.com"), indent=2))
