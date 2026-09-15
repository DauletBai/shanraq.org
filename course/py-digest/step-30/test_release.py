import json
import release

def test_build_writes_complete_manifest(tmp_path):
    public = tmp_path / "public"
    release.build(public)
    manifest = json.loads((public / "manifest.json").read_text(encoding="utf-8"))
    assert set(manifest["files"]) == {"index.html", "report.csv", "inflation.png"}
    assert all(len(value) == 64 for value in manifest["files"].values())
