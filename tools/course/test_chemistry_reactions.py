"""Scientific, rendering and staged-release gates for chemistry block five."""
from collections import Counter
from decimal import Decimal
from functools import reduce
from math import gcd
import json
from pathlib import Path
import re
import unittest
import xml.etree.ElementTree as ET
from .chemistry_reactions_data import LESSONS, EQUATIONS
from .build_chemistry_reactions import render, visual, map_svg, bridge_body, LANGS, LESSON_DIR, MAP_DIR
from .build_chemistry_bonds import render as previous_render
from .build_chemistry_particles import slug
from .prepare_chemistry import read_lesson
from .prepare_chemistry_bonds import prepare as prepare_previous
from .prepare_chemistry_reactions import prepare

ROOT=Path(__file__).resolve().parents[2]


def atoms(formula):
    """Independent recursive parser: reject unconsumed or malformed formulas."""
    tokens=re.findall(r'[A-Z][a-z]?|[1-9][0-9]*|[()]',formula)
    if ''.join(tokens)!=formula:
        raise ValueError(formula)
    pos=0
    def group(nested=False):
        nonlocal pos
        count=Counter()
        while pos<len(tokens) and tokens[pos]!=')':
            token=tokens[pos];pos+=1
            if token=='(':
                sub=group(True)
                if pos>=len(tokens) or tokens[pos]!=')':raise ValueError(formula)
                pos+=1
            elif re.fullmatch('[A-Z][a-z]?',token):sub=Counter({token:1})
            else:raise ValueError(formula)
            multiplier=1
            if pos<len(tokens) and tokens[pos].isdigit():multiplier=int(tokens[pos]);pos+=1
            for element,quantity in sub.items():count[element]+=quantity*multiplier
        if not count:raise ValueError(formula)
        return count
    result=group()
    if pos!=len(tokens):raise ValueError(formula)
    return result


def side_counts(side):
    result=Counter()
    for coefficient,formula in side:
        for element,count in atoms(formula).items():result[element]+=coefficient*count
    return result


