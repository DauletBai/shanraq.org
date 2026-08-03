# Changelog

All notable changes to this project are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.7.1] — 2026-08-03

### Added
- SMS gateway can now deliver verification codes over **Telegram** (SMSC `tg=1`)
  via a new `SHANRAQ_SMS_CHANNEL=telegram` setting — no paid operator sender
  name, far cheaper than SMS, and it works immediately once the SMSC account is
  active.

## [0.7.0] — 2026-08-03

### Added
- Audience analytics now shows the **device, OS and browser mix** (mobile /
  tablet / desktop; Android / iOS / Windows / macOS / Linux; Chrome / Safari /
  Firefox / Edge / …) as aggregate 30-day panels — coarse, no per-visitor
  profiling. The User-Agent is classified and discarded, like the bot check.

### Changed
- The team's own visits no longer count as audience: page views from a logged-in
  admin/operator are skipped, so the owner's constant browsing stops inflating
  the guest and "Direct" counts.

## [0.6.5] — 2026-08-02

### Changed
- Dependency maintenance (Dependabot): bumped the Go dependency group (pgx,
  viper, zap, prometheus client, jwt, golang.org/x/crypto, goose, validator,
  anthropic-sdk and others — all minor/patch) and CI GitHub Actions
  (checkout@v7, setup-go@v7, docker/setup-buildx-action@v4). Build and full
  test suite green.

## [0.6.4] — 2026-08-02

### Changed
- Commercial SEO scanners (already turned away in robots.txt) are no longer
  recorded in analytics at all — they neither count as guests nor clutter the
  bot panel. Existing historical rows are cleared. Yandex source detection now
  also catches the short domain `ya.ru`.

## [0.6.3] — 2026-08-02

### Changed
- Telegram auto-posts now tag the article link with `?utm_source=telegram`, so
  visits from the channel are attributed to Telegram in analytics instead of
  falling into "Direct" (the messenger strips the referrer). No manual work per
  post.

## [0.6.2] — 2026-08-02

### Added
- UTM attribution: traffic analytics now reads `?utm_source=` first (mapped to a
  known source label) and falls back to the Referer host. Links shared with
  `?utm_source=telegram` / `instagram` / etc. are attributed correctly even when
  the browser strips the referrer (messengers, in-app browsers) — no third-party
  service required.

### Changed
- robots.txt now turns away commercial SEO scanners (Ahrefs, Semrush, MJ12, Dot,
  BLEX, DataForSeo, Petal, MegaIndex) — the heaviest crawlers that return no
  value. Search engines and AI crawlers stay allowed (AI-answer discovery is a
  channel worth keeping).

## [0.6.1] — 2026-08-02

### Fixed
- Field-help tooltip was unreadable on the dark theme (white text on the light
  tooltip background). Dark theme now uses dark text; light theme keeps its dark
  tooltip with white text.
- Audience tiles used calendar windows, so at the start of a month the "week"
  count could exceed the "month" count (the current week reached back into the
  previous month). The tiles now use rolling windows — today, last 7, 30 and 365
  days — so they always nest: today ≤ week ≤ month ≤ year.

## [0.6.0] — 2026-08-02

