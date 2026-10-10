"""Transparent toy text classifier for fictional tasks; suggestions only.

A small labelled file is not evidence of accuracy on real learners. No network
request, account, or automatic change to a task is made.
"""
import csv
from collections import Counter, defaultdict
from pathlib import Path
import re

LABELS = ('math', 'reading')
TOKEN = re.compile(r'\w+', re.UNICODE)


def words(title):
    return [word.lower() for word in TOKEN.findall(title)]


def read_examples(path):
    with Path(path).open(encoding='utf-8', newline='') as stream:
        reader = csv.DictReader(stream)
        if reader.fieldnames != ['title', 'label', 'split']:
            raise ValueError('wrong training file headers')
        rows = list(reader)
    if any(row['label'] not in LABELS or row['split'] not in ('train', 'test') or not words(row['title']) for row in rows):
        raise ValueError('invalid labelled example')
    if not any(row['split'] == 'train' for row in rows) or not any(row['split'] == 'test' for row in rows):
        raise ValueError('need separate train and test examples')
    return rows


def train(rows):
    """Count words in the training split only; a tie means human review."""
    counts = {label: Counter() for label in LABELS}
    for row in rows:
        if row['split'] == 'train':
            counts[row['label']].update(set(words(row['title'])))
    return counts


def suggest(title, counts):
    tokens = set(words(title))
    scores = {label: sum(counts[label][token] for token in tokens) for label in LABELS}
    if scores['math'] == scores['reading']:
        return 'review'
    return max(LABELS, key=lambda label: scores[label])


def evaluate(rows, counts):
    matrix = defaultdict(int)
    for row in rows:
        if row['split'] == 'test':
            matrix[(row['label'], suggest(row['title'], counts))] += 1
    return dict(matrix)


def main():
    import argparse
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('title', nargs='?', default='Кітап оқу')
    args = parser.parse_args()
    rows = read_examples(Path(__file__).with_name('labelled_tasks.csv'))
    counts = train(rows)
    print('suggestion:', suggest(args.title, counts))
    print('held-out:', evaluate(rows, counts))


if __name__ == '__main__':
    main()
