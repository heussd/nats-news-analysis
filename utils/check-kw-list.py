#!/usr/bin/python3
import re
import sys

with open(sys.argv[1]) as f:
    lines = f.readlines()
    for line in lines:
        print("Parsing", line.strip(), end=' ')
        re.compile(line)
        print("... OK")
