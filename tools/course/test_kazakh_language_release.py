"""Structural, localization, and publication checks for the Kazakh course."""
import json
from pathlib import Path
import re
import unittest
import xml.etree.ElementTree as ET
import wave

from tools.course import prepare_kazakh_language

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/kazakh-language"
MAPS = ROOT / "web/static/course/kazakh-language"
PUBLISHED_STEMS = tuple(stem for stem, _, _ in prepare_kazakh_language.ROUTE)
LESSON_STEMS = PUBLISHED_STEMS
MAP_STEMS = (
    "01-first-contact", "02-nine-sounds", "03-harmony", "04-sentence",
    "05-introduction", "06-questions", "07-repair", "08-mastery",
    "09-short-vowels", "10-family", "11-possessives", "12-plurals",
    "13-numbers", "14-home", "15-existence", "16-mastery",
    "17-school-map", "18-days-subjects", "19-clock-time", "20-directions",
    "21-from-to", "22-transport", "23-rules", "24-mastery",
    "25-daily-routine", "26-sequence", "27-shop", "28-cafe",
    "29-appointment", "30-health", "31-service-error", "32-mastery",
    "33-case-map", "34-genitive", "35-accusative", "36-dative",
    "37-locative", "38-ablative", "39-instrumental", "40-mastery",
    "41-time-map", "42-habit", "43-progressive", "44-past",
    "45-future", "46-questions", "47-story", "48-mastery",
    "49-request", "50-agreement", "51-possibility", "52-cause",
    "53-condition", "54-problem", "55-choice", "56-mastery",
    "57-chunks", "58-intonation", "59-gist", "60-repair",
    "61-story", "62-context", "63-retell", "64-mastery",
    "65-fact-opinion", "66-source-date", "67-small-chart", "68-baseline",
    "69-claim-evidence", "70-dialogue", "71-public-card", "72-mastery",
)


