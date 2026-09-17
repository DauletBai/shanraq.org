# Free Go book sample

The product page `/shop/go-book` links to the free preface and chapters 1–12 in Russian. Readers can open the HTML book without signing in or download PDF, EPUB, offline HTML and the example code. The last chapter links back to the full product.

The checked static package is embedded from `web/static/shop/go-book-sample/0.21.0-sample/`. Its `sample.json` declares the chapter boundary and file hashes. Chapters 13–42 and the full paid book are not included. Build the application again after updating these static files.

In the shop admin, use `/static/shop/go-book-sample/0.21.0-sample/index.html` for the preview URL and `/static/shop/go-book-sample/0.21.0-sample/read/cover/go-book-cover-v3.png` for the cover. The current Russian title is “Go: от первой строки до интернет-магазина”. The product remains in its existing sale state; a free sample does not enable payment or announce availability.

To regenerate, use the public gopher-shop repository: run the book checks, build and validate the sample, then `python book/tools/prepare_sample_site.py`. Inspect the website ZIP before replacing this version. Keep each new edition in a separate directory so cached files from different editions do not mix.

Revision 0.21 contains 116 PDF pages and 18 runnable learning labs in the first twelve chapters. The sidebar banner uses the same versioned cover as the sample; changing the URL avoids a stale browser/CDN image. Migration 20260917000100 updates the product cover and preview URL, preserving its price and sale state. Earlier sample assets remain available for existing links.

After deployment, check the product's five sample links without signing in, open chapter 12, and verify chapter 13 returns 404. The running binary must be rebuilt: uploading files beside an old binary does not replace embedded static assets.

Migration 20260917000200 replaces obsolete seeded titles, edition 0.4 and the old eight-chapter progress paragraph with the current 42-chapter draft and twelve-chapter sample. Custom descriptions and sale settings are preserved.
