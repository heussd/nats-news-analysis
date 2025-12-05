import json
import os

from fastmcp import FastMCP
import requests

import ai_search

mcp = FastMCP("News Desk")


@mcp.tool()
async def news_start_page() -> str:
    """
        Shows a news start page as RSS feed. A nice entry point for all kinds of news interests.
        Use `retrieve_article_full_text` to retrieve the full text for articles.
        Use `search` to find related news articles based on keywords or topics.

    Args:
        None
    """
    url = os.getenv(
        "RAINDROP_RSS_FEED", "https://www.abc.net.au/news/feed/2942460/rss.xml"
    )
    spoiler = requests.get(url, verify=False)
    spoiler.raise_for_status()
    data = spoiler.text
    return data


@mcp.tool
def search(query: str, top: int = 10, baseUrl: str = None) -> str:
    """
    Find related news to a query with a multi-language, semantic search. Do not search for generic terms such as "latest technology trends 2025", but search for specific technologies, keywords, names, events, etc.

    Args:
        query (str): The search query. Can be in any language. Search will be keyword- and embedding-based.
        top (int): The number of top results to return (default 10).
    """
    return json.dumps(ai_search.search(query, top=top), indent=0)


@mcp.tool
def latest_on(query: str, top: int = 200, baseUrl: str = None) -> str:
    """
    Find the latest things related to a query, based on a baseUrl such as github.com. Do not search for generic terms such as "latest technology trends 2025", but search for specific technologies, keywords, names, events, etc.

    Args:
        query (str): The search query. Can be in any language. Search will be keyword- and embedding-based.
        top (int): The number of top results to return (default 200).
        baseUrl (str): The base URL to filter the search results.
    """
    return json.dumps(ai_search.search(query, top=top, baseUrl=baseUrl), indent=0)


@mcp.tool()
async def retrieve_article_full_text(url: str) -> any:
    """
        Retrieves the full text content and additional metadata of an article from the internet.

    Args:
        URL: The URL of the article to retrieve the full text of.
    """
    try:
        response = requests.get(
            f"http://fullfeedrss:80/extract.php", params={"url": url}
        )
        response.raise_for_status()
        data = response.json()
        return data
    except Exception as e:
        from readabilipy import simple_json_from_html_string

        req = requests.get(url)
        article = simple_json_from_html_string(req.text, use_readability=True)

        return article


if __name__ == "__main__":
    mcp.run(transport="streamable-http", host="0.0.0.0")
