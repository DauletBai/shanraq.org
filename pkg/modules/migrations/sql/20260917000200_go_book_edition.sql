-- +goose Up
-- Update obsolete seeded copy while preserving operator-written descriptions and sale settings.
UPDATE shop_products
SET edition = CASE WHEN edition IN ('0.4', '0.19.0-draft', '0.20.0-draft') THEN '0.21.0-draft' ELSE edition END,
    title_kz = CASE WHEN title_kz = 'Go: бірінші жолдан кітап дүкеніне дейін' THEN 'Go: бірінші жолдан интернет-дүкенге дейін' ELSE title_kz END,
    title_ru = CASE WHEN title_ru = 'Go: от первой строки до книжного магазина' THEN 'Go: от первой строки до интернет-магазина' ELSE title_ru END,
    title_en = CASE WHEN title_en = 'Go: from the first line to a bookshop' THEN 'Go: from the first line to an online shop' ELSE title_en END,
    body_kz = replace(body_kz, 'Қазір кіріспе және 1–8 тараулар дайын: жұмыс орны, алғашқы бағдарлама, мәндер, шарттар, функциялар мен алғашқы тест, жолдар, каталог, құрылымдар мен сілтемелер. Қалған тараулар жазылып жатыр.', 'Толық жұмыс редакциясы 42 тараудан тұрады. Кіріспе мен 1–12 тарауларды тегін онлайн оқуға немесе PDF, EPUB және HTML түрінде жүктеуге болады. Тегін үзіндіде 18 шағын іске қосылатын тәжірибе бар. Кітап орыс тілінде; оқу жобасындағы төлем нақты ақша алмайтын тест режимінде жұмыс істейді.'),
    body_ru = replace(body_ru, 'Сейчас готовы вступление и главы 1–8: рабочее место, первая программа, значения, условия, функции и первый тест, строки, каталог, структуры и указатели. Остальные главы пишутся.', 'Полная рабочая редакция содержит 42 главы. Бесплатно доступны предисловие и главы 1–12: их можно читать на сайте или скачать в PDF, EPUB и HTML. В ознакомительном фрагменте — 18 небольших запускаемых опытов. Оплата в учебном проекте работает в тестовом режиме без списания реальных денег.'),
    body_en = replace(body_en, 'The preface and chapters 1–8 are ready: the workspace, the first program, values, conditions, functions and the first test, strings, the catalogue, structs and pointers. The remaining chapters are being written.', 'The complete working draft contains 42 chapters. The preface and chapters 1–12 are free to read online or download as PDF, EPUB and HTML. The sample includes 18 small runnable experiments. The book is in Russian; payments in its teaching project use a sandbox without charging real money.'),
    updated_at = now()
WHERE slug = 'go-book';

-- +goose Down
-- Editorial corrections remain accurate when rolling back application code.
SELECT 1;