### Added
- Address ↔ map sync on the listing form. A new "Find on the map by address"
  button geocodes the entered address (cascade + street/house) and drops a
  precise pin; placing or dragging the pin reverse-geocodes it and fills in the
  street, house number and microdistrict — filling only what the geocoder
  returns, never clearing what the author typed. Geocoding is proxied
  server-side over OpenStreetMap Nominatim (the browser CSP blocks direct calls)
  with an in-memory cache; the country/region/city cascade stays manual by
  design (it is bound to the platform's own geo database).

## [0.5.1] — 2026-08-02

### Changed
- Listing form now reflects the currency before submit: picking Russia in the
  location cascade flips the price hint to rubles and shows a ₽ chip next to the
  price, while Kazakhstan shows ₸. The geo API now returns each node's country so
  the form can react to it. (The stored currency was already correct in 0.5.0;
  this removes the confusing "in tenge" hint for Russian listings.)

## [0.5.0] — 2026-08-02

### Added
- Listings for property in Russia, not just Kazakhstan. The location cascade
  already carried both countries; now each listing also carries its own currency
  — tenge (₸) for Kazakh addresses, ruble (₽) for Russian ones — chosen
  automatically from the selected location and shown on cards, the listing page,
  "my listings", favorites and JSON-LD. Posting needs only a verified email (no
  phone/SMS), so Russian users can list right away.

### Changed
- The mandatory Kazakh title is waived for Russian listings: when the location is
  in Russia, only Russian and English titles are required (Kazakh optional). The
  price-field label no longer bakes in ₸, since the currency now follows country.

## [0.4.1] — 2026-08-01

### Fixed
- Info-bar exchange rates could stay blank for up to 6 hours after a single
  transient fetch failure at boot. The National Bank fetch now retries with a
  short backoff, the HTTP timeout is more forgiving (10s), and empty rates are
  re-fetched on the 30-minute tick instead of only every 6 hours.

## [0.4.0] — 2026-08-01

### Added
- SMS phone verification for agents/authors: a provider-agnostic SMS gateway
  (`pkg/modules/sms`) with Mobizon.kz and SMSC.kz backends, chosen by config
  (`SHANRAQ_SMS_PROVIDER` + credentials). The verification flow already existed
  (code mint, hash, rate-limit, confirm); this wires real delivery. With no
  provider set, SMS stays off and codes are dev-logged — switching provider is a
  single environment variable, so onboarding friction with one aggregator never
  blocks the platform.

## [0.3.1] — 2026-08-01

### Changed
- Listing form: each language tab's placeholders now read in that tab's own
  language (Russian tab in Russian, Kazakh in Kazakh, English in English),
  regardless of the interface language — so the hint matches the content asked for.

## [0.3.0] — 2026-08-01

### Added
- Trilingual listings (KZ/RU/EN) — the flagship, mandatory feature: title and
  description in all three languages via a tabbed form; each reader sees the
  listing in their language with fallback. Script sanity on submit — English
  must be Latin, Russian Cyrillic, Kazakh either (Cyrillic or the new Latin).

## [0.2.0] — 2026-08-01

### Added
- Listing documents: agents can attach PDF plans/passports/contracts and image
  schemes to a listing (`/media/upload-doc`), shown in a "Documents / floor plan"
  block on the listing page.
- Traffic-source analytics (referrer → Google/Yandex/Telegram/Facebook/direct/…).
- Separate bot vs human classification so audience counts reflect real people.
- Dedicated real-estate sitemap (`/sitemap-listings.xml`) for Search Console.
- Branded PNG Open Graph card and descriptive homepage title for link previews.
- Brand watermark overlay on article covers and listing photos.
- GitHub icon in the footer; Telegram and Facebook social links from config.

### Changed
- Telegram bot token / chat id now bind from environment (secret never in config).

## [0.1.0] — 2026-07-30

First tagged release. Live in closed beta at [shanraq.org](https://shanraq.org).

### Added
- Trilingual (kk/ru/en) publishing with per-language SEO (hreflang, sitemaps, JSON-LD).
- Real-estate classifieds with photos, geo, amenities, and promotion/feature tariffs.
- Block-resilient syndication: always-on RSS and automatic Telegram posting on publish.
- Optional Claude-powered AI co-editor and trilingual auto-translation (off by default).
- Media pipeline: upload, EXIF strip, brand watermark; pluggable storage backend.
- Ratings with weighted author karma and anti-brigading.
- Referral loop: invite links, attribution, promotion-credit rewards.
- Operator admin panel: editable legal/info pages, tariffs, service flags, payment
  provider, and AI settings — no redeploy required.
- Aggregate-only, privacy-respecting audience analytics (no per-visitor profiling).
- Secure auth: refresh-token rotation, RBAC, password-reset flows, CSRF protection.
- Production stack: Docker Compose + Caddy automatic HTTPS, embedded Goose migrations.

[Unreleased]: https://github.com/DauletBai/shanraq.org/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/DauletBai/shanraq.org/releases/tag/v0.1.0
