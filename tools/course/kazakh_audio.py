"""Approved audio assets used by the Kazakh-language lesson generators."""
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
AUDIO_DIR = ROOT / "web/static/course/kazakh-language/audio"
PUBLIC_PREFIX = "/static/course/kazakh-language/audio/"


def approved_audio_by_phrase():
    manifest = json.loads((AUDIO_DIR / "manifest.json").read_text(encoding="utf-8"))
    result = {}
    for recording in manifest["recordings"]:
        if recording["status"] != "approved":
            raise ValueError(f"audio is not approved: {recording['id']}")
        phrase = recording["phrase"]
        if phrase in result:
            raise ValueError(f"duplicate audio phrase: {phrase}")
        path = AUDIO_DIR / recording["file"]
        if not path.is_file():
            raise FileNotFoundError(path)
        result[phrase] = PUBLIC_PREFIX + recording["file"]
    return result


AUDIO_BY_PHRASE = approved_audio_by_phrase()


def audio_url(phrase):
    try:
        return AUDIO_BY_PHRASE[phrase]
    except KeyError as exc:
        raise ValueError(f"approved Kazakh audio is missing for: {phrase}") from exc
