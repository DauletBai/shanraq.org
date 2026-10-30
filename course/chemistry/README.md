# Chemistry / Химия / Химия

The first block has five lessons and fifteen published language versions.
The second block (particles and atoms) has eight lessons and twenty-four
published language versions.
The third block (elements and the periodic table, lessons 14–22) has nine
published lessons and twenty-seven language versions.
The fourth block (bonds and formulas, lessons 23–31) is published: nine lessons,
27 language versions, 27 localized maps and 27 original teaching diagrams.
`curriculum.json` sketches the 100-lesson route. See `docs/chemistry-course.md`
for teaching method, evidence limits and safety rules.

Generate and validate from the repository root:

```sh
python3 tools/course/build_chemistry_block.py
python3 tools/course/build_chemistry_particles.py
python3 tools/course/build_chemistry_elements.py
python3 tools/course/build_chemistry_bonds.py
python3 -m unittest tools.course.test_chemistry_release
python3 tools/course/prepare_chemistry.py --sql /tmp/chemistry-course.sql --expected /tmp/chemistry-expected.json
python3 tools/course/prepare_chemistry_particles.py --sql /tmp/chemistry-particles.sql --expected /tmp/chemistry-particles-expected.json
python3 tools/course/prepare_chemistry_elements.py --sql /tmp/chemistry-elements.sql --expected /tmp/chemistry-elements-expected.json
python3 tools/course/prepare_chemistry_bonds.py --sql /tmp/chemistry-bonds.sql --expected /tmp/chemistry-bonds-expected.json
```

The prepare command never writes to production. Back up the database before
using its guarded transaction, deploy static assets, then verify every live
lesson, language and course navigation. The second release transaction guards
the first block, article authors and occupied course positions. Lessons must
not enter the news Telegram queue. Chemistry has no audio tracks yet; confirm
that no lesson references a nonexistent track before publication.

The third release also guards the thirteen existing lessons and their 39
translations, updates lesson 13’s onward link, and attaches all nine new
lessons to the course before commit. Its expected manifest includes 30 pages:
the new 27 plus the previous lesson in three languages.

The fourth release guards 22 prior lessons and 66 translations, bridges lesson
22 and forces deferred triggers before checking that no new lesson entered
the news queue. Its manifest contains the 27 new pages and lesson 22 in all
three languages. Stage the transaction on a restored backup before production.

Block six (41–49) was published on 28 October 2026: moles, molar mass, entity counts, elemental
mass fractions, stoichiometry, limiting reactants and product yield. It adds
27 pages, 27 retrieval maps and 27 localized vector diagram covers. Run:

```sh
python3 tools/course/build_chemistry_quantities.py
python3 -m unittest tools.course.test_chemistry_release tools.course.test_chemistry_reactions tools.course.test_chemistry_quantities
python3 tools/course/prepare_chemistry_quantities.py --sql /tmp/chemistry-quantities.sql --expected /tmp/chemistry-quantities-expected.json
```

Its transaction requires the 40 existing lessons and 120 ready translations.
It attaches lessons before forcing publication triggers, checks the absence
of news jobs, and adds lesson 40's onward link in the same transaction.
Each translation receives its own localized cover. Numeric tests independently
check formula masses, conservation, excess remaining, purity and yield. Audio
is still unavailable; no missing audio links are introduced.

Release verification: 49 published lessons and 147 ready translations;
30 expected translation hashes matched; 87 live page, asset and course checks
passed, as did 30 sitemap URL checks. There were no chemistry news jobs.
Rehearsal used a schema-only production dump and synthetic author/course
fixtures, without copying production user data. Backup remains on the server.

Block five (32–40, reactions and conservation) adds 27 language pages, 27
retrieval maps and 27 original diagram covers. Regenerate and validate with:

```sh
python3 tools/course/build_chemistry_reactions.py
python3 -m unittest tools.course.test_chemistry_release tools.course.test_chemistry_reactions
python3 tools/course/prepare_chemistry_reactions.py --sql /tmp/chemistry-reactions.sql --expected /tmp/chemistry-reactions-expected.json
```

Release block four first. Block five requires 31 published lessons and 93 ready
translations, and bridges lesson 31 only inside its release transaction. The
on-disk lesson 31 remains safe to release without block five. The manifest has
30 pages. Tests independently count atoms in all reference equations, check
minimal coefficients, mass and energy examples, particle diagrams, localized
routes and assets, and the staged bridge. No lesson may enter news Telegram.

Block seven (50–58), water and solutions, was published on 30 October 2026 with 27 language pages,
27 localized retrieval maps and 27 original vector diagram covers. Regenerate:

```sh
python3 tools/course/build_chemistry_solutions.py
python3 -m unittest tools.course.test_chemistry_solutions
python3 tools/course/prepare_chemistry_solutions.py --sql /tmp/chemistry-solutions.sql --expected /tmp/chemistry-solutions-expected.json
```

Its release requires 49 prior lessons and 147 ready translations. It bridges
lesson 49 only inside the transaction, attaches all nine lessons before
publication constraints run, and rejects any news Telegram job. Numerical
tests check separate liquid/solid balances, percentages, millilitres, molarity,
unequal-mass mixing and final-solution volume. Only paper and supplied data
are used; no physical preparation or electrical experiments are assigned.

Release verification: 58 published lessons and 174 ready translations; exact
matches for the 30-page manifest; no news Telegram jobs for this release.
The release was rehearsed on a restored backup, including rollback on a
missing prior translation and an unchanged second run. All seven PR workflows
passed, including Go race tests and template snapshots. Live verification
passed 87 page/asset/course checks and 60 sitemap checks. Browser inspection
covered Russian desktop, English tablet and Russian/Kazakh mobile layouts,
including tables, navigation and answer disclosure. The database backup and
previous application image were retained on the server.
