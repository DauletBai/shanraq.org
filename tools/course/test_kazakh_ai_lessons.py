#!/usr/bin/env python3
"""Run every printed program in Kazakh AI lessons 6–30 with novice inputs."""

import re
import json
from datetime import date
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

LESSONS = Path(__file__).resolve().parents[2] / "course/lessons/kazakh-ai"
STEPS = Path(__file__).resolve().parents[2] / "course/kazakh-ai"
STEP_FILES = {
    "06-cyrillic": "step-06/letters.py",
    "07-words": "step-07/words.py",
    "08-roots": "step-08/roots.py",
    "09-plurals": "step-09/plural.py",
    "10-order": "step-10/order.py",
    "11-cases": "step-11/case.py",
    "12-ambiguity": "step-12/meanings.py",
    "13-entities": "step-13/entities.py",
    "14-intent": "step-14/intent.py",
    "15-abstain": "step-15/request.py",
    "16-catalog": "step-16/catalog.py",
    "17-lookup": "step-17/search.py",
    "18-provenance": "step-18/source.py",
    "19-template": "step-19/answer.py",
    "20-expiry": "step-20/model.py",
    "21-labels": "step-21/labels.py",
    "22-splits": "step-22/splits.py",
    "23-classifier": "step-23/classifier.py",
    "24-gate": "step-24/gate.py",
    "25-metrics": "step-25/evaluate.py",
    "26-qazaq-ir": "step-26/compare.py",
    "27-boundary": "step-27/checks.py",
    "28-benchmark": "step-28/bench.py",
    "29-audit": "step-29/engine.py",
    "30-cli": "step-30/main.py",
}
CASES = {
    "06-cyrillic": [("үй\n", "Сөз: үй\n")],
    "07-words": [("Үйлерде, мектептерде?\n", "Сұрақ: ['үйлерде', 'мектептерде']\n"),
                 ("\n", "Сұрақ: []\n")],
    "08-roots": [("ҮЙ\n", "Негіз: баспана\n"),
                 ("мектеп\n", "Негіз: оқу орны\n"),
                 ("үйлер\n", "Негіз: білмеймін\n")],
    "09-plurals": [(word + "\n", "Негіз: " + result + "\n") for word, result in (
        ("қала", "қалалар"), ("үй", "үйлер"), ("қалам", "қаламдар"),
        ("тіл", "тілдер"), ("кітап", "кітаптар"),
        ("мектеп", "мектептер"), ("жол", "білмеймін"))],
    "10-order": [("", "үйлерімізге\n")],
    "11-cases": [("үйлерімізге\n", "Сөз: үй | лер | іміз | барыс\n"),
                 ("мектептерімізде\n", "Сөз: мектеп | тер | іміз | жатыс\n"),
                 ("МЕКТЕПТЕРІМІЗГЕ\n", "Сөз: мектеп | тер | іміз | барыс\n"),
                 ("үйімізге\n", "Сөз: білмеймін\n"),
                 ("үйлерімізді\n", "Сөз: білмеймін\n")],
    "12-ambiguity": [("ат\n", "Сөз: ['жылқы', 'есім']\n"),
                     ("АТ\n", "Сөз: ['жылқы', 'есім']\n"),
                     ("үй\n", "Сөз: ['баспана']\n"),
                     ("жол\n", "Сөз: білмеймін\n")],
    "13-entities": [("Шахмат қашан?\n", "Сұрақ: ['шахмат']\n"),
                    ("Шахмат, сурет қашан?\n", "Сұрақ: ['шахмат', 'сурет']\n"),
                    ("СУРЕТ қашан?\n", "Сұрақ: ['сурет']\n"),
                    ("Шахмат шахмат қашан?\n", "Сұрақ: ['шахмат', 'шахмат']\n"),
                    ("Робот қайда?\n", "Сұрақ: []\n")],
    "14-intent": [("Шахмат қашан?\n", "Сұрақ: ['уақыт']\n"),
                  ("Шахмат қайда?\n", "Сұрақ: ['орын']\n"),
                  ("Шахмат қашан, қайда?\n", "Сұрақ: ['уақыт', 'орын']\n"),
                  ("ҚАЙДА?\n", "Сұрақ: ['орын']\n"),
                  ("\n", "Сұрақ: []\n")],
    "15-abstain": [("Шахмат қашан?\n", "Сұрақ: іздеу: шахмат уақыт\n"),
                   ("СУРЕТ қайда?\n", "Сұрақ: іздеу: сурет орын\n"),
                   ("Робот қашан?\n", "Сұрақ: білмеймін: қай үйірме?\n"),
                   ("Шахмат, сурет қайда?\n", "Сұрақ: нақтылаңыз: бір үйірмені атаңыз\n"),
                   ("Шахмат қашан қайда?\n", "Сұрақ: нақтылаңыз: бір сұрақты таңдаңыз\n"),
                   ("Шахмат неге?\n", "Сұрақ: білмеймін: уақыт па, орын ба?\n"),
                   ("\n", "Сұрақ: білмеймін: қай үйірме?\n")],
    "16-catalog": [("", "Жазба саны: 2\nшахмат уақыт бейсенбі, 15:00\nшахмат орын 203-бөлме\n")],
    "17-lookup": [("шахмат\nуақыт\n", "Үйірме: Мәлімет түрі: Табылды: бейсенбі, 15:00\n"),
                   ("сурет\nорын\n", "Үйірме: Мәлімет түрі: білмеймін: дерек жоқ\n")],
    "18-provenance": [("шахмат\nорын\n", "Үйірме: Мәлімет түрі: Дерек: 203-бөлме\nДереккөз: club-sheet-01\nТексерілген күні: 2026-09-01\n")],
    "19-template": [("шахмат\nуақыт\n", "Үйірме: Мәлімет түрі: шахмат үйірмесінің уақыты: бейсенбі, 15:00\nДереккөз: club-sheet-01 | тексерілген күні: 2026-09-01\n")],
    "20-expiry": [("Шахмат қашан?\n", "Сұрақ: шахмат үйірмесінің уақыты: бейсенбі, 15:00\nДереккөз: club-sheet-01 | тексерілген күні: 2026-09-01\n" if date.today() <= date(2026, 12, 31) else "Сұрақ: білмеймін: дерек ескірген\n"),
                  ("Сурет қайда?\n", "Сұрақ: білмеймін: дерек жоқ\n"),
                  ("Шахмат қайда қашан?\n", "Сұрақ: нақтылаңыз: бір үйірме және бір сұрақ түрі керек\n")],
    "21-labels": [("", "Уақыт: 6\nОрын: 5\nБелгісіз: 2\n")],
    "22-splits": [("", "Оқыту: 8\nБаптау: 5\nБақылау: 0\nЖинақ жарамды: True\n")],
    "23-classifier": [("", "Шахмат қашан өтеді? → уақыт\nСурет қайда өтеді? → орын\nШахмат нешеде? → білмеймін\nШахмат неге? → білмеймін\nШахмат қашан және қайда өтеді? Уақыты қандай? → уақыт\n")],
    "24-gate": [("", "Шахмат қашан өтеді? → уақыт\nСурет қайда өтеді? → орын\nШахмат нешеде? → білмеймін\nШахмат неге? → білмеймін\nШахмат қашан және қайда өтеді? Уақыты қандай? → білмеймін\n")],
    "25-metrics": [("", "Сурет қашан басталады? → уақыт | уақыт\nШахмат қайда болады? → орын | орын\nШахмат қай күні? → білмеймін | уақыт\nСурет қай жерде? → білмеймін | орын\nСурет неге? → білмеймін | белгісіз\nШахмат кімге? → білмеймін | белгісіз\nЖауап: 2 Дұрыс: 2 Белгілі: 4\nБас тарту: 4 Белгісізге дұрыс бас тарту: 2\nДәлдік: 1.0\nТолықтық: 0.5\nБас тарту үлесі: 0.67\n")],
    "26-qazaq-ir": [("", "Оқу үлгісі: мектептерімізде → мектеп ['тер', 'іміз', 'де']\nqazaq-ir: қосымша салыстыру іске қосылмады\n")],
    "27-boundary": [("", "Шахмат қашан? → {'status': 'ready', 'club': 'шахмат', 'kind': 'уақыт'}\nШахмаат қашан? → {'status': 'refuse', 'reason': 'білмеймін: таныс емес сөз'}\nШахмат неге? → {'status': 'refuse', 'reason': 'білмеймін: таныс емес сөз'}\nШахмат қашан интернет? → {'status': 'refuse', 'reason': 'білмеймін: таныс емес сөз'}\nШахмат қашан қайда? → {'status': 'refuse', 'reason': 'нақтылаңыз: бір сұрақ түрі керек'}\n")],
    "28-benchmark": [],
    "29-audit": [("", "Шахмат қашан? → шахмат үйірмесінің уақыты: бейсенбі, 15:00 | 2026-09-24 күнгі дерек | дереккөз: club-sheet-01\nІз: ['Кілт: шахмат / уақыт', 'Жазба саны: 1', 'Дереккөз: club-sheet-01', 'Тексерілген: 2026-09-01', 'Жарамды: 2026-12-31 дейін']\nСурет қайда? → білмеймін: дерек жоқ\nІз: ['Кілт: сурет / орын', 'Жазба саны: 0']\nШахмаат қашан? → білмеймін: таныс емес сөз\nІз: ['Сұрақ түрі анықталмады']\n")],
    "30-cli": [],
}