class ReactionTests(unittest.TestCase):
    def test_equations_preserve_each_element_and_minimal_coefficients(self):
        for name,(left,right) in EQUATIONS.items():
            with self.subTest(name=name):
                self.assertEqual(side_counts(left),side_counts(right))
                self.assertEqual(reduce(gcd,[n for n,_ in left+right]),1)
        self.assertEqual(atoms('Fe2(SO4)3'),{'Fe':2,'S':3,'O':12})
        self.assertEqual(side_counts(EQUATIONS['brackets'][0]),{'Fe':2,'S':3,'Na':6,'O':18,'H':6})
        self.assertEqual(side_counts(EQUATIONS['calcium_carbonate'][0]),{'Ca':1,'C':1,'O':3,'H':2,'Cl':2})
        for bad in ('H0','2H2O','Ca(OH','H2O)','()','Na+','H2O garbage'):
            with self.assertRaises(ValueError):atoms(bad)

    def test_all_locales_routes_assets_and_navigation(self):
        self.assertEqual(len(LESSONS),9)
        curriculum=json.loads((ROOT/'course/chemistry/curriculum.json').read_text())
        self.assertIn(curriculum['blocks'][4]['status'],('prepared','published'))
        routes=json.loads((ROOT/'tools/course/lesson-slugs.json').read_text())
        for n,row in enumerate(LESSONS,32):
            for lang in LANGS:
                with self.subTest(n=n,lang=lang):
                    self.assertEqual(len(row[lang]),10)
                    suffix='' if lang=='ru' else '-'+lang
                    path=LESSON_DIR/f'{row["stem"]}{suffix}.md'
                    self.assertEqual(path.read_text(),render(n,lang))
                    title,summary,body=read_lesson(path)
                    self.assertTrue(title and summary)
                    self.assertGreater(len(body.split()),300)
                    self.assertEqual(routes[f'chemistry/{row["stem"]}{suffix}'],[slug(n,row['stem']),lang])
                    self.assertIn(f'/read/chemistry-{n-1:02d}-',body)
                    if n<40:self.assertIn(f'/read/chemistry-{n+1:02d}-',body)
                    else:self.assertNotIn('/read/chemistry-41-',body)
                    for heading in {'ru':('Задание','Ответы'),'kz':('Тапсырма','Жауаптар'),'en':('Exercise','Answers')}[lang]:
                        self.assertIn('## '+heading+'\n',body)
                    self.assertNotIn('<audio',body)
                    self.assertNotIn('<details',body)
                    self.assertNotRegex(body,r'\d, \d{2}\b')
                    self.assertIn('|---|---|---|',body)
                    for linked in re.findall(r'!\[[^]]+\]\((/static/[^)]+)\)',body):
                        asset=ROOT/'web'/linked.lstrip('/')
                        ET.parse(asset)
                    self.assertEqual((MAP_DIR/f'map-{n:02d}-{lang}.svg').read_text(),map_svg(n,lang))
                    self.assertEqual((MAP_DIR/f'visuals/{n:02d}-reactions-{lang}.svg').read_text(),visual(n,lang))

    def test_particle_diagram_atom_counts_are_not_molecule_counts(self):
        for lang in LANGS:
            root=ET.fromstring(visual(34,lang))
            for side in ('before','after'):
                counts=Counter(node.get('data-element') for node in root if node.get('data-side')==side)
                self.assertEqual(counts,{'H':4,'O':2})

    def test_mass_and_energy_examples(self):
        D=Decimal
        self.assertEqual(D('97')+D('3'),D('100'))
        self.assertEqual(D('10')+D('2'),D('12'))
        self.assertEqual(D('150.00')-D('149.56'),D('0.44'))
        self.assertEqual(D('64.00')-D('61.50'),D('2.50'))
        self.assertEqual(D('80.00')-D('79.20'),D('0.80'))
        self.assertLess(abs(D('150.01')-D('150.00')),D('0.02'))
        self.assertEqual(100-140,-40)
        self.assertEqual(120-80,40)
        self.assertEqual(50-65,-15)
        for lang,phrase in [('ru','Вымышленные'),('kz','Ойдан берілген'),('en','Invented')]:
            self.assertIn(phrase,visual(38,lang))
            self.assertIn('0.44' if lang=='en' else '0,44',render(40,lang))

    def test_release_bridges_only_when_new_block_is_attached(self):
        previous_sql,previous_manifest=prepare_previous()
        self.assertNotIn('/read/chemistry-32-',previous_sql)
        sql,manifest=prepare()
        self.assertEqual(len(manifest),30)
        self.assertEqual(len({(x['slug'],x['lang']) for x in manifest}),30)
        self.assertEqual(sql.count('INSERT INTO article_series_items'),9)
        self.assertEqual(sql.count('INSERT INTO article_translations'),27)
        self.assertEqual(sql.count('UPDATE article_translations SET body_md='),3)
        self.assertTrue(sql.startswith('BEGIN;'))
        self.assertTrue(sql.rstrip().endswith('COMMIT;'))
        for clause in ('<> 31 THEN','<> 93 THEN','BETWEEN 320 AND 400','Reactions lesson position occupied',
                       'Reactions lesson belongs to another course','Reactions lesson slug belongs to another author',
                       'SET CONSTRAINTS ALL IMMEDIATE','Course lessons entered news queue'):
            self.assertIn(clause,sql)
        self.assertLess(sql.index('INSERT INTO article_series_items'),sql.index('SET CONSTRAINTS ALL IMMEDIATE'))
        for lang in LANGS:
            self.assertNotIn('/read/chemistry-32-',previous_render(31,lang))
            self.assertIn('/read/chemistry-32-',bridge_body(lang))
            # Bridge contains the complete previous body, not a replacement lesson.
            self.assertTrue(bridge_body(lang).startswith(previous_render(31,lang).split('\n\n',2)[2].rstrip()))

if __name__=='__main__':unittest.main()
