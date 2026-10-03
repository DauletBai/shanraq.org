import json
from pathlib import Path
import re
import unittest

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/kazakh-language"
MAPS = ROOT / "web/static/course/kazakh-language"


class KazakhLanguageReleaseTest(unittest.TestCase):
    def test_architecture_has_ten_blocks_and_eighty_lessons(self):
        data = json.loads((ROOT / "course/kazakh-language/curriculum.json").read_text())
        self.assertEqual(len(data["blocks"]), 10)
        self.assertEqual(data["blocks"][-1]["range"], "73-80")
        self.assertEqual(len(data["first_block_lessons"]), 8)
        self.assertEqual(data["blocks"][3]["lexical_target"], 800)
        self.assertEqual(data["blocks"][6]["lexical_target"], 1783)
        self.assertEqual(data["blocks"][-1]["lexical_target"], 2229)

    def test_started_lesson_is_manual_three_language_content(self):
        names = ["01-map-first-conversation.md", "01-map-first-conversation-kz.md", "01-map-first-conversation-en.md"]
        for name in names:
            text = (LESSONS / name).read_text()
            self.assertIn("data-speak-kz=", text)
            self.assertIn("## Тапсырма" if "-kz" in name else "## Exercise" if "-en" in name else "## Задание", text)
            self.assertGreater(len(text.split()), 450)
            self.assertNotRegex(text, r"\b(TODO|TBD)\b")

    def test_all_first_block_maps_exist_and_fit_the_canvas(self):
        for number, stem in enumerate([
            "01-first-contact", "02-nine-sounds", "03-harmony", "04-sentence",
            "05-introduction", "06-questions", "07-repair", "08-mastery",
        ], 1):
            for lang in ("ru", "kz", "en"):
                path = MAPS / f"map-{stem}-{lang}.svg"
                text = path.read_text()
                self.assertIn('viewBox="0 0 1600 900"', text, path)
                self.assertNotIn("TODO", text, path)
                self.assertFalse(re.search(r'<text[^>]*x="(?:1[6-9]\d\d|[2-9]\d{3})"', text), path)


if __name__ == "__main__":
    unittest.main()
