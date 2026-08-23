package articles

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

// Пустой ответ банка — это не поломка. За выходные и за даты вне его окна
// приходит документ без единой валюты, и если считать это ошибкой, догрузка
// будет вечно спотыкаться о каждую субботу.
func TestAnEmptyDayFromTheBankIsNotAnError(t *testing.T) {
	rates, err := parseFxRates([]byte(`<?xml version="1.0"?><rates></rates>`), time.Now())
	if err != nil {
		t.Fatalf("пустой день признан ошибкой: %v", err)
	}
	if len(rates) != 0 {
		t.Fatalf("из пустого документа взялось %d курсов", len(rates))
	}
}

// Банк объявляет иену сотнями, и кратность надо сохранить: без неё курс иены
// окажется в сто раз больше, чем он есть.
func TestTheBanksMultiplierSurvivesParsing(t *testing.T) {
	body := []byte(`<rates>
		<item><fullname>ДОЛЛАР США</fullname><title>USD</title><description>456.88</description><quant>1</quant></item>
		<item><fullname>ЯПОНСКАЯ ЙЕНА</fullname><title>JPY</title><description>304.12</description><quant>100</quant></item>
	</rates>`)
	rates, err := parseFxRates(body, time.Date(2026, 8, 22, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if len(rates) != 2 {
		t.Fatalf("разобрано %d курсов вместо двух", len(rates))
	}
	by := map[string]FxRate{}
	for _, r := range rates {
		by[r.Code] = r
	}
	if by["JPY"].Quant != 100 {
		t.Errorf("кратность иены потеряна: %d", by["JPY"].Quant)
	}
	if by["USD"].Value != 456.88 || by["USD"].Name == "" {
		t.Errorf("доллар разобран неверно: %+v", by["USD"])
	}
}

// В выгрузке BIS все валюты меряются долларом. Тенге за единицу валюты — это
// частное двух её строк, и если считать иначе, глубокая история будет врать.
func TestTheDeepArchiveIsCrossedThroughTheDollar(t *testing.T) {
	csv := `FREQ,REF_AREA,CURRENCY,COLLECTION,TIME_PERIOD,OBS_VALUE
M,KZ,KZT,E,2026-06,480.72
M,XM,EUR,E,2026-06,0.8777
M,KZ,KZT,E,1993-11,4.70
M,XM,EUR,E,1993-11,NaN
`
	rows, err := parseBISMonthly(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]float64{}
	for _, r := range rows {
		got[r.Code+"@"+r.Month.Format("2006-01")] = r.Value
	}
	// Доллар — база ряда, его курс к тенге берётся напрямую.
	if v := got["USD@2026-06"]; v < 480.7 || v > 480.8 {
		t.Errorf("доллар посчитан как %.2f вместо 480.72", v)
	}
	if v := got["EUR@2026-06"]; v < 547 || v > 549 {
		t.Errorf("евро посчитано как %.2f, а 480.72/0.8777 это около 547.7", v)
	}
	// Месяц рождения тенге доллар покрывает, а евро тогда ещё не было.
	if v := got["USD@1993-11"]; v < 4.69 || v > 4.71 {
		t.Errorf("ноябрь 1993 потерян: %.2f", v)
	}
	if _, ok := got["EUR@1993-11"]; ok {
		t.Error("евро появилось в 1993 году, хотя в выгрузке там пусто")
	}
}

// Месяц без курса тенге нельзя пересчитать ни во что: делить не на что.
// Такой месяц должен выпадать целиком, а не превращаться в ноль.
func TestAMonthWithoutTheTengeIsDropped(t *testing.T) {
	rows, err := parseBISMonthly(strings.NewReader(
		"FREQ,REF_AREA,CURRENCY,COLLECTION,TIME_PERIOD,OBS_VALUE\nM,XM,EUR,E,1998-01,0.9\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("без курса тенге получилось %d строк", len(rows))
	}
}

// Прореживание длинного ряда не должно терять его концы: именно по ним
// страница считает изменение за период.
func TestThinningKeepsBothEndsOfTheSeries(t *testing.T) {
	pts := make([]FxPoint, 3000)
	base := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range pts {
		pts[i] = FxPoint{Day: base.AddDate(0, 0, i), Value: float64(i)}
	}
	out := fxThin(pts, fxMaxPoints)
	if len(out) != fxMaxPoints {
		t.Fatalf("после прореживания %d точек вместо %d", len(out), fxMaxPoints)
	}
	if out[0] != pts[0] {
		t.Error("потеряно начало ряда")
	}
	if out[len(out)-1] != pts[len(pts)-1] {
		t.Error("потерян конец ряда — по нему считается текущий курс")
	}
}

// Мелкие валюты вроде узбекского сума меньше единицы. Два знака после запятой
// превратили бы весь их ряд в одинаковые нули.
func TestSmallRatesKeepEnoughDigits(t *testing.T) {
	if got := fxNum(0.0384); got != "0,0384" {
		t.Errorf("мелкий курс напечатан как %q", got)
	}
	if got := fxNum(1234.5); got != "1\u202f234,50" {
		t.Errorf("крупный курс напечатан как %q", got)
	}
	if got := fxDelta(2.5, 456.0); !strings.HasPrefix(got, "+") {
		t.Errorf("рост напечатан без плюса: %q", got)
	}
	// Точность изменения задаёт величина курса, а не самого изменения: иначе в
	// одном столбце встают «+2,25» и «+0,5900», и колонка перестаёт читаться.
	if got := fxDelta(0.59, 456.0); got != "+0,59" {
		t.Errorf("мелкое изменение крупного курса напечатано как %q", got)
	}
	if got := fxDelta(0.0031, 0.0384); got != "+0,0031" {
		t.Errorf("изменение мелкого курса напечатано как %q", got)
	}
	if got := fxPct(-1.9); got != "−1,90%" {
		t.Errorf("проценты напечатаны как %q", got)
	}
}

// seedFx кладёт короткую историю: месяц дневных курсов и одну глубокую точку.
func seedFx(app *testApp) {
	app.exec(`DELETE FROM fx_rates WHERE code IN ('USD','JPY')`)
	app.exec(`DELETE FROM fx_monthly WHERE code IN ('USD','JPY')`)
	app.exec(`INSERT INTO fx_rates (day, code, value, quant, name)
		SELECT g::date, 'USD', 450 + (g::date - CURRENT_DATE + 30), 1, 'ДОЛЛАР США'
		  FROM generate_series(CURRENT_DATE - 30, CURRENT_DATE, interval '1 day') g
		ON CONFLICT (day, code) DO NOTHING`)
	app.exec(`INSERT INTO fx_monthly (month, code, value) VALUES
		(DATE '1993-11-01', 'USD', 4.70),
		(DATE '1999-04-01', 'USD', 88.00)
		ON CONFLICT (month, code) DO NOTHING`)
}

// Страница курса публичная и открывается без входа: она затем и нужна, чтобы
// на неё приходили посторонние.
func TestTheRatesPageIsOpenAndShowsTheRate(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	seedFx(app)

	w := app.do(http.MethodGet, "/rates?c=USD&p=month", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница курса отдала %d без входа", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, T(LangRU, "fx.title")) {
		t.Error("на странице нет её собственного заголовка")
	}
	if strings.Contains(body, T(LangRU, "fx.empty")) {
		t.Fatal("страница говорит, что данных нет, хотя они засеяны")
	}
	// Линия рисуется на сервере: без неё останутся одни числа.
	if !strings.Contains(body, "fx-svg__line") {
		t.Error("на странице нет линии курса")
	}
	// Источник обязан быть назван на самой странице.
	if !strings.Contains(body, T(LangRU, "fx.src_bis")) {
		t.Error("страница не называет источник глубокой истории")
	}
}

// Неизвестная валюта в адресе не должна ронять страницу или показывать пустоту:
// её место занимает доллар.
func TestAnUnknownCurrencyFallsBackToTheDollar(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	seedFx(app)

	w := app.do(http.MethodGet, "/rates?c=ZZZ&p=year", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("неизвестная валюта уронила страницу: %d", w.Code)
	}
	if strings.Contains(w.Body.String(), T(LangRU, "fx.empty")) {
		t.Error("вместо подстановки доллара страница осталась пустой")
	}
}

// Весь период должен начинаться с ноября 1993 года — с месяца, когда тенге
// появился. Ради этого глубокий архив и заводился.
func TestTheWholePeriodReachesTheBirthOfTheTenge(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	seedFx(app)

	pts, err := app.module().fx.SeriesMonthly(context.Background(), "USD", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(pts) < 3 {
		t.Fatalf("месячный ряд собрал %d точек", len(pts))
	}
	if got := pts[0].Day.Format("2006-01"); got != "1993-11" {
		t.Errorf("ряд начинается с %s, а не с рождения тенге", got)
	}
	// Хвост, до которого глубокий источник не дошёл, берётся из дневного
	// архива — иначе страница показывала бы курс двухмесячной давности.
	last := pts[len(pts)-1].Day
	if last.Before(monthStart(time.Now().UTC().AddDate(0, -1, 0))) {
		t.Errorf("ряд обрывается на %s — свежий хвост не подставлен", last.Format("2006-01"))
	}
}

// Догрузка не должна переспрашивать банк про дни, о которых уже спрашивала:
// пять лет выходных — это тысяча запросов ни за чем при каждом перезапуске.
func TestAProbedDayIsNeverAskedTwice(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	ctx := context.Background()
	fx := app.module().fx

	// Окно готовим сами: тест, который берёт первый попавшийся неопрошенный
	// день, зависит от того, что оставила в базе предыдущая работа.
	app.exec(`INSERT INTO fx_probed (day, found)
	          SELECT g::date, 0 FROM generate_series(CURRENT_DATE - 5, CURRENT_DATE, interval '1 day') g
	          ON CONFLICT (day) DO NOTHING`)
	app.exec(`DELETE FROM fx_probed WHERE day = CURRENT_DATE - 3`)

	floor := time.Now().UTC().AddDate(0, 0, -5)
	day, ok, err := fx.NextToProbe(ctx, floor)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("единственный неопрошенный день окна не найден")
	}
	if got := day.Format("2006-01-02"); got != time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02") {
		t.Fatalf("к опросу предложен %s, а не единственный пропуск", got)
	}

	if err := fx.MarkProbed(ctx, day, 0); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := fx.NextToProbe(ctx, floor); err != nil {
		t.Fatal(err)
	} else if ok {
		t.Error("после отметки догрузка снова нашла, что опрашивать в закрытом окне")
	}
}

// Длинная цепочка пустых дней — это край чужого архива, и упираться в него
// бесконечно незачем. Праздники дают до девяти пустых дней подряд, поэтому
// порог должен быть заметно выше.
func TestTheEndOfTheSourcesWindowIsRecognised(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	ctx := context.Background()
	fx := app.module().fx

	app.exec(`DELETE FROM fx_probed WHERE day BETWEEN DATE '2019-01-01' AND DATE '2019-03-01'`)
	app.exec(`INSERT INTO fx_probed (day, found)
	          SELECT g::date, 0 FROM generate_series(DATE '2019-01-01', DATE '2019-02-25', interval '1 day') g`)
	app.exec(`UPDATE fx_probed SET found = 40 WHERE day = DATE '2019-01-01'`)

	run, err := fx.EmptyRunBelow(ctx, time.Date(2019, 2, 25, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if run < fxEmptyRun {
		t.Errorf("подряд пустых дней насчитано %d — край окна не распознан", run)
	}
	// Строка с курсами обрывает цепочку: ниже неё считать нечего.
	if run > 55 {
		t.Errorf("цепочка в %d дней перешагнула день с курсами", run)
	}
}

// Названия банк пишет прописными, а аббревиатуры в них нельзя опускать в
// строчные: «ДОЛЛАР США» должен стать «Доллар США», а не «Доллар сша».
func TestTheBanksShoutingNamesAreCalmedDown(t *testing.T) {
	for in, want := range map[string]string{
		"ДОЛЛАР США":       "Доллар США",
		"ЕВРО":             "Евро",
		"РОССИЙСКИЙ РУБЛЬ": "Российский рубль",
		"ДИРХАМ ОАЭ":       "Дирхам ОАЭ",
		"СДР":              "СДР",
		"":                 "",
	} {
		if got := fxTidyName(in); got != want {
			t.Errorf("%q превратилось в %q, а ожидалось %q", in, got, want)
		}
	}
}

// У курса каждой валюты свой адрес. Свернуть их в один значит, что человек,
// ищущий курс рубля, нашу страницу курса рубля не найдёт: поисковик знает
// только тот адрес, который сайт объявил своим.
func TestEachCurrencyIsItsOwnPage(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	seedFx(app)
	app.exec(`INSERT INTO fx_rates (day, code, value, quant, name)
		SELECT g::date, 'EUR', 520, 1, 'ЕВРО'
		  FROM generate_series(CURRENT_DATE - 30, CURRENT_DATE, interval '1 day') g
		ON CONFLICT (day, code) DO NOTHING`)

	w := app.do(http.MethodGet, "/rates?c=EUR&p=year", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница евро отдала %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `href="/rates?c=EUR&amp;lang=ru"`) {
		t.Error("canonical евро не указывает на собственный адрес валюты")
	}
	// Период в canonical не входит: это тот же курс на другую глубину.
	if strings.Contains(body, "c=EUR&amp;p=year&amp;lang=ru") {
		t.Error("период попал в canonical и размножил страницу на четыре копии")
	}
	if !strings.Contains(body, "EUR — курс к тенге") {
		t.Error("у страницы евро общий заголовок вместо собственного")
	}
	if !strings.Contains(body, "Евро (EUR)") {
		t.Error("в описании страницы нет названия валюты")
	}

	// Доллар — это и есть /rates, второй адрес ему не нужен.
	w = app.do(http.MethodGet, "/rates?c=USD", nil)
	if b := w.Body.String(); strings.Contains(b, `href="/rates?c=USD&amp;lang=ru"`) {
		t.Error("доллар получил отдельный адрес, хотя он и есть страница курсов")
	}
}

// Страница, которой нет в карте сайта, для поисковика не существует.
func TestEveryCurrencyIsOfferedInTheSitemap(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	seedFx(app)
	app.exec(`INSERT INTO fx_rates (day, code, value, quant, name)
		SELECT g::date, 'EUR', 520, 1, 'ЕВРО'
		  FROM generate_series(CURRENT_DATE - 3, CURRENT_DATE, interval '1 day') g
		ON CONFLICT (day, code) DO NOTHING`)

	w := app.do(http.MethodGet, "/sitemap.xml", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("карта сайта отдала %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "/rates?c=EUR&amp;lang=ru") {
		t.Error("курса евро нет в карте сайта")
	}
	if !strings.Contains(body, "/rates?lang=ru") {
		t.Error("самой страницы курсов нет в карте сайта")
	}
	// Доллар не должен предлагаться дважды — своим адресом и как /rates.
	if strings.Contains(body, "/rates?c=USD") {
		t.Error("доллар предложен вторым адресом — это дубль страницы курсов")
	}
}

// Постоянные страницы должны уходить в IndexNow на всех трёх языках и ровно
// теми адресами, которые сайт объявляет своими: заявка на чужой адрес
// поисковику ничего не даёт.
func TestPermanentPagesAreAnnouncedInEveryLanguage(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	urls := app.module().pageURLs([]string{"/rates", "/analytics"})
	if len(urls) != 2*len(Langs) {
		t.Fatalf("адресов %d, а страниц две на %d языка", len(urls), len(Langs))
	}
	for _, want := range []string{"/rates?lang=ru", "/rates?lang=kz", "/rates?lang=en", "/analytics?lang=ru"} {
		found := false
		for _, u := range urls {
			if strings.HasSuffix(u, want) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("в заявке нет адреса %s", want)
		}
	}
	for _, u := range urls {
		if !strings.HasPrefix(u, "http") {
			t.Errorf("адрес %q не абсолютный — IndexNow такой отвергнет", u)
		}
	}
}
