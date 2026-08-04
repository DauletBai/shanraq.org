# Changelog

All notable changes to this project are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- A "follow us" card in the article aside, under the table of contents. The
  aside is sticky on desktop, so the ask now travels with the reader for the
  whole piece instead of waiting at the bottom, where only those who finish ever
  see it. Only configured profiles are listed — a `#` placeholder (YouTube today)
  would be a dead click on a card whose single job is to be clicked. Clicks
  report an aggregate `follow_social` event to the existing counter, so the card
  can be judged on numbers rather than on taste.

### Changed
- The home sidebar's email newsletter block is replaced by the same follow card
  (accent variant, keeping that column's one point of colour). The form was not
  broken — addresses were stored and the weekly digest ran — but it never sent a
  confirmation, so subscribing felt like nothing happened, and in a week it
  collected no one. The `/subscribe` route, the `subscribers` table and the
  digest job are deliberately left intact: this is a UI decision, reversible by
  restoring the block.

### Added
- Article cards now carry the view counter, which until now only existed on the
  article page itself: `👍 1 · 👁 115 · D. Baimurza · 04.08.26`. Counters above
  999 collapse to a localized short form — `1,2 мың` / `1,2 тыс.` / `1.2k`,
  truncated rather than rounded so a card never claims more views than there
  were.

### Changed
- Bylines on cards are abbreviated to initial + family name ("Daulet Baimurza"
  → "D. Baimurza"). The given name is the part that gets cut — the byline
  convention everywhere space is tight — freeing the width the counter needed.
- The card footer (rating, views, byline, date) now comes from one shared
  `post_meta` partial instead of three near-copies in home / favourites /
  author templates, which had already drifted (favourites still drew the rating
  as a bare "▲").

## [0.10.2] — 2026-08-04

### Changed
- The team's own browsing no longer pollutes analytics. Staff (admin/operator)
  visits were already skipped, but the owner's publish routine — logging into the
  admin account, then a test account to like the new article — still generated
  guest page views and login clicks while logged out. A device is now flagged
  with a persistent opt-out cookie the moment it loads a page as staff (or as a
  configured `SHANRAQ_ANALYTICS_EXCLUDE_EMAILS` address), so it stays out of the
  counts even after logging out or switching accounts. Click events honour the
  same rule.

### Added
- A "Data center / VPN — reading language" panel that splits the masked VPN
  bucket by language, because that split identifies who is behind the VPN:
  Russian ≈ Russia/CIS, while English ≈ genuine international readers (China,
  Iran, the West) who bridge through English rather than Russian.

### Changed
- Reworded the "English readers — origin" note: English from the Data center /
  VPN bucket is now described as *most likely genuine foreigners* (China, Iran,
  the West), not vaguely "masked" — they read English precisely because they
  don't read Russian.

## [0.10.0] — 2026-08-03

### Added
- Audience analytics now tracks the **reading language** (Kazakh / Russian /
  English) and crosses it with visitor origin, to answer whether the English
  audience is genuine foreign readers or something else. Two new 30-day panels:
  "Reading language" (the overall kk/ru/en mix) and "English readers — origin"
  (where the English-version visits come from). The origin panel is the sharp
  signal a VPN cannot mask: English from a real foreign country is a genuine
  foreigner; English from Kazakhstan is a curious local; English from
  "Data center / VPN" is masked traffic — e.g. a reader on a VPN from a country
  with a restricted internet (Russia, China, Iran). Aggregate-only as always.

## [0.9.0] — 2026-08-03

### Added
- The country panel now separates **hosting/cloud/VPN traffic from real
  readers**. An optional ASN database (DB-IP ASN Lite) flags visits from cloud
  networks (AWS, Google, Azure, Cloudflare, OVH, Hetzner, VPN exits, …) and
  buckets them as "ЦОД / VPN" instead of a country — so the geographic rows
  reflect actual eyeballs, not US-hosted infrastructure that otherwise dominates
  by IP. Only the coarse label is counted; the IP is still discarded. Optional
  and graceful: without the ASN file every IP is bucketed by country as before.

### Changed
- Bot detection now also recognises common scraper/library agents (Scrapy,
  axios, node-fetch, Guzzle, aiohttp, Postman, …), so fewer automated hits leak
  into the human audience counts.

## [0.8.0] — 2026-08-03

### Added
- Audience analytics now shows a **visitor-country** panel (30-day) so the team
  can tell domestic (Kazakhstan) readers from a genuine foreign audience — the
  English content shared on LinkedIn draws both, and by-language view counts
  alone couldn't say which. Countries are resolved from the visitor IP with a
  local DB-IP Lite database (bind-mounted read-only, refreshed monthly on the
  host); only the coarse country code is counted and the IP is immediately
  discarded, exactly like the User-Agent in the device/OS panels — no
  per-visitor profiling. The feature is optional: with no database present the
  panel simply stays empty and nothing else changes.

## [0.7.2] — 2026-08-03

### Changed
- The author phone-verification form now states the one-time code arrives **via
  Telegram** (with a nudge to have Telegram on that number), so a user without
  Telegram isn't sent down a dead end. Labels updated from "SMS" to "Telegram".

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
