# Kazakh language course sources

The course is planned as ten blocks of eight lessons. A block becomes visible on
the site only after all eight lessons exist in Russian, Kazakh, and English and
its support maps, exercises, links, and mobile layout pass review.

Current state: the full 80-lesson curriculum is fixed, the first lesson of block
one is implemented in all three languages, and all eight block-one map templates
are available for content development. The public footer link remains disabled
until block one is complete.

## Regenerate and verify

```sh
python3 tools/course/generate_kazakh_language_foundations.py
python3 tools/course/generate_kazakh_language_maps.py
python3 -m unittest tools.course.test_kazakh_language_release
```

The generators contain explicit human-authored localization. They do not call a
translation service during the build.
