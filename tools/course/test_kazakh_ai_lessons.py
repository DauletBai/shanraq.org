#!/usr/bin/env python3
"""Run every printed program in Kazakh AI lessons 6–10 with novice inputs."""

import re
import subprocess
import sys
import unittest
from pathlib import Path

LESSONS = Path(__file__).resolve().parents[2] / "course/lessons/kazakh-ai"
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
}


class KazakhAILessonsTest(unittest.TestCase):
    def test_printed_programs_and_outputs_in_all_locales(self):
        for stem, cases in CASES.items():
            for suffix in ("", "-kz", "-en"):
                path = LESSONS / f"{stem}{suffix}.md"
                blocks = re.findall(r"^```python\n(.*?)^```", path.read_text(encoding="utf-8"), re.M | re.S)
                self.assertEqual(len(blocks), 1, path)
                for user_input, expected in cases:
                    with self.subTest(path=path.name, user_input=user_input):
                        run = subprocess.run([sys.executable, "-c", blocks[0]], input=user_input,
                                             text=True, capture_output=True, timeout=5, check=False)
                        self.assertEqual(run.returncode, 0, run.stderr)
                        self.assertEqual(run.stdout, expected)

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
