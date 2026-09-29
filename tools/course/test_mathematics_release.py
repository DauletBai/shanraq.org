"""Structural and pedagogical checks for the mathematics course."""
import json
import re
import unittest
from pathlib import Path
import xml.etree.ElementTree as ET

import prepare_mathematics


ROOT = Path(__file__).resolve().parents[2]
CURRICULUM = ROOT / "course/mathematics/curriculum.json"


class MathematicsReleaseTests(unittest.TestCase):
    def test_curriculum_is_a_connected_dag(self):
        data = json.loads(CURRICULUM.read_text(encoding="utf-8"))
        modules = data["modules"]
        ids = [m["id"] for m in modules]
        self.assertEqual(len(ids), len(set(ids)))
        known = set(ids)
        for module in modules:
            self.assertTrue({"title_ru", "title_kz", "title_en"} <= module.keys())
            self.assertTrue(set(module["prerequisites"]) <= known)

        visiting, visited = set(), set()

        def walk(node):
            if node in visiting:
                self.fail(f"curriculum cycle at {node}")
            if node in visited:
                return
            visiting.add(node)
            current = next(m for m in modules if m["id"] == node)
            for dependency in current["prerequisites"]:
                walk(dependency)
            visiting.remove(node)
            visited.add(node)

        for node in ids:
            walk(node)
        self.assertEqual(visited, known)
        self.assertIn("capstone-modeling", known)

    def test_published_route_is_atomic_and_non_destructive(self):
        sql, expected = prepare_mathematics.prepare()
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.endswith("COMMIT;\n"))
        self.assertEqual(len(expected), 156)
        self.assertEqual(sql.count("INSERT INTO article_series_items("), 52)
        self.assertIn("'society','education'", sql)
        self.assertIn("'published','math'", sql)
        self.assertNotIn("DELETE FROM", sql)
        self.assertNotIn("TRUNCATE", sql)
        positions = [position for _, _, position in prepare_mathematics.ROUTE]
        self.assertEqual(positions, sorted(positions))
        self.assertEqual(len(positions), len(set(positions)))
        self.assertLess(
            positions[prepare_mathematics.STEMS.index("18-foundations-mastery")],
            positions[prepare_mathematics.STEMS.index("02-fraction-meaning")],
        )
        self.assertEqual({item["lang"] for item in expected}, {"ru", "kz", "en"})

    def test_every_route_has_three_localized_pages_and_maps(self):
        headings = {
            "kz": ("## Картадағы орнымыз", "## Сабақтың тірек сигналы", "## Тапсырма"),
            "en": ("## Where we are on the map", "## The lesson's support signal", "## Exercise"),
        }
        for stem in prepare_mathematics.STEMS:
            for lang in ("kz", "en"):
                path = prepare_mathematics.LESSONS / f"{stem}-{lang}.md"
                self.assertTrue(path.is_file(), path)
                text = path.read_text(encoding="utf-8")
                for heading in headings[lang]:
                    self.assertIn(heading, text, path)
                found = re.findall(r"/static/course/mathematics/([^\s)]+\.svg)", text)
                self.assertEqual(len(found), 1, path)
                self.assertTrue(found[0].endswith(f"-{lang}.svg"), found[0])
                target = ROOT / "web/static/course/mathematics" / found[0]
                self.assertTrue(target.is_file(), target)
                ET.parse(target)
                self.assertIn(f"/course/mathematics?lang={lang}", text)

    def test_every_lesson_is_a_full_learning_cycle(self):
        sections = (
            "## Где мы на карте",
            "## Опорный сигнал",
            "## Воспроизведите без подсказки",
            "## Найдите и исправьте ошибку",
            "## Задание",
        )
        for stem in prepare_mathematics.TEACHING_STEMS:
            path = prepare_mathematics.LESSONS / f"{stem}.md"
            text = path.read_text(encoding="utf-8")
            for heading in sections:
                self.assertIn(heading, text, path)
            self.assertIn("/static/course/mathematics/", text, path)
            self.assertGreaterEqual(len(re.findall(r"\b[\w/-]+\b", text)), 500, path)
        for stem in prepare_mathematics.MASTERY_STEMS:
            mastery = (prepare_mathematics.LESSONS / f"{stem}.md").read_text(encoding="utf-8")
            self.assertRegex(mastery, r"8 (?:из 10|из десяти)")
            self.assertRegex(mastery, r"[Чч]ерез семь дней")

    def test_each_lesson_uses_its_own_map(self):
        maps = []
        for stem in prepare_mathematics.STEMS[1:]:
            text = (prepare_mathematics.LESSONS / f"{stem}.md").read_text(encoding="utf-8")
            found = re.findall(r"/static/course/mathematics/([^\s)]+\.svg)", text)
            self.assertEqual(len(found), 1, stem)
            maps.extend(found)
        self.assertEqual(len(maps), len(set(maps)))
        for name in maps:
            self.assertTrue((ROOT / "web/static/course/mathematics" / name).is_file(), name)
            ET.parse(ROOT / "web/static/course/mathematics" / name)
        ET.parse(ROOT / "web/static/course/mathematics/map-full-ru.svg")

    def test_course_map_direction_and_percent_cells_are_explicit(self):
        maps = ROOT / "web/static/course/mathematics"
        course_map = (maps / "map-full-ru.svg").read_text(encoding="utf-8")
        self.assertEqual(course_map.count('class="course-step"'), 8)
        self.assertEqual(course_map.count('class="course-arrow"'), 7)
        self.assertNotIn("M1135 465", course_map)

        percent_map = (maps / "map-07-percent.svg").read_text(encoding="utf-8")
        self.assertEqual(percent_map.count('class="percent-cell'), 100)
        self.assertEqual(percent_map.count("percent-cell--filled"), 37)
        self.assertEqual(percent_map.count("percent-cell--final-seven"), 7)
        self.assertIn("30 + 7 = 37", percent_map)

    def test_lesson_maps_encode_the_examples_they_explain(self):
        maps = ROOT / "web/static/course/mathematics"
        lessons = prepare_mathematics.LESSONS

        diagnostic = (maps / "map-01-diagnostic.svg").read_text(encoding="utf-8")
        self.assertIn('data-diagnostic-route="12|3/8|2:3|25%|x=5"', diagnostic)

        expected = (
            ("02-fraction-meaning.md", "map-02-fraction.svg", "3/8", 'data-fraction="3/8"'),
            ("03-equivalent-fractions.md", "map-03-equivalence.svg", "1/2", 'data-equivalence="1/2=2/4"'),
            ("04-compare-fractions.md", "map-04-compare.svg", "2/3", 'data-comparison="1/2&lt;2/3"'),
            ("05-fraction-operations.md", "map-05-operations.svg", "5/6", 'data-equation="1/2+1/3=5/6"'),
            ("06-ratio.md", "map-06-ratio.svg", "4:6", 'data-ratios="2:3=4:6"'),
            ("07-percent.md", "map-07-percent.svg", "37/100", 'data-fraction="37/100"'),
            ("08-proportion.md", "map-08-proportion.svg", "6/9", 'data-equivalence="2/3=6/9"'),
        )
        for lesson_name, map_name, example, marker in expected:
            lesson = (lessons / lesson_name).read_text(encoding="utf-8")
            illustration = (maps / map_name).read_text(encoding="utf-8")
            self.assertIn(example, lesson, lesson_name)
            self.assertIn(marker, illustration, map_name)

        fraction = ET.parse(maps / "map-02-fraction.svg").getroot()
        namespace = {"svg": "http://www.w3.org/2000/svg"}
        spokes = fraction.find(".//svg:g[@class='fraction-spokes']", namespace)
        ticks = fraction.find(".//svg:g[@class='eighth-ticks']", namespace)
        self.assertEqual(len(spokes.findall("svg:path", namespace)), 8)
        self.assertEqual(len(ticks.findall("svg:path", namespace)), 9)

        comparison = ET.parse(maps / "map-04-compare.svg").getroot()
        sixth_ticks = comparison.find(".//svg:g[@class='sixth-ticks']", namespace)
        self.assertEqual(len(sixth_ticks.findall("svg:path", namespace)), 7)

        mastery = (maps / "map-09-mastery.svg").read_text(encoding="utf-8")
        self.assertEqual(mastery.count('class="mastery-skill"'), 5)
        self.assertIn("8/10", mastery)
        self.assertIn("7/10", mastery)

        ratio = (maps / "map-06-ratio.svg").read_text(encoding="utf-8")
        self.assertIn("×2", ratio)
        proportion = (maps / "map-08-proportion.svg").read_text(encoding="utf-8")
        self.assertIn("2 × 3 = 6; 3 × 3 = 9", proportion)

        foundations = (
            ("10-quantity-counting.md", "map-10-quantity-counting.svg", "4 × 6 + 3 = 27", 'data-counting="4x6+3=27"'),
            ("11-place-value.md", "map-11-place-value.svg", "4 072", 'data-place-value="4072=4000+70+2"'),
            ("12-addition-subtraction.md", "map-12-addition-subtraction.svg", "268 + 157 = 425", 'data-family="268+157=425"'),
            ("13-multiplication-division.md", "map-13-multiplication-division.svg", "4 × 6", 'data-array="4x6=24"'),
            ("14-order-estimation.md", "map-14-order-estimation.svg", "240 − 6 × (18 + 7)", 'data-expression="240-6*(18+7)=90"'),
            ("15-negative-numbers.md", "map-15-negative-numbers.svg", "−3 + 7", 'data-integers="-3+7=4"'),
            ("16-divisibility-primes.md", "map-16-divisibility-primes.svg", "84 = 2 × 2 × 3 × 7", 'data-factorization="84=2^2*3*7"'),
            ("17-decimal-fractions.md", "map-17-decimal-fractions.svg", "3/8 = 375/1000", 'data-decimal="3/8=375/1000=0.375"'),
        )
        for lesson_name, map_name, example, marker in foundations:
            lesson = (lessons / lesson_name).read_text(encoding="utf-8")
            illustration = (maps / map_name).read_text(encoding="utf-8")
            self.assertIn(example, lesson, lesson_name)
            self.assertIn(marker, illustration, map_name)

        counting = (maps / "map-10-quantity-counting.svg").read_text(encoding="utf-8")
        self.assertEqual(counting.count("<circle"), 27)
        array = (maps / "map-13-multiplication-division.svg").read_text(encoding="utf-8")
        self.assertEqual(array.count("<circle"), 24)
        foundation_mastery = (maps / "map-18-foundations-mastery.svg").read_text(encoding="utf-8")
        self.assertEqual(foundation_mastery.count('class="foundation-skill"'), 8)
        self.assertIn("8/10", foundation_mastery)
        self.assertIn("7/10", foundation_mastery)

        prealgebra = (
            ("19-variables.md", "map-19-variables.svg", "C = 700 + 120d", 'data-variable="C=700+120d;d=5;C=1300"'),
            ("20-expressions.md", "map-20-expressions.svg", "3(2x + 5) − 4x", 'data-expression="3(2x+5)-4x=2x+15;x=4;value=23"'),
            ("21-equations-basic.md", "map-21-equations-basic.svg", "3x + 5 = 26", 'data-equation="3x+5=26;x=7"'),
            ("22-coordinates.md", "map-22-coordinates.svg", "A(−3; 2)", 'data-coordinates="A(-3,2);B(4,2);distance=7"'),
            ("23-inequalities.md", "map-23-inequalities.svg", "3x + 2 ≤ 14", 'data-inequality="3x+2&lt;=14;x&lt;=4"'),
        )
        for lesson_name, map_name, example, marker in prealgebra:
            lesson = (lessons / lesson_name).read_text(encoding="utf-8")
            illustration = (maps / map_name).read_text(encoding="utf-8")
            self.assertIn(example, lesson, lesson_name)
            self.assertIn(marker, illustration, map_name)

        prealgebra_mastery = (maps / "map-24-prealgebra-mastery.svg").read_text(encoding="utf-8")
        self.assertEqual(prealgebra_mastery.count('class="prealgebra-skill"'), 5)
        self.assertIn("8/10", prealgebra_mastery)
        self.assertIn("7/10", prealgebra_mastery)

        algebra = (
            ("25-linear-functions.md", "map-25-linear-functions.svg", "y = 2x + 1", 'data-linear="y=2x+1;(-1,-1);(0,1);(2,5)"'),
            ("26-systems.md", "map-26-systems.svg", "x + y = 10", 'data-system="x+y=10;2x+y=16;(6,4)"'),
            ("27-powers-roots.md", "map-27-powers-roots.svg", "2³ × 2⁴", 'data-powers="2^3*2^4=2^7=128;sqrt(144)=12"'),
            ("28-polynomials.md", "map-28-polynomials.svg", "(2x + 3)(x − 4)", 'data-polynomial="(2x+3)(x-4)=2x^2-5x-12"'),
            ("29-quadratics.md", "map-29-quadratics.svg", "x² − 4x + 3", 'data-quadratic="x^2-4x+3=(x-1)(x-3)=(x-2)^2-1"'),
            ("30-exponential-log.md", "map-30-exponential-log.svg", "N = 500 × 2ᵗ", 'data-exponential="N=500*2^t;log_2(8)=3"'),
            ("31-sequences.md", "map-31-sequences.svg", "5, 8, 11, 14", 'data-sequences="5,8,11,14;d=3|2,6,18,54;q=3"'),
        )
        for lesson_name, map_name, example, marker in algebra:
            lesson = (lessons / lesson_name).read_text(encoding="utf-8")
            illustration = (maps / map_name).read_text(encoding="utf-8")
            self.assertIn(example, lesson, lesson_name)
            self.assertIn(marker, illustration, map_name)

        algebra_mastery = (maps / "map-32-algebra-mastery.svg").read_text(encoding="utf-8")
        self.assertEqual(algebra_mastery.count('class="algebra-skill"'), 7)
        self.assertIn("8/10", algebra_mastery)
        self.assertIn("7/10", algebra_mastery)

        geometry = (
            ("33-geometry-language.md", "map-33-geometry-language.svg", "125°", 'data-geometry="adjacent:125+55=180"'),
            ("34-triangles-congruence.md", "map-34-triangles-congruence.svg", "5, 7, 8", 'data-triangles="SSS:5,7,8"'),
            ("35-similarity-scale.md", "map-35-similarity-scale.svg", "3,4,5", 'data-similarity="3-4-5~6-8-10;k=2"'),
            ("36-area-volume.md", "map-36-area-volume.svg", "40 см²", 'data-measure="S=8*5=40;V=40*3=120"'),
            ("37-pythagoras.md", "map-37-pythagoras.svg", "3²+4²", 'data-pythagoras="3^2+4^2=5^2"'),
            ("38-circles.md", "map-38-circles.svg", "8π", 'data-circle="r=4;d=8;C=8pi;S=16pi"'),
            ("39-trigonometry.md", "map-39-trigonometry.svg", "sin α=3/5", 'data-trigonometry="3-4-5;sin=3/5;cos=4/5;tan=3/4"'),
            ("40-vectors-geometry.md", "map-40-vectors-geometry.svg", "a=(3,4)", 'data-vectors="a=(3,4);b=(-1,2);a+b=(2,6);|a|=5"'),
        )
        for lesson_name, map_name, example, marker in geometry:
            lesson = (lessons / lesson_name).read_text(encoding="utf-8")
            illustration = (maps / map_name).read_text(encoding="utf-8")
            self.assertIn(example, lesson, lesson_name)
            self.assertIn(marker, illustration, map_name)

        geometry_mastery = (maps / "map-41-geometry-mastery.svg").read_text(encoding="utf-8")
        self.assertEqual(geometry_mastery.count('class="geometry-skill"'), 8)
        self.assertIn("8/10", geometry_mastery)
        self.assertIn("7/10", geometry_mastery)

        calculus = (
            ("42-limits.md", "map-42-limits.svg", "(x²−4)/(x−2)", 'data-limit="x->2;(x^2-4)/(x-2)->4"'),
            ("43-derivative.md", "map-43-derivative.svg", "f′(x)=2x", 'data-derivative="f=x^2;x=2;fprime=4"'),
            ("44-derivative-applications.md", "map-44-derivative-applications.svg", "A=x(10−x)", 'data-optimization="A=x(10-x);max=A(5)=25"'),
            ("45-integral.md", "map-45-integral.svg", "∫[0,3]2x dx=3²−0²=9", 'data-integral="integral_0^3 2x dx=9"'),
            ("46-differential-equations.md", "map-46-differential-equations.svg", "y′=0,05y", 'data-diffeq="yprime=.05y;y(0)=100;y=100e^(.05t)"'),
        )
        for lesson_name, map_name, example, marker in calculus:
            lesson = (lessons / lesson_name).read_text(encoding="utf-8")
            illustration = (maps / map_name).read_text(encoding="utf-8")
            self.assertIn(example, lesson, lesson_name)
            self.assertIn(marker, illustration, map_name)

        calculus_mastery = (maps / "map-47-calculus-mastery.svg").read_text(encoding="utf-8")
        self.assertEqual(calculus_mastery.count('class="calculus-skill"'), 5)
        self.assertIn("8/10", calculus_mastery)
        self.assertIn("7/10", calculus_mastery)

        linear = (
            ("48-matrices.md", "map-48-matrices.svg", "Ax=(10,10)", 'data-matrix="[[2,1],[1,3]]*(4,2)=(10,10)"'),
            ("49-linear-spaces.md", "map-49-linear-spaces.svg", "v=4b₁+2b₂", 'data-basis="v=(6,2)=4(1,1)+2(1,-1)"'),
            ("50-eigen.md", "map-50-eigen.svg", "Av=λv", 'data-eigen="A=[[2,1],[1,2]];lambda1=3;v1=(1,1);lambda2=1;v2=(1,-1)"'),
        )
        for lesson_name, map_name, example, marker in linear:
            lesson = (lessons / lesson_name).read_text(encoding="utf-8")
            illustration = (maps / map_name).read_text(encoding="utf-8")
            self.assertIn(example, lesson, lesson_name)
            self.assertIn(marker, illustration, map_name)

        linear_mastery = (maps / "map-51-linear-mastery.svg").read_text(encoding="utf-8")
        self.assertEqual(linear_mastery.count('class="linear-skill"'), 3)
        self.assertIn("8/10", linear_mastery)
        self.assertIn("7/10", linear_mastery)

        sql, _ = prepare_mathematics.prepare()
        for cover in prepare_mathematics.LESSON_COVERS.values():
            self.assertIn(cover, sql)
            self.assertTrue((ROOT / "web" / cover.lstrip("/")).is_file(), cover)

    def test_claims_do_not_turn_targets_into_results(self):
        preface = (prepare_mathematics.LESSONS / "preface.md").read_text(encoding="utf-8")
        self.assertIn("не станем выдавать", preface)
        self.assertIn("Числа приблизительны", preface)
        self.assertNotIn("КПД образования — 10%", preface)
        self.assertNotIn("95% успеваемости", preface)

    def test_sql_delimiter_is_rejected(self):
        with self.assertRaises(ValueError):
            prepare_mathematics.literal("bad $mathematics_course$ value")


if __name__ == "__main__":
    unittest.main()
