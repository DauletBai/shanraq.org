"""Local classroom security checkpoint: verify, back up and restore fictional data.

The database remains loopback-only and read-only in the web view. The backup is
NOT encrypted and must stay on a trusted local device. This is not account auth.
"""
import argparse
from contextlib import closing
import hashlib
import json
import os
from pathlib import Path
import sqlite3
import tempfile

from data_store import connect_readonly


def digest(path):
    hasher = hashlib.sha256()
    with Path(path).open('rb') as source:
        for chunk in iter(lambda: source.read(65536), b''):
            hasher.update(chunk)
    return hasher.hexdigest()


def verify(path):
    with closing(connect_readonly(path)) as db:
        result = db.execute('PRAGMA integrity_check').fetchone()[0]
        tasks = db.execute('SELECT COUNT(*) FROM tasks').fetchone()[0]
        sessions = db.execute('SELECT COUNT(*) FROM study_sessions').fetchone()[0]
        links = db.execute('PRAGMA foreign_key_check').fetchall()
    if result != 'ok' or links:
        raise ValueError('database integrity check failed')
    return tasks, sessions


def backup(source, destination):
    source, destination = Path(source), Path(destination)
    if destination.exists() or destination.with_suffix(destination.suffix + '.sha256').exists():
        raise ValueError('backup or checksum exists; refusing to overwrite')
    expected = verify(source)
    destination.parent.mkdir(parents=True, exist_ok=True)
    fd, name = tempfile.mkstemp(prefix='.backup-', suffix='.db', dir=destination.parent)
    os.close(fd)
    temporary = Path(name)
    try:
        with closing(connect_readonly(source)) as original, closing(sqlite3.connect(temporary)) as copy:
            original.backup(copy)
        if verify(temporary) != expected:
            raise ValueError('backup counts differ from source')
        checksum = digest(temporary)
        os.link(temporary, destination)
        destination.with_suffix(destination.suffix + '.sha256').write_text(checksum + '\n', encoding='ascii')
        return checksum, expected
    finally:
        temporary.unlink(missing_ok=True)


def restore(source, destination):
    source, destination = Path(source), Path(destination)
    if destination.exists():
        raise ValueError('destination exists; refusing to overwrite')
    recorded = source.with_suffix(source.suffix + '.sha256').read_text(encoding='ascii').strip()
    if digest(source) != recorded:
        raise ValueError('backup checksum mismatch')
    expected = verify(source)
    destination.parent.mkdir(parents=True, exist_ok=True)
    fd, name = tempfile.mkstemp(prefix='.restore-', suffix='.db', dir=destination.parent)
    os.close(fd)
    temporary = Path(name)
    try:
        with closing(connect_readonly(source)) as original, closing(sqlite3.connect(temporary)) as copy:
            original.backup(copy)
        if verify(temporary) != expected:
            raise ValueError('restored counts differ')
        os.link(temporary, destination)
        return expected
    finally:
        temporary.unlink(missing_ok=True)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('action', choices=('verify', 'backup', 'restore'))
    parser.add_argument('source', type=Path)
    parser.add_argument('destination', nargs='?', type=Path)
    args = parser.parse_args()
    if args.action == 'verify':
        print('tasks=%d sessions=%d' % verify(args.source))
        return
    if args.destination is None:
        parser.error('destination is required')
    result = backup(args.source, args.destination) if args.action == 'backup' else restore(args.source, args.destination)
    print(json.dumps({'action': args.action, 'result': result}, ensure_ascii=False))


if __name__ == '__main__':
    main()
