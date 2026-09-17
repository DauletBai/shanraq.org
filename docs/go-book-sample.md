# Free Go book sample

The product page `/shop/go-book` links to the free preface and chapters 1–12 in Russian. Readers can open the HTML book without signing in or download PDF, EPUB, offline HTML and the example code. The last chapter links back to the full product.

The checked static package is embedded from `web/static/shop/go-book-sample/0.19.0-sample/`. Its `sample.json` declares the chapter boundary and file hashes. Chapters 13–42 and the full paid book are not included. Build the application again after updating these static files.

In the shop admin, use `/static/shop/go-book-sample/0.19.0-sample/index.html` for the preview URL and `/static/shop/go-book-sample/0.19.0-sample/read/cover/go-book-cover-v3.png` for the cover. The current Russian title is “Go: от первой строки до интернет-магазина”. The product remains in its existing sale state; a free sample does not enable payment or announce availability.

To regenerate, use the public gopher-shop repository: run the book checks, build and validate the sample, then `python book/tools/prepare_sample_site.py`. Inspect the website ZIP before replacing this version. Keep each new edition in a separate directory so cached files from different editions do not mix.
