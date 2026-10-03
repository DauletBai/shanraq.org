# Kazakh language course sources

The course is planned as ten blocks of eight lessons. A block becomes visible on
the site only after all eight lessons exist in Russian, Kazakh, and English and
its support maps, exercises, links, and mobile layout pass review.

Current state: the full 80-lesson curriculum is fixed and the first two blocks,
16 lessons, are complete in Russian, Kazakh, and English. Their 48 localized
support maps, listening controls, exercises, mastery gates, block covers, and
atomic SQL publication are checked together.

## Regenerate and verify

```sh
python3 tools/course/generate_kazakh_language_foundations.py
python3 tools/course/generate_kazakh_language_maps.py
python3 tools/course/generate_kazakh_language_people_home.py
python3 tools/course/generate_kazakh_language_people_home_maps.py
python3 tools/course/prepare_kazakh_language.py --sql /tmp/kazakh-language.sql --expected /tmp/kazakh-language.json
python3 -m unittest tools.course.test_kazakh_language_release
```

The generators contain explicit human-authored localization. They do not call a
translation service during the build.
