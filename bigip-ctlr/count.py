#!/usr/bin/python3

import sys
import os


def is_code(x):
    return x.endswith(".go")  and "/vendor/" not in x

def count_lines(x):
    with open(x, 'r') as file:
        line_count = sum(1 for line in file)
    return line_count

directory = os.path.dirname(os.path.abspath(__file__))
file_list = []
for root, dirs, files in os.walk(directory):
    for file in files:
        file_list.append(os.path.join(root, file))

code_file_list = list(filter(is_code, file_list))

lines = 0
for file in code_file_list:
    print(file)
    lines += count_lines(file)

print("total files:", len(code_file_list), "  total lines:", lines)
