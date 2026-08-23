-- +goose Up
-- Курсы валют, которые мы храним сами.
--
-- Нацбанк отдаёт исторические курсы примерно за последние пять лет: запрос за
-- 2020 год и раньше возвращает пустой ответ. Это значит, что источник — не
-- архив, а окно, и оно едет вперёд вместе с сегодняшним днём. Поэтому мы
-- складываем каждый день к себе: через год наша таблица будет глубже, чем то,
-- что банк готов отдать, а «за весь период» перестанет упираться в чужое окно.
--
-- quant — за сколько единиц валюты дан курс (за 100 йен, за 10 гривен).
-- Хранится рядом со значением, потому что без него сравнивать валюты нельзя.
CREATE TABLE IF NOT EXISTS fx_rates (
    day   date          NOT NULL,
    code  text          NOT NULL,
    value numeric(14,4) NOT NULL,
    quant integer       NOT NULL DEFAULT 1,
    name  text          NOT NULL DEFAULT '',
    PRIMARY KEY (day, code)
);

-- Ряд по одной валюте за период — основной запрос страницы.
CREATE INDEX IF NOT EXISTS idx_fx_rates_code_day ON fx_rates (code, day);

-- Журнал опросов: какой день мы уже спрашивали и сколько курсов получили.
--
-- Без него догрузка после каждого перезапуска заново спрашивала бы банк про
-- каждые выходные за пять лет — тысячу запросов ни за чем. Ноль в found это
-- ответ «в этот день курсов нет», и он такой же результат, как и число: он
-- показывает и выходные, и границу, за которой источник молчит.
CREATE TABLE IF NOT EXISTS fx_probed (
    day        date        NOT NULL PRIMARY KEY,
    found      integer     NOT NULL DEFAULT 0,
    probed_at  timestamptz NOT NULL DEFAULT now()
);

-- Месячный архив с ноября 1993 года — от рождения тенге.
--
-- Нацбанк отдаёт дневной курс лет на пять назад и молчит про всё, что раньше.
-- Поэтому глубину берём у Банка международных расчётов: у него месячный ряд
-- курсов к доллару с 1993-11, и из него кросс-курсом получается тенге за любую
-- из наших валют. Хранится тенге за ОДНУ единицу; привычные «за 100 иен»
-- делает уже отображение.
CREATE TABLE IF NOT EXISTS fx_monthly (
    month date          NOT NULL,
    code  text          NOT NULL,
    value numeric(18,6) NOT NULL,
    PRIMARY KEY (month, code)
);
CREATE INDEX IF NOT EXISTS idx_fx_monthly_code_month ON fx_monthly (code, month);

-- +goose Down
DROP TABLE IF EXISTS fx_monthly;
DROP TABLE IF EXISTS fx_probed;
DROP TABLE IF EXISTS fx_rates;
