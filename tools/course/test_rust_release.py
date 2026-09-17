"""Publication boundaries, source coverage, and SQL quoting checks."""
import unittest
from prepare_rust import literal, prepare


class RustReleaseTests(unittest.TestCase):
    def test_first_batch_is_complete_and_atomic(self):
        sql, expected = prepare()
        self.assertTrue(sql.startswith('BEGIN;'))
        self.assertTrue(sql.endswith('COMMIT;\n'))
        self.assertEqual(len(expected), 18)
        self.assertEqual({p['lang'] for p in expected}, {'kz', 'ru', 'en'})
        self.assertEqual(len({p['slug'] for p in expected}), 6)
        self.assertEqual(sql.count('INSERT INTO article_series_items('), 6)
        self.assertNotIn('rust-06-', sql)
        self.assertNotIn('DELETE FROM', sql)
        self.assertNotIn('TRUNCATE', sql)
        self.assertNotIn('Редакционный черновик', sql)
        self.assertNotIn('Editorial draft.', sql)
        self.assertNotIn('Редакциялық жоба.', sql)

    def test_quote_preserves_text_and_rejects_delimiter(self):
        value = "Reader's task; $HOME — менің ісім"
        self.assertEqual(literal(value), '$rust_course$'+value+'$rust_course$')
        with self.assertRaises(ValueError):
            literal('$rust_course$; DROP TABLE articles;')


if __name__ == '__main__':
    unittest.main()
