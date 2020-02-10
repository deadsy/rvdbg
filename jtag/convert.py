#!/usr/bin/env python3
# convert the idcodes file taken from OpenOCD to Go code.

import re

fname = "jep106.inc"

def main():
  f = open(fname)
  x = f.readlines()
  f.close()
  for l in x:
    l = l.strip()
    if l.startswith("/*"):
      continue
    y, name = l.split("=")
    m = re.match(r"\[?(\d+)\]\[0x?([\dA-Fa-f]+) - 1\]", y)
    page = int(m.group(1), 16)
    code = int(m.group(2), 16)
    idcode = (page << 7) | code
    print("0x%03x: %s" % (idcode, name))

main()
