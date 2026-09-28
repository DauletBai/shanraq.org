# Mathematics course

The curriculum is a dependency graph, not a list of school years.
`curriculum.json` is the source of truth for the full 50-node route from
quantity to mathematical modelling. Every prerequisite must name another node,
and the graph must remain acyclic; `tools/course/test_mathematics_release.py`
checks both properties.

The first published route is the nine-lesson fractions, ratios, percent, and
proportion pilot in `course/lessons/mathematics`. It tests the lesson cycle and
the mastery gate before the remaining nodes are written:

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

The Russian edition is published first. Kazakh and English lesson files should
be added as reviewed translations, with `-kz` and `-en` suffixes, before those
editions are described as available. Course metadata already has all three
languages so the hub is never blank while localization is in progress.

Generate the original support maps and run the release checks:

```sh
python3 tools/course/generate_math_maps.py
cd tools/course && python3 -m unittest test_mathematics_release.py
```

Prepare a reviewable transaction without touching production:

```sh
python3 tools/course/prepare_mathematics.py \
  --sql /tmp/mathematics-course.sql \
  --expected /tmp/mathematics-expected.json
```

