# Kazakh AI course, first release

On 2026-09-24 the preface and lessons 1–5 were published in kz/ru/en as one transaction. Source commit: `411cb425c54e83f122a0a69251773da3a0cda25f`. The series and all 18 translations became visible together.

## Preflight

- Manual novice walkthrough of the five tasks in each language; expected outputs checked against the given records and the Python program.
- `langcheck.py` passed all 544 course Markdown files; offline links passed (456 external targets recorded, 1306 internal links checked).
- `go test ./pkg/site ./pkg/modules/articles` and the release manifest test passed. The preview had 72 browser checks at 390, 768, 844 and 1366 pixels with no page overflow.
- The production backup `/var/backups/shanraq/kazakh-ai-before-20260924.dump` is 4,648,394 bytes; its archive listing was checked. Its SHA-256 is `c6fa2a5f332d36fd7f8cd1426c1457445a9ba425838f091c`.
- The publication SQL SHA-256 was `6352d410a48e5d0082eb4e660d46ee8c84d41249826bf6c9986165a3b1fe1682` on both the workstation and the server.

## Activation and checks

The production source fast-forwarded to the commit above. The Docker image built successfully, and the app was recreated after the owner-authorized five-minute CI observation window. The `Course` and `docker-smoke` workflows had succeeded; `CI` and `Rust course checks` were still running without a reported failure at activation. `/readyz` returned 200 and the new footer link was present before the content transaction.

After publication, all 18 titles, bodies and summaries matched the release SHA-256 manifest. The 18 article URLs, three course pages, footer links in three languages, health endpoints and checks that new lessons stay out of the general feed passed. [Live browser checks](kazakh-ai-live-browser-checks.json) covered all 18 article pages at four viewports without horizontal overflow.

The first live browser pass exposed a display issue: the small course mark was stretched into a large article cover above the title. The six article covers were cleared in a separate transaction (`UPDATE 6`), while the series mark and footer mark remained. A second live 72-page browser pass and screenshot inspection passed. The publication script now creates future article rows without that cover.

The long CI workflow later failed only at `gofmt` for `pkg/site/i18n.go`; its unit and snapshot jobs passed. The file was formatted locally and `gofmt -l .` returned no files. This correction was committed separately as `89305f3b7626e67458138444719ec365f2550772`; its `CI`, `Course`, and `docker-smoke` workflows completed successfully.

## Second release: lessons 6–10

On 2026-09-24, lessons 6–10 were published in ru/kz/en as one database transaction from source commit `8571588e021de7a21cd8df4b5578ee52a6043b0f`. The release left the preface and lessons 1–5 in the same series. Production now has 11 series items and 33 article translations. The five runnable snapshots live under `course/kazakh-ai/step-06` through `step-10`; the course test compares each file with code in all three lesson versions.

Preflight: all 559 course Markdown files passed language checks; offline link checks covered 1306 internal and 465 external links. Lesson code tests, the 33-entry SQL release manifest, Go site/article tests, manual novice walkthrough, and 132 local browser checks at 390, 768, 844 and 1366 pixels passed. The SQL SHA-256 was `b69ef317aefd0f4bf139b7f45973e1274eabd1eadba1a4c45b2ce8a87322a690` on both workstation and server. The verified database backup is `/var/backups/shanraq/kazakh-ai-before-6-10-20260924.dump` (4,674,362 bytes; SHA-256 `dcdfb8fcf8fbec252a96d043f09d15efcaaa029b3bcc0611aeb842e600bf6db8`).

After the owner's five-minute observation window, `Course` and `docker-smoke` had succeeded; the main `checks` job had succeeded and the long `CI` unit job was still running without errors. The server source fast-forwarded to the commit above. The content transaction completed, and all 33 production titles, bodies and summaries matched the SHA-256 release manifest. All 33 article URLs, three course pages, footer links, feeds, `/readyz` and `/healthz` passed HTTP checks. [The live browser report](kazakh-ai-live-browser-checks-10.json) contains 144 page/viewport checks with no horizontal overflow. Kazakh mobile, English desktop and course mobile screenshots were reviewed. The full `CI` workflow for the publication commit later completed successfully, including the Go race test.

## Third release: lessons 11–15

