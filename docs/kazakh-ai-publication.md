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

The long CI workflow later failed only at `gofmt` for `pkg/site/i18n.go`; its unit and snapshot jobs passed. The file was formatted locally and `gofmt -l .` returned no files. This correction is committed separately and its CI result must be recorded here after completion.
