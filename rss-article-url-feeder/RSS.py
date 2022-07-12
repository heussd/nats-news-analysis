import requests
from lxml import html


def retrieve_article_links(feedurl):
    article_urls = []
    page = requests.get(feedurl)
    tree = html.fromstring(page.content)
    links = tree.xpath('//item/link')

    for link in links:
        article_urls.append(link.tail)

    return article_urls


if __name__ == "__main__":
    print(retrieve_article_links('https://hnrss.org/newest?count=50'))
