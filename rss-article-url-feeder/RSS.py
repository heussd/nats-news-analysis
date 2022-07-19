import requests
from django.core.validators import URLValidator
from lxml import html

validate = URLValidator()


def retrieve_article_links(feedurl):
    article_urls = []

    try:
        page = requests.get(feedurl)
        tree = html.fromstring(page.content)
        links = tree.xpath('//item/link')

        for link in links:
            the_link = link.tail
            try:
                validate(the_link)
                article_urls.append(the_link)
            except:
                print("Invalid URL", the_link)

    except:
        print("Retrieval failed:", feedurl)

    return article_urls


if __name__ == "__main__":
    print(retrieve_article_links('https'))
