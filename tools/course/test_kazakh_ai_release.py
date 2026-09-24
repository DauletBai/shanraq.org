#!/usr/bin/env python3
"""Check that the Kazakh AI release is complete and atomic."""

import unittest
import prepare_kazakh_ai as release


class KazakhAIReleaseTest(unittest.TestCase):
    def test_first_block_is_trilingual_and_transactional(self):
        sql, expected = release.prepare()
        self.assertEqual(len(expected), 18)
        self.assertEqual(len({(x["slug"], x["lang"]) for x in expected}), 18)
        self.assertEqual({x["lang"] for x in expected}, {"ru", "kz", "en"})
        self.assertEqual(sql.count("INSERT INTO article_series_items("), 6)
        self.assertEqual(sql.count("INSERT INTO article_translations("), 18)
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.endswith("COMMIT;\n"))
        self.assertLess(sql.index("INSERT INTO article_series("),
                        sql.index("INSERT INTO articles("))
        for slug in release.SLUGS:
            self.assertIn(slug, sql)


if __name__ == "__main__":
    unittest.main()
