import requests
from lxml import html


def retrieve_article_links(feedurl):
    article_urls = []

    try:
        page = requests.get(feedurl)
        tree = html.fromstring(page.content)
        links = tree.xpath('//item/link')

        for link in links:
            article_urls.append(link.tail)

    except:
        print("Retrieval failed:", feedurl)

    return article_urls


if __name__ == "__main__":
    print(retrieve_article_links('https://hnrss9999.org/newest?count=50'))
