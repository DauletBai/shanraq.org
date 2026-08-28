# The speech synthesiser, in its own image.
#
# It is separate from the site for one reason: the site is a small static Go
# binary, and this is Python with an ONNX runtime and three voice models. Fusing
# them would multiply the size of every deploy by a component that changes once
# a year.
#
# The models are not baked in. They are 240 MB together, they are licensed
# separately (CC-BY-4.0 for the Kazakh voice, from IS2AI/Kazakh_TTS), and an
# image is the wrong place for either. They arrive as a mounted volume, filled
# once by deploy/tts/fetch-voices.sh.
FROM python:3.12-slim

# ffmpeg encodes the finished audio to Opus. espeak-ng is what turns text into
# the phonemes Piper speaks; the wheel bundles its data, but the system package
# is what makes the loader find it reliably on this base image.
RUN apt-get update \
 && apt-get install -y --no-install-recommends ffmpeg espeak-ng-data \
 && rm -rf /var/lib/apt/lists/*

RUN pip install --no-cache-dir piper-tts

WORKDIR /app
COPY server.py /app/server.py

ENV TTS_PORT=8081
EXPOSE 8081

# The check loads nothing: it reports which voices are on disk, which is what
# distinguishes "starting" from "started without its models".
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD python3 -c "import urllib.request,sys; sys.exit(0 if urllib.request.urlopen('http://127.0.0.1:8081/health',timeout=3).status==200 else 1)"

CMD ["python3", "/app/server.py"]
