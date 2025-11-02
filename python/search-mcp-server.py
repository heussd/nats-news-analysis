import json

from fastmcp import FastMCP

import ai_search

mcp = FastMCP("News Analysis Search")

@mcp.tool
def search(query: str) -> str:
    """
        Multi-language, semantic search for news articles related to the query.

        Args:
            query (str): The search query. Can be in any language. Search will be keyword- and embedding-based.
    """
    return json.dumps(ai_search.search(query), indent=0)

if __name__ == "__main__":
    mcp.run(transport="streamable-http", host="0.0.0.0")
