import os
import sys
# https://docs.python-guide.org/writing/structure/#test-suite
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

import unittest
import Keywords


class KeywordsUnitTest(unittest.TestCase):

    def test_ai(self):
        self.assertFalse(Keywords.match('Human Intelligence')[0])
        self.assertTrue(Keywords.match('Artificial Intelligence are things that fail where')[0])

    def test_ml(self):
        self.assertTrue(Keywords.match("Machine Learning like Bert")[0])
        self.assertFalse(Keywords.match("Machine Learning")[0])

    def test_nlp(self):
        self.assertTrue(Keywords.match("nlp bla bla rule-engine")[0])
        self.assertTrue(Keywords.match("nlp bla bla rule engine")[0])

    def test_pofalla(self):
        self.assertTrue(Keywords.match("pofalla-wende")[0])
        self.assertTrue(Keywords.match("Pofalla-Wende")[0])
        self.assertTrue(Keywords.match("Pofalla Wende")[0])


if __name__ == '__main__':
    unittest.main()
