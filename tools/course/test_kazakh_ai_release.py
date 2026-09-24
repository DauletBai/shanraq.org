#!/usr/bin/env python3
"""Check that the Kazakh AI release is complete and atomic."""

import unittest
import prepare_kazakh_ai as release


class KazakhAIReleaseTest(unittest.TestCase):
    def test_complete_course_is_trilingual_and_transactional(self):
        sql, expected = release.prepare()
        self.assertEqual(release.META["ru"][0], "AI без LLM: создаем свою модель ИИ")
        self.assertEqual(release.META["kz"][0], "Үлкен тілдік үлгісіз ЖИ: өз үлгімізді жасаймыз")
        self.assertEqual(release.META["en"][0], "AI without an LLM: build your own AI model")
        self.assertEqual(len(expected), 93)
        self.assertEqual(len({(x["slug"], x["lang"]) for x in expected}), 93)
        self.assertEqual({x["lang"] for x in expected}, {"ru", "kz", "en"})
        self.assertEqual(sql.count("INSERT INTO article_series_items("), 31)
        self.assertEqual(sql.count("INSERT INTO article_translations("), 93)
        self.assertEqual(sql.count("ON CONFLICT(slug) DO UPDATE SET cover_url='',status='published'"), 31)
        self.assertIn("/static/brand/shanraq.svg", sql)
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.endswith("COMMIT;\n"))
        self.assertLess(sql.index("INSERT INTO article_series("),
                        sql.index("INSERT INTO articles("))
        for slug in release.SLUGS:
            self.assertIn(slug, sql)


if __name__ == "__main__":
    unittest.main()
