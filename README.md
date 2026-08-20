<p align="center">
  <img src="web/static/brand/shanraq.svg" alt="Shanraq" width="120">
</p>

<h1 align="center">Shanraq</h1>

<p align="center">
  <b>A censorship-resistant, trilingual publishing &amp; classifieds platform for Kazakhstan.</b><br>
  Independent journalism and real-estate listings — in Kazakh, Russian and English — under one roof.
</p>

<p align="center">
  <a href="https://github.com/DauletBai/shanraq.org/actions/workflows/ci.yml"><img src="https://github.com/DauletBai/shanraq.org/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26"></a>
  <a href="https://www.postgresql.org"><img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL 16"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-Proprietary-red" alt="License: Proprietary"></a>
  <a href="https://shanraq.org"><img src="https://img.shields.io/badge/live-shanraq.org-e53935" alt="Live: shanraq.org"></a>
</p>

---

**Shanraq** ([shanraq.org](https://shanraq.org)) is a modular Go platform that fuses two things regional media has always kept apart: **trustworthy local journalism** and a **modern classifieds marketplace** (real estate first, autos and more to come). It is built to keep reaching readers even when a single channel is blocked — through automatic syndication to Telegram and RSS — and to run three languages as one publication, not three separate sites.

The name refers to the *shanyrak* — the crown of the Kazakh yurt — the living centre that holds the whole roof together.

## Why it exists

- **Resilience over dependence.** Content is pushed to Telegram and RSS automatically on publish, so the audience stays reachable even if the website is unavailable.
- **One roof, not five apps.** News, real estate, and (next) autos and a marketplace live in a single product instead of scattered across separate sites and apps.
- **AI that removes drudgery, not judgement.** Optional machine translation turns one journalist's work into a trilingual, SEO-ready, auto-distributed publication. The platform does not write or edit for the author: the text is theirs, and so is the byline.
- **Built for a small, open economy.** Trilingual (kk/ru/en) with per-language SEO (hreflang, sitemaps, structured data), data stored in-country, self-hostable on a single VPS.

## Features

- 📰 **Trilingual publishing** — one article, three language variants; per-language `hreflang`, sitemaps and JSON-LD.
- 🏠 **Real-estate classifieds** — listings with photos, geo, amenities, promotion/feature tariffs, and a dedicated listings sitemap.
- 🤖 **AI translation & moderation** — provider-agnostic (Claude, ChatGPT or Kimi), off unless a key is set. Translation keeps the author's Markdown and links; the platform then compares the result against the original and reports anything mechanically lost, so an author can check a language they do not read.
- 📣 **Block-resilient syndication** — RSS always on; automatic Telegram posting on publish.
- 🖼️ **Media pipeline** — upload, EXIF-strip, brand-watermark; pluggable storage (filesystem now, S3/MinIO ready).
- ⭐ **Ratings & author karma** — weighted voting with anti-brigading.
- ◎ **Prediction ledger** — every forecast recorded with its date and judged in public at `/predictions`, misses included; the score is allowed to go down.
- 🔎 **Built to be cited** — sitemaps, IndexNow, conditional GET, a ready-made reference under every article, and an `llms.txt` that offers assistants only the human-written work.
- 🔗 **Referral loop** — invite links, attribution, and promotion-credit rewards.
- 🛠️ **Operator admin panel** — edit legal/info pages, tariffs, service flags, payment provider, and AI settings without a redeploy.
- 📊 **Privacy-respecting analytics** — aggregate-only audience metrics; bots separated from humans; traffic sources; **no per-visitor profiling**.
- 🔐 **Secure auth** — refresh-token rotation, RBAC, password-reset flows, CSRF protection.

## Tech stack

| Layer | Choice |
|------|--------|
| Language | **Go 1.26** (modular monolith) |
| Database | **PostgreSQL 16** via `pgxpool`, embedded Goose migrations |
| Web | Server-rendered `html/template`, progressive enhancement |
| Config | Viper — typed, env-first (`SHANRAQ_*`), no secrets in git |
| Observability | Zap structured logging, Prometheus metrics, health/readiness probes |
| Delivery | Docker Compose + Caddy (automatic HTTPS) on a single VPS |
| AI | Anthropic Claude (optional, provider-agnostic) |

## Quick start (Docker Compose + automatic HTTPS)

```bash
git clone https://github.com/DauletBai/shanraq.org.git
cd shanraq.org
cp .env.example .env
# edit .env: POSTGRES_PASSWORD, DOMAIN, SHANRAQ_PUBLIC_BASE_URL,
#            SHANRAQ_AUTH_TOKEN_SECRET (>=32 random chars), operator details.
docker compose -f docker-compose.prod.yml up -d --build
```

Migrations run on first boot; only Caddy is exposed publicly, and it fetches a TLS
certificate automatically. Secrets live only in `.env` (git-ignored) — in
`environment=production` the app refuses to start with a weak token secret or a
non-HTTPS base URL. Full guide: **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)**.

Local development:

```bash
make run      # start the app
make test     # unit tests
make smoke    # smoke test
```

## Architecture

A single Go binary composed of independent modules that implement the
`shanraq.Module` contract (init / start / routes):

```
articles · auth · ai · media · syndicate · notifier · jobs · ratings · migrations
```

Each module wires its own routes, migrations, and background workers, so features
compose without a service mesh. See **[docs/](docs/)** for the configuration
reference, module guides, and deployment runbook.

## Documentation

- [Deployment guide](docs/DEPLOYMENT.md) — VPS, Docker, backups, data migration.
- [Configuration reference](docs/CONFIGURATION.md) — every `SHANRAQ_*` setting.
- [Backups & resilience](docs/BACKUPS.md) — what runs, how to restore, and which single points of failure are still open.
- [Editorial & crawling policy](docs/EDITORIAL-POLICY.md) — who writes what, and which crawlers may read it.
- In-app docs at `/docs` on a running instance.

## Security

Secrets never live in the repository — configuration is env-first and `.env` is
git-ignored. To report a vulnerability, see **[SECURITY.md](SECURITY.md)**
(please do not open a public issue for security reports).

## Status & versioning

Live in **closed beta** at [shanraq.org](https://shanraq.org). Releases follow
[Semantic Versioning](https://semver.org); see **[CHANGELOG.md](CHANGELOG.md)**.

Cutting a release is one command — pushing an annotated `vX.Y.Z` tag. The
`Release` workflow then publishes the GitHub Release with the notes taken from
the matching CHANGELOG section, and fails loudly if that section is missing, so
a tag can never ship without its notes.

## License

**Proprietary — all rights reserved.** The source is published for transparency
and review; it is **not** open-source and may not be used, deployed, or
redistributed without written permission. See **[LICENSE](LICENSE)**.
For partnership or commercial-use inquiries: **shanirak.org@gmail.com**.
