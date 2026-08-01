# Changelog

All notable changes to this project are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
