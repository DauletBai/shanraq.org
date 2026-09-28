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

    def test_course_map_direction_and_percent_cells_are_explicit(self):
        maps = ROOT / "web/static/course/mathematics"
        course_map = (maps / "map-full-ru.svg").read_text(encoding="utf-8")
        self.assertEqual(course_map.count('class="course-step"'), 8)
        self.assertEqual(course_map.count('class="course-arrow"'), 7)
        self.assertNotIn("M1135 465", course_map)

        percent_map = (maps / "map-07-percent.svg").read_text(encoding="utf-8")
        self.assertEqual(percent_map.count('class="percent-cell'), 100)
        self.assertEqual(percent_map.count("percent-cell--filled"), 37)
        self.assertEqual(percent_map.count("percent-cell--final-seven"), 7)
        self.assertIn("30 + 7 = 37", percent_map)

    def test_lesson_maps_encode_the_examples_they_explain(self):
        maps = ROOT / "web/static/course/mathematics"
        lessons = prepare_mathematics.LESSONS

        diagnostic = (maps / "map-01-diagnostic.svg").read_text(encoding="utf-8")
        self.assertIn('data-diagnostic-route="12|3/8|2:3|25%|x=5"', diagnostic)

        expected = (
            ("02-fraction-meaning.md", "map-02-fraction.svg", "3/8", 'data-fraction="3/8"'),
            ("03-equivalent-fractions.md", "map-03-equivalence.svg", "1/2", 'data-equivalence="1/2=2/4"'),
            ("04-compare-fractions.md", "map-04-compare.svg", "2/3", 'data-comparison="1/2&lt;2/3"'),
            ("05-fraction-operations.md", "map-05-operations.svg", "5/6", 'data-equation="1/2+1/3=5/6"'),
            ("06-ratio.md", "map-06-ratio.svg", "4:6", 'data-ratios="2:3=4:6"'),
            ("07-percent.md", "map-07-percent.svg", "37/100", 'data-fraction="37/100"'),
            ("08-proportion.md", "map-08-proportion.svg", "6/9", 'data-equivalence="2/3=6/9"'),
        )
        for lesson_name, map_name, example, marker in expected:
            lesson = (lessons / lesson_name).read_text(encoding="utf-8")
            illustration = (maps / map_name).read_text(encoding="utf-8")
            self.assertIn(example, lesson, lesson_name)
            self.assertIn(marker, illustration, map_name)

        fraction = ET.parse(maps / "map-02-fraction.svg").getroot()
        namespace = {"svg": "http://www.w3.org/2000/svg"}
        spokes = fraction.find(".//svg:g[@class='fraction-spokes']", namespace)
        ticks = fraction.find(".//svg:g[@class='eighth-ticks']", namespace)
        self.assertEqual(len(spokes.findall("svg:path", namespace)), 8)
        self.assertEqual(len(ticks.findall("svg:path", namespace)), 9)

        comparison = ET.parse(maps / "map-04-compare.svg").getroot()
        sixth_ticks = comparison.find(".//svg:g[@class='sixth-ticks']", namespace)
        self.assertEqual(len(sixth_ticks.findall("svg:path", namespace)), 7)

        mastery = (maps / "map-09-mastery.svg").read_text(encoding="utf-8")
        self.assertEqual(mastery.count('class="mastery-skill"'), 5)
        self.assertIn("8/10", mastery)
        self.assertIn("7/10", mastery)

        ratio = (maps / "map-06-ratio.svg").read_text(encoding="utf-8")
        self.assertIn("×2", ratio)
        proportion = (maps / "map-08-proportion.svg").read_text(encoding="utf-8")
        self.assertIn("2 × 3 = 6; 3 × 3 = 9", proportion)

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
