"""Independent acid/base balances, model limits and complete release checks."""
from decimal import Decimal as D
import json
import math
from pathlib import Path
import unittest
import xml.etree.ElementTree as ET
from .chemistry_acids_data import LESSONS
from .build_chemistry_acids import render, visual, map_svg, bridge_body, LANGS, LESSON_DIR, MAP_DIR
from .build_chemistry_solutions import render as previous_render
from .prepare_chemistry_acids import prepare
from .prepare_chemistry_solutions import prepare as previous_prepare
from .build_chemistry_particles import slug
from .prepare_chemistry import read_lesson

ROOT=Path(__file__).resolve().parents[2]
NS={'s':'http://www.w3.org/2000/svg'}

class AcidsTests(unittest.TestCase):
    def test_complete_localized_block_and_navigation(self):
        routes=json.loads((ROOT/'tools/course/lesson-slugs.json').read_text())
        curriculum=json.loads((ROOT/'course/chemistry/curriculum.json').read_text())
        self.assertIn(curriculum['blocks'][7]['status'],('prepared','published'))
        self.assertEqual(len(LESSONS),9)
        for n,row in enumerate(LESSONS,59):
            for lang in LANGS:
                with self.subTest(number=n,lang=lang):
                    self.assertEqual(len(row[lang]),10)
                    suffix='' if lang=='ru' else '-'+lang
                    path=LESSON_DIR/f'{row["stem"]}{suffix}.md'
                    self.assertEqual(path.read_text(),render(n,lang))
                    title,summary,body=read_lesson(path)
                    self.assertTrue(title and summary)
                    self.assertGreater(len(body.split()),300)
                    self.assertEqual(routes['chemistry/'+row['stem']+suffix],[slug(n,row['stem']),lang])
                    self.assertNotIn('<audio',body)
                    self.assertNotIn('.mp3',body)
                    self.assertNotIn('/read/chemistry-68',body)
                    previous='58-solution-investigation-review' if n==59 else LESSONS[n-60]['stem']
                    self.assertIn('/read/'+slug(n-1,previous)+'?lang='+lang,body)
                    if n<67:self.assertIn('/read/'+slug(n+1,LESSONS[n-58]['stem'])+'?lang='+lang,body)
                    for path,text in [(MAP_DIR/f'map-{n}-{lang}.svg',map_svg(n,lang)),(MAP_DIR/f'visuals/{n}-acids-{lang}.svg',visual(n,lang))]:
                        self.assertEqual(path.read_text(),text)
                        root=ET.fromstring(text)
                        self.assertTrue(root.findall('.//s:text',NS))
                        self.assertTrue(root.find('s:title',NS).text)
                        width,height=map(float,root.get('viewBox').split()[2:])
                        for el in root.findall('.//s:text',NS):
                            self.assertGreaterEqual(float(el.get('x')),0)
                            self.assertLessEqual(float(el.get('x')),width)
                            self.assertGreaterEqual(float(el.get('y')),0)
                            self.assertLess(float(el.get('y')),height)

    def test_atoms_and_charge_in_proton_transfer_and_neutralisation(self):
        # Explicit species compositions count H, O, N, Cl, Na and charge;
        # these balances are independent of the rendered equation strings.
        species={'HCl':(1,0,0,1,0,0),'H2O':(2,1,0,0,0,0),'H3O':(3,1,0,0,0,1),
                 'Cl':(0,0,0,1,0,-1),'NH3':(3,0,1,0,0,0),'NH4':(4,0,1,0,0,1),
                 'OH':(1,1,0,0,0,-1),'HNO3':(1,3,1,0,0,0),'NO3':(0,3,1,0,0,-1),
                 'Na':(0,0,0,0,1,1)}
        equations=[({'HCl':1,'H2O':1},{'H3O':1,'Cl':1}),({'NH3':1,'H2O':1},{'NH4':1,'OH':1}),
                   ({'HNO3':1,'H2O':1},{'H3O':1,'NO3':1}),({'H3O':1,'OH':1},{'H2O':2}),
                   ({'H3O':1,'Cl':1,'Na':1,'OH':1},{'Na':1,'Cl':1,'H2O':2})]
        def total(side):return tuple(sum(species[s][i]*n for s,n in side.items()) for i in range(6))
        for left,right in equations:self.assertEqual(total(left),total(right))

    def test_neutralisation_amounts_and_excess_ph(self):
        cases=[('0.025','0.20','0.040','0.10','0.065','0.0010','0.0040',1.81),
               ('0.020','0.10','0.030','0.10','0.050','-0.0010','0.0020',12.30),
               ('0.020','0.10','0.018','0.10','0.038','0.0002','0.0018',2.28)]
        for va,ca,vb,cb,vfinal,excess,salt,ph in cases:
            a,b=D(va)*D(ca),D(vb)*D(cb)
            self.assertEqual(a-b,D(excess));self.assertEqual(min(a,b),D(salt))
            concentration=float(abs(a-b)/D(vfinal))
            actual=-math.log10(concentration) if a>b else 14+math.log10(concentration)
            self.assertAlmostEqual(actual,ph,places=2)
        self.assertEqual(D('0.0020')*D('58.5'),D('0.11700'))
        self.assertEqual(D('0.0018')*D('58.5'),D('0.10530'))
        self.assertEqual(D('0.0020')/D('0.040'),D('0.050'))
        # Equal volumes do not imply equivalence when concentrations differ.
        self.assertEqual(D('.010')*D('.20')-D('.010')*D('.10'),D('.0010'))

    def test_dilution_strength_and_water_limit(self):
        for v,c,final,expected_ph in [('.010','.010','.100',3),('.020','.001','.200',4),('.010','.001','.100',4)]:
            amount=D(v)*D(c);after=amount/D(final)
            self.assertEqual(after*D(final),amount)
            self.assertAlmostEqual(-math.log10(float(after)),expected_ph)
        self.assertEqual(D('.030')*D('.020'),D('.00060'))
        self.assertEqual(D('.010')*D('.10')/D('.00010'),D(10))
        self.assertEqual(D(10)**(-4)/D(10)**(-6),D(100))
        # Electroneutrality plus Kw for very dilute fully dissociated HCl.
        acid,kw=1e-8,1e-14
        h=(acid+math.sqrt(acid*acid+4*kw))/2
        self.assertAlmostEqual(h-kw/h,acid)
        self.assertLess(-math.log10(h),7)
        self.assertGreater(-math.log10(h),6.9)
        self.assertNotEqual(round(-math.log10(h)),8)

    def test_particle_diagram_and_proportional_bars(self):
        root=ET.fromstring(visual(62,'en'))
        units=root.findall('.//*[@data-acid-unit]')
        self.assertEqual(len(units),10)
        labels=[e.text for e in root.findall('.//s:text',NS)]
        self.assertEqual(labels.count('HA'),9)
        self.assertEqual(labels.count('A⁻'),1)
        self.assertEqual(len(root.findall('.//*[@data-hydronium]')),1)
        root=ET.fromstring(visual(64,'en'))
        bars=root.findall('.//*[@data-amount-mmol]')
        self.assertEqual([int(e.get('data-amount-mmol')) for e in bars],[5,4,1])
        self.assertEqual([float(e.get('width'))/float(e.get('data-amount-mmol')) for e in bars],[140]*3)
        labels=[e.text for e in ET.fromstring(visual(61,'en')).findall('.//s:text',NS)]
        self.assertTrue(all(str(i) in labels for i in range(15)))

    def test_transaction_guards_bridge_and_news_exclusion(self):
        sql,manifest=prepare()
        self.assertTrue(sql.startswith('BEGIN;'));self.assertTrue(sql.endswith('COMMIT;\n'))
        self.assertEqual(len(manifest),30)
        self.assertEqual(sql.count('INSERT INTO article_translations'),27)
        self.assertEqual(sql.count('INSERT INTO article_series_items'),9)
        self.assertIn('<> 58',sql);self.assertIn('<> 174',sql)
        self.assertIn('position BETWEEN 590 AND 670',sql)
        self.assertIn('another author',sql);self.assertIn('another course',sql)
        self.assertIn('SET CONSTRAINTS ALL IMMEDIATE',sql)
        self.assertLess(sql.rfind('INSERT INTO article_series_items'),sql.index('SET CONSTRAINTS ALL IMMEDIATE'))
        self.assertIn("name='syndicate_telegram'",sql)
        self.assertIn('IS DISTINCT FROM',sql)
        old_sql,_=previous_prepare()
        self.assertNotIn('/read/chemistry-59',old_sql)
        for lang in LANGS:
            old=previous_render(58,lang).split('\n\n',2)[2].rstrip()
            self.assertTrue(bridge_body(lang).startswith(old+'\n\n'))
            self.assertIn('/read/chemistry-59',bridge_body(lang))

if __name__=='__main__':unittest.main()
