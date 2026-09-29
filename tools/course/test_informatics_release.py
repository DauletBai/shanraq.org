"""Structural and content checks for the Informatics and AI course."""
import json
import re
import unittest
from pathlib import Path
import xml.etree.ElementTree as ET


ROOT = Path(__file__).resolve().parents[2]
CURRICULUM = ROOT / "course/informatics/curriculum.json"
LESSONS = ROOT / "course/lessons/informatics"
MAPS = ROOT / "web/static/course/informatics"
RELEASE_STEMS = (
    "01-whole-map",
    "02-diagnostic-product",
    "03-device-system",
    "04-input-output-sensors",
    "05-cpu-memory-storage",
    "06-os-process-app",
    "07-files-folders-paths",
    "08-formats-software-licenses",
    "09-versions-collaboration-accessibility",
    "10-systems-mastery",
)


class InformaticsReleaseTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.data = json.loads(CURRICULUM.read_text(encoding="utf-8"))

    def test_curriculum_has_eight_atomic_blocks_and_72_lessons(self):
        blocks = self.data["blocks"]
        lessons = self.data["lessons"]
        self.assertEqual(self.data["publication_policy"], "complete-block-only")
        self.assertEqual([b["lesson_count"] for b in blocks], [10, 8, 8, 12, 8, 8, 8, 10])
        self.assertEqual(sum(b["lesson_count"] for b in blocks), 72)
        self.assertEqual([lesson["number"] for lesson in lessons], list(range(1, 73)))
        self.assertEqual(len({lesson["id"] for lesson in lessons}), 72)

        for block in blocks:
            start, end = map(int, block["range"].split("-"))
            members = [lesson for lesson in lessons if lesson["block"] == block["id"]]
            self.assertEqual([lesson["number"] for lesson in members], list(range(start, end + 1)))
            self.assertEqual(len(members), block["lesson_count"])
            self.assertIn(block["status"], {"planned", "in_development", "ready", "published"})

    def test_curriculum_is_an_ordered_dag_with_three_titles(self):
        lessons = self.data["lessons"]
        by_id = {lesson["id"]: lesson for lesson in lessons}
        for lesson in lessons:
            for key in ("title_ru", "title_kz", "title_en"):
                self.assertTrue(lesson[key].strip(), (lesson["id"], key))
            for prerequisite in lesson["prerequisites"]:
                self.assertIn(prerequisite, by_id)
                self.assertLess(by_id[prerequisite]["number"], lesson["number"])

        visiting, visited = set(), set()

        def walk(lesson_id):
            if lesson_id in visiting:
                self.fail(f"curriculum cycle at {lesson_id}")
            if lesson_id in visited:
                return
            visiting.add(lesson_id)
            for prerequisite in by_id[lesson_id]["prerequisites"]:
                walk(prerequisite)
            visiting.remove(lesson_id)
            visited.add(lesson_id)

        for lesson_id in by_id:
            walk(lesson_id)
        self.assertEqual(visited, set(by_id))

    def test_project_and_mastery_policy_are_explicit(self):
        project = self.data["project"]
        self.assertTrue(project["uses_fictional_data_only"])
        self.assertEqual(project["title_ru"], "Мой цифровой помощник")
        self.assertEqual(project["title_kz"], "Менің цифрлық көмекшім")
        gate = self.data["mastery_gate"]
        self.assertEqual(gate["immediate"], "8/10")
        self.assertEqual(gate["delayed_7_days"], "7/10")
        self.assertTrue(gate["delayed_30_days"])
        self.assertTrue(gate["requires_unseen_transfer"])
        self.assertTrue(gate["requires_project_evidence"])

        first_cover = ROOT / "web" / self.data["blocks"][0]["cover"].lstrip("/")
        self.assertTrue(first_cover.is_file(), first_cover)
        self.assertEqual(first_cover.suffix, ".webp")
        checkpoint = ROOT / "course/informatics-assistant/step-00/project-passport.template.md"
        self.assertTrue(checkpoint.is_file(), checkpoint)
        checkpoint_text = checkpoint.read_text(encoding="utf-8")
        self.assertIn("Three fictional records", checkpoint_text)
        self.assertIn("Decision log", checkpoint_text)

    def test_release_has_three_complete_localizations(self):
        required = {
            "ru": ("## Где мы на карте", "## Опорный сигнал", "## Задание"),
            "kz": ("## Картадағы орнымыз", "## Сабақтың тірек сигналы", "## Тапсырма"),
            "en": ("## Where we are on the map", "## The lesson's support signal", "## Exercise"),
        }
        forbidden = {
            "kz": ("## Где мы", "## Задание", "Следующий урок"),
            "en": ("## Где мы", "## Задание", "Келесі сабақ"),
        }
        minimum_words = {"ru": 500, "kz": 440, "en": 500}
        for stem in RELEASE_STEMS:
            for lang in ("ru", "kz", "en"):
                suffix = "" if lang == "ru" else f"-{lang}"
                path = LESSONS / f"{stem}{suffix}.md"
                self.assertTrue(path.is_file(), path)
                text = path.read_text(encoding="utf-8")
                for heading in required[lang]:
                    self.assertIn(heading, text, path)
                for phrase in forbidden.get(lang, ()):
                    self.assertNotIn(phrase, text, path)
                self.assertGreaterEqual(
                    len(re.findall(r"\b[\w'-]+\b", text)), minimum_words[lang], path
                )
                self.assertIn("1, 7", text, path)

    def test_diagnostic_has_same_ten_tasks_and_examples(self):
        paths = (
            LESSONS / "02-diagnostic-product.md",
            LESSONS / "02-diagnostic-product-kz.md",
            LESSONS / "02-diagnostic-product-en.md",
        )
        for path in paths:
            text = path.read_text(encoding="utf-8")
            numbers = [int(value) for value in re.findall(r"^### (\d+)\.", text, re.MULTILINE)]
            self.assertEqual(numbers, list(range(1, 11)), path)
            for marker in ("plan.txt", "10 × 10", "100 × 100", "3 → 2 → 4 → 1", "tasks.json"):
                self.assertIn(marker, text, path)

    def test_all_release_maps_are_valid_and_exact(self):
        for lang in ("ru", "kz", "en"):
            whole = MAPS / f"map-01-whole-map-{lang}.svg"
            diagnostic = MAPS / f"map-02-diagnostic-{lang}.svg"
            ET.parse(whole)
            ET.parse(diagnostic)
            whole_text = whole.read_text(encoding="utf-8")
            diagnostic_text = diagnostic.read_text(encoding="utf-8")
            self.assertEqual(whole_text.count('class="course-step"'), 8)
            self.assertEqual(whole_text.count('class="course-arrow"'), 7)
            self.assertEqual(diagnostic_text.count('class="diagnostic-group"'), 4)
            for span in ("1–3", "4–6", "7–8", "9–10"):
                self.assertIn(f'data-questions="{span}"', diagnostic_text)

            system = MAPS / f"map-03-device-system-{lang}.svg"
            io = MAPS / f"map-04-input-output-{lang}.svg"
            ET.parse(system)
            ET.parse(io)
            system_text = system.read_text(encoding="utf-8")
            io_text = io.read_text(encoding="utf-8")
            self.assertEqual(system_text.count('class="system-role"'), 5)
            self.assertEqual(system_text.count('class="device-system"'), 2)
            self.assertEqual(io_text.count('class="io-stage"'), 5)
            self.assertEqual(io_text.count('class="io-example"'), 3)
            self.assertIn("x=412, y=728", io_text)

        map_names = (
            "map-05-memory", "map-06-os-process", "map-07-paths",
            "map-08-formats", "map-09-versions", "map-10-mastery",
        )
        for name in map_names:
            for lang in ("ru", "kz", "en"):
                path = MAPS / f"{name}-{lang}.svg"
                root = ET.parse(path).getroot()
                self.assertEqual(root.attrib.get("viewBox"), "0 0 1600 900", path)
                text = path.read_text(encoding="utf-8")
                self.assertGreaterEqual(text.count('class="foundation-stage"'), 2, path)
                self.assertIn("filter=\"url(#shadow)\"", text, path)

        exact_labels = {
            "map-05-memory": {
                "ru": ("накопитель", "ОЗУ", "процессор", "сохранение"),
                "kz": ("сақтау құрылғысы", "жедел жад", "процессор"),
                "en": ("storage", "RAM", "processor"),
            },
            "map-06-os-process": {
                "ru": ("файл программы", "ОС", "процесс"),
                "kz": ("бағдарлама файлы", "ОЖ", "процесс"),
                "en": ("program file", "OS", "process"),
            },
            "map-10-mastery": {
                "ru": ("найти ошибку", "8/10 сейчас"),
                "kz": ("қатені табу", "қазір 8/10"),
                "en": ("diagnose", "8/10 now"),
            },
        }
        for name, translations in exact_labels.items():
            for lang, labels in translations.items():
                text = (MAPS / f"{name}-{lang}.svg").read_text(encoding="utf-8")
                for label in labels:
                    self.assertIn(label, text, (name, lang, label))

    def test_publication_is_atomic_and_complete(self):
        from tools.course.prepare_informatics import prepare

        first = self.data["blocks"][0]
        self.assertEqual(first["status"], "published")
        sql, expected = prepare()
        self.assertEqual(len(expected), 30)
        self.assertEqual(len({item["slug"] for item in expected}), 10)
        self.assertEqual({item["lang"] for item in expected}, {"ru", "kz", "en"})
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.rstrip().endswith("COMMIT;"))
        self.assertEqual(sql.count("INSERT INTO article_series_items"), 10)
        self.assertIn("shanraq-informatics-course", sql)

    def test_first_release_has_no_missing_or_extra_pages(self):
        actual = {path.name for path in LESSONS.glob("*.md")}
        expected = {
            f"{stem}{suffix}.md"
            for stem in RELEASE_STEMS
            for suffix in ("", "-kz", "-en")
        }
        self.assertEqual(actual, expected)


if __name__ == "__main__":
    unittest.main()
