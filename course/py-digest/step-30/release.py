"""Build a complete static artifact before replacing the public directory."""
import hashlib
import json
import shutil
import tempfile
from datetime import UTC, datetime
from pathlib import Path

ROOT = Path(__file__).parent
DATA = ROOT / "data"
PUBLIC = ROOT / "public"
FILES = ("report.html", "report.csv", "inflation.png")

def digest(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()

def build(public: Path = PUBLIC) -> None:
    """Build beside the destination and replace it only when complete."""
    public.parent.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory(dir=public.parent) as raw:
        staged = Path(raw) / "artifact"
        staged.mkdir()
        for name in FILES:
            source = DATA / name
            if not source.is_file() or source.stat().st_size == 0:
                raise RuntimeError(f"missing artifact input: {name}")
            target = "index.html" if name == "report.html" else name
            shutil.copy2(source, staged / target)
        manifest = {
            "built_at": datetime.now(UTC).isoformat(),
            "source_ids": ["FP.CPI.TOTL.ZG", "FM.LBL.BMNY.GD.ZS"],
            "files": {path.name: digest(path) for path in sorted(staged.iterdir())},
        }
        (staged / "manifest.json").write_text(
            json.dumps(manifest, ensure_ascii=False, indent=2) + "\n", encoding="utf-8"
        )
        backup = public.with_name(public.name + ".previous")
        if backup.exists():
            shutil.rmtree(backup)
        if public.exists():
            public.replace(backup)
        staged.replace(public)

if __name__ == "__main__":
    build()
    print("checks: passed")
    print("artifact: complete")
    print("publish: ready")
