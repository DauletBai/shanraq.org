# Обложка книги по Go

Здесь лежит исходник, а на сайт уходит растровая обложка
`web/static/shop/go-book-cover.jpg`. Так сделано потому, что маскот —
растровый рисунок: SVG, показанный через `<img>`, не имеет права подгружать
внешнюю картинку, а вшитая в него base64-копия сделала бы карточку товара
тяжелее полумегабайта.

- `go-book-cover.svg` — исходник: поле, логотип Go, заголовок, подвал.
- `gopher-cart.png` — маскот книги (оригинал живёт в репозитории книги
  `gopher-shop`, здесь уменьшенная до 860 px копия).

Пересобрать обложку после правки исходника:

```sh
"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" \
  --headless=new --disable-gpu --force-device-scale-factor=2 \
  --window-size=720,960 --screenshot=/tmp/cover@2x.png \
  file://$PWD/web/design/shop/go-book-cover.svg
sips -Z 1080 /tmp/cover@2x.png --out /tmp/cover.png
sips -s format jpeg -s formatOptions 86 /tmp/cover.png \
  --out web/static/shop/go-book-cover.jpg
```

Адрес обложки хранится у товара в поле `cover_url` — его можно поменять в
панели `/admin/shop`, не трогая код.
