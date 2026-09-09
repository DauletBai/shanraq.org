"""The answer to lesson 18's exercise, Kazakh: the dollar rate for one fixed day.

The day is fixed on purpose: a rate for yesterday changes every night, and a
lesson cannot promise an output that does. Everything else is what the exercise
is about -- a request with a header of its own, a timeout, the status read
before the body, and the bytes decoded with the charset the server named.

The rate is pulled out with string methods, and that is deliberately the crude
way: the next lesson takes the same answer apart with a parser that knows what
a tag is.
"""

import urllib.error
import urllib.request

TEMPLATE = "https://nationalbank.kz/rss/get_rates.cfm?fdate={}"
URL = TEMPLATE.format("15.01.2026")

request = urllib.request.Request(URL, headers={"User-Agent": "shanraq-course/1.0"})
try:
    with urllib.request.urlopen(request, timeout=30) as answer:
        status = answer.status
        charset = answer.headers.get_content_charset() or "utf-8"
        body = answer.read()
except urllib.error.HTTPError as error:
    raise SystemExit(f"сервер {error.code} {error.reason} деп жауап берді")

print("статус:", status, "| алынған байт:", len(body))

text = body.decode(charset)
print("жауаптағы күн:", text.split("<date>")[1].split("</date>")[0])

item = text.split("<title>USD</title>")[1]
rate = item.split("<description>")[1].split("</description>")[0]
change = item.split("<change>")[1].split("</change>")[0]
print("доллар:", rate, "теңге | күндегі өзгеріс:", change)
