import re

import Config
from model.Match import Match
from Utils import clean_content
from model.RSSFullTextResponse import RSSFullTextResponse


class Keywords:

    def __init__(self):
        self.keywords = []

    def init(self):
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
                            self.keywords.append(expr)
                            #print("Added keyword id", (len(self.keywords) - 1), "expr", expr.pattern)
                        except re.error:
                            raise Exception("Error parsing keyword regexp: " + one_line)

                line = keywords_file.readline()
            keywords_file.close()
        except FileNotFoundError:
            print("Keywords file", Config.KEYWORDS_FILE, "not found, terminating")
            exit(1)

        assert len(self.keywords) > 0

        print(len(self.keywords), "keywords initialized")

    def match(self, news: RSSFullTextResponse) -> Match:
        news_text = ""
        news_text = news_text + " {}".format(news.excerpt)
        news_text = news_text + " {}".format(news.title)
        news_text = news_text + " {}".format(news.content)
        news_text = clean_content(news_text)
        news_text = news_text.lower()
        news_text = news_text[:Config.MAX_ARTICLE_LENGTH]

        print("Analysing news...")
        for i in range(len(self.keywords)):
            if i % 25. == 0:
                print("... keywords >= " + str(i))

            keyword = self.keywords[i]

            if keyword.search(news_text):
                m = Match(kw_id=str(i),
                          kw_pattern=keyword.pattern,
                          news=news)
                print("Match", m)
                return m

        return None


if __name__ == "__main__":
    print("JAO")
    n = News()
    n.content = "sydas red hat developer roundup asdas"
    n.title = "BEST TITLE"
    n.url = "https://www.tagesschau.de/newsticker/liveblog-ukraine-mittwoch-121.html"

    print(match(n))
