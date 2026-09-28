"""Publication and source-boundary checks for the Cloud & DevOps course."""
import tempfile
import unittest
from pathlib import Path
import re
import prepare_cloud_devops


class CloudDevOpsReleaseTests(unittest.TestCase):
    def test_release_is_complete_atomic_and_non_destructive(self):
        sql, expected = prepare_cloud_devops.prepare()
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.endswith("COMMIT;\n"))
        self.assertEqual(len(expected), 75)
        self.assertEqual(len({x["slug"] for x in expected}), 25)
        self.assertEqual({x["lang"] for x in expected}, {"ru", "kz", "en"})
        self.assertEqual(sql.count("INSERT INTO article_series_items("), 25)
        self.assertIn("'it','devops'", sql)
        self.assertIn("'published','shell'", sql)
        self.assertNotIn("DELETE FROM", sql)
        self.assertNotIn("TRUNCATE", sql)

    def test_every_lesson_has_exercise_map_and_safe_source(self):
        for stem in prepare_cloud_devops.STEMS[1:]:
            for lang, suffix in prepare_cloud_devops.LANGS.items():
                path = prepare_cloud_devops.LESSONS / f"{stem}{suffix}.md"
                text = path.read_text(encoding="utf-8")
                head = {"ru": "## Задание", "kz": "## Тапсырма", "en": "## Exercise"}[lang]
                self.assertIn(head, text, path)
                self.assertIn("/static/course/cloud-devops/", text, path)

    def test_every_lesson_teaches_a_beginner_before_giving_commands(self):
        sections = {
            "ru": ("## Сначала знакомый образ", "## Опорный сигнал урока",
                   "## Новые слова простыми словами", "## Сначала предскажите результат",
                   "## Воспроизведите без подсказки", "## Найдите и исправьте ошибку"),
            "kz": ("## Алдымен таныс бейне", "## Сабақтың тірек сигналы",
                   "## Жаңа сөздер қарапайым тілмен", "## Алдымен нәтижені болжаңыз",
                   "## Көмексіз жаңғыртыңыз", "## Қатені тауып түзетіңіз"),
            "en": ("## Begin with a familiar image", "## The lesson's support signal",
                   "## New words in plain language", "## Predict the result first",
                   "## Recall without a prompt", "## Find and correct the mistake"),
        }
        for stem in prepare_cloud_devops.STEMS[1:]:
            for lang, suffix in prepare_cloud_devops.LANGS.items():
                path = prepare_cloud_devops.LESSONS / f"{stem}{suffix}.md"
                text = path.read_text(encoding="utf-8")
                for heading in sections[lang]:
                    self.assertIn(heading, text, path)
                # The earlier draft averaged roughly 250 words per lesson.  A
                # floor guards the explanatory bridge without rewarding filler.
                words = re.findall(r"\b[\w&/-]+\b", text)
                self.assertGreaterEqual(len(words), 500, path)

    def test_beginner_dependency_order(self):
        self.assertEqual(
            prepare_cloud_devops.STEMS[:7],
            ("preface", "01-why-now", "03-terminal-git", "04-linux-access",
             "05-processes-logs", "06-networks-dns-http", "02-service-map"),
        )

    def test_preface_contains_zero_setup_and_teaching_method(self):
        for lang, suffix in prepare_cloud_devops.LANGS.items():
            text = (prepare_cloud_devops.LESSONS / f"preface{suffix}.md").read_text(encoding="utf-8")
            self.assertIn("docker version", text)
            self.assertIn("course/cloud-devops-lab", text)
            method = {"ru": "методика опорных сигналов", "kz": "тірек сигналдар әдісі",
                      "en": "support-signal method"}[lang]
            self.assertIn(method, text)

    def test_sql_delimiter_is_rejected(self):
        with self.assertRaises(ValueError):
            prepare_cloud_devops.literal("bad $cloud_devops_course$ value")


if __name__ == "__main__":
    unittest.main()
