"""Publication boundaries, source coverage, and SQL quoting checks."""
import unittest
import json
import tempfile
from pathlib import Path
from unittest.mock import patch
import prepare_rust
from prepare_rust import literal, prepare


class RustReleaseTests(unittest.TestCase):
    def test_seventh_batch_is_complete_and_atomic(self):
        sql, expected = prepare()
        self.assertTrue(sql.startswith('BEGIN;'))
        self.assertTrue(sql.endswith('COMMIT;\n'))
        self.assertEqual(len(expected), 108)
        self.assertEqual({p['lang'] for p in expected}, {'kz', 'ru', 'en'})
        self.assertEqual(len({p['slug'] for p in expected}), 36)
        self.assertEqual(sql.count('INSERT INTO article_series_items('), 36)
        self.assertIn('rust-35-arguments', sql)
        self.assertNotIn('rust-36-', sql)
        self.assertNotIn('DELETE FROM', sql)
        self.assertNotIn('TRUNCATE', sql)
        self.assertNotIn('Редакционный черновик', sql)
        self.assertNotIn('Editorial draft.', sql)
        self.assertNotIn('Редакциялық жоба.', sql)

    def test_incomplete_or_misnumbered_group_is_rejected(self):
        original = json.loads((prepare_rust.ROOT / 'tools/course/rust-release.json').read_text(encoding='utf-8'))
        for defect in ('partial_group', 'missing_item', 'wrong_batch'):
            with self.subTest(defect=defect), tempfile.TemporaryDirectory() as name:
                manifest = json.loads(json.dumps(original))
                if defect == 'partial_group':
                    manifest['lessons'] = 11
                elif defect == 'missing_item':
                    manifest['items'].pop()
                else:
                    manifest['batch'] = 1
                root = Path(name)
                path = root / 'tools/course/rust-release.json'
                path.parent.mkdir(parents=True)
                path.write_text(json.dumps(manifest), encoding='utf-8')
                with patch.object(prepare_rust, 'ROOT', root), self.assertRaises(AssertionError):
                    prepare()

    def test_quote_preserves_text_and_rejects_delimiter(self):
        value = "Reader's task; $HOME — менің ісім"
        self.assertEqual(literal(value), '$rust_course$'+value+'$rust_course$')
        with self.assertRaises(ValueError):
            literal('$rust_course$; DROP TABLE articles;')


if __name__ == '__main__':
    unittest.main()
