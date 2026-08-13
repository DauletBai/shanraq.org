# Editorial and crawling policy

Who writes what on Shanraq.org, and which crawlers are allowed to read it.
This is an operational document: the rules below are enforced in code, and each
section names the switch that enforces it.

## Two kinds of article

**Authored work.** Written by a person, published under their name. This is the
site: 18 pieces as of August 2026, and everything published from here on.

**AI columns.** 90 pieces published between July and early August 2026 under
the "AI Dake" byline, labelled *Мнение ИИ / ИИ пікірі / AI opinion*. They are a
closed set — no more are produced.

## What happened to the AI columns

On 2026-08-06 all 90 were taken out of the search index. Ninety machine-written
columns against eighteen authored pieces is the profile of a content farm, and
that is what a search engine was being shown.

They were **not deleted**, and deleting them is refused. Ninety URLs would start
answering 404, nothing could be undone, and a column that later turns out to be
worth keeping would be gone. The pages stay published, readable and linkable.

The switch is one column: `articles.indexable`. Setting it false puts a page
behind a `noindex` tag, out of `sitemap.xml`, out of the news sitemap, out of
`llms.txt`, and into the AI-crawler exclusion list. Setting it true undoes all
five at once. There is no second flag to forget.

## Why the AI columns are also closed to AI crawlers

`noindex` is a directive to *search engines*. AI crawlers do not act on it, and
`robots.txt` said `Allow: /` to everything, so assistants went on reading the
columns after the de-indexing and could quote them back with this site's name
attached as the source.

That is worse than being ranked badly. A search engine ranks us; an assistant
**cites** us, and a citation is an assertion of authority. The *Мнение ИИ* badge
protects a human reader, who can see it — but it is page furniture, and it does
not reliably survive into the chunk a model retrieves. So the labelling that
makes the columns honest to a reader does not make them honest to a model.

And it cannot be undone. A web page can be corrected or withdrawn; a claim that
has entered an answer engine's index or a training corpus cannot be recalled.
Weighed against that, the cost of being wrong in the other direction — some lost
crawler traffic, which is not readers and not revenue — is small.

Three mechanisms, because no single one is honoured by everybody:

| Mechanism | Reaches |
|---|---|
| `<meta name="robots" content="noindex, follow">` | search engines that parse HTML |
| `X-Robots-Tag: noindex, follow, noai, noimageai` | crawlers that read headers only, and HEAD requests |
| `robots.txt` group naming ~22 AI agents | crawlers that fetch robots.txt and nothing else |

The `robots.txt` group repeats the private paths (`/admin`, `/studio`, `/api/`,
`/jobs`) from the wildcard group. A crawler obeys exactly one group — the most
specific that names it — and stops reading `User-agent: *` once it finds itself,
so the repeats are what keep those paths closed to the agents named.

If the query that builds the exclusion list fails, `/read` is closed to AI
crawlers **entirely** rather than a short list being served. Hiding the authored
articles for one crawl costs a little reach and is undone by the next fetch;
exposing one machine-written column cannot be undone at all.

## What is offered instead

`/llms.txt` lists only the authored articles, with dates, summaries, the
citation format and an explicit request not to quote the machine-written
columns. Every authored article carries a ready-made reference with a copy
button. The reasoning is not politeness: attribution follows the path of least
resistance, and a reference someone has to assemble themselves usually becomes
"источник — интернет".

## Forecasts

Forecasts made in articles go on the ledger at `/predictions` with the date they
were made, and are judged there in public — hit, partial or miss. The rules that
keep it honest:

- The wording of a forecast is never edited after it is judged.
- A partial hit counts as half. Scoring it as a success flatters the number;
  scoring it as a failure punishes admitting the middle ground, which is where
  most forecasts land.
- Open forecasts are excluded from the score entirely, so an unjudged ledger
  reads "nothing judged yet" rather than 0%.
- An open forecast past its own deadline is flagged in the admin. Never judging
  anything is the way a ledger like this gets quietly gamed.
- Deletion exists for a duplicate or a typo. Deleting a forecast because it went
  badly empties the ledger of the only thing it is for, and the confirmation
  dialog says so.
