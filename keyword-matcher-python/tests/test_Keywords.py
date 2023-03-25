import os
import sys

import Keywords
from Keywords import Keywords

# https://docs.python-guide.org/writing/structure/#test-suite
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))


import unittest


class KeywordsUnitTest(unittest.TestCase):

    def first(self, ret) -> bool:
        (flag, id) = ret
        return flag

    def testLocalIT(self):
        keywords = Keywords()
        keywords.init()

        self.assertTrue(self.first(keywords.match("Peach")))
        self.assertFalse(self.first(keywords.match("Pineapple")))
        self.assertFalse(self.first(keywords.match("Hamburger")))
        self.assertTrue(self.first(keywords.match("Apple")))
        self.assertTrue(self.first(keywords.match("Banana split")))
        self.assertTrue(self.first(keywords.match("Delicious pineapple recipes")))
        self.assertTrue(self.first(keywords.match("Delicious recipes")))
        self.assertTrue(self.first(keywords.match("Delicious pineapple pies")))
        self.assertTrue(self.first(keywords.match("Delicious dark-chocolate pineapple pies")))
        self.assertTrue(self.first(keywords.match("ICE cold cream")))
        self.assertFalse(self.first(keywords.match("ICE and also some cream")))
        self.assertFalse(self.first(keywords.match("whipped cream")))
        self.assertFalse(self.first(keywords.match("# Should not match")))

        self.assertFalse(self.first(keywords.match("Mister Cool")))
        self.assertFalse(self.first(keywords.match("Miss Gray")))
        self.assertTrue(self.first(keywords.match("Mississippi")))

        self.assertFalse(self.first(keywords.match("Bias")))
        self.assertTrue(self.first(keywords.match("as")))

        self.assertTrue(self.first(keywords.match("All of us")))
        self.assertTrue(self.first(keywords.match("All-of-us")))
        self.assertFalse(self.first(keywords.match("Alloofuus")))

        self.assertFalse(self.first(keywords.match("I drink cold beer. I eat hot pizza.")))
        self.assertTrue(self.first(keywords.match("I ate cold yummy pizza yesterday afternoon. I drink hot chocolate.")))

        self.assertTrue(self.first(keywords.match("The king lived long and prosper.")))
        self.assertTrue(self.first(keywords.match("Long live the king.")))
        self.assertTrue(self.first(keywords.match("The queen lived long and prosper.")))
        self.assertTrue(self.first(keywords.match("Long live the queen.")))

        self.assertFalse(self.first(keywords.match("Like king and queen.")))

    def test_human_readable(self):
        keywords = Keywords()
        keywords.init()

        self.assertEqual("delicious pie recipes", keywords.human_readable("(?i)(delicious).*(pie|recipes)"))

    def test_string_match_return(self):
        keywords = Keywords()
        keywords.init()

        (_, text) = keywords.match("A little Peach acc day")
        self.assertEqual("Apple peach", text)

        (_, text) = keywords.match("I like to eat delicious original organic-sourced pineapple pies twice a day")
        self.assertEqual("delicious pie recipes", text)

        (_, text) = keywords.match("Long live the queen. Something else")
        self.assertEqual("king queen long", text)


if __name__ == '__main__':
    unittest.main()