On 2026-09-24 lessons 11–15 were published in ru/kz/en as one transaction from source commit `802ea706ff44718b234f5c92d8582dc94f6f042a`. Production has 16 course items and 48 article translations. The five runnable snapshots in `course/kazakh-ai/step-11` through `step-15` match the code printed in all three lesson versions.

Preflight: manual novice walkthrough checked the new tasks, including comma handling, an ambiguous word, an unknown club, and multiple requested facts. All 574 course Markdown files passed language checks; offline links covered 1306 internal and 471 external links, and all seven unique external links of the new block were reachable. Lesson output tests, the 48-entry release manifest, Go site/article tests, and 192 local browser checks at 390, 768, 844, and 1366 pixels passed. The SQL SHA-256 was `21f18fe776270d5f3a49fd959c32f1281fa84d19c3a592b4b32bb4c521c5392f` on workstation and server. The verified backup `/var/backups/shanraq/kazakh-ai-before-11-15-20260924.dump` is 4,695,124 bytes, SHA-256 `57ad98e1e3d8a7ca02447e32c4984f3b0f2bdbe9bc4cbb47320267f9237a4297`.

After the owner-authorized five-minute observation window, `Course` and `docker-smoke` had succeeded and other workflows had no reported failure. The content transaction completed. All 48 production titles, bodies, and summaries matched the source SHA-256 manifest. Every article URL, the three course pages, footer links, feeds, `/readyz`, and `/healthz` passed HTTP checks. [The live browser report](kazakh-ai-live-browser-checks-15.json) contains 204 page/viewport checks with no horizontal overflow. The full `CI`, `Rust course checks`, `Course`, and `docker-smoke` workflows for the lesson commit later completed successfully.

Live screenshot inspection found that the course sidebar held all 16 items but its 352-pixel scrolling list started at the top, hiding the current lesson after lesson 10. Before the fix, a browser measurement of lesson 11 returned `items: 16`, `scrollTop: 0`, `visible: false`. Commit `429a990bd3e3f106e111d33068ba0e077c8d8ade` scrolls only the inner list to its current item. The app image was rebuilt and replaced; the previous image was tagged `shanraq-app:pre-kazakh-15`. After deployment, lesson 11 returned `scrollTop: 209`, `visible: true`; lesson 15 in kz and lesson 1 in ru were also visible. `docker-smoke` and the main `checks` job for this fix succeeded before deployment; its full `CI` workflow later completed successfully.

## Fourth release: lessons 16–20

On 2026-09-24 the fourth block was prepared in ru/kz/en. Its five standalone Python snapshots and JSON catalogs live in `course/kazakh-ai/step-16` through `step-20`. The lesson test compares both the printed Python and JSON with these files in every locale. A novice walkthrough covered a normal answer, an absent fact, a duplicate or conflicting fact, and an expired fact. The final example joins the short Kazakh question parser to sourced lookup and refusal. All activities and source IDs in this block are fictional.

Preflight: 589 Markdown files passed language checks; offline link checks covered 1306 internal and 476 external references. The lesson tests, 63-entry transactional release manifest, Go site/article tests, coursecheck self-tests and `gofmt` passed. Local browser preview checked 63 pages at 390, 768, 844 and 1366 pixels (252 page/viewport checks) without horizontal overflow; new Kazakh mobile, Russian tablet and English desktop pages were visually reviewed.

Publication used source commit `96b86efcda631cd42160921312fb4bd17e19d29a` after the owner's five-minute CI observation window. `Course` and `docker-smoke` had succeeded; `CI` and `Rust course checks` had no failure at activation. The verified backup is `/var/backups/shanraq/kazakh-ai-before-16-20-20260924.dump` (4,723,745 bytes, SHA-256 `4b1d08e29212088a81bd087acf42986c25abc8cd2a4b8c49c9a5b3a5cc8c8e16`). The publication SQL SHA-256 was `51894952c5f6c3469fc09d2bb20ac9606790431200ac54b2ac0ba5d74568bed1` on both workstation and server. The production source fast-forwarded, and the SQL transaction completed successfully without an application rebuild.

