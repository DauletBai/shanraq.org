import shutil
import tempfile
import unittest
from pathlib import Path

from data_store import create_database
from security_assistant import backup, digest, restore, verify

ROOT = Path(__file__).resolve().parents[1]

class SecurityTests(unittest.TestCase):
    def test_backup_restore_and_tamper_detection(self):
        with tempfile.TemporaryDirectory() as folder:
            directory = Path(folder)
            live = directory / 'assistant.db'
            saved = directory / 'backup.db'
            restored = directory / 'restored.db'
            create_database(live, ROOT / 'tasks.json', ROOT / 'study_sessions.csv')
            self.assertEqual(verify(live), (3, 4))
            backup(live, saved)
            self.assertEqual(restore(saved, restored), (3, 4))
            self.assertEqual(digest(saved), digest(restored))
            with self.assertRaisesRegex(ValueError, 'exists'):
                restore(saved, restored)
            with saved.open('ab') as file:
                file.write(b'tampered')
            with self.assertRaisesRegex(ValueError, 'checksum'):
                restore(saved, directory / 'another.db')

if __name__ == '__main__':
    unittest.main()
