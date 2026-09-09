"""The answer to lesson 19's exercise, English: the same answer, read by a parser.

Every value is taken by the name of its tag rather than by cutting the string,
and the price of one unit is the rate divided by quant -- the bank quotes some
currencies per ten, per hundred or per thousand, and a comparison that ignores
that is wrong by two orders of magnitude.
"""

import urllib.request
import xml.etree.ElementTree as ET

URL = "https://nationalbank.kz/rss/get_rates.cfm?fdate=15.01.2026"
WANTED = ("USD", "RUB", "UZS")

request = urllib.request.Request(URL, headers={"User-Agent": "shanraq-course/1.0"})
with urllib.request.urlopen(request, timeout=30) as answer:
    charset = answer.headers.get_content_charset() or "utf-8"
    root = ET.fromstring(answer.read().decode(charset))

print("date:", root.findtext("date"), "| currencies in the answer:", len(root.findall("item")))

for item in root.findall("item"):
    code = item.findtext("title")
    if code not in WANTED:
        continue
    rate = float(item.findtext("description"))
    quant = int(item.findtext("quant"))
    print(f"{code}: {rate} per {quant} → {rate / quant:.4f} tenge for one")

first = root.find("item")
print("the answer has no author:", first.findtext("author", "no author tag"))
