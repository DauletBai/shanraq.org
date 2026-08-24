package articles

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

// An empty answer from the bank is not a breakage. For weekends and for dates
// outside its window a document with no currencies arrives, and treating that as
// an error would have the backfill trip over every Saturday for ever.
func TestAnEmptyDayFromTheBankIsNotAnError(t *testing.T) {
	rates, err := parseFxRates([]byte(`<?xml version="1.0"?><rates></rates>`), time.Now())
	if err != nil {
		t.Fatalf("пустой день признан ошибкой: %v", err)
	}
	if len(rates) != 0 {
		t.Fatalf("из пустого документа взялось %d курсов", len(rates))
	}
}

// The bank quotes the yen in hundreds, and that multiple has to be kept: without
// it the yen's rate comes out a hundred times larger than it is.
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

// In the BIS export every currency is measured in dollars. Tenge per unit of a
// currency is the quotient of two of its rows, and computing it any other way makes
// the deep history lie.
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
	// The dollar is the series' base, so its rate against the tenge is taken
	// directly.
	if v := got["USD@2026-06"]; v < 480.7 || v > 480.8 {
		t.Errorf("доллар посчитан как %.2f вместо 480.72", v)
	}
	if v := got["EUR@2026-06"]; v < 547 || v > 549 {
		t.Errorf("евро посчитано как %.2f, а 480.72/0.8777 это около 547.7", v)
	}
	// The dollar covers the tenge's birth month; the euro did not exist yet.
	if v := got["USD@1993-11"]; v < 4.69 || v > 4.71 {
		t.Errorf("ноябрь 1993 потерян: %.2f", v)
	}
	if _, ok := got["EUR@1993-11"]; ok {
		t.Error("евро появилось в 1993 году, хотя в выгрузке там пусто")
	}
}

// A month with no tenge rate cannot be converted into anything: there is nothing
// to divide by. Such a month must drop out whole rather than become a zero.
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

// Thinning a long series must not lose its ends: they are what the page computes
// the change over the period from.
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

// Small currencies like the Uzbek sum are worth less than one. Two decimal places
// would turn their whole series into identical zeros.
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
	// A change's precision is set by the size of the rate, not of the change itself:
	// otherwise "+2,25" and "+0,5900" stand in one column and it stops reading.
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

// seedFx lays down a short history: a month of daily rates and one deep point.
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

// The rates page is public and opens without a login: its whole purpose is that
// strangers arrive on it.
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
	// The line is drawn on the server: without it only the numbers remain.
	if !strings.Contains(body, "fx-svg__line") {
		t.Error("на странице нет линии курса")
	}
	// The source has to be named on the page itself.
	if !strings.Contains(body, T(LangRU, "fx.src_bis")) {
		t.Error("страница не называет источник глубокой истории")
	}
}

// An unknown currency in the address must not break the page or show an emptiness:
// the dollar takes its place.
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

// "All time" has to begin in November 1993 — the month the tenge appeared. That is
// what the deep archive was built for.
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
	// The tail the deep source has not reached comes from the daily archive —
	// otherwise the page would show a rate two months old.
	last := pts[len(pts)-1].Day
	if last.Before(monthStart(time.Now().UTC().AddDate(0, -1, 0))) {
		t.Errorf("ряд обрывается на %s — свежий хвост не подставлен", last.Format("2006-01"))
	}
}

// The backfill must not ask the bank again about days it has already asked about:
// five years of weekends is a thousand pointless requests on every restart.
func TestAProbedDayIsNeverAskedTwice(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	ctx := context.Background()
	fx := app.module().fx

	// The window is staged here: a test that takes whichever unprobed day comes
	// first depends on what an earlier run left in the database.
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

// A long run of empty days is the edge of somebody else's archive, and there is no
// point pushing against it for ever. Holidays give up to nine empty days in a row,
// so the threshold has to sit well above that.
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
	// A row holding rates breaks the run: below it there is nothing to count.
	if run > 55 {
		t.Errorf("цепочка в %d дней перешагнула день с курсами", run)
	}
}

// The bank writes the names in capitals, and abbreviations inside them must not be
// lower-cased: "ДОЛЛАР США" has to become "Доллар США", not "Доллар сша".
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

// Each currency's rate has its own address. Folding them into one means that
// someone searching for the rouble rate will not find our rouble rate page: a
// search engine only knows the address the site declared as its own.
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
	// The period stays out of the canonical: it is the same rate at another depth.
	if strings.Contains(body, "c=EUR&amp;p=year&amp;lang=ru") {
		t.Error("период попал в canonical и размножил страницу на четыре копии")
	}
	if !strings.Contains(body, "EUR — курс к тенге") {
		t.Error("у страницы евро общий заголовок вместо собственного")
	}
	if !strings.Contains(body, "Евро (EUR)") {
		t.Error("в описании страницы нет названия валюты")
	}

	// The dollar is /rates, and it needs no second address.
	w = app.do(http.MethodGet, "/rates?c=USD", nil)
	if b := w.Body.String(); strings.Contains(b, `href="/rates?c=USD&amp;lang=ru"`) {
		t.Error("доллар получил отдельный адрес, хотя он и есть страница курсов")
	}
}

// A page missing from the sitemap does not exist as far as a search engine is
// concerned.
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
	// The dollar must not be offered twice — once by its own address and once as
	// /rates.
	if strings.Contains(body, "/rates?c=USD") {
		t.Error("доллар предложен вторым адресом — это дубль страницы курсов")
	}
}

// The standing pages have to go to IndexNow in all three languages and by exactly
// the addresses the site declares as its own: a submission for someone else's
// address gives a search engine nothing.
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
