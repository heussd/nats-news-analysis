import re

import Config
from model.Match import Match
from Utils import clean_content
from model.RSSFullTextResponse import RSSFullTextResponse


MAX_ARTICLE_LENGTH = 6000

keywords = []

try:
    keywords_file = open(Config.KEYWORDS_FILE, "r")
    line = keywords_file.readline()
    while line:
        one_line = line.strip()
        one_line = one_line.lower()

        if not one_line == "":
            if not one_line.startswith("#"):
                try:
                    expr = re.compile(one_line)
                    keywords.append(expr)
                    print("Added keyword id", (len(keywords) - 1), "expr", expr.pattern)
                except re.error:
                    raise Exception("Error parsing keyword regexp: " + one_line)

        line = keywords_file.readline()
    keywords_file.close()
except FileNotFoundError:
    print("Keywords file", Config.KEYWORDS_FILE, "not found, terminating")
    exit(1)

assert len(keywords) > 0

print("Keywords initialized")


def match(news: RSSFullTextResponse) -> Match:
    news_text = news.excerpt + news.content
    news_text = clean_content(news_text)
    news_text = news.title + "\n\n" + news_text
    news_text = news_text.lower()
    #news_text = news_text[:MAX_ARTICLE_LENGTH]

    for i in range(len(keywords)):
        keyword = keywords[i]

        if keyword.search(news_text):
            m = Match(kw_id=str(i),
                      kw_pattern=keyword.pattern,
                      news=news)
            print(m)
            return m

    return None


if __name__ == "__main__":
    print("JAO")
    n = News()
    n.content = "sydas red hat developer roundup asdas"
    n.title = "BEST TITLE"
    n.url = "https://www.tagesschau.de/newsticker/liveblog-ukraine-mittwoch-121.html"

    print(match(n))

