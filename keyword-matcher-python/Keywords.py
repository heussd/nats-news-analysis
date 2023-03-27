import re
from re import Pattern
from typing import List

from pydantic import BaseModel

import Config


class KeywordEntry(BaseModel):
    id: str
    text: str
    regexp: Pattern

    class Config:
        arbitrary_types_allowed = True


clean_up_regexes = [
    re.compile("[^a-zA-Z]"),
    re.compile("\\s\\S\\s"),
    re.compile("\\s\\s+")
]


class Keywords:

    def __init__(self):
        self.keywords: List[KeywordEntry] = []

    def human_readable(self, regex: str):
        s = regex
        for r in clean_up_regexes:
            s = r.sub(" ", s, 0)

        return s.strip()

    def init(self):
        try:
            keywords_file = open(Config.KEYWORDS_FILE, "r")
            line = keywords_file.readline()
            while line:
                one_line = line.strip()

                if not one_line == "":
                    if not one_line.startswith("#"):
                        try:
                            expr = re.compile(one_line)
                            self.keywords.append(KeywordEntry(
                                regexp=expr,
                                id=self.human_readable(one_line),
                                text=one_line
                            ))
                        except re.error:
                            raise Exception("Error parsing keyword regexp: " + one_line)

                line = keywords_file.readline()
            keywords_file.close()
        except FileNotFoundError:
            print("Keywords file", Config.KEYWORDS_FILE, "not found, terminating")
            exit(1)

        assert len(self.keywords) > 0

        print(len(self.keywords), "keywords initialized")

    def match(self, news: str) -> (bool, str):
        for i in range(len(self.keywords)):
            keyword = self.keywords[i]

            if keyword.regexp.match(news) is not None:
                return True, keyword.id

        return False, ""


if __name__ == "__main__":
    content = "Delicious dark-chocolate pineapple pies"

    keywords = Keywords()
    keywords.init()
    print(keywords.match(content))

