import json

from fastmcp import FastMCP

import ai_search

mcp = FastMCP("News Desk")


@mcp.tool
def search(query: str, top: int = 10, baseUrl: str = None) -> str:
    """
    Find related news to a query with a multi-language, semantic search. Do not search for generic terms such as "latest technology trends 2025", but search for specific technologies, keywords, names, events, etc.

    Args:
        query (str): The search query. Can be in any language. Search will be keyword- and embedding-based.
        top (int): Optional: The number of top results to return (default 10).
        baseUrl (str): Optional: The base URL to filter the search results. If baseUrl is provided, search result order will switch from relevance based (default) to recency based, meaning that the latest things related to the query will be returned. This is useful for finding the latest news on a specific topic from a specific source, such as github.com, twitter.com, etc.
    """
    return json.dumps(ai_search.search(query, top=top, baseUrl=baseUrl), indent=0)


if __name__ == "__main__":
    mcp.run(transport="streamable-http", host="0.0.0.0")
