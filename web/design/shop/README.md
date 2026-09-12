# Обложка книги по Go

`go-book-cover.svg` — исходник обложки, самодостаточный: маскот вшит в файл,
поэтому обложка видна в любом просмотрщике, а не только в браузере. Оригинал
рисунка живёт в репозитории книги — `gopher-shop/docs/public/gopher.png`
(сделан нейросетью, дата и запрос — у владельца).

На сайт уходит не он, а растровая копия `web/static/shop/go-book-cover.jpg`
(1080 px, ~137 КБ). Причина: SVG, показанный через `<img>`, отдельной картинки
не подгружает, а восемьсот килобайт base64 в карточке товара — дорого за
миниатюру.

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
