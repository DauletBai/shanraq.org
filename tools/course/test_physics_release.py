"""Checks that protect the first physics block's scientific and publishing contract."""
import json
from pathlib import Path
import re
import unittest
import xml.etree.ElementTree as ET

from tools.course.build_physics_block import DATA, LESSONS, MAPS, ROOT
from tools.course.physics_motion_data import MOTION
from tools.course.physics_forces_data import FORCES
from tools.course.physics_fluids_data import FLUIDS
from tools.course.build_physics_fluids import DIAGRAM
from tools.course.physics_heat_data import HEAT
from tools.course.physics_electricity_data import ELECTRICITY
from tools.course.physics_magnetism_data import MAGNETISM
from tools.course.physics_waves_data import WAVES
from tools.course.physics_astronomy_data import ASTRONOMY
from tools.course.build_physics_magnetism import NOTES
from tools.course.prepare_physics import prepare, read_lesson


class PhysicsReleaseTests(unittest.TestCase):
    def test_curriculum_has_contiguous_route_and_complete_first_block(self):
        data = json.loads((ROOT / "course/physics/curriculum.json").read_text())
        self.assertEqual(data["lesson_count"], 100)
        self.assertEqual(data["publication_policy"], "complete-block-only")
        self.assertEqual(data["languages"], ["ru", "kz", "en"])
        next_number = 1
        for block in data["blocks"]:
            first, last = map(int, block["range"].split("-"))
            self.assertEqual(first, next_number)
            self.assertEqual(last - first + 1, block["lesson_count"])
            next_number = last + 1
        self.assertEqual(next_number, 101)
        self.assertEqual([row["number"] for row in data["first_block_lessons"]], list(range(1, 9)))
        self.assertEqual([row["number"] for row in data["second_block_lessons"]], list(range(9, 19)))
        self.assertEqual(data["blocks"][1]["status"], "published")
        self.assertEqual([row["number"] for row in data["third_block_lessons"]], list(range(19, 29)))
        self.assertEqual(data["blocks"][2]["status"], "published")
        self.assertEqual([row["number"] for row in data["fourth_block_lessons"]], list(range(29, 37)))
        self.assertEqual(data["blocks"][3]["status"], "published")
        self.assertEqual([row["number"] for row in data["fifth_block_lessons"]], list(range(37, 47)))
        self.assertEqual(data["blocks"][4]["status"], "published")
        self.assertTrue((ROOT / "web" / data["blocks"][4]["cover"].lstrip("/")).is_file())
        self.assertEqual([row["number"] for row in data["sixth_block_lessons"]], list(range(47, 57)))
        self.assertEqual(data["blocks"][5]["status"], "published")
        self.assertTrue((ROOT / "web" / data["blocks"][5]["cover"].lstrip("/")).is_file())
        self.assertEqual([row["number"] for row in data["seventh_block_lessons"]], list(range(57, 67)))
        self.assertEqual(data["blocks"][6]["status"], "published")
        self.assertTrue((ROOT / "web" / data["blocks"][6]["cover"].lstrip("/")).is_file())
        self.assertEqual([row["number"] for row in data["eighth_block_lessons"]], list(range(67, 75)))
        self.assertEqual(data["blocks"][7]["status"], "published")
        self.assertTrue((ROOT / "web" / data["blocks"][7]["cover"].lstrip("/")).is_file())
        self.assertEqual([row["number"] for row in data["ninth_block_lessons"]], list(range(75, 85)))
        self.assertEqual(data["blocks"][8]["status"], "published")
        self.assertTrue((ROOT / "web" / data["blocks"][8]["cover"].lstrip("/")).is_file())

    def test_every_language_has_the_same_experiment_and_working_assets(self):
        slugs = json.loads((ROOT / "tools/course/lesson-slugs.json").read_text())
        route = [*DATA, *MOTION, *FORCES, *FLUIDS, *HEAT, *ELECTRICITY, *MAGNETISM, *WAVES, *ASTRONOMY]
        for index, row in enumerate(route, 1):
            for lang in ("ru", "kz", "en"):
                suffix = "" if lang == "ru" else f"-{lang}"
                path = LESSONS / f"{row['stem']}{suffix}.md"
                self.assertTrue(path.is_file(), path)
                slug = f"physics-{index:02d}-{row['stem'][3:]}"
                self.assertEqual(slugs.get(f"physics/{row['stem']}{suffix}"), [slug, lang])
                title, summary, body = read_lesson(path)
                self.assertTrue(title and summary)
                self.assertGreater(len(body.split()), 280, path)
                ref = re.search(r"!\[[^]]+\]\((/static/course/physics/[^)]+)\)", body)
                self.assertIsNotNone(ref, path)
                asset = ROOT / "web" / ref.group(1).lstrip("/")
                self.assertTrue(asset.is_file(), asset)
                ET.parse(asset)
                if index < 84:
                    next_slug = f"physics-{index+1:02d}-{route[index]['stem'][3:]}"
                    self.assertIn(f"/read/{next_slug}?lang={lang}", body)

    def test_motion_maps_match_worked_examples(self):
        expected = {
            10: {"10", "2"}, 11: {"7", "2", "5"},
            12: {"12", "4", "3"}, 15: {"0", "2", "4", "8"},
            17: {"0", "2", "4", "6"}, 18: {"0", "2", "4", "6"},
        }
        for number, values in expected.items():
            row = MOTION[number - 9]
            for lang in ("ru", "kz", "en"):
                path = MAPS / f"map-{row['stem']}-{lang}.svg"
                svg = path.read_text()
                for value in values:
                    self.assertIn(value, svg, (path, value))
                ET.parse(path)

    def test_forces_work_and_energy_examples_are_consistent(self):
        self.assertEqual(3 - 1, 2)  # horizontal net force, N
        self.assertEqual((3 - 1) / 0.5, 4)  # acceleration, m/s²
        self.assertEqual(3 * 2 - 1 * 2, 4)  # total work, J
        self.assertEqual(0.5 * 0.5 * 4**2, 4)  # final kinetic energy, J
        self.assertEqual(0.5 * 10 * 0.8, 4)  # potential energy, J
        for row in FORCES:
            for lang in ("ru", "kz", "en"):
                text = (MAPS / f"map-{row['stem']}-{lang}.svg").read_text()
                self.assertIn(row[lang][0], text)
                ET.fromstring(text)

    def test_fluid_calculations_and_localized_maps(self):
        self.assertEqual(20 / .02, 1000)
        self.assertEqual(4 / .0002, 20000)
        self.assertEqual(1000 * 10 * .2, 2000)
        self.assertEqual(1000 * 10 * .0001, 1)
        self.assertAlmostEqual(3 / (1000 * 10), .0003)
        self.assertEqual(4 * 10 / 2, 20)
        self.assertAlmostEqual(4 * .05, 20 * .01)
        for row in FLUIDS:
            for lang in ('ru', 'kz', 'en'):
                svg = MAPS / f"map-{row['stem']}-{lang}.svg"
                root = ET.parse(svg).getroot()
                self.assertEqual(root.attrib['viewBox'], '0 0 1600 900')
                self.assertIn(row[lang][0], ''.join(root.itertext()))
                map_text=' '.join(root.itertext())
                for label in DIAGRAM[row['kind']][lang]:
                    for token in label.split():
                        self.assertIn(token, map_text, (svg, token))

    def test_heat_models_and_translations_match_their_maps(self):
        self.assertEqual(50 + 20 - 30, 40)
        self.assertEqual(.2 * 4200 * 5, 4200)
        self.assertAlmostEqual(.1 * 334, 33.4)
        self.assertEqual(600 - 200 - 100, 300)
        self.assertEqual(600 - 150 - 100, 350)
        for row in HEAT:
            for lang in ("ru", "kz", "en"):
                path = MAPS / f"map-{row['stem']}-{lang}.svg"
                root = ET.parse(path).getroot()
                self.assertEqual(root.attrib["viewBox"], "0 0 1600 900")
                visible = " ".join(root.itertext())
                self.assertIn(row[lang][0], visible)
                for label in row["map"][lang]:
                    self.assertIn(label, visible)

    def test_electricity_models_and_translations_match_maps(self):
        self.assertEqual(6/3, 2)
        self.assertEqual(12/4, 3)
        self.assertAlmostEqual(3/300, .01)
        self.assertAlmostEqual(3/300+3/600, .015)
        self.assertAlmostEqual(3*.015*60, 2.7)
        for row in ELECTRICITY:
            for lang in ("ru", "kz", "en"):
                path=MAPS/f"map-{row['stem']}-{lang}.svg"
                root=ET.parse(path).getroot()
                self.assertEqual(root.attrib["viewBox"], "0 0 1600 900")
                self.assertIn(row[lang][0], " ".join(root.itertext()))
                body=(LESSONS/f"{row['stem']}{'' if lang=='ru' else '-'+lang}.md").read_text()
                self.assertIn(row[lang][0], body)

    def test_magnetism_block_distinguishes_examples_from_observations(self):
        self.assertEqual(10, 6 + 4)
        self.assertEqual(100, 60 + 40)
        self.assertEqual('10' + '01', '1001')
        for row in MAGNETISM:
            for lang in ('ru', 'kz', 'en'):
                svg = MAPS / f"map-{row['stem']}-{lang}.svg"
                root = ET.parse(svg).getroot()
                self.assertEqual(root.attrib['viewBox'], '0 0 1600 900')
                visible = ' '.join(root.itertext())
                for note in NOTES[row['kind']][lang]:
                    for part in note.split():
                        self.assertIn(part, visible, (svg, part))
                lesson = (LESSONS / f"{row['stem']}{'' if lang=='ru' else '-'+lang}.md").read_text()
                self.assertIn(row[lang][0], lesson)
        # Voltage may exist in an open induction coil, but current in the
        # capstone requires a closed path and the correct reed-switch state.
        for lang, words in {'ru':('разомкнутой', 'нормально замкнутого'), 'kz':('ашық', 'жабық'), 'en':('open coil', 'normally closed')}.items():
            suffix = '' if lang == 'ru' else '-'+lang
            self.assertIn(words[0], (LESSONS / f'63-changing-field-induction{suffix}.md').read_text())
            self.assertIn(words[1], (LESSONS / f'66-magnetism-communication-checkpoint{suffix}.md').read_text())

    def test_waves_block_numbers_and_visual_models(self):
        self.assertEqual(8 / 4, 2)  # cycles per second
        self.assertEqual(1.5 * 2, 3)  # wavelength times frequency
        self.assertEqual(1 / 10 - 1 / 20, 1 / 20)  # thin-lens model
        self.assertEqual(2 * 10, 20)  # reflected ray turns twice mirror angle
        for row in WAVES:
            for lang in ("ru", "kz", "en"):
                image = MAPS / f"map-{row['stem']}-{lang}.svg"
                root = ET.parse(image).getroot()
                self.assertEqual(root.attrib["viewBox"], "0 0 1600 900")
                visible = " ".join(root.itertext())
                for label in row["map"][lang]:
                    self.assertIn(label, visible)
                lesson = (LESSONS / f"{row['stem']}{'' if lang == 'ru' else '-'+lang}.md").read_text()
                self.assertIn(row[lang][0], lesson)
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else "-" + lang
            link = f"/read/physics-67-disturbance-and-wave?lang={lang}"
            self.assertIn(link, (LESSONS / f"66-magnetism-communication-checkpoint{suffix}.md").read_text())

    def test_astronomy_block_numbers_and_support_maps(self):
        self.assertAlmostEqual(384000 / 12700, 30.24, places=2)
        self.assertEqual(360 / 24, 15)
        self.assertAlmostEqual(4 * .24, .96)
        self.assertEqual(10 * 5.2, 52)
        self.assertEqual(150_000_000 / 300_000, 500)
        for row in ASTRONOMY:
            for lang in ('ru', 'kz', 'en'):
                svg = MAPS / f"map-{row['stem']}-{lang}.svg"
                root = ET.parse(svg).getroot()
                self.assertEqual(root.attrib['viewBox'], '0 0 1600 900')
                body = (LESSONS / f"{row['stem']}{'' if lang == 'ru' else '-'+lang}.md").read_text()
                self.assertIn(row[lang][0], body)
                self.assertIn('/static/course/physics/'+svg.name, body)
        for lang in ('ru','kz','en'):
            suffix = '' if lang == 'ru' else '-'+lang
            self.assertIn('/read/physics-75-cosmic-scale?lang='+lang,
                          (LESSONS/f'74-waves-light-investigation{suffix}.md').read_text())

    def test_experiment_numbers_and_limits_match_across_languages(self):
        # These are worked numbers, never represented as a learner's observation.
        short = [10.9, 11.1, 11.0]
        long = [15.5, 15.7, 15.6]
        self.assertAlmostEqual(sum(short) / len(short) / 10, 1.10)
        self.assertAlmostEqual(sum(long) / len(long) / 10, 1.56)
        self.assertAlmostEqual(1.56 - 1.10, 0.46)
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else f"-{lang}"
            text = (LESSONS / f"08-investigation-mastery{suffix}.md").read_text()
            values = ("10,9", "11,1", "11,0", "15,5", "15,7", "15,6", "1,10", "1,56", "0,46") if lang != "en" else (
                "10.9", "11.1", "11.0", "15.5", "15.7", "15.6", "1.10", "1.56", "0.46")
            for value in values:
                self.assertIn(value, text, (lang, value))
            self.assertIn("| 30 |", text)
            self.assertIn("| 60 |", text)

    def test_guarded_publication_is_atomic_and_complete(self):
        sql, expected = prepare()
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.rstrip().endswith("COMMIT;"))
        self.assertIn("Physics slug belongs to another author", sql)
        self.assertIn("Physics lesson belongs to another course", sql)
        self.assertEqual(len(expected), 252)
        self.assertEqual({e["lang"] for e in expected}, {"ru", "kz", "en"})
        self.assertEqual({e["slug"] for e in expected}, {
            f"physics-{i:02d}-{row['stem'][3:]}" for i, row in enumerate([*DATA, *MOTION, *FORCES, *FLUIDS, *HEAT, *ELECTRICITY, *MAGNETISM, *WAVES, *ASTRONOMY], 1)})
        electricity_sql, electricity_expected = prepare(start=46)
        self.assertEqual(len(electricity_expected), 117)
        self.assertEqual({e["slug"] for e in electricity_expected}, {
            f"physics-{i:02d}-{row['stem'][3:]}" for i, row in enumerate([*HEAT[-1:], *ELECTRICITY, *MAGNETISM, *WAVES, *ASTRONOMY], 46)})
        astronomy_sql, astronomy_expected = prepare(start=74)
        self.assertEqual(len(astronomy_expected), 33)
        self.assertIn('physics-75-cosmic-scale', astronomy_sql)
        self.assertNotIn("physics-45-room-energy-budget", electricity_sql)


if __name__ == "__main__":
    unittest.main()
