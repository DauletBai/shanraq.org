# Kazakh language course sources

The course is planned as ten blocks of eight lessons. A block becomes visible on
the site only after all eight lessons exist in Russian, Kazakh, and English and
its support maps, exercises, links, and mobile layout pass review.

Current state: the full 80-lesson curriculum is fixed and the first two blocks,
16 lessons, are published in Russian, Kazakh, and English. Blocks 3–4, School
and City and Daily Life and Services (lessons 17–32), have complete manually
localized pages, 48 localized support maps, 4K covers, exercises, and mastery
gates. Their 95 new Kazakh recordings remain a review set until the owner
approves them; the atomic publisher therefore still contains only the first 16
lessons.

## Regenerate and verify

```sh
python3 tools/course/generate_kazakh_language_foundations.py
python3 tools/course/generate_kazakh_language_maps.py
python3 tools/course/generate_kazakh_language_people_home.py
python3 tools/course/generate_kazakh_language_people_home_maps.py
python3 tools/course/generate_kazakh_language_school_city.py
python3 tools/course/generate_kazakh_language_school_city_maps.py
python3 tools/course/generate_kazakh_language_day_services.py
python3 tools/course/generate_kazakh_language_day_services_maps.py
python3 tools/course/prepare_kazakh_language.py --sql /tmp/kazakh-language.sql --expected /tmp/kazakh-language.json
python3 -m unittest tools.course.test_kazakh_language_release
```

The generators contain explicit human-authored localization. They do not call a
translation service during the build.
