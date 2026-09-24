#!/usr/bin/env python3
"""Run every printed program in Kazakh AI lessons 6–25 with novice inputs."""

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
    "21-labels": [("", "Уақыт: 8\nОрын: 7\nБелгісіз: 4\n")],
    "22-splits": [("", "Оқыту: 8\nБаптау: 5\nБақылау: 6\nЖинақ жарамды: True\n")],
    "23-classifier": [("", "Шахмат қашан өтеді? → уақыт\nСурет қайда өтеді? → орын\nШахмат нешеде? → білмеймін\nШахмат неге? → білмеймін\nШахмат қашан және қайда өтеді? Уақыты қандай? → уақыт\n")],
    "24-gate": [("", "Шахмат қашан өтеді? → уақыт\nСурет қайда өтеді? → орын\nШахмат нешеде? → білмеймін\nШахмат неге? → білмеймін\nШахмат қашан және қайда өтеді? Уақыты қандай? → білмеймін\n")],
    "25-metrics": [("", "Сурет қашан басталады? → уақыт | уақыт\nШахмат қайда болады? → орын | орын\nШахмат қай күні? → білмеймін | уақыт\nСурет қай жерде? → білмеймін | орын\nСурет неге? → білмеймін | белгісіз\nШахмат кімге? → білмеймін | белгісіз\nЖауап: 2 Дұрыс: 2 Белгілі: 4\nБас тарту: 4 Белгісізге дұрыс бас тарту: 2\nДәлдік: 1.0\nТолықтық: 0.5\nБас тарту үлесі: 0.67\n")],
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
        files = [STEPS / f"step-{n}" / "examples.json" for n in range(22, 26)]
        datasets = [json.loads(path.read_text(encoding="utf-8")) for path in files]
        self.assertTrue(all(data == datasets[0] for data in datasets), "lesson datasets diverged")
        rows = datasets[0]
        self.assertEqual(len(rows), 19)
        self.assertEqual({split: sum(row["split"] == split for row in rows)
                          for split in ("train", "tune", "test")},
                         {"train": 8, "tune": 5, "test": 6})
        self.assertEqual(len({row["text"] for row in rows}), len(rows))
        self.assertEqual({row["label"] for row in rows if row["split"] == "train"},
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
