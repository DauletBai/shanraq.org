import unittest
from pathlib import Path
from classifier import read_examples, train, suggest, evaluate

DATA = Path(__file__).resolve().parents[1] / 'labelled_tasks.csv'

class ClassifierTests(unittest.TestCase):
    def test_holdout_never_enters_training_and_unknown_is_review(self):
        rows = read_examples(DATA)
        counts = train(rows)
        self.assertEqual(suggest('Непонятная новая задача', counts), 'review')
        self.assertEqual(suggest('Кітап оқу', counts), 'reading')
        self.assertEqual(suggest('Геометрия есебі', counts), 'math')
        held_out = evaluate(rows, counts)
        self.assertEqual(sum(held_out.values()), 4)
        self.assertEqual(held_out.get(('reading', 'reading')), 2)
        self.assertEqual(held_out.get(('math', 'math')), 2)

if __name__ == '__main__':
    unittest.main()
