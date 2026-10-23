# Chemistry / Химия / Химия

The first block has five lessons and fifteen published language versions.
The second block (particles and atoms) has eight lessons and twenty-four
published language versions.
The third block (elements and the periodic table, lessons 14–22) has nine
prepared lessons and twenty-seven language versions.
`curriculum.json` sketches the 100-lesson route. See `docs/chemistry-course.md`
for teaching method, evidence limits and safety rules.

Generate and validate from the repository root:

```sh
python3 tools/course/build_chemistry_block.py
python3 tools/course/build_chemistry_particles.py
python3 tools/course/build_chemistry_elements.py
python3 -m unittest tools.course.test_chemistry_release
python3 tools/course/prepare_chemistry.py --sql /tmp/chemistry-course.sql --expected /tmp/chemistry-expected.json
python3 tools/course/prepare_chemistry_particles.py --sql /tmp/chemistry-particles.sql --expected /tmp/chemistry-particles-expected.json
python3 tools/course/prepare_chemistry_elements.py --sql /tmp/chemistry-elements.sql --expected /tmp/chemistry-elements-expected.json
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
