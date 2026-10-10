"""Structural and content checks for the Informatics and AI course."""
import json
import re
import subprocess
import sys
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
    "11-bits-states",
    "12-binary-numbers",
    "13-text-unicode",
    "14-pixels-color",
    "15-sound-video-sampling",
    "16-compression",
    "17-integrity-errors",
    "18-representation-mastery",
    "19-problem-decomposition",
    "20-state-variables",
    "21-sequence-tracing",
    "22-conditions-boundaries",
    "23-loops-invariants",
    "24-functions-contracts",
    "25-correctness-efficiency",
    "26-algorithms-mastery",
    "27-python-first-state",
    "28-python-input-output-types",
    "29-python-expressions",
    "30-python-conditions",
    "31-python-loops",
    "32-python-collections",
    "33-python-text-dates",
    "34-python-functions",
    "35-python-files-json",
    "36-python-errors-debugging",
    "37-python-tests-modules",
    "38-python-cli-release",
    "39-network-message-journey",
    "40-ip-router-packets",
    "41-dns-names",
    "42-tcp-udp-delivery",
    "43-tls-trust",
    "44-http-browser-server",
    "45-html-css-accessibility",
    "46-web-cloud-release",
    "47-observations-data-schema",
    "48-spreadsheets-formulas",
    "49-cleaning-provenance",
    "50-charts-honesty",
    "51-relational-keys",
    "52-sql-queries",
    "53-joins-reports",
    "54-data-release",
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

        for block in self.data["blocks"][:6]:
            cover = ROOT / "web" / block["cover"].lstrip("/")
            self.assertTrue(cover.is_file(), cover)
            self.assertEqual(cover.suffix, ".webp")
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
                headings = required[lang] if stem in RELEASE_STEMS[:18] else (required[lang][0], required[lang][1],
                    {"ru": "## Задание и доказательство", "kz": "## Тапсырма және дәлел", "en": "## Task and evidence"}[lang])
                for heading in headings:
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

    def test_representation_maps_and_examples_are_exact(self):
        for lang in ("ru", "kz", "en"):
            for number, stem in enumerate(RELEASE_STEMS[10:], 11):
                map_stem = "map-" + stem
                path = MAPS / f"{map_stem}-{lang}.svg"
                root = ET.parse(path).getroot()
                self.assertEqual(root.attrib.get("viewBox"), "0 0 1600 900", path)
            cases = {
                12: ("1101", "8 + 4 + 0 + 1"),
                13: ("U+04D8", "D3 98", "U+042F", "D0 AF"),
                14: ("255, 0, 0", "0, 255, 0", "0, 0, 255", "255, 255, 255"),
                15: ("0.25", "0.75", "0.5"),
                16: ("ААААБББ", "4А3Б", "201", "200"),
                17: ("4 + 7 + 2 = 13", "4 + 8 + 2 = 14", "5 + 6 + 2 = 13"),
                18: ("t-01", "t-02", "t-03", "false", "true"),
            }
            for number, labels in cases.items():
                content = (MAPS / f"map-{RELEASE_STEMS[number-1]}-{lang}.svg").read_text()
                for label in labels:
                    self.assertIn(label, content)

    def test_project_checkpoint_matches_lesson_eighteen(self):
        sample = ROOT / "course/informatics-assistant/step-01/data/tasks.json"
        data = json.loads(sample.read_text(encoding="utf-8"))
        self.assertEqual(data["version"], "0.2")
        self.assertEqual([item["id"] for item in data["tasks"]], ["t-01", "t-02", "t-03"])
        self.assertEqual([item["done"] for item in data["tasks"]], [False, True, False])
        for lang in ("ru", "kz", "en"):
            self.assertTrue((ROOT / f"course/informatics-assistant/step-01/FORMAT-{lang}.md").is_file())
            suffix = "" if lang == "ru" else f"-{lang}"
            lesson = (LESSONS / f"18-representation-mastery{suffix}.md").read_text()
            self.assertIn('"version":"0.2"', lesson)
            self.assertIn('"title":"Кітап оқу"', lesson)
            self.assertIn('8/10', lesson)
            self.assertIn('7/10', lesson)

    def test_algorithm_release_keeps_contract_and_examples_aligned(self):
        cases = json.loads((ROOT / "course/informatics-assistant/step-02/cases.json").read_text())
        self.assertEqual(len(cases["reminder_cases"]), 8)
        for case in cases["reminder_cases"]:
            expected = ("DONE" if case["done"] else "NO_DATE" if case["days_to_due"] is None
                        else "OVERDUE" if case["days_to_due"] < 0
                        else "REMIND" if case["days_to_due"] <= 2 else "NOT_YET")
            self.assertEqual(case["expected"], expected, case)
        for case in cases["count_cases"]:
            self.assertEqual(case["expected"], sum(case["done_values"]))
        for case in cases["duplicate_cases"]:
            self.assertEqual(case["has_duplicate"], len(case["ids"]) != len(set(case["ids"])))
        for lang in ("ru", "kz", "en"):
            self.assertTrue((ROOT / f"course/informatics-assistant/step-02/ALGORITHM-{lang}.md").is_file())
            for stem in RELEASE_STEMS[18:26]:
                path = MAPS / f"map-{stem}-{lang}.svg"
                root = ET.parse(path).getroot()
                self.assertEqual(root.attrib.get("viewBox"), "0 0 1600 900", path)
                self.assertEqual(path.read_text().count('class="card-text"'), 4, path)
        source = json.loads((ROOT / "course/informatics-assistant/step-01/data/tasks.json").read_text())
        self.assertTrue(all("due" not in task for task in source["tasks"]))

    def test_publication_is_atomic_and_complete(self):
        from tools.course.prepare_informatics import prepare

        self.assertEqual([block["status"] for block in self.data["blocks"][:6]], ["published"] * 6)
        sql, expected = prepare()
        self.assertEqual(len(expected), 162)
        self.assertEqual(len({item["slug"] for item in expected}), 54)
        self.assertEqual({item["lang"] for item in expected}, {"ru", "kz", "en"})
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.rstrip().endswith("COMMIT;"))
        self.assertEqual(sql.count("INSERT INTO article_series_items"), 54)
        self.assertIn("shanraq-informatics-course", sql)
        web_sql, web_expected = prepare(39)
        self.assertEqual(len(web_expected), 24)
        self.assertEqual(web_sql.count("INSERT INTO article_series_items"), 8)
        self.assertNotIn("informatics-38-python-cli-release", web_sql)
        data_sql, data_expected = prepare(47)
        self.assertEqual(len(data_expected), 24)
        self.assertEqual(data_sql.count("INSERT INTO article_series_items"), 8)
        self.assertNotIn("informatics-46-web-cloud-release", data_sql)

    def test_first_six_releases_have_no_missing_or_extra_pages(self):
        actual = {path.name for path in LESSONS.glob("*.md")}
        expected = {
            f"{stem}{suffix}.md"
            for stem in RELEASE_STEMS
            for suffix in ("", "-kz", "-en")
        }
        self.assertEqual(actual, expected)

    def test_python_examples_run_and_match_all_three_pages(self):
        project = ROOT / "course/informatics-assistant/step-03"
        for stem in RELEASE_STEMS[26:38]:
            number = int(stem[:2])
            for lang in ("ru", "kz", "en"):
                suffix = "" if lang == "ru" else f"-{lang}"
                page = (LESSONS / f"{stem}{suffix}.md").read_text(encoding="utf-8")
                code = re.search(r"```python\n(.*?)\n```", page, re.S)
                output = re.search(r"```text\n(.*?)\n```", page, re.S)
                self.assertIsNotNone(code, (number, lang))
                self.assertIsNotNone(output, (number, lang))
                result = subprocess.run([sys.executable, "-c", code.group(1)], cwd=project,
                                        capture_output=True, text=True, timeout=10)
                self.assertEqual(result.returncode, 0, (number, lang, result.stderr))
                self.assertEqual(result.stdout.strip(), output.group(1), (number, lang))
                map_path = MAPS / f"map-{stem}-{lang}.svg"
                root = ET.parse(map_path).getroot()
                self.assertEqual(root.attrib["viewBox"], "0 0 1600 900")
                self.assertEqual(map_path.read_text().count('class="python-stage"'), 3)
        sample = json.loads((project / "tasks.json").read_text(encoding="utf-8"))
        source = json.loads((ROOT / "course/informatics-assistant/step-01/data/tasks.json").read_text(encoding="utf-8"))
        self.assertEqual(sample["version"], "1.0")
        self.assertEqual([tuple(task[key] for key in ("id", "title", "done")) for task in sample["tasks"]],
                         [tuple(task[key] for key in ("id", "title", "done")) for task in source["tasks"]])

    def test_web_examples_match_three_languages_and_the_checkpoint(self):
        project = ROOT / "course/informatics-assistant/step-04"
        self.assertEqual((project / "assistant_core.py").read_bytes(),
                         (ROOT / "course/informatics-assistant/step-03/assistant_core.py").read_bytes())
        self.assertEqual((project / "tasks.json").read_bytes(),
                         (ROOT / "course/informatics-assistant/step-03/tasks.json").read_bytes())
        for stem in RELEASE_STEMS[38:46]:
            for lang in ("ru", "kz", "en"):
                suffix = "" if lang == "ru" else f"-{lang}"
                page = (LESSONS / f"{stem}{suffix}.md").read_text(encoding="utf-8")
                code = re.search(r"```python\n(.*?)\n```", page, re.S)
                output = re.search(r"```text\n(.*?)\n```", page, re.S)
                self.assertIsNotNone(code, (stem, lang))
                self.assertIsNotNone(output, (stem, lang))
                result = subprocess.run([sys.executable, "-c", code.group(1)], cwd=project,
                                        capture_output=True, text=True, timeout=10)
                self.assertEqual(result.returncode, 0, (stem, lang, result.stderr))
                self.assertEqual(result.stdout.strip(), output.group(1), (stem, lang))
                map_path = MAPS / f"map-{stem}-{lang}.svg"
                root = ET.parse(map_path).getroot()
                self.assertEqual(root.attrib["viewBox"], "0 0 1600 900")
                self.assertEqual(map_path.read_text().count('class="web-stage"'), 3)
        for lang in ("ru", "kz", "en"):
            self.assertTrue((project / f"README-{lang}.md").is_file())

    def test_data_examples_match_three_languages_and_sqlite_checkpoint(self):
        project = ROOT / "course/informatics-assistant/step-05"
        self.assertEqual((project / "assistant_core.py").read_bytes(),
                         (ROOT / "course/informatics-assistant/step-04/assistant_core.py").read_bytes())
        self.assertEqual((project / "tasks.json").read_bytes(),
                         (ROOT / "course/informatics-assistant/step-04/tasks.json").read_bytes())
        for stem in RELEASE_STEMS[46:54]:
            for lang in ("ru", "kz", "en"):
                suffix = "" if lang == "ru" else f"-{lang}"
                page = (LESSONS / f"{stem}{suffix}.md").read_text(encoding="utf-8")
                code = re.search(r"```python\n(.*?)\n```", page, re.S)
                output = re.search(r"```text\n(.*?)\n```", page, re.S)
                self.assertIsNotNone(code, (stem, lang))
                self.assertIsNotNone(output, (stem, lang))
                result = subprocess.run([sys.executable, "-c", code.group(1)], cwd=project,
                                        capture_output=True, text=True, timeout=10)
                self.assertEqual(result.returncode, 0, (stem, lang, result.stderr))
                self.assertEqual(result.stdout.strip(), output.group(1), (stem, lang))
                map_path = MAPS / f"map-{stem}-{lang}.svg"
                root = ET.parse(map_path).getroot()
                self.assertEqual(root.attrib["viewBox"], "0 0 1600 900")
                self.assertEqual(map_path.read_text().count('class="data-stage"'), 3)
        for lang in ("ru", "kz", "en"):
            self.assertTrue((project / f"README-{lang}.md").is_file())


if __name__ == "__main__":
    unittest.main()
