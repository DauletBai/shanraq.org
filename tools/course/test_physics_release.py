"""Checks that protect the first physics block's scientific and publishing contract."""
import json
from pathlib import Path
import re
import unittest
import xml.etree.ElementTree as ET

from tools.course.build_physics_block import DATA, LESSONS, MAPS, ROOT
from tools.course.prepare_physics import prepare, read_lesson


class PhysicsReleaseTests(unittest.TestCase):
    def test_curriculum_has_contiguous_route_and_complete_first_block(self):
        data = json.loads((ROOT / "course/physics/curriculum.json").read_text())
        self.assertEqual(data["lesson_count"], 100)
        self.assertEqual(data["publication_policy"], "complete-block-only")
        self.assertEqual(data["languages"], ["ru", "kz", "en"])
        next_number = 1
        for block in data["blocks"]:
            first, last = map(int, block["range"].split("-"))
            self.assertEqual(first, next_number)
            self.assertEqual(last - first + 1, block["lesson_count"])
            next_number = last + 1
        self.assertEqual(next_number, 101)
        self.assertEqual([row["number"] for row in data["first_block_lessons"]], list(range(1, 9)))

    def test_every_language_has_the_same_experiment_and_working_assets(self):
        for index, row in enumerate(DATA, 1):
            for lang in ("ru", "kz", "en"):
                suffix = "" if lang == "ru" else f"-{lang}"
                path = LESSONS / f"{row['stem']}{suffix}.md"
                self.assertTrue(path.is_file(), path)
                title, summary, body = read_lesson(path)
                self.assertTrue(title and summary)
                self.assertGreater(len(body.split()), 250, path)
                ref = re.search(r"!\[[^]]+\]\((/static/course/physics/[^)]+)\)", body)
                self.assertIsNotNone(ref, path)
                asset = ROOT / "web" / ref.group(1).lstrip("/")
                self.assertTrue(asset.is_file(), asset)
                ET.parse(asset)
                if index < 8:
                    next_slug = f"physics-{index+1:02d}-{DATA[index]['stem'][3:]}"
                    self.assertIn(f"/read/{next_slug}?lang={lang}", body)

    def test_experiment_numbers_and_limits_match_across_languages(self):
        # These are worked numbers, never represented as a learner's observation.
        short = [10.9, 11.1, 11.0]
        long = [15.5, 15.7, 15.6]
        self.assertAlmostEqual(sum(short) / len(short) / 10, 1.10)
        self.assertAlmostEqual(sum(long) / len(long) / 10, 1.56)
        self.assertAlmostEqual(1.56 - 1.10, 0.46)
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else f"-{lang}"
            text = (LESSONS / f"08-investigation-mastery{suffix}.md").read_text()
            values = ("10,9", "11,1", "11,0", "15,5", "15,7", "15,6", "1,10", "1,56", "0,46") if lang != "en" else (
                "10.9", "11.1", "11.0", "15.5", "15.7", "15.6", "1.10", "1.56", "0.46")
            for value in values:
                self.assertIn(value, text, (lang, value))
            self.assertIn("| 30 |", text)
            self.assertIn("| 60 |", text)

    def test_guarded_publication_is_atomic_and_only_first_block(self):
        sql, expected = prepare()
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.rstrip().endswith("COMMIT;"))
        self.assertIn("Physics slug belongs to another author", sql)
        self.assertIn("Physics lesson belongs to another course", sql)
        self.assertEqual(len(expected), 24)
        self.assertEqual({e["lang"] for e in expected}, {"ru", "kz", "en"})
        self.assertEqual({e["slug"] for e in expected}, {
            f"physics-{i:02d}-{row['stem'][3:]}" for i, row in enumerate(DATA, 1)})


if __name__ == "__main__":
    unittest.main()
