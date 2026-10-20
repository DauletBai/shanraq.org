# Chemistry / Химия / Химия

The complete first block has five lessons and fifteen language versions.
`curriculum.json` sketches the 100-lesson route. See `docs/chemistry-course.md`
for teaching method, evidence limits and safety rules.

Generate and validate from the repository root:

```sh
python3 tools/course/build_chemistry_block.py
python3 -m unittest tools.course.test_chemistry_release
python3 tools/course/prepare_chemistry.py --sql /tmp/chemistry-course.sql --expected /tmp/chemistry-expected.json
```

The prepare command never writes to production. Back up the database before
using its guarded transaction, deploy static assets, then verify every live
lesson, language and course navigation. Lessons must not enter the news
Telegram queue.
