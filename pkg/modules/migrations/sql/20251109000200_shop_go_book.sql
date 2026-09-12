-- +goose Up
-- The first product: the Go book the site is writing.
--
-- It is seeded as a draft, which means only staff can open its page. The panel
-- at /admin/shop is where it is announced, priced and put on sale -- that is an
-- editorial decision, not a migration.
INSERT INTO shop_products (slug, kind, status, price, edition, cover_url, title_kz, title_ru, title_en,
                           summary_kz, summary_ru, summary_en, body_kz, body_ru, body_en, position)
VALUES (
    'go-book', 'book', 'draft', 0, '0.2', '/static/shop/go-book-cover.svg',
    'Go: бірінші жолдан кітап дүкеніне дейін',
    'Go: от первой строки до книжного магазина',
    'Go: from the first line to a bookshop',
    'Go тілін нөлден үйрететін практикалық кітап: бір мысал бірінші жолдан жұмыс істейтін интернет-дүкенге дейін өседі. Үш деңгей: junior, middle, senior.',
    'Практическая книга по Go с нуля: один пример вырастает от первой строки до работающего интернет-магазина. Три уровня: junior, middle, senior.',
    'A practical Go book from scratch: one example grows from the first line into a working online shop. Three levels: junior, middle, senior.',
    'Кітап Go тіліндегі алғашқы жолдан бастап жұмыс істейтін интернет-дүкенге дейін жеткізеді. Бұл — оқырманмен бірге өсетін бір мысал: әуелі бағдарлама терминалда кітап атауын басып шығарады, сосын баға мен күй, каталог, дерекқор, сервер, төлем және серверге орналастыру қосылады.

Кітаптың бағдарламасы — үш деңгейдегі 42 тарау. **Junior** — тіл және алғашқы бағдарламалар. **Middle** — дерекқор, HTTP, тесттер, құралдар. **Senior** — орнықтылық, өнімділік, қауіпсіздік және пайдалану.

Кітаптағы код — шынымен орындалатын код. Әр орындалатын тараудың өз бақылау нүктесі бар: жеке модуль, оны жинап, іске қосуға болады, ал мәтіндегі шығыс кітап жиналған сайын бағдарламаның нақты шығысымен салыстырылады.

Қазір кіріспе және 1–8 тараулар дайын: жұмыс орны, алғашқы бағдарлама, мәндер, шарттар, функциялар мен алғашқы тест, жолдар, каталог, құрылымдар мен сілтемелер. Қалған тараулар жазылып жатыр.

Кітап орыс тілінде жазылған. Пішімдері: PDF, EPUB және HTML.
',
    'Книга ведёт от первой строки на Go до работающего интернет-магазина. Это один пример, который растёт вместе с читателем: сначала программа печатает название книги в терминале, потом появляются цена и состояние, каталог, база данных, сервер, оплата и выкладка на сервер.

Программа книги — 42 главы на трёх уровнях. **Junior** — язык и первые программы. **Middle** — база данных, HTTP, тесты, инструменты. **Senior** — устойчивость, производительность, безопасность и эксплуатация.

Код в книге — это код, который запускается. У каждой исполняемой главы есть своя контрольная точка: отдельный модуль, который можно собрать и запустить, а вывод в тексте сверяется с настоящим выводом программы при каждой сборке книги.

Сейчас готовы вступление и главы 1–8: рабочее место, первая программа, значения, условия, функции и первый тест, строки, каталог, структуры и указатели. Остальные главы пишутся.

Книга написана на русском языке. Форматы: PDF, EPUB и HTML.
',
    'The book goes from the first line of Go to a working online shop. It is one example that grows with the reader: first the program prints a book''s title in the terminal, then come the price and the condition, the catalogue, the database, the server, payment, and a deployment of its own.

The programme is 42 chapters on three levels. **Junior** — the language and the first programs. **Middle** — database, HTTP, tests, tooling. **Senior** — resilience, performance, security and operations.

The code in the book is code that runs. Every executable chapter has a checkpoint of its own: a separate module you can build and run, and the output printed in the text is checked against the program''s real output every time the book is built.

The preface and chapters 1–8 are ready: the workspace, the first program, values, conditions, functions and the first test, strings, the catalogue, structs and pointers. The remaining chapters are being written.

The book is written in Russian. Formats: PDF, EPUB and HTML.
',
    10
)
ON CONFLICT (slug) DO NOTHING;

-- +goose Down
DELETE FROM shop_products WHERE slug = 'go-book';