Production has 21 series items and 63 translations. All 63 titles, bodies, and summaries match the release SHA-256 manifest. All article URLs, three course pages, footer links, feeds, `/readyz`, and `/healthz` passed HTTP checks. [Live browser checks](kazakh-ai-live-browser-checks-20.json) covered all 63 articles plus the three course pages at four widths (264 checks), with no horizontal overflow. Kazakh mobile, Russian tablet, English desktop, and Kazakh course mobile screenshots were reviewed. After publication `Course`, `docker-smoke`, and `Rust course checks` had succeeded; the long `CI` workflow was still running without a reported failure.

## Fifth release: lessons 21–25 and shorter names

On 2026-09-24 the fifth block was prepared in ru/kz/en. The footer course label is now “Курс ИИ” in Russian, “ЖИ курсы” in Kazakh, and “AI course” in English. The Russian course title is “ИИ без LLM: создаем свою модель ИИ”; Kazakh and English titles carry the same meaning. The course URLs and lesson slugs are unchanged.

The five standalone Python programs and their JSON datasets live in `course/kazakh-ai/step-21` through `step-25`. Thirteen fictional questions are visible during training and tuning (8/5); six additional labeled test questions are opened only in lesson 25. A word and word-prefix count classifier demonstrates a dual-intent mistake, then an explicit cue rule stops it. The test reports 2/2 correct predictions among answered questions, 2/4 known-type questions covered, and 4/6 refusals. These are exercise results on a tiny dataset, not a product benchmark or a claim about fact accuracy.

Preflight: novice walkthrough exercised an added label, duplicate question, the classifier's dual-intent error, its correction, and all six held-out examples. All 604 course Markdown files passed language checks; offline links covered 1306 internal and 484 external references. Lesson tests, the 78-entry release manifest, Go site/article tests, coursecheck self-tests and `gofmt` passed. Local browser preview checked 78 pages at 390, 768, 844, and 1366 pixels (312 page/viewport checks) without horizontal overflow; new Kazakh mobile, Russian tablet and English desktop pages were visually reviewed.

Publication used corrected source commit `2a2b3a96a9103d9001598f8c192cd82a00500d86` after the owner-authorized five-minute CI observation window. `Course` and `docker-smoke` succeeded; long `CI` had no reported failure at activation. The verified backup is `/var/backups/shanraq/kazakh-ai-before-21-25-20260924.dump` (4,730,188 bytes, SHA-256 `0aa7574a6f367e3fc2291dc6c41d6a389b4f2d8113376da7878c2a6929035cb5`). The final SQL SHA-256 was `64a96841b5e777b25da1f7d27c341864fd35679cbe9f287e85ec99ce1a78d61f` on workstation and server. The new application image was built, activated, and passed `/readyz` before the content transaction.

Production has 26 series items and 78 translations. All 78 titles, bodies, and summaries match the release SHA-256 manifest. All article URLs, three course pages, footer labels and links, feeds, `/readyz`, and `/healthz` passed HTTP checks. [Live browser checks](kazakh-ai-live-browser-checks-25.json) covered all 78 articles plus three course pages at four widths (324 checks), with no horizontal overflow. Kazakh mobile, Russian tablet, English desktop, and Kazakh course mobile screenshots were reviewed. The long `CI` workflow was still running without a reported failure at the last check.

## Lessons 26–30: final block

The final five steps connect the educational parser to an optional `qazaq-ir` comparison, guard unfamiliar or ambiguous questions, measure the local program, trace every fact-backed answer, and expose a reproducible command-line interface. The optional comparison was exercised with local `qazaq-ir` 0.31.0: the root matched, while `analysis_status` was `partial`, so the lesson reports only that limited result. The exercise does not claim a measured thousandfold advantage over an LLM or zero risk of false source data.

Preflight: a novice walkthrough found and corrected the placement of `--trace` in all three lesson texts. All five runnable steps match the printed code in ru/kz/en. Seven lesson checks, the 93-entry release manifest, Go site/article tests, coursecheck self-tests, 620-file language check, and 1306-internal/490-external offline link check passed. Local browser preview covered 93 pages at four widths (372 checks) without horizontal overflow; new Kazakh mobile, English desktop, and Russian tablet pages were visually reviewed. The prepared publication SQL SHA-256 is `61604b53e9a7164fb4beb0df41486183ec57f4edb90d5181b437bd9ed12d3698`.
