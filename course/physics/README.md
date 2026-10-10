# Physics / Физика / Физика

The 100-lesson route is described in `curriculum.json` and `docs/physics-course.md`.
Only complete blocks are published. The first block contains eight lessons in
Russian, Kazakh, and English. The Russian filename is the base stem; Kazakh
and English append `-kz` and `-en`. Maps use `-ru`, `-kz`, and `-en`.

The learner maintains a laboratory notebook of fictional or self-measured,
non-personal data. The first release compares the period of a lightweight
pendulum at two string lengths. It does not claim to have derived the formula
for the period or established a universal law from two measurements.

From the repository root, regenerate and validate the first block:

```sh
python3 tools/course/build_physics_block.py
python3 -m unittest tools.course.test_physics_release
python3 tools/course/prepare_physics.py \
  --sql /tmp/physics-course.sql \
  --expected /tmp/physics-course-expected.json
```

The prepare command writes a guarded SQL transaction locally. It does not
modify the live site. The public course uses the shared prose checker mode
only after a separately authored exercise and suitable physics review prompt
are added; the first block instead uses its explicit notebook self-checks.
