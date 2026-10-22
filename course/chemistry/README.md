# Chemistry / Химия / Химия

The first block has five lessons and fifteen published language versions.
The second block (particles and atoms) has eight lessons and twenty-four
language versions prepared for release.
`curriculum.json` sketches the 100-lesson route. See `docs/chemistry-course.md`
for teaching method, evidence limits and safety rules.

Generate and validate from the repository root:

```sh
python3 tools/course/build_chemistry_block.py
python3 tools/course/build_chemistry_particles.py
python3 -m unittest tools.course.test_chemistry_release
python3 tools/course/prepare_chemistry.py --sql /tmp/chemistry-course.sql --expected /tmp/chemistry-expected.json
python3 tools/course/prepare_chemistry_particles.py --sql /tmp/chemistry-particles.sql --expected /tmp/chemistry-particles-expected.json
```

The prepare command never writes to production. Back up the database before
using its guarded transaction, deploy static assets, then verify every live
lesson, language and course navigation. The second release transaction guards
the first block, article authors and occupied course positions. Lessons must
not enter the news Telegram queue. Chemistry has no audio tracks yet; confirm
that no lesson references a nonexistent track before publication.
