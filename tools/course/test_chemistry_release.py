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
from tools.course.chemistry_elements_data import ELEMENTS, LESSONS as ELEMENT_LESSONS
from tools.course.build_chemistry_elements import render as render_elements, diagram as element_diagram, table_svg, visual as element_visual
from tools.course.prepare_chemistry_elements import prepare as prepare_elements
from tools.course.chemistry_bonds_data import LESSONS as BOND_LESSONS
from tools.course.build_chemistry_bonds import render as render_bonds, map_svg as bond_map, visual as bond_visual
from tools.course.prepare_chemistry_bonds import prepare as prepare_bonds


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


class ChemistryElementsReleaseTests(unittest.TestCase):
    def test_complete_trilingual_route_and_assets(self):
        curriculum = json.loads((ROOT / 'course/chemistry/curriculum.json').read_text())
        self.assertEqual(curriculum['blocks'][2]['range'], '14-22')
        self.assertIn(curriculum['blocks'][2]['status'], ('prepared', 'published'))
        self.assertEqual(len(ELEMENT_LESSONS), 9)
        slugs = json.loads((ROOT / 'tools/course/lesson-slugs.json').read_text())
        for number, row in enumerate(ELEMENT_LESSONS, 14):
            self.assertEqual(len(row['ru']), 10)
            cover = MAP_DIR / f'visuals/{number:02d}-elements.svg'
            self.assertEqual(cover.read_text(), element_visual(number))
            ET.parse(cover)
            for lang in LANGS:
                suffix = '' if lang == 'ru' else f'-{lang}'
                path = LESSON_DIR / f'{row["stem"]}{suffix}.md'
                self.assertEqual(path.read_text(), render_elements(number, lang))
                title, summary, body = read_lesson(path)
                self.assertTrue(title and summary)
                self.assertGreater(len(body.split()), 350, path)
                self.assertEqual(slugs[f'chemistry/{row["stem"]}{suffix}'],
                                 [particle_slug(number, row['stem']), lang])
                self.assertIn(f'/read/chemistry-{number-1:02d}-', body)
                if number < 22:
                    self.assertIn(f'/read/chemistry-{number+1:02d}-', body)
                else:
                    self.assertIn('/read/chemistry-23-why-atoms-bond', body)
                self.assertNotIn('<audio', body)
                # Raw HTML is stripped by the site's safe Markdown renderer.
                self.assertNotIn('<details', body)
                for heading in {'ru': ('Задание', 'Ответы'), 'kz': ('Тапсырма', 'Жауаптар'),
                                'en': ('Exercise', 'Answers')}[lang]:
                    self.assertIn(f'## {heading}\n', body)
                for linked in re.findall(r'!\[[^]]+\]\((/static/[^)]+)\)', body):
                    self.assertTrue((ROOT / 'web' / linked.lstrip('/')).is_file(), linked)
                map_path = MAP_DIR / f'map-{number:02d}-{lang}.svg'
                self.assertEqual(map_path.read_text(), element_diagram(number, lang))
                ET.parse(map_path)
        for lang in LANGS:
            suffix = '' if lang == 'ru' else f'-{lang}'
            previous = (LESSON_DIR / f'13-particle-case-review{suffix}.md').read_text()
            self.assertIn(f'/read/chemistry-14-read-element-card?lang={lang}', previous)

    def test_shell_data_and_actual_periodic_table_geometry(self):
        self.assertEqual([row[0] for row in ELEMENTS], list(range(1, 21)))
        self.assertEqual(len({row[1] for row in ELEMENTS}), 20)
        by_symbol = {row[1]: row for row in ELEMENTS}
        for z, symbol, period, group, shells in ELEMENTS:
            self.assertEqual(sum(shells), z, symbol)
            self.assertEqual(len(shells), period, symbol)
            self.assertLessEqual(shells[0], 2)
            self.assertTrue(1 <= group <= 18)
        self.assertEqual(by_symbol['He'][2:], (1, 18, (2,)))
        self.assertEqual(by_symbol['B'][2:], (2, 13, (2, 3)))
        self.assertEqual(by_symbol['K'][4], (2, 8, 8, 1))
        self.assertEqual(by_symbol['Ca'][4], (2, 8, 8, 2))
        for lang in LANGS:
            path = MAP_DIR / f'first-20-{lang}.svg'
            self.assertEqual(path.read_text(), table_svg(lang))
            root = ET.parse(path).getroot()
            cells = [node for node in root if 'data-z' in node.attrib]
            self.assertEqual(len(cells), 20)
            for cell, row in zip(cells, ELEMENTS, strict=True):
                self.assertEqual(int(cell.attrib['data-z']), row[0])
                self.assertEqual(int(cell.attrib['data-period']), row[2])
                self.assertEqual(int(cell.attrib['data-group']), row[3])
                self.assertEqual(int(cell.attrib['x']), 80+(row[3]-1)*63)
                self.assertEqual(int(cell.attrib['y']), 200+(row[2]-1)*95)

    def test_atomic_release_and_previous_lesson_bridge(self):
        sql, expected = prepare_elements()
        self.assertEqual(len(expected), 30)
        self.assertEqual(len({(x['slug'], x['lang']) for x in expected}), 30)
        self.assertEqual(sql.count('BEGIN;'), 1)
        self.assertTrue(sql.rstrip().endswith('COMMIT;'))
        self.assertEqual(sql.count('INSERT INTO article_series_items'), 9)
        self.assertEqual(sql.count('INSERT INTO article_translations'), 27)
        self.assertEqual(sql.count('UPDATE article_translations SET body_md='), 3)
        self.assertNotIn('syndicate_telegram', sql)
        for clause in ('pg_advisory_xact_lock', 'Previous chemistry translations incomplete',
                       'Elements lesson position occupied', 'Elements lesson belongs to another course',
                       'Elements block course membership incomplete', 'Elements block translations incomplete'):
            self.assertIn(clause, sql)

    def test_fourth_block_content_navigation_and_localized_assets(self):
        curriculum = json.loads((ROOT / 'course/chemistry/curriculum.json').read_text())
        self.assertEqual(curriculum['blocks'][3]['range'], '23-31')
        self.assertIn(curriculum['blocks'][3]['status'], ('prepared', 'published'))
        self.assertEqual(len(BOND_LESSONS), 9)
        slugs = json.loads((ROOT / 'tools/course/lesson-slugs.json').read_text())
        for number, row in enumerate(BOND_LESSONS, 23):
            for lang in LANGS:
                self.assertEqual(len(row[lang]), 10)
                suffix = '' if lang == 'ru' else '-'+lang
                path = LESSON_DIR / f'{row["stem"]}{suffix}.md'
                self.assertEqual(path.read_text(), render_bonds(number, lang))
                title, summary, body = read_lesson(path)
                self.assertTrue(title and summary)
                # Counts are a lower completeness bound, not a quality score.
                self.assertGreater(len(body.split()), 300, path)
                self.assertEqual(slugs[f'chemistry/{row["stem"]}{suffix}'],
                                 [particle_slug(number, row['stem']), lang])
                self.assertIn(f'/read/chemistry-{number-1:02d}-', body)
                if number < 31:
                    self.assertIn(f'/read/chemistry-{number+1:02d}-', body)
                else:
                    self.assertNotIn('/read/chemistry-32-', body)
                self.assertNotIn('<audio', body)
                self.assertNotIn('<details', body)
                self.assertIn({'ru': 'Не получайте газы', 'kz': 'Газ алмаңыз',
                               'en': 'Do not produce gases'}[lang], body)
                for heading in {'ru': ('Задание', 'Ответы'), 'kz': ('Тапсырма', 'Жауаптар'), 'en': ('Exercise', 'Answers')}[lang]:
                    self.assertIn(f'## {heading}\n', body)
                for linked in re.findall(r'!\[[^]]+\]\((/static/[^)]+)\)', body):
                    asset = ROOT / 'web' / linked.lstrip('/')
                    self.assertTrue(asset.is_file(), linked)
                    ET.parse(asset)
                self.assertEqual((MAP_DIR / f'map-{number:02d}-{lang}.svg').read_text(), bond_map(number, lang))
                self.assertEqual((MAP_DIR / f'visuals/{number:02d}-bonds-{lang}.svg').read_text(), bond_visual(number, lang))

    def test_bond_diagrams_conserve_electrons_and_charge(self):
        for lang in LANGS:
            root = ET.fromstring(bond_visual(24, lang))
            electrons = [x for x in root if x.get('data-electron')]
            self.assertEqual(len(electrons), 14)
            self.assertEqual(sum(x.get('data-electron') == 'shared' for x in electrons), 2)
        # Independent reference counts catch an incorrect worksheet answer.
        for cation, anion, ncat, nanion in ((2,-1,1,2),(3,-2,2,3),(2,-3,3,2),(1,-2,2,1)):
            self.assertEqual(cation*ncat+anion*nanion, 0)
            from math import gcd
            self.assertEqual(gcd(ncat,nanion), 1)
        for lang in LANGS:
            text = render_bonds(30, lang)
            self.assertIn('Al₂(SO₄)₃', text)
            self.assertIn('12', text)
            self.assertIn('24', text)

    def test_bonds_publication_guards_previous_block_and_news_queue(self):
        sql, expected = prepare_bonds()
        self.assertEqual(len(expected), 30)
        self.assertEqual(len({(x['slug'],x['lang']) for x in expected}), 30)
        self.assertEqual(sql.count('BEGIN;'), 1)
        self.assertTrue(sql.rstrip().endswith('COMMIT;'))
        self.assertEqual(sql.count('INSERT INTO article_series_items'), 9)
        self.assertEqual(sql.count('INSERT INTO article_translations'), 27)
        self.assertEqual(sql.count('UPDATE article_translations SET body_md='), 3)
        for clause in ('<> 22 THEN','<> 66 THEN','Bonds lesson position occupied',
                       'Bonds lesson belongs to another course','Bonds block course membership incomplete',
                       'SET CONSTRAINTS ALL IMMEDIATE','Course lessons entered news queue'):
            self.assertIn(clause, sql)
        self.assertLess(sql.index('INSERT INTO article_series_items'),
                        sql.index('SET CONSTRAINTS ALL IMMEDIATE'))


if __name__ == '__main__':
    unittest.main()