class KazakhAILessonsTest(unittest.TestCase):
    def test_printed_programs_and_outputs_in_all_locales(self):
        for stem, cases in CASES.items():
            for suffix in ("", "-kz", "-en"):
                path = LESSONS / f"{stem}{suffix}.md"
                blocks = re.findall(r"^```python\n(.*?)^```", path.read_text(encoding="utf-8"), re.M | re.S)
                self.assertEqual(len(blocks), 1, path)
                self.assertEqual(blocks[0], (STEPS / STEP_FILES[stem]).read_text(encoding="utf-8"), path)
                if stem.startswith(("16-", "17-", "18-", "19-", "20-", "21-", "22-")):
                    data = re.findall(r"^```json\n(.*?)^```", path.read_text(encoding="utf-8"), re.M | re.S)
                    self.assertEqual(len(data), 1, path)
                    self.assertEqual(data[0], (STEPS / STEP_FILES[stem].split("/")[0] / ("facts.json" if stem[:2] <= "20" else "examples.json")).read_text(encoding="utf-8"))
                for user_input, expected in cases:
                    with self.subTest(path=path.name, user_input=user_input):
                        run = subprocess.run([sys.executable, "-c", blocks[0]], input=user_input, cwd=STEPS / STEP_FILES[stem].split("/")[0],
                                             text=True, capture_output=True, timeout=5, check=False)
                        self.assertEqual(run.returncode, 0, run.stderr)
                        self.assertEqual(run.stdout, expected)

    def test_learning_splits_keep_test_questions_unseen(self):
        import json
        early_files = [STEPS / f"step-{n}" / "examples.json" for n in range(21, 25)]
        early = [json.loads(path.read_text(encoding="utf-8")) for path in early_files]
        self.assertEqual(len(early[0]), 13)
        self.assertEqual(len(early[1]), 13)
        self.assertTrue(all(data == early[1] for data in early[2:]), "tuning datasets diverged")
        self.assertTrue(all(row["split"] != "test" for row in early[1]))
        final = json.loads((STEPS / "step-25/examples.json").read_text(encoding="utf-8"))
        self.assertEqual(len(final), 19)
        self.assertEqual(final[:13], early[1])
        self.assertEqual({split: sum(row["split"] == split for row in final)
                          for split in ("train", "tune", "test")},
                         {"train": 8, "tune": 5, "test": 6})
        self.assertEqual(len({row["text"] for row in final}), len(final))
        self.assertEqual({row["label"] for row in final if row["split"] == "train"},
                         {"уақыт", "орын"})

    def test_conflict_and_expiry_stop_before_answer(self):
        folder = STEPS / "step-20"
        program = (folder / "model.py").read_text(encoding="utf-8")
        facts = json.loads((folder / "facts.json").read_text(encoding="utf-8"))
        with tempfile.TemporaryDirectory() as tmp:
            work = Path(tmp)
            def run(records):
                (work / "facts.json").write_text(json.dumps(records, ensure_ascii=False), encoding="utf-8")
                return subprocess.run([sys.executable, "-c", program], input="Шахмат қашан?\n",
                                      cwd=work, text=True, capture_output=True, timeout=5, check=True).stdout
            conflict = facts + [dict(facts[0], value="жұма, 16:00")]
            self.assertEqual(run(conflict), "Сұрақ: тоқта: бірнеше дерек табылды\n")
            expired = [dict(facts[0], valid_until="2000-01-01"), facts[1]]
            self.assertEqual(run(expired), "Сұрақ: білмеймін: дерек ескірген\n")

    def test_final_project_reuses_its_checked_components(self):
        for n in (28, 29, 30):
            folder = STEPS / f"step-{n}"
            self.assertEqual((folder / "checks.py").read_bytes(),
                             (STEPS / "step-27/checks.py").read_bytes())
            self.assertEqual((folder / "examples.json").read_bytes(),
                             (STEPS / "step-27/examples.json").read_bytes())
        self.assertEqual((STEPS / "step-30/engine.py").read_bytes(),
                         (STEPS / "step-29/engine.py").read_bytes())
        self.assertEqual((STEPS / "step-30/facts.json").read_bytes(),
                         (STEPS / "step-29/facts.json").read_bytes())
        import json
        rows = json.loads((STEPS / "step-30/examples.json").read_text(encoding="utf-8"))
        self.assertEqual(len(rows), 13)
        self.assertNotIn("test", {row["split"] for row in rows})

    def test_optional_comparison_and_measurements(self):
        compare = subprocess.run([sys.executable, "compare.py", "/no/such/qazaq-ir"],
                                 cwd=STEPS / "step-26", text=True, capture_output=True, check=True)
        self.assertIn("бағдарлама файлы табылмады", compare.stdout)
        bench = subprocess.run([sys.executable, "bench.py"], cwd=STEPS / "step-28",
                               text=True, capture_output=True, timeout=10, check=True)
        self.assertIn("Сұрау саны: 200\n", bench.stdout)
        self.assertRegex(bench.stdout, r"Жаңа іске қосу, мс: [0-9]+\.[0-9]+")
        self.assertRegex(bench.stdout, r"Ең көп бақыланған бөлу, байт: [1-9][0-9]*")
        self.assertIn("құрылғы мен еңбек құны есептелмеді", bench.stdout)

    def test_final_cli_answers_and_refuses_with_reason(self):
        folder = STEPS / "step-30"
        cases = [
            (["2026-09-24", "Шахмат қашан?"], "шахмат үйірмесінің уақыты: бейсенбі, 15:00 | 2026-09-24 күнгі дерек | дереккөз: club-sheet-01\n", 0),
            (["2026-09-24", "Сурет қайда?"], "білмеймін: дерек жоқ\n", 0),
            (["2026-09-24", "Шахмаат қашан?"], "білмеймін: таныс емес сөз\n", 0),
            (["2026-09-24", "Шахмат қашан қайда?"], "нақтылаңыз: бір сұрақ түрі керек\n", 0),
            (["2027-01-01", "Шахмат қашан?"], "білмеймін: бұл күнге жарамды дерек жоқ\n", 0),
            (["2026-01-01", "Шахмат қашан?"], "білмеймін: бұл күнге жарамды дерек жоқ\n", 0),
            (["bad", "Шахмат қашан?"], "Қате күн: YYYY-MM-DD түрінде жазыңыз\n", 2),
        ]
        for args, expected, code in cases:
            with self.subTest(args=args):
                run = subprocess.run([sys.executable, "main.py", *args], cwd=folder,
                                     text=True, capture_output=True, timeout=5, check=False)
                self.assertEqual(run.returncode, code)
                self.assertEqual(run.stdout, expected)
        traced = subprocess.run([sys.executable, "main.py", "--trace", "2026-09-24", "Шахмат қашан?"],
                                cwd=folder, text=True, capture_output=True, timeout=5, check=True)
        self.assertEqual(traced.stdout.count("Із:"), 5)
        import json
        with tempfile.TemporaryDirectory() as tmp:
            work = Path(tmp)
            for name in ("main.py", "engine.py", "checks.py", "examples.json"):
                (work / name).write_bytes((folder / name).read_bytes())
            facts = json.loads((folder / "facts.json").read_text(encoding="utf-8"))
            facts.append(dict(facts[0], value="жұма, 16:00"))
            (work / "facts.json").write_text(json.dumps(facts, ensure_ascii=False), encoding="utf-8")
            conflict = subprocess.run([sys.executable, "main.py", "2026-09-24", "Шахмат қашан?"],
                                      cwd=work, text=True, capture_output=True, timeout=5, check=True)
            self.assertEqual(conflict.stdout, "тоқта: бірнеше дерек табылды\n")

    def test_wrong_order_is_rejected_in_all_locales(self):
        for suffix in ("", "-kz", "-en"):
            path = LESSONS / f"10-order{suffix}.md"
            block = re.findall(r"^```python\n(.*?)^```", path.read_text(encoding="utf-8"), re.M | re.S)[0]
            changed = block.replace('chosen = ["көптік", "тәуелдік", "барыс"]',
                                    'chosen = ["тәуелдік", "көптік", "барыс"]')
            self.assertNotEqual(changed, block)
            run = subprocess.run([sys.executable, "-c", changed],
                                 text=True, capture_output=True, timeout=5, check=False)
            self.assertEqual(run.returncode, 0, run.stderr)
            self.assertEqual(run.stdout, "реті қате\n")


if __name__ == "__main__":
    unittest.main()
