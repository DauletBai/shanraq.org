-- +goose Up
-- Lending programmes for the loan calculator.
--
-- A rate is not a constant: it is a claim about a lender's offer on a given
-- day, and a mortgage is a twenty-year commitment, so a stale figure here
-- misleads someone about the largest purchase of their life. Every row carries
-- the source it came from and the date it was last checked, and the page prints
-- both -- the discipline the exchange-rate page already follows.
--
-- Rows seeded without a source and without a checked date are deliberately
-- unverified: the page shows them as unchecked so the operator can see at a
-- glance what still has to be confirmed. Better a visible gap than a confident
-- wrong number.
--
-- Rates are stored in basis points so the arithmetic stays in integers:
-- 700 = 7.00 %. Terms are in months, because a consumer loan is not counted in
-- years and a mortgage should not force the form to change units.
CREATE TABLE IF NOT EXISTS loan_programs (
    code         text PRIMARY KEY,
    kind         text        NOT NULL,
    sort         integer     NOT NULL DEFAULT 0,
    active       boolean     NOT NULL DEFAULT TRUE,
    lender       text        NOT NULL DEFAULT '',
    name_kz      text        NOT NULL DEFAULT '',
    name_ru      text        NOT NULL DEFAULT '',
    name_en      text        NOT NULL DEFAULT '',
    note_kz      text        NOT NULL DEFAULT '',
    note_ru      text        NOT NULL DEFAULT '',
    note_en      text        NOT NULL DEFAULT '',
    rate_bp      integer     NOT NULL,
    max_months   integer     NOT NULL,
    min_down_pct integer     NOT NULL DEFAULT 0,
    source_url   text        NOT NULL DEFAULT '',
    checked_on   date,
    updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS loan_programs_kind_idx ON loan_programs (kind, sort) WHERE active;

INSERT INTO loan_programs
    (code, kind, sort, lender, name_kz, name_ru, name_en, note_kz, note_ru, note_en,
     rate_bp, max_months, min_down_pct, source_url, checked_on)
VALUES
    -- Mortgages. Checked against public sources on 26 August 2026.
    ('mortgage-7-20-25', 'mortgage', 10, 'Отбасы банк',
     '«7-20-25»', '«7-20-25»', '"7-20-25"',
     'Тек жаңа тұрғын үй. Жылына 100 млрд теңге бөлінеді, кезек ұзақ.',
     'Только новостройки. На программу выделяют 100 млрд тенге в год, очередь длинная.',
     'New builds only. Funded at 100 bn tenge a year; the queue is long.',
     700, 300, 20, 'https://krisha.kz/content/articles/2026/2026-ipoteka-na-novostroyki-2026-usloviya-vseh-bankov', DATE '2026-08-26'),

    ('mortgage-nauryz', 'mortgage', 20, 'Отбасы банк',
     '«Наурыз»', '«Наурыз»', '"Nauryz"',
     'Соңғы бес жылда тұрғын үйі болмағандарға. Отбасы банкте кемінде 2 млн теңге жинақ керек. Әлеуметтік осал топтарға — 7 %.',
     'Для тех, у кого не было жилья последние пять лет. Нужны накопления от 2 млн тенге в Отбасы банке. Социально уязвимым категориям — 7 %.',
     'For those without property for five years. Requires at least 2 m tenge saved at Otbasy Bank. 7 % for socially vulnerable categories.',
     900, 228, 20, 'https://informburo.kz/novosti/otbasy-bank-opublikoval-spisok-ucastnikov-programmy-nauryz-2026', DATE '2026-08-26'),

    ('mortgage-jana', 'mortgage', 30, 'Отбасы банк',
     'JAÑA ипотека', 'JAÑA ипотека', 'JAÑA mortgage',
     'Ақша кепілімен, қайталама нарық. Мөлшерлеме кепіл көлеміне қарай.',
     'Под заклад денег, вторичное жильё. Ставка зависит от размера заклада.',
     'Against a cash pledge, secondary market. The rate depends on the pledge.',
     500, 180, 30, 'https://krisha.kz/content/articles/2026/2026-ipoteka-na-vtorichku-2026-usloviya-vseh-bankov', DATE '2026-08-26'),

    ('mortgage-market', 'mortgage', 40, '',
     'Коммерциялық ипотека', 'Коммерческая ипотека', 'Commercial mortgage',
     'Нарықтық мөлшерлеме: банкке және жарнаға қарай 14–18 %.',
     'Рыночная ставка: 14–18 % в зависимости от банка и взноса.',
     'Market rate: 14–18 %, depending on the bank and the down payment.',
     1800, 240, 20, 'https://informburo.kz/cards/ipoteka-2026-stavki-snizili-trebovaniia-uzestocili', DATE '2026-08-26'),

    -- Everything below ships unverified: plausible starting values, no source
    -- and no date, so the page marks them unchecked until an operator confirms.
    ('auto-state', 'auto', 10, '', 'Жеңілдікті автонесие', 'Льготный автокредит', 'Subsidised car loan',
     'Отандық құрастырудағы жаңа көлікке.', 'На новые автомобили отечественной сборки.', 'For new, locally assembled cars.',
     400, 84, 20, '', NULL),
    ('auto-market', 'auto', 20, '', 'Банк автонесиесі', 'Автокредит банка', 'Bank car loan',
     'Кепілге көліктің өзі қалады.', 'Автомобиль остаётся в залоге.', 'The car itself stands as collateral.',
     2100, 84, 20, '', NULL),

    ('consumer-secured', 'consumer', 10, '', 'Кепілді тұтыну несиесі', 'Потребительский под залог', 'Secured consumer loan',
     'Кепілмен: ЖТСМ шегі — 40 %.', 'С залогом: предельная ГЭСВ — 40 %.', 'Secured: the effective-rate cap is 40 %.',
     2200, 120, 0, '', NULL),
    ('consumer-plain', 'consumer', 20, '', 'Кепілсіз тұтыну несиесі', 'Потребительский без залога', 'Unsecured consumer loan',
     'Кепілсіз: ЖТСМ шегі — 56 %.', 'Без залога: предельная ГЭСВ — 56 %.', 'Unsecured: the effective-rate cap is 56 %.',
     2800, 60, 0, '', NULL),

    ('installment-zero', 'installment', 10, '', 'Пайызсыз бөліп төлеу', 'Рассрочка без переплаты', 'Interest-free instalments',
     'Мөлшерлеме нөл: бағаға сатушы қосады.', 'Ставка ноль: наценку закладывает продавец в цену.', 'Zero rate: the seller builds the margin into the price.',
     0, 24, 0, '', NULL),

    ('business-state', 'business', 10, '', 'Кәсіпкерлікті қолдау', 'Господдержка бизнеса', 'State business support',
     'Субсидияланатын мөлшерлеме, талаптары бар.', 'Субсидируемая ставка, есть условия отбора.', 'Subsidised rate, with eligibility conditions.',
     700, 120, 0, '', NULL),
    ('business-market', 'business', 20, '', 'ШОБ несиесі', 'Кредит для МСБ', 'SME loan',
     'Айналым қаражатына немесе инвестицияға.', 'На пополнение оборотных средств или инвестиции.', 'For working capital or investment.',
     2000, 84, 0, '', NULL)
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS loan_programs;
