"""Release gates for a complete, safe, trilingual first chemistry block."""
import json
from pathlib import Path
import re
import unittest
import xml.etree.ElementTree as ET

from tools.course.build_chemistry_block import LESSONS, LANGS, LESSON_DIR, MAP_DIR, ROOT, VISUAL_FILES, map_svg, render
from tools.course.prepare_chemistry import prepare, read_lesson
from tools.course.chemistry_particles_data import LESSONS as PARTICLE_LESSONS
from tools.course.build_chemistry_particles import diagram, render as render_particles, slug as particle_slug, visual
from tools.course.prepare_chemistry_particles import prepare as prepare_particles


class ChemistryReleaseTests(unittest.TestCase):
    def test_route_and_each_language_are_complete(self):
        curriculum = json.loads((ROOT / 'course/chemistry/curriculum.json').read_text())
        self.assertEqual(curriculum['publication_policy'], 'complete-block-only')
        self.assertEqual(curriculum['languages'], list(LANGS))
        self.assertEqual(curriculum['blocks'][0]['range'], '1-5')
        self.assertEqual(len(LESSONS), 5)
        slugs = json.loads((ROOT / 'tools/course/lesson-slugs.json').read_text())
        for number, row in enumerate(LESSONS, 1):
            slug = f'chemistry-{number:02d}-{row["stem"].split("-",1)[1]}'
            for lang in LANGS:
                suffix = '' if lang == 'ru' else '-' + lang
                path = LESSON_DIR / f'{row["stem"]}{suffix}.md'
                self.assertEqual(path.read_text(), render(number, lang), path)
                self.assertEqual(slugs[f'chemistry/{row["stem"]}{suffix}'], [slug, lang])
                title, summary, body = read_lesson(path)
                self.assertTrue(title and summary)
                self.assertGreater(len(body.split()), 350, path)
                self.assertIn(VISUAL_FILES[number-1], body)
                for linked in re.findall(r'!\[[^]]+\]\((/static/[^)]+)\)', body):
                    self.assertTrue((ROOT / 'web' / linked.lstrip('/')).is_file(), linked)
                map_path = MAP_DIR / f'map-{number:02d}-{lang}.svg'
                self.assertEqual(map_path.read_text(), map_svg(number, lang))
                ET.parse(map_path)
                self.assertNotIn('textLength=', map_path.read_text())
                if number < 5:
                    self.assertIn(f'/read/chemistry-{number+1:02d}-', body)

    def test_illustrations_and_safety_gate(self):
        for asset in VISUAL_FILES:
            self.assertGreater((ROOT / 'web' / asset.lstrip('/')).stat().st_size, 100_000)
        for lang in LANGS:
            first = (LESSON_DIR / f'01-question-and-safety{("-"+lang) if lang != "ru" else ""}.md').read_text().lower()
            last = (LESSON_DIR / f'05-separate-and-explain{("-"+lang) if lang != "ru" else ""}.md').read_text().lower()
            self.assertTrue(any(term in first for term in ('не пробуйте', 'дәмін татпаңыз', 'do not heat or taste')))
            self.assertTrue(any(term in last for term in ('не пробуйте фильтрат', 'дәмін татпаңыз', 'do not heat anything or taste the filtrate')))

    def test_publication_is_guarded_and_has_no_news_syndication(self):
        sql, expected = prepare()
        self.assertEqual(len(expected), 15)
        self.assertIn('BEGIN;', sql)
        self.assertIn('COMMIT;', sql)
        self.assertIn('pg_advisory_xact_lock', sql)
        self.assertIn("s.slug='chemistry'", sql)
        self.assertNotIn('syndicate_telegram', sql)

    def test_second_block_is_complete_and_navigable(self):
        curriculum = json.loads((ROOT / 'course/chemistry/curriculum.json').read_text())
        self.assertEqual(curriculum['blocks'][1]['range'], '6-13')
        self.assertEqual(curriculum['blocks'][1]['status'], 'published')
        self.assertEqual(len(PARTICLE_LESSONS), 8)
        slugs = json.loads((ROOT / 'tools/course/lesson-slugs.json').read_text())
        self.assertEqual(len({row['stem'] for row in PARTICLE_LESSONS}), 8)
        for number in range(6, 14):
            row = PARTICLE_LESSONS[number - 6]
            article_slug = particle_slug(number, row['stem'])
            cover = ROOT / f'web/static/course/chemistry/visuals/{number:02d}-particles.svg'
            self.assertEqual(cover.read_text(), visual(number))
            ET.parse(cover)
            for lang in LANGS:
                suffix = '' if lang == 'ru' else '-' + lang
                path = LESSON_DIR / f'{row["stem"]}{suffix}.md'
                body = path.read_text()
                self.assertEqual(body, render_particles(number, lang))
                self.assertGreater(len(body.split()), 350, path)
                self.assertEqual(slugs[f'chemistry/{row["stem"]}{suffix}'], [article_slug, lang])
                self.assertIn(f'/read/chemistry-{number-1:02d}-', body)
                if number < 13:
                    self.assertIn(f'/read/chemistry-{number+1:02d}-', body)
                self.assertIn('![' , body)
                for linked in re.findall(r'!\[[^]]+\]\((/static/[^)]+)\)', body):
                    self.assertTrue((ROOT / 'web' / linked.lstrip('/')).is_file(), linked)
                map_path = MAP_DIR / f'map-{number:02d}-{lang}.svg'
                self.assertEqual(map_path.read_text(), diagram(number, lang))
                ET.parse(map_path)
                self.assertNotIn('<audio', body)

    def test_particle_science_and_release_guards(self):
        for lang, terms in {
            'ru': ('NaCl', 'H₂O', '155 г'),
            'kz': ('NaCl', 'H₂O', '155 г'),
            'en': ('NaCl', 'H₂O', '155 g'),
        }.items():
            endings = '' if lang == 'ru' else '-' + lang
            combined = '\n'.join((LESSON_DIR / f'{row["stem"]}{endings}.md').read_text()
                                 for row in PARTICLE_LESSONS)
            for term in terms:
                self.assertIn(term, combined)
        sql, expected = prepare_particles()
        self.assertEqual(len(expected), 27)
        self.assertEqual(len({(item['slug'], item['lang']) for item in expected}), 27)
        for clause in ('BEGIN;', 'COMMIT;', 'pg_advisory_xact_lock',
                       'First chemistry block incomplete', 'Particle lesson position occupied',
                       "s.slug='chemistry'"):
            self.assertIn(clause, sql)
        self.assertNotIn('syndicate_telegram', sql)


if __name__ == '__main__':
    unittest.main()
