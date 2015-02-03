#!/usr/bin/env python

import fileinput
import json

if __name__ == '__main__':
    for i, line in enumerate(fileinput.input()):
        _ = json.loads(line)
        if i % 1000000 == 0:
            print(i)
