"""Structural, localization, and publication checks for the Kazakh course."""
import json
from pathlib import Path
import re
import unittest
import xml.etree.ElementTree as ET

from tools.course import prepare_kazakh_language

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/kazakh-language"
MAPS = ROOT / "web/static/course/kazakh-language"
LESSON_STEMS = tuple(stem for stem, _, _ in prepare_kazakh_language.ROUTE)
MAP_STEMS = (
    "01-first-contact", "02-nine-sounds", "03-harmony", "04-sentence",
    "05-introduction", "06-questions", "07-repair", "08-mastery",
)


class KazakhLanguageReleaseTest(unittest.TestCase):
    def test_architecture_has_ten_blocks_and_eighty_lessons(self):
        data = json.loads((ROOT / "course/kazakh-language/curriculum.json").read_text())
        self.assertEqual(data["publication_policy"], "complete-block-only")
        self.assertEqual(len(data["blocks"]), 10)
        self.assertEqual(data["blocks"][-1]["range"], "73-80")
        self.assertEqual(len(data["first_block_lessons"]), 8)
        self.assertEqual(data["blocks"][0]["status"], "ready")
        self.assertEqual(data["blocks"][3]["lexical_target"], 800)
        self.assertEqual(data["blocks"][6]["lexical_target"], 1783)
        self.assertEqual(data["blocks"][-1]["lexical_target"], 2229)

    def test_first_block_is_complete_in_three_languages(self):
        headings = {"ru": "## Задание", "kz": "## Тапсырма", "en": "## Exercise"}
        minimum = {"ru": 600, "kz": 520, "en": 680}
        forbidden = {
            "kz": ("## Задание", "Следующий урок", "Проверка переноса"),
            "en": ("## Задание", "Келесі сабақ", "Проверка переноса"),
        }
        for index, stem in enumerate(LESSON_STEMS):
            for lang in ("ru", "kz", "en"):
                suffix = "" if lang == "ru" else f"-{lang}"
                path = LESSONS / f"{stem}{suffix}.md"
                self.assertTrue(path.is_file(), path)
                text = path.read_text(encoding="utf-8")
                self.assertIn("data-speak-kz=", text, path)
                self.assertGreaterEqual(text.count("data-speak-kz="), 5, path)
                self.assertIn("```kazakh\n", text, path)
                self.assertIn(headings[lang], text, path)
                self.assertGreaterEqual(len(text.split()), minimum[lang], path)
                self.assertNotRegex(text, r"\b(TODO|TBD)\b")
                for phrase in forbidden.get(lang, ()):
                    self.assertNotIn(phrase, text, path)
                if index + 1 < len(LESSON_STEMS):
                    next_slug = LESSON_STEMS[index + 1]
                    self.assertIn(f"/read/kazakh-language-{next_slug}?lang={lang}", text, path)
                else:
                    self.assertIn(f"/course/kazakh-language?lang={lang}", text, path)

    def test_all_first_block_maps_are_valid_and_fit_the_canvas(self):
        for stem in MAP_STEMS:
            for lang in ("ru", "kz", "en"):
                path = MAPS / f"map-{stem}-{lang}.svg"
                root = ET.parse(path).getroot()
                self.assertEqual(root.attrib.get("viewBox"), "0 0 1600 900", path)
                text = path.read_text(encoding="utf-8")
                self.assertNotIn("TODO", text, path)
                self.assertEqual(text.count('filter="url(#shadow)"'), 4, path)
                self.assertFalse(
                    re.search(r'<text[^>]*x="(?:1[6-9]\d\d|[2-9]\d{3})"', text), path
                )

    def test_atomic_sql_contains_the_full_route_and_expected_hashes(self):
        sql, expected = prepare_kazakh_language.prepare()
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.rstrip().endswith("COMMIT;"))
        self.assertEqual(sql.count("INSERT INTO articles("), 8)
        self.assertEqual(sql.count("INSERT INTO article_translations("), 24)
        self.assertEqual(sql.count("INSERT INTO article_series_items("), 8)
        self.assertEqual(len(expected), 24)
        self.assertEqual(
            [position for _, _, position in prepare_kazakh_language.ROUTE],
            list(range(10, 81, 10)),
        )
        self.assertIn("'kazakh-language'", sql)
        self.assertIn("'kazakh'", sql)


if __name__ == "__main__":
    unittest.main()
