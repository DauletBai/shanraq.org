# Mathematics course

The curriculum is a dependency graph, not a list of school years.
`curriculum.json` is the source of truth for the full 50-node route from
quantity to mathematical modelling. Every prerequisite must name another node,
and the graph must remain acyclic; `tools/course/test_mathematics_release.py`
checks both properties.

The prepared route now contains seven connected blocks: numbers and operations;
fractions, ratios, percent, and proportion; variables and prealgebra; then
algebra and functions; geometry and space; calculus; then matrices, vector
spaces, bases, and eigenvectors. Together they provide fifty-one lessons plus the course
orientation and test the lesson cycle and mastery gate before the remaining
nodes are written:

1. whole-map orientation;
2. familiar physical image;
3. exact definition;
4. compact support signal;
5. fully worked example;
6. faded example;
7. independent problem;
8. explanation in the learner's words;
9. error diagnosis;
10. transfer to an unseen context;
11. retrieval after a delay;
12. mastery gate.

Every released page and support map is available in Russian, Kazakh, and
English. Russian source files have no suffix; reviewed Kazakh and English
versions use `-kz` and `-en`. The language versions keep the same formulas,
examples, checks, and lesson order.

Generate the original support maps and run the release checks:

```sh
python3 tools/course/generate_math_maps.py
python3 tools/course/generate_math_geometry.py
python3 tools/course/generate_math_calculus.py
python3 tools/course/generate_math_linear.py
python3 tools/course/generate_math_localizations.py
python3 tools/course/generate_math_map_localizations.py
cd tools/course && python3 -m unittest test_mathematics_release.py
```

Prepare a reviewable transaction without touching production:

```sh
python3 tools/course/prepare_mathematics.py \
  --sql /tmp/mathematics-course.sql \
  --expected /tmp/mathematics-expected.json
```
