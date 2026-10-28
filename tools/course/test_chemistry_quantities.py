"""Numerical, entity-count, navigation and release checks for lessons 41–49."""
from decimal import Decimal as D
import json
from pathlib import Path
import re
import unittest
import xml.etree.ElementTree as ET

from .chemistry_quantities_data import LESSONS
from .build_chemistry_quantities import render, visual, map_svg, bridge_body, LANGS, LESSON_DIR, MAP_DIR
from .build_chemistry_reactions import render as previous_render
from .prepare_chemistry_quantities import prepare
from .prepare_chemistry_reactions import prepare as prepare_previous
from .prepare_chemistry import read_lesson
from .build_chemistry_particles import slug
from .test_chemistry_reactions import atoms, side_counts

ROOT=Path(__file__).resolve().parents[2]


class QuantitiesTests(unittest.TestCase):
    def test_all_translations_and_assets_form_a_complete_block(self):
        self.assertEqual(len(LESSONS),9)
        routes=json.loads((ROOT/'tools/course/lesson-slugs.json').read_text())
        curriculum=json.loads((ROOT/'course/chemistry/curriculum.json').read_text())
        self.assertEqual(curriculum['blocks'][5]['range'],'41-49')
        self.assertIn(curriculum['blocks'][5]['status'],('prepared','published'))
        for number,row in enumerate(LESSONS,41):
            for lang in LANGS:
                with self.subTest(number=number,lang=lang):
                    self.assertEqual(len(row[lang]),10)
                    suffix='' if lang=='ru' else '-'+lang
                    path=LESSON_DIR/f'{row["stem"]}{suffix}.md'
                    self.assertEqual(path.read_text(),render(number,lang))
                    title,summary,body=read_lesson(path)
                    self.assertTrue(title and summary)
                    self.assertGreater(len(body.split()),250)
                    self.assertEqual(routes[f'chemistry/{row["stem"]}{suffix}'],[slug(number,row['stem']),lang])
                    self.assertIn(f'/read/chemistry-{number-1:02d}-',body)
                    if number<49:self.assertIn(f'/read/chemistry-{number+1:02d}-',body)
                    else:self.assertNotIn('/read/chemistry-50-',body)
                    self.assertNotIn('<audio',body)
                    self.assertNotIn('<details',body)
                    self.assertIn('|---|---|---|',body)
                    linked=re.findall(r'!\[[^]]+\]\((/static/[^)]+)\)',body)
                    self.assertEqual(len(linked),2)
                    for link in linked:
                        self.assertIn(f'-{lang}.svg',link)
                        ET.parse(ROOT/'web'/link.lstrip('/'))
                    self.assertEqual((MAP_DIR/f'map-{number:02d}-{lang}.svg').read_text(),map_svg(number,lang))
                    self.assertEqual((MAP_DIR/f'visuals/{number:02d}-quantities-{lang}.svg').read_text(),visual(number,lang))
                    for heading in {'ru':('Задание','Ответы'),'kz':('Тапсырма','Жауаптар'),'en':('Exercise','Answers')}[lang]:
                        self.assertIn('## '+heading+'\n',body)

    def test_entity_counts_and_exact_definition(self):
        exact=D('6.02214076e23')
        self.assertEqual(exact,D(602214076000000000000000))
        self.assertEqual(D('0.5')*D('6.02e23'),D('3.01e23'))
        self.assertEqual(D('0.25')*2,D('0.50'))
        self.assertEqual(atoms('CaCl2')['Cl'],2)
        for lang in LANGS:
            root=ET.fromstring(visual(41,lang))
            elements=[node.get('data-element') for node in root if node.get('data-element')]
            self.assertEqual(elements.count('H'),2)
            self.assertEqual(elements.count('O'),1)
            self.assertIn('6.02214076 × 10²³ mol⁻¹',visual(41,lang))
            self.assertIn('CaCl₂',render(41,lang))

    def test_molar_masses_fractions_and_unit_conversions(self):
        weights={k:D(v) for k,v in {'H':'1','C':'12','O':'16','Na':'23','Cl':'35.5','Ca':'40'}.items()}
        def mass(formula):return sum(weights[e]*n for e,n in atoms(formula).items())
        for f,expected in {'H2O':'18','CO2':'44','Ca(OH)2':'74','NaCl':'58.5','CaCO3':'100','Na2CO3':'106','Na2O':'62','CaCl2':'111'}.items():
            self.assertEqual(mass(f),D(expected),f)
        self.assertEqual(D(2)/18+D(16)/18,1)
        self.assertEqual(D(36)*2/18,4)
        self.assertEqual(D(22)*12/44,6)
        self.assertEqual(D(25)*48/100,12)
        self.assertEqual(D(500)/1000,D('0.500'))
        self.assertEqual(D(9)/18,D('0.5'))
        self.assertEqual(D('0.20')*44,D('8.8'))
        self.assertEqual(D('1.204e23')/D('6.02e23'),D('0.20'))

    def test_equations_mass_balance_and_limiting_reactants(self):
        equations=[([(2,'H2'),(1,'O2')],[(2,'H2O')]),
                   ([(1,'CaCO3')],[(1,'CaO'),(1,'CO2')]),
                   ([(1,'Na2CO3'),(2,'HCl')],[(2,'NaCl'),(1,'H2O'),(1,'CO2')])]
        for left,right in equations:self.assertEqual(side_counts(left),side_counts(right))
        self.assertEqual(D(10)/100*44,D('4.4'))
        self.assertEqual(D(10)/100*56+D(10)/100*44,D(10))
        self.assertEqual(D('5.30')/106*2*D('58.5'),D('5.85'))
        for h,o,water,left_h,left_o in [('3','2','3','0','0.5'),('4','1','2','2','0'),('2','1','2','0','0'),('1','1','1','0','0.5')]:
            extent=min(D(h)/2,D(o))
            self.assertEqual(2*extent,D(water))
            self.assertEqual(D(h)-2*extent,D(left_h))
            self.assertEqual(D(o)-extent,D(left_o))
            self.assertGreaterEqual(D(h)-2*extent,0)
            self.assertGreaterEqual(D(o)-extent,0)

    def test_yield_and_final_case_use_purity_and_both_reactants(self):
        carbonate=D('13.25')*D('0.80')/106
        extent=min(carbonate,D('0.150')/2)
        theory=2*extent*D('58.5')
        actual=D('8.00')*D('0.8775')
        self.assertEqual(carbonate,D('0.100'))
        self.assertEqual(theory,D('8.775'))
        self.assertEqual(actual,D('7.020'))
        self.assertEqual(actual/theory*100,80)
        self.assertEqual((carbonate-extent)*106,D('2.65'))
        # The transfer case reverses the limiting reactant; it is not a scaled copy.
        carbonate2=D('6.625')*D('0.80')/106
        extent2=min(carbonate2,D('0.200')/2)
        self.assertEqual(extent2,carbonate2)
        self.assertEqual(2*extent2*D('58.5'),D('5.85'))
        self.assertEqual(D('5.00')*D('0.936')/D('5.85')*100,80)
        self.assertEqual(D('0.200')-2*extent2,D('0.100'))
        self.assertEqual(D('4.68')/D('5.85')*100,80)
        self.assertGreater(D('6.20')/D('5.85')*100,100)

    def test_transaction_is_complete_guarded_and_does_not_release_news(self):
        sql,expected=prepare()
        previous_sql,_=prepare_previous()
        self.assertNotIn('/read/chemistry-41-',previous_sql)
        self.assertEqual(len(expected),30)
        self.assertEqual(len({(x['slug'],x['lang']) for x in expected}),30)
        self.assertEqual(sql.count('INSERT INTO article_series_items'),9)
        self.assertEqual(sql.count('INSERT INTO article_translations'),27)
        self.assertEqual(sql.count('UPDATE article_translations SET body_md='),3)
        for clause in ('<> 40 THEN','<> 120 THEN','BETWEEN 410 AND 490',
                       'Quantities lesson position occupied','Quantities lesson belongs to another course',
                       'Quantities lesson slug belongs to another author','SET CONSTRAINTS ALL IMMEDIATE',
                       'Course lessons entered news queue'):
            self.assertIn(clause,sql)
        self.assertTrue(sql.startswith('BEGIN;'))
        self.assertTrue(sql.rstrip().endswith('COMMIT;'))
        self.assertLess(sql.index('INSERT INTO article_series_items'),sql.index('SET CONSTRAINTS ALL IMMEDIATE'))
        for lang in LANGS:
            self.assertNotIn('/read/chemistry-41-',previous_render(40,lang))
            self.assertTrue(bridge_body(lang).startswith(previous_render(40,lang).split('\n\n',2)[2].rstrip()))
            self.assertIn(f'/read/chemistry-41-mole-and-specified-entities?lang={lang}',bridge_body(lang))
            self.assertIn(f'41-quantities-{lang}.svg',sql)


if __name__=='__main__':unittest.main()
