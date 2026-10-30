"""Independent balances, localization, diagrams and release gates for 50–58."""
from decimal import Decimal as D
import json
from pathlib import Path
import re
import unittest
import xml.etree.ElementTree as ET
from .chemistry_solutions_data import LESSONS
from .build_chemistry_solutions import render,visual,map_svg,bridge_body,LANGS,LESSON_DIR,MAP_DIR
from .build_chemistry_quantities import render as previous_render
from .prepare_chemistry_solutions import prepare
from .prepare_chemistry_quantities import prepare as previous_prepare
from .prepare_chemistry import read_lesson
from .build_chemistry_particles import slug

ROOT=Path(__file__).resolve().parents[2]

class SolutionsTests(unittest.TestCase):
    def test_complete_localized_block_navigation_assets(self):
        self.assertEqual(len(LESSONS),9)
        routes=json.loads((ROOT/'tools/course/lesson-slugs.json').read_text())
        curriculum=json.loads((ROOT/'course/chemistry/curriculum.json').read_text())
        self.assertEqual(curriculum['blocks'][6]['range'],'50-58')
        self.assertIn(curriculum['blocks'][6]['status'],('prepared','published'))
        for n,row in enumerate(LESSONS,50):
            self.assertEqual(len({row[lang][0] for lang in LANGS}),3)
            for lang in LANGS:
                with self.subTest(n=n,lang=lang):
                    self.assertEqual(len(row[lang]),10)
                    suffix='' if lang=='ru' else '-'+lang
                    p=LESSON_DIR/f'{row["stem"]}{suffix}.md'
                    self.assertEqual(p.read_text(),render(n,lang))
                    title,summary,body=read_lesson(p)
                    self.assertTrue(title and summary)
                    self.assertGreater(len(body.split()),250)
                    self.assertEqual(routes[f'chemistry/{row["stem"]}{suffix}'],[slug(n,row['stem']),lang])
                    self.assertIn(f'/read/chemistry-{n-1:02d}-',body)
                    if n<58:self.assertIn(f'/read/chemistry-{n+1:02d}-',body)
                    else:self.assertNotIn('/read/chemistry-59-',body)
                    self.assertNotIn('<audio',body)
                    self.assertNotIn('<details',body)
                    self.assertIn('|---|---|---|',body)
                    assets=re.findall(r'!\[[^]]+\]\((/static/[^)]+)\)',body)
                    self.assertEqual(len(assets),2)
                    for a in assets:
                        self.assertIn(f'-{lang}.svg',a)
                        ET.parse(ROOT/'web'/a.lstrip('/'))
                    self.assertEqual((MAP_DIR/f'map-{n:02d}-{lang}.svg').read_text(),map_svg(n,lang))
                    self.assertEqual((MAP_DIR/f'visuals/{n:02d}-solutions-{lang}.svg').read_text(),visual(n,lang))

    def test_mass_fractions_and_component_balances(self):
        for solute,water,fraction in [('12','188','0.06'),('12','88','0.12'),('7.5','142.5','0.05'),('9','171','0.05')]:
            s,w=D(solute),D(water)
            self.assertEqual(s/(s+w),D(fraction))
            self.assertEqual(s/(s+w)+w/(s+w),1)
        self.assertEqual(D('0.08')*250,20)
        self.assertEqual(D(250)-20,230)
        self.assertEqual(D(8)/D('0.05'),160)
        self.assertEqual(D('0.04')*300,12)
        self.assertEqual(D(20)/270*100,D('7.407407407407407407407407407'))

    def test_solubility_solid_and_liquid_are_separate(self):
        capacity=D(30)*50/100
        added=D(20)
        self.assertEqual(capacity,15)
        self.assertEqual(added-capacity,5)
        self.assertEqual(50+capacity,65)
        self.assertEqual(50+added,70)
        self.assertEqual(D(50)-30,20)
        self.assertEqual(D(200)*30/100,60)
        self.assertEqual(D(70)-60,10)

    def test_molarity_conversions_dilution_and_unequal_mixing(self):
        molar_mass=D('58.5')
        self.assertEqual(D('5.85')/molar_mass,D('0.100'))
        self.assertEqual((D('2.925')/molar_mass)/(D(250)/1000),D('0.2'))
        self.assertEqual(D('0.30')*D('0.20')*molar_mass,D('3.5100'))
        self.assertEqual(D('0.50')*D('0.100')/D('0.250'),D('0.20'))
        self.assertEqual(D('0.75')*D('0.040')/D('0.120'),D('0.25'))
        # Unequal masses rule out the incorrect arithmetic mean of percentages.
        self.assertEqual((50*D('0.10')+150*D('0.04'))/200,D('0.055'))
        self.assertEqual(80*D('0.15')/120,D('0.10'))
        self.assertEqual(D(50)*D('0.12')/150,D('0.04'))
        self.assertEqual(D(12)/200*50/75,D('0.04'))
        self.assertLess(abs(D(10)/molar_mass/D('0.195')-D('0.877')),D('0.001'))

    def test_particle_counts_and_supplied_graph_not_fabricated_measurement(self):
        for lang in LANGS:
            nodes=list(ET.fromstring(visual(52,lang)))
            ions=[n.get('data-ion') for n in nodes if n.get('data-ion')]
            self.assertEqual(ions.count('Na⁺'),4)
            self.assertEqual(ions.count('Cl⁻'),4)
            points=[(n.get('data-temperature'),n.get('data-solubility')) for n in ET.fromstring(visual(53,lang)) if n.get('data-temperature')]
            self.assertEqual(points,[('20','30'),('40','50')])
            dots=[n for n in ET.fromstring(visual(53,lang)) if n.get('data-temperature')]
            # With the chart origin at (180,540), both axes have a linear scale.
            xdist=[D(n.get('cx'))-180 for n in dots]
            ydist=[540-D(n.get('cy')) for n in dots]
            self.assertEqual(xdist[1]/xdist[0],D(40)/20)
            self.assertEqual(ydist[1]/ydist[0],D(50)/30)
            water=[n.get('data-element') for n in ET.fromstring(visual(50,lang)) if n.get('data-element')]
            self.assertEqual(water.count('O'),1)
            self.assertEqual(water.count('H'),2)
            self.assertIn('q = 0',visual(50,lang))
        for n in (53,58):
            self.assertIn('вымышлен',render(n,'ru'))
            self.assertIn('ойдан',render(n,'kz').lower())
            self.assertIn('imaginary',render(n,'en'))

    def test_guarded_transaction_and_previous_release_unchanged(self):
        sql,expected=prepare();oldsql,_=previous_prepare()
        self.assertNotIn('/read/chemistry-50-',oldsql)
        self.assertEqual(len(expected),30)
        self.assertEqual(len({(x['slug'],x['lang']) for x in expected}),30)
        self.assertEqual(sql.count('INSERT INTO article_series_items'),9)
        self.assertEqual(sql.count('INSERT INTO article_translations'),27)
        self.assertEqual(sql.count('UPDATE article_translations SET body_md='),3)
        for clause in ('<> 49 THEN','<> 147 THEN','BETWEEN 500 AND 580',
                       'Solutions lesson position occupied','Solutions lesson belongs to another course',
                       'Solutions lesson slug belongs to another author','SET CONSTRAINTS ALL IMMEDIATE',
                       'Course lessons entered news queue'):
            self.assertIn(clause,sql)
        self.assertTrue(sql.startswith('BEGIN;'))
        self.assertTrue(sql.rstrip().endswith('COMMIT;'))
        self.assertLess(sql.index('INSERT INTO article_series_items'),sql.index('SET CONSTRAINTS ALL IMMEDIATE'))
        for lang in LANGS:
            self.assertNotIn('/read/chemistry-50-',previous_render(49,lang))
            self.assertTrue(bridge_body(lang).startswith(previous_render(49,lang).split('\n\n',2)[2].rstrip()))
            self.assertIn(f'/read/chemistry-50-water-and-polarity?lang={lang}',bridge_body(lang))
            self.assertIn(f'50-solutions-{lang}.svg',sql)

if __name__=='__main__':unittest.main()
