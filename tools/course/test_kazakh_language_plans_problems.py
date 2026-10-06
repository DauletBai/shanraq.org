"""Checks for the published Kazakh Plans and Problems block."""
import json
from pathlib import Path
import re
import unittest
import xml.etree.ElementTree as ET

from tools.course import generate_kazakh_language_plans_problems as draft
from tools.course import prepare_kazakh_language


class PlansProblemsTest(unittest.TestCase):
    def test_eight_trilingual_lessons_are_complete_and_linked(self):
        self.assertEqual(len(draft.LESSONS), 8)
        self.assertEqual([int(row["stem"][:2]) for row in draft.LESSONS], list(range(49, 57)))
        manifest = json.loads((draft.AUDIO / "manifest.json").read_text(encoding="utf-8"))
        by_phrase = {row["phrase"]: row for row in manifest["recordings"]}
        for index, lesson in enumerate(draft.LESSONS):
            self.assertEqual(len(lesson["phrases"]), 6)
            for lang, minimum in (("ru", 470), ("kz", 440), ("en", 570)):
                suffix = "" if lang == "ru" else f"-{lang}"
                page = draft.OUT / f"{lesson['stem']}{suffix}.md"
                self.assertTrue(page.is_file(), page)
                body = page.read_text(encoding="utf-8")
                self.assertGreaterEqual(len(body.split()), minimum, page)
                self.assertEqual(body.count("](/static/course/kazakh-language/audio/"), 6, page)
                self.assertIn(f"map-{lesson['map']}-{lang}.svg", body, page)
                self.assertNotIn("#speak-kz=", body, page)
                self.assertNotRegex(body, r"\b(?:TODO|TBD)\b")
                for phrase, *_ in lesson["phrases"]:
                    self.assertIn(phrase, body, page)
                    self.assertIn(by_phrase[phrase]["file"], body, page)
                if index < 7:
                    self.assertIn(f"/read/kazakh-language-{draft.LESSONS[index+1]['stem']}?lang={lang}", body)
                else:
                    self.assertIn(f"/course/kazakh-language?lang={lang}", body)

    def test_maps_parse_and_match_their_explained_phrases(self):
        for lesson in draft.LESSONS:
            for lang in draft.LANGS:
                map_path = draft.MAPS / f"map-{lesson['map']}-{lang}.svg"
                root = ET.parse(map_path).getroot()
                self.assertEqual(root.attrib["viewBox"], "0 0 1600 900", map_path)
                text = map_path.read_text(encoding="utf-8")
                self.assertIn(lesson["title"][draft.LANGS.index(lang)], text, map_path)
                self.assertGreaterEqual(text.count('<rect x='), 5, map_path)
                self.assertGreaterEqual(text.count('filter="url(#shadow)"'), 5, map_path)
                if lesson["stem"].startswith("53-"):
                    self.assertIn('M446 440H520V365H608', text)
                    self.assertIn('M446 505H520V635H608', text)
                if lesson["stem"].startswith("51-"):
                    self.assertNotIn('M418 453h21', text)

    def test_approved_block_is_in_published_route(self):
        self.assertEqual(len(prepare_kazakh_language.ROUTE), 56)
        self.assertEqual(
            [stem for stem, _, _ in prepare_kazakh_language.ROUTE[-8:]],
            [lesson["stem"] for lesson in draft.LESSONS],
        )
        manifest = json.loads((draft.AUDIO / "manifest.json").read_text(encoding="utf-8"))
        published = [
            row for row in manifest["recordings"]
            if any(int(name[:2]) <= 56 for name in row["lessons"])
        ]
        self.assertTrue(published)
        self.assertTrue(all(row["status"] == "approved" for row in published))
        staged = [row for row in manifest["recordings"] if row["status"] == "candidate_review"]
        self.assertEqual(len(staged), 48)
        self.assertTrue(all(
            all(57 <= int(name[:2]) <= 64 for name in row["lessons"])
            for row in staged
        ))
        review = (draft.AUDIO / "review.html").read_text(encoding="utf-8")
        self.assertIn("уроков 57–64", review)
        self.assertEqual(review.count('<audio controls'), 48)


if __name__ == "__main__":
    unittest.main()
