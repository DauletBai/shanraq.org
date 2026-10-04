# Kazakh language course sources

The course is planned as ten blocks of eight lessons. A block becomes visible on
the site only after all eight lessons exist in Russian, Kazakh, and English and
its support maps, exercises, links, and mobile layout pass review.

Current state: the full 80-lesson curriculum is fixed and the first six blocks,
48 lessons, are published in Russian, Kazakh, and English. School and City,
Daily Life and Services, Space and Cases, and Action and Time (lessons 17–48)
have complete manually localized pages, 144 localized support maps, 4K covers,
exercises, mastery gates, and 273 approved Kazakh recordings.

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
python3 tools/course/generate_kazakh_language_space_cases.py
python3 tools/course/generate_kazakh_language_space_cases_maps.py
python3 tools/course/generate_kazakh_language_action_time.py
python3 tools/course/generate_kazakh_language_action_time_maps.py
python3 tools/course/prepare_kazakh_language.py --sql /tmp/kazakh-language.sql --expected /tmp/kazakh-language.json
python3 -m unittest tools.course.test_kazakh_language_release
```

The generators contain explicit human-authored localization. They do not call a
translation service during the build.
