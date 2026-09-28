# Cloud & DevOps course

## Scope

The published course is one preface and 24 lessons in Kazakh, Russian, and
English. It follows one runnable service from a local process through Docker,
Compose, CI/CD, a secured server, HTTPS, OpenTofu, Ansible, Kubernetes,
observability, backup, and a recovery drill.

The preface states the course boundary explicitly. It teaches software and
platform operations. It does not claim to train electrical, cooling, or
facilities engineers for a physical data centre, and it does not turn
infrastructure investment into a guaranteed employment forecast.

## Source of truth

- Lessons: `course/lessons/cloud-devops/`
- Runnable project: `course/cloud-devops-lab/`
- Support maps: `web/static/course/cloud-devops/`
- Cover: `web/static/covers/it/devops/cloud-devops-course.webp`
- URL mapping: `tools/course/lesson-slugs.json`
- SQL builder: `tools/course/prepare_cloud_devops.py`

Regenerate the localized lesson pages and maps with:

```sh
python3 tools/course/generate_cloud_devops.py
```

Prepare a reviewable publication without contacting production:

```sh
PYTHONPATH=tools/course python3 -m unittest tools/course/test_cloud_devops_release.py
python3 tools/course/prepare_cloud_devops.py \
  --sql /tmp/cloud-devops-release.sql \
  --expected /tmp/cloud-devops-expected.json
```

The SQL uses one transaction, an advisory lock, author and ownership guards,
and upserts. It never deletes course content.

## Local verification

```sh
cd course/cloud-devops-lab
go test ./...
sh -n scripts/*.sh
docker compose config
docker compose up --build -d
./scripts/smoke.sh
docker compose down
```

The 4K cover was generated with the built-in image generation tool. Final
prompt: a documentary-quality, physically plausible modern data centre with
one Central Asian infrastructure engineer, real racks, cabling and cooling,
natural lighting, no logos, text, holograms, fantasy architecture, or abstract
clouds.
