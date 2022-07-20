import feedparser


def retrieve_article_links(feedurl):
    article_urls = []

    try:
        feed = feedparser.parse(feedurl)

        for entry in feed['entries']:
            article_urls.append(entry.link)

    except Exception as e:
        print("Retrieval failed:", feedurl)
        print(e)

    article_urls = list(filter(None, article_urls))
    article_urls = list(filter(lambda item: item.startswith('http'), article_urls))
    article_urls = [i.strip() for i in article_urls]

    if len(article_urls) == 0:
        print("WARNING: No article URLs found in feed", feedurl)
        if (feedurl.startswith("http://")):
            print("Retrying with https")
            return retrieve_article_links(feedurl.replace("http:", "https:"))

    return article_urls


if __name__ == "__main__":
    print(retrieve_article_links('https://www.hessenschau.de/index.rss'))
    print(retrieve_article_links('https://katapult-magazin.de/feed/rss'))
    'https://waylonwalker.com/rss.xml'

