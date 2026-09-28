"""Publication and source-boundary checks for the Cloud & DevOps course."""
import tempfile
import unittest
from pathlib import Path
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

    def test_sql_delimiter_is_rejected(self):
        with self.assertRaises(ValueError):
            prepare_cloud_devops.literal("bad $cloud_devops_course$ value")


if __name__ == "__main__":
    unittest.main()
