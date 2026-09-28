"""Structural and pedagogical checks for the mathematics course."""
import json
import re
import unittest
from pathlib import Path
import xml.etree.ElementTree as ET

import prepare_mathematics


ROOT = Path(__file__).resolve().parents[2]
CURRICULUM = ROOT / "course/mathematics/curriculum.json"


class MathematicsReleaseTests(unittest.TestCase):
    def test_curriculum_is_a_connected_dag(self):
        data = json.loads(CURRICULUM.read_text(encoding="utf-8"))
        modules = data["modules"]
        ids = [m["id"] for m in modules]
        self.assertEqual(len(ids), len(set(ids)))
        known = set(ids)
        for module in modules:
            self.assertTrue({"title_ru", "title_kz", "title_en"} <= module.keys())
            self.assertTrue(set(module["prerequisites"]) <= known)

        visiting, visited = set(), set()

        def walk(node):
            if node in visiting:
                self.fail(f"curriculum cycle at {node}")
            if node in visited:
                return
            visiting.add(node)
            current = next(m for m in modules if m["id"] == node)
            for dependency in current["prerequisites"]:
                walk(dependency)
            visiting.remove(node)
            visited.add(node)

        for node in ids:
            walk(node)
        self.assertEqual(visited, known)
        self.assertIn("capstone-modeling", known)

    def test_published_route_is_atomic_and_non_destructive(self):
        sql, expected = prepare_mathematics.prepare()
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.endswith("COMMIT;\n"))
        self.assertEqual(len(expected), 10)
        self.assertEqual(sql.count("INSERT INTO article_series_items("), 10)
        self.assertIn("'society','education'", sql)
        self.assertIn("'published','math'", sql)
        self.assertNotIn("DELETE FROM", sql)
        self.assertNotIn("TRUNCATE", sql)

    def test_every_lesson_is_a_full_learning_cycle(self):
        sections = (
            "## Где мы на карте",
            "## Опорный сигнал",
            "## Воспроизведите без подсказки",
            "## Найдите и исправьте ошибку",
            "## Задание",
        )
        for stem in prepare_mathematics.STEMS[1:-1]:
            path = prepare_mathematics.LESSONS / f"{stem}.md"
            text = path.read_text(encoding="utf-8")
            for heading in sections:
                self.assertIn(heading, text, path)
            self.assertIn("/static/course/mathematics/", text, path)
            self.assertGreaterEqual(len(re.findall(r"\b[\w/-]+\b", text)), 500, path)
        mastery = (prepare_mathematics.LESSONS / "09-mastery.md").read_text(encoding="utf-8")
        self.assertIn("8 из 10", mastery)
        self.assertIn("Через семь дней", mastery)

    def test_each_lesson_uses_its_own_map(self):
        maps = []
        for stem in prepare_mathematics.STEMS[1:]:
            text = (prepare_mathematics.LESSONS / f"{stem}.md").read_text(encoding="utf-8")
            found = re.findall(r"/static/course/mathematics/([^\s)]+\.svg)", text)
            self.assertEqual(len(found), 1, stem)
            maps.extend(found)
        self.assertEqual(len(maps), len(set(maps)))
        for name in maps:
            self.assertTrue((ROOT / "web/static/course/mathematics" / name).is_file(), name)
            ET.parse(ROOT / "web/static/course/mathematics" / name)
        ET.parse(ROOT / "web/static/course/mathematics/map-full-ru.svg")

    def test_claims_do_not_turn_targets_into_results(self):
        preface = (prepare_mathematics.LESSONS / "preface.md").read_text(encoding="utf-8")
        self.assertIn("не станем выдавать", preface)
        self.assertIn("Числа приблизительны", preface)
        self.assertNotIn("КПД образования — 10%", preface)
        self.assertNotIn("95% успеваемости", preface)

    def test_sql_delimiter_is_rejected(self):
        with self.assertRaises(ValueError):
            prepare_mathematics.literal("bad $mathematics_course$ value")


if __name__ == "__main__":
    unittest.main()
