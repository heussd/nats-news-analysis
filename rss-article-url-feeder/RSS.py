from xml.etree import ElementTree
import requests


def retrieve_article_links(feedurl):
    article_urls = []

    try:
        page = requests.get(feedurl)
        assert page.ok

        tree = ElementTree.fromstring(page.content)
        links = tree.findall('.//item/link')

        for link in links:
            article_urls.append(link.text)
            article_urls.append(link.tail)

    except Exception as e:
        print("Retrieval failed:", feedurl)
        print(e)

    article_urls = list(filter(None, article_urls))
    article_urls = list(filter(lambda item: item.startswith('http'), article_urls))
    article_urls = [i.strip() for i in article_urls]

    if len(article_urls) == 0:
        print("WARNING: No article URLs found in feed", feedurl)

    return article_urls


if __name__ == "__main__":
    print(retrieve_article_links('https://hnrss.org/newest?count=50'))