class KazakhLanguageReleaseTest(unittest.TestCase):
    def test_architecture_has_ten_blocks_and_eighty_lessons(self):
        data = json.loads((ROOT / "course/kazakh-language/curriculum.json").read_text())
        self.assertEqual(data["publication_policy"], "complete-block-only")
        self.assertEqual(len(data["blocks"]), 10)
        self.assertEqual(data["blocks"][-1]["range"], "73-80")
        self.assertEqual(len(data["first_block_lessons"]), 8)
        self.assertEqual(len(data["second_block_lessons"]), 8)
        self.assertEqual(len(data["third_block_lessons"]), 8)
        self.assertEqual(len(data["fourth_block_lessons"]), 8)
        self.assertEqual(len(data["fifth_block_lessons"]), 8)
        self.assertEqual(len(data["sixth_block_lessons"]), 8)
        self.assertEqual(len(data["seventh_block_lessons"]), 8)
        self.assertEqual(len(data["eighth_block_lessons"]), 8)
        self.assertEqual(len(data["ninth_block_lessons"]), 8)
        self.assertEqual(data["blocks"][0]["status"], "ready")
        self.assertEqual(data["blocks"][1]["status"], "ready")
        self.assertEqual(data["blocks"][2]["status"], "ready")
        self.assertEqual(data["blocks"][3]["status"], "ready")
        self.assertEqual(data["blocks"][4]["status"], "ready")
        self.assertEqual(data["blocks"][5]["status"], "ready")
        self.assertEqual(data["blocks"][6]["status"], "ready")
        self.assertEqual(data["blocks"][7]["status"], "ready")
        self.assertEqual(data["blocks"][8]["status"], "ready")
        self.assertEqual(data["blocks"][3]["lexical_target"], 800)
        self.assertEqual(data["blocks"][6]["lexical_target"], 1783)
        self.assertEqual(data["blocks"][-1]["lexical_target"], 2229)

    def test_completed_blocks_are_complete_in_three_languages(self):
        headings = {"ru": "## Задание", "kz": "## Тапсырма", "en": "## Exercise"}
        # Kazakh carries more meaning inside each inflected word, so its honest
        # word count is lower than Russian or English for equivalent content.
        minimum = {"ru": 470, "kz": 440, "en": 570}
        forbidden = {
            "kz": ("## Задание", "Следующий урок", "Проверка переноса"),
            "en": ("## Задание", "Келесі сабақ", "Проверка переноса"),
        }
        for index, stem in enumerate(LESSON_STEMS):
            for lang in ("ru", "kz", "en"):
                suffix = "" if lang == "ru" else f"-{lang}"
                path = LESSONS / f"{stem}{suffix}.md"
                self.assertTrue(path.is_file(), path)
                text = path.read_text(encoding="utf-8")
                audio_prefix = "](/static/course/kazakh-language/audio/kz-"
                self.assertIn(audio_prefix, text, path)
                self.assertGreaterEqual(text.count(audio_prefix), 5, path)
                self.assertNotIn("#speak-kz=", text, path)
                self.assertIn("```kazakh\n", text, path)
                self.assertIn(headings[lang], text, path)
                self.assertGreaterEqual(len(text.split()), minimum[lang], path)
                self.assertNotRegex(text, r"\b(TODO|TBD)\b")
                for phrase in forbidden.get(lang, ()):
                    self.assertNotIn(phrase, text, path)
                if index + 1 < len(LESSON_STEMS):
                    next_slug = LESSON_STEMS[index + 1]
                    self.assertIn(f"/read/kazakh-language-{next_slug}?lang={lang}", text, path)
                else:
                    self.assertIn(f"/course/kazakh-language?lang={lang}", text, path)

    def test_all_completed_block_maps_are_valid_and_fit_the_canvas(self):
        for stem in MAP_STEMS:
            for lang in ("ru", "kz", "en"):
                path = MAPS / f"map-{stem}-{lang}.svg"
                root = ET.parse(path).getroot()
                self.assertEqual(root.attrib.get("viewBox"), "0 0 1600 900", path)
                text = path.read_text(encoding="utf-8")
                self.assertNotIn("TODO", text, path)
                expected_filters = 4 if stem[:2].isdigit() and int(stem[:2]) <= 16 else 1
                self.assertGreaterEqual(text.count('filter="url(#'), expected_filters, path)
                self.assertFalse(
                    re.search(r'<text[^>]*x="(?:1[6-9]\d\d|[2-9]\d{3})"', text), path
                )

    def test_atomic_sql_contains_the_full_route_and_expected_hashes(self):
        sql, expected = prepare_kazakh_language.prepare()
        self.assertTrue(sql.startswith("BEGIN;"))
        self.assertTrue(sql.rstrip().endswith("COMMIT;"))
        self.assertEqual(sql.count("INSERT INTO articles("), 72)
        self.assertEqual(sql.count("INSERT INTO article_translations("), 216)
        self.assertEqual(sql.count("INSERT INTO article_series_items("), 72)
        self.assertEqual(len(expected), 216)
        self.assertEqual(
            [position for _, _, position in prepare_kazakh_language.ROUTE],
            list(range(10, 721, 10)),
        )
        self.assertIn("'kazakh-language'", sql)
        self.assertIn("'kazakh'", sql)

    def test_phonetic_laboratory_distinguishes_sound_reduction_and_deletion(self):
        required = {
            "ru": ("не означает немую букву", "ауыз + ы → аузы", "орын + ы → орны"),
            "kz": ("дыбыссыз әріп емес", "ауыз + ы → аузы", "орын + ы → орны"),
            "en": ("not a silent letter", "ауыз + ы → аузы", "орын + ы → орны"),
        }
        for lang, phrases in required.items():
            suffix = "" if lang == "ru" else f"-{lang}"
            text = (LESSONS / f"09-short-vowels-y-i{suffix}.md").read_text(encoding="utf-8")
            for phrase in phrases:
                self.assertIn(phrase, text)
            self.assertIn("1929", text)
            self.assertIn("cambridge.org", text)

    def test_second_block_cover_is_present_and_large(self):
        cover = ROOT / "web/static/covers/school/kazakh-language/people-home/02-people-home.webp"
        self.assertTrue(cover.is_file())
        self.assertGreater(cover.stat().st_size, 300_000)
        self.assertEqual(cover.read_bytes()[:4], b"RIFF")

    def test_third_block_cover_is_present_and_large(self):
        cover = ROOT / "web/static/covers/school/kazakh-language/school-city/03-school-city.webp"
        self.assertTrue(cover.is_file())
        self.assertGreater(cover.stat().st_size, 500_000)
        self.assertEqual(cover.read_bytes()[:4], b"RIFF")

    def test_fourth_block_cover_is_present_and_4k(self):
        cover = ROOT / "web/static/covers/school/kazakh-language/day-services/04-day-services.webp"
        self.assertTrue(cover.is_file())
        self.assertGreater(cover.stat().st_size, 500_000)
        self.assertEqual(cover.read_bytes()[:4], b"RIFF")
        self.assertIn(b"VP8", cover.read_bytes()[:32])

    def test_fifth_block_cover_is_present_and_4k(self):
        cover = ROOT / "web/static/covers/school/kazakh-language/space-cases/05-space-cases.webp"
        self.assertTrue(cover.is_file())
        self.assertGreater(cover.stat().st_size, 500_000)
        self.assertEqual(cover.read_bytes()[:4], b"RIFF")
        self.assertIn(b"VP8", cover.read_bytes()[:32])

    def test_sixth_block_cover_is_present_and_4k(self):
        cover = ROOT / "web/static/covers/school/kazakh-language/action-time/06-action-time.webp"
        self.assertTrue(cover.is_file())
        self.assertGreater(cover.stat().st_size, 500_000)
        self.assertEqual(cover.read_bytes()[:4], b"RIFF")
        self.assertIn(b"VP8", cover.read_bytes()[:32])

    def test_seventh_block_cover_is_present_and_4k(self):
        cover = ROOT / "web/static/covers/school/kazakh-language/plans-problems/07-plans-problems.webp"
        self.assertTrue(cover.is_file())
        self.assertGreater(cover.stat().st_size, 500_000)
        self.assertEqual(cover.read_bytes()[:4], b"RIFF")
        self.assertIn(b"VP8", cover.read_bytes()[:32])

    def test_every_spoken_phrase_has_a_valid_recording_and_review_state(self):
        audio = MAPS / "audio"
        manifest = json.loads((audio / "manifest.json").read_text(encoding="utf-8"))
        recordings = manifest["recordings"]

        lesson_files = []
        for path in LESSONS.glob("*.md"):
            lesson_files.extend(re.findall(
                r"/static/course/kazakh-language/audio/(kz-\d{3}\.wav)",
                path.read_text(encoding="utf-8"),
            ))

        self.assertEqual(len(recordings), 411)
        self.assertEqual(len(lesson_files), 1278)
        self.assertTrue(all(row["status"] == "approved" for row in recordings))
        self.assertEqual({row["file"] for row in recordings}, set(lesson_files))
        review = (audio / "review.html").read_text(encoding="utf-8")
        frame_counts = set()
        for row in recordings:
            path = audio / row["file"]
            self.assertTrue(path.is_file(), path)
            with wave.open(str(path), "rb") as wav:
                self.assertEqual(wav.getnchannels(), 1, path)
                self.assertEqual(wav.getframerate(), 48_000, path)
                self.assertEqual(wav.getsampwidth(), 3, path)
                self.assertGreater(wav.getnframes(), 36_000, path)
                frame_counts.add(wav.getnframes())
                raw = wav.readframes(wav.getnframes())
                peak = max(
                    abs(int.from_bytes(raw[i : i + 3], "little", signed=True))
                    for i in range(0, len(raw) - 2, 3)
                )
                self.assertGreater(peak, 10_000, path)
            if int(row["id"][3:]) >= 317:
                self.assertIn(f'src="{row["file"]}"', review)
            else:
                self.assertNotIn(f'src="{row["file"]}"', review)
        # MP3 source duration is quantised in codec frames, so several short
        # phrases legitimately share a length. The failed local renderer made
        # every asset identical; a few dozen distinct durations catches that case.
        self.assertGreater(len(frame_counts), 30)

    def test_eighth_block_is_complete_and_approved(self):
        from tools.course.generate_kazakh_language_connected_speech import LESSONS as BLOCK_LESSONS

        curriculum = json.loads((ROOT / "course/kazakh-language/curriculum.json").read_text())
        self.assertEqual(curriculum["blocks"][7]["status"], "ready")
        self.assertEqual(len(curriculum["eighth_block_lessons"]), 8)
        self.assertEqual([int(lesson["stem"][:2]) for lesson in BLOCK_LESSONS], list(range(57, 65)))
        for lesson in BLOCK_LESSONS:
            for lang, minimum in (("ru", 470), ("kz", 440), ("en", 570)):
                suffix = "" if lang == "ru" else "-" + lang
                path = LESSONS / f"{lesson['stem']}{suffix}.md"
                text = path.read_text(encoding="utf-8")
                self.assertGreaterEqual(len(text.split()), minimum, path)
                self.assertEqual(text.count("/static/course/kazakh-language/audio/kz-"), 6, path)
                self.assertIn("```kazakh\n", text, path)
                self.assertIn("## " + {"ru": "Задание", "kz": "Тапсырма", "en": "Exercise"}[lang], text, path)
                diagram = MAPS / f"map-{lesson['map']}-{lang}.svg"
                root = ET.parse(diagram).getroot()
                self.assertEqual(root.attrib["viewBox"], "0 0 1600 900")
        cover = ROOT / "web/static/covers/school/kazakh-language/connected-speech/08-connected-speech.webp"
        self.assertGreater(cover.stat().st_size, 500_000)

    def test_ninth_block_numbers_and_localized_assets(self):
        from tools.course.generate_kazakh_language_information_opinion import LESSONS as BLOCK

        self.assertEqual([d['n'] for d in BLOCK], list(range(65, 73)))
        for d in BLOCK:
            stem=f"{d['n']:02d}-{d['slug']}"
            for lang, minimum in (('ru', 470), ('kz', 440), ('en', 570)):
                suffix='' if lang=='ru' else '-'+lang
                text=(LESSONS/f'{stem}{suffix}.md').read_text(encoding='utf-8')
                self.assertGreaterEqual(len(text.split()), minimum)
                self.assertEqual(text.count('/static/course/kazakh-language/audio/kz-'), 6)
                self.assertIn('```kazakh\n', text)
                self.assertIn('## '+{'ru':'Задание','kz':'Тапсырма','en':'Exercise'}[lang], text)
                map_path=MAPS/f"map-{d['map']}-{lang}.svg"
                self.assertEqual(ET.parse(map_path).getroot().attrib['viewBox'],'0 0 1600 900')
                svg=map_path.read_text(encoding='utf-8')
                if d['n'] not in (67,68):
                    from tools.course.generate_kazakh_language_information_opinion import MAP_EXAMPLES
                    for label in d['labels'][{'ru':1,'kz':0,'en':2}[lang]]:
                        self.assertIn(f'>{label}</text>',svg,map_path)
                    for example in MAP_EXAMPLES[d['n']][lang]:
                        self.assertIn(f'>{example}</text>',svg,map_path)
        cover=ROOT/'web/static/covers/school/kazakh-language/information-opinion/09-information-opinion.webp'
        self.assertGreater(cover.stat().st_size,500_000)


if __name__ == "__main__":
    unittest.main()
