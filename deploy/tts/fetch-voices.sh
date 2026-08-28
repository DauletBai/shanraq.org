#!/bin/sh
# Fill the voice volume. Run once on the host, before the synthesiser starts.
#
# The models are not in the image and not in git: 240 MB together, they change
# about never, and the Kazakh one carries its own licence. Keeping them out of
# both means a deploy stays small and the licence stays visible.
#
#   deploy/tts/fetch-voices.sh /var/lib/shanraq/voices
#
set -eu

DIR="${1:-/var/lib/shanraq/voices}"
BASE="https://huggingface.co/rhasspy/piper-voices/resolve/main"

mkdir -p "$DIR"

# Kazakh is the reason any of this exists: no browser ships a Kazakh voice, so
# before this the button offered a Russian one and read Kazakh in it.
# "high" rather than "medium" because this is the language we cannot get wrong;
# it costs about five times the compute of the other two and is worth it.
fetch() {
  path="$1"; name="$2"
  if [ -s "$DIR/$name.onnx" ]; then
    echo "have $name"
    return
  fi
  echo "fetching $name"
  curl -fsSL "$BASE/$path/$name.onnx"      -o "$DIR/$name.onnx.part"
  curl -fsSL "$BASE/$path/$name.onnx.json" -o "$DIR/$name.onnx.json"
  # Renamed only once the whole file is down, so an interrupted fetch cannot
  # leave a truncated model that loads and then produces noise.
  mv "$DIR/$name.onnx.part" "$DIR/$name.onnx"
}

fetch "kk/kk_KZ/issai/high"    "kk_KZ-issai-high"
fetch "ru/ru_RU/dmitri/medium" "ru_RU-dmitri-medium"
fetch "en/en_US/lessac/medium" "en_US-lessac-medium"

# The Kazakh voice is trained on IS2AI/Kazakh_TTS under CC-BY-4.0: commercial
# use is allowed, attribution is required. Keep this next to the files so the
# obligation travels with them.
cat > "$DIR/ATTRIBUTION.txt" <<'TXT'
kk_KZ-issai-high
  Trained from scratch on IS2AI/Kazakh_TTS — https://github.com/IS2AI/Kazakh_TTS
  Dataset licence: CC-BY-4.0. Attribution required: ISSAI, Nazarbayev University.

Piper itself is MIT — https://github.com/rhasspy/piper
TXT

ls -la "$DIR"
