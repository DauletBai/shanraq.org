package articles

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

// Раздел «Как формируется курс тенге и инфляция».
//
// Страница не пишется — она считается. Нацбанк раз в месяц публикует, сколько
// в стране тенге и сколько долларов в золотовалютных резервах; Всемирный банк
// раз в год — инфляцию. Всё, что здесь делается, это кладёт одно рядом с
// другим: количество денег, количество валюты, за которую эти деньги можно
// обменять, и курс, который из этого соотношения выходит.
//
// Никаких утверждений о намерениях: показаны ряды и их сопоставление, а вывод
// из них читатель делает сам. Цифра, за которой стоит источник, убеждает
// сильнее любого эпитета.

// MacroChart — две линии в одной системе координат.
type MacroChart struct {
	A      string
	B      string
	Area   string
	Y      []FxTick
	X      []FxAxis
	Width  int
	Height int
	// Log отмечает логарифмическую шкалу — читателю об этом надо сказать.
	Log bool
}

// MacroYear — строка годовой сводки.
type MacroYear struct {
	Year  int
	CPI   string
	Money string
	Rate  string
	Fund  string
	// FundDown отмечает год, когда Национальный фонд уменьшился.
	FundDown bool
}

// MacroErode — во что превратилась тысяча тенге того года.
type MacroErode struct {
	Year int
	// Kept — сколько сегодняшних тенге стоит та тысяча по покупательной
	// способности.
	Kept string
	// Then и Now — сколько долларов давали за тысячу тенге тогда и сейчас.
	Then string
	Now  string
}

// MacroBlock — весь раздел целиком.
type MacroBlock struct {
	Has bool

	M3      FxChart
	M3Last  string
	M3Month string
	M3Year  string
	M3Up    bool

	Cmp      MacroChart
	CmpFrom  string
	MoneyMul string
	RateMul  string

	Cover      FxChart
	CoverLast  string
	CoverMonth string
	RateAt     string

	Res      MacroChart
	ResLast  string
	FundLast string
	ResMonth string

	CPI     FxChart
	CPILast string
	CPIYear string

	Years []MacroYear
	Erode []MacroErode
	// NowDollars — сколько долларов дают за тысячу тенге сегодня.
	NowDollars string
}

// buildMacro собирает раздел. Любая недостающая часть просто не показывается:
// половина разбора полезнее пустой страницы.
func (m *Module) buildMacro(ctx context.Context, lang string) MacroBlock {
	var b MacroBlock
	if m.macro == nil || m.fx == nil {
		return b
	}

	m3, _ := m.macro.Series(ctx, MacroM3)
	res, _ := m.macro.Series(ctx, MacroReserves)
	fund, _ := m.macro.Series(ctx, MacroFund)
	cpi, _ := m.macro.Series(ctx, MacroCPI)
	rate, _ := m.fx.SeriesMonthly(ctx, fxDefaultCode, time.Time{})
	if len(m3) < 24 {
		return b
	}
	b.Has = true

	// 1. Сколько тенге в стране.
	b.M3 = fxBuildChart(macroScaled(m3, 1e6), "all", lang) // триллионы тенге
	last := m3[len(m3)-1]
	b.M3Last = macroTenge(last.Value)
	b.M3Month = macroMonth(last.Period, lang)
	if prev, ok := macroAt(m3, last.Period.AddDate(-1, 0, 0)); ok && prev > 0 {
		d := (last.Value/prev - 1) * 100
		b.M3Year, b.M3Up = fxPct(d), d > 0
	}

	// 2. Денежная масса и курс доллара на одной шкале. Обе линии приведены к
	//    ста в первый общий месяц: у одной единица измерения — тенге, у другой
	//    тенге за доллар, и сравнивать их можно только в разах.
	if len(rate) > 24 {
		b.Cmp, b.CmpFrom = macroCompare(m3, rate, lang)
		b.MoneyMul = macroTimes(m3, lang)
		b.RateMul = macroTimes(rate, lang)
	}

	// 3. Сколько тенге приходится на доллар золотовалютных резервов. Это не
	//    курс и не прогноз: это отношение двух опубликованных чисел.
	if len(res) > 24 {
		cover := macroCover(m3, res)
		if len(cover) > 24 {
			b.Cover = fxBuildChart(cover, "all", lang)
			c := cover[len(cover)-1]
			b.CoverLast, b.CoverMonth = fxNum(c.Value), macroMonth(c.Day, lang)
			if v, ok := macroRateAt(rate, c.Day); ok {
				b.RateAt = fxNum(v)
			}
		}
	}

	// 4. Золотовалютные резервы и Национальный фонд.
	if len(res) > 24 && len(fund) > 24 {
		b.Res = macroTwoLines(macroScaled(res, 1e3), macroScaled(fund, 1e3), lang, false) // млрд долларов
		b.ResLast = macroDollars(res[len(res)-1].Value)
		b.FundLast = macroDollars(fund[len(fund)-1].Value)
		b.ResMonth = macroMonth(res[len(res)-1].Period, lang)
	}

	// 5. Инфляция по годам.
	if len(cpi) > 5 {
		b.CPI = fxBuildChart(macroPoints(cpi), "all", lang)
		c := cpi[len(cpi)-1]
		b.CPILast, b.CPIYear = fxPct(c.Value), fmt.Sprintf("%d", c.Period.Year())
	}

	b.Years = macroYears(cpi, m3, rate, fund)
	b.Erode, b.NowDollars = macroErosion(cpi, rate)
	return b
}

// macroErosion считает, во что превратилась тысяча тенге разных лет.
//
// Две меры сразу, потому что они отвечают на разные вопросы. Покупательная
// способность — сколько сегодняшних тенге нужно, чтобы купить то же самое:
// накопленная инфляция за годы. Доллары — сколько валюты давали за ту же
// тысячу тогда и дают сейчас. Первое про цены внутри страны, второе про то,
// чего стоят эти цены снаружи.
func macroErosion(cpi []MacroPoint, rate []FxPoint) ([]MacroErode, string) {
	if len(cpi) < 5 || len(rate) < 12 {
		return nil, ""
	}
	byYear := map[int]float64{}
	for _, c := range cpi {
		byYear[c.Period.Year()] = c.Value
	}
	last := rate[len(rate)-1]
	nowRate := last.Value
	if nowRate <= 0 {
		return nil, ""
	}
	rateAtYear := map[int]float64{}
	for _, p := range rate {
		if _, ok := rateAtYear[p.Day.Year()]; !ok {
			rateAtYear[p.Day.Year()] = p.Value
		}
	}

	out := []MacroErode{}
	for _, y := range []int{1995, 2000, 2005, 2010, 2015, 2020} {
		if y >= last.Day.Year() {
			continue
		}
		// Накопленный рост цен от начала года y до последнего года с данными.
		growth := 1.0
		ok := true
		for yy := y; yy < last.Day.Year(); yy++ {
			v, has := byYear[yy]
			if !has {
				ok = false
				break
			}
			growth *= 1 + v/100
		}
		r, hasRate := rateAtYear[y]
		if !ok || !hasRate || r <= 0 || growth <= 0 {
			continue
		}
		out = append(out, MacroErode{
			Year: y,
			Kept: fxFormat(1000/growth, 0),
			Then: fxFormat(1000/r, 2),
			Now:  fxFormat(1000/nowRate, 2),
		})
	}
	return out, fxFormat(1000/nowRate, 2)
}

// macroPoints переводит макроряд в точки графика.
func macroPoints(pts []MacroPoint) []FxPoint {
	out := make([]FxPoint, 0, len(pts))
	for _, p := range pts {
		out = append(out, FxPoint{Day: p.Period, Value: p.Value})
	}
	return out
}

// macroAt возвращает значение ряда за конкретный месяц.
func macroAt(pts []MacroPoint, when time.Time) (float64, bool) {
	for _, p := range pts {
		if p.Period.Equal(when) {
			return p.Value, true
		}
	}
	return 0, false
}

// macroRateAt возвращает курс за месяц точки, беря ближайший не позже него.
func macroRateAt(rate []FxPoint, when time.Time) (float64, bool) {
	var v float64
	var ok bool
	for _, p := range rate {
		if p.Day.After(when.AddDate(0, 1, 0)) {
			break
		}
		v, ok = p.Value, true
	}
	return v, ok
}

// macroCover считает, сколько тенге приходится на доллар золотовалютных
// резервов: денежная масса в миллионах тенге, делённая на резервы в миллионах
// долларов, даёт тенге за доллар.
func macroCover(m3 []MacroPoint, res []MacroPoint) []FxPoint {
	byMonth := make(map[string]float64, len(res))
	for _, r := range res {
		byMonth[r.Period.Format("2006-01")] = r.Value
	}
	out := make([]FxPoint, 0, len(m3))
	for _, p := range m3 {
		r, ok := byMonth[p.Period.Format("2006-01")]
		if !ok || r <= 0 {
			continue
		}
		out = append(out, FxPoint{Day: p.Period, Value: p.Value / r})
	}
	return out
}

// macroTimes — во сколько раз ряд вырос от первого значения к последнему,
// вместе со словом «раз» в нужной форме.
func macroTimes[T any](pts []T, lang string) string {
	first, last, ok := macroEnds(pts)
	if !ok || first <= 0 {
		return ""
	}
	return macroMul(last/first, lang)
}

// macroEnds достаёт первое и последнее значение ряда любого из наших видов.
func macroEnds[T any](pts []T) (float64, float64, bool) {
	if len(pts) < 2 {
		return 0, 0, false
	}
	val := func(x any) (float64, bool) {
		switch v := x.(type) {
		case MacroPoint:
			return v.Value, true
		case FxPoint:
			return v.Value, true
		}
		return 0, false
	}
	f, ok1 := val(any(pts[0]))
	l, ok2 := val(any(pts[len(pts)-1]))
	return f, l, ok1 && ok2
}

// macroMul печатает кратность вместе со словом: «6 776 раз», «2,4 раза».
//
// Слово приходится склонять: по-русски после 2, 3 и 4 — «раза», после
// остальных — «раз», а дробное число всегда берёт «раза». Без этого в тексте
// появляется «в 97,2 раз», и читатель спотыкается о фразу вместо того, чтобы
// смотреть на число.
func macroMul(x float64, lang string) string {
	whole := x >= 10
	n := fxFormat(x, 1)
	if whole {
		n = fxFormat(math.Round(x), 0)
	}
	switch lang {
	case LangKZ:
		return n + " есе"
	case LangEN:
		return n + " times"
	}
	if !whole {
		return n + " раза"
	}
	return n + " " + ruTimesWord(int64(math.Round(x)))
}

// ruTimesWord выбирает форму слова «раз» для целого числа.
func ruTimesWord(n int64) string {
	if n < 0 {
		n = -n
	}
	if t := n % 100; t >= 11 && t <= 14 {
		return "раз"
	}
	switch n % 10 {
	case 2, 3, 4:
		return "раза"
	}
	return "раз"
}

// macroCompare строит две линии, приведённые к ста в первый общий месяц.
func macroCompare(m3 []MacroPoint, rate []FxPoint, lang string) (MacroChart, string) {
	a := macroPoints(m3)
	from := a[0].Day
	if rate[0].Day.After(from) {
		from = rate[0].Day
	}
	a = macroSince(a, from)
	b := macroSince(rate, from)
	if len(a) < 2 || len(b) < 2 {
		return MacroChart{}, ""
	}
	return macroTwoLines(macroIndex(a), macroIndex(b), lang, true), macroMonthIn(from, lang)
}

// macroScaled переводит ряд в удобные единицы: миллионы тенге в триллионы,
// миллионы долларов в миллиарды. Ось с подписью «60 000 000» не читается и не
// помещается — а это те же 60 триллионов.
func macroScaled(pts []MacroPoint, div float64) []FxPoint {
	out := make([]FxPoint, 0, len(pts))
	for _, p := range pts {
		out = append(out, FxPoint{Day: p.Period, Value: p.Value / div})
	}
	return out
}

// macroSince отбрасывает всё раньше указанного месяца.
func macroSince(pts []FxPoint, from time.Time) []FxPoint {
	out := make([]FxPoint, 0, len(pts))
	for _, p := range pts {
		if p.Day.Before(from) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// macroIndex приводит ряд к ста в первой точке.
func macroIndex(pts []FxPoint) []FxPoint {
	if len(pts) == 0 || pts[0].Value == 0 {
		return pts
	}
	base := pts[0].Value
	out := make([]FxPoint, 0, len(pts))
	for _, p := range pts {
		out = append(out, FxPoint{Day: p.Day, Value: p.Value / base * 100})
	}
	return out
}

// macroTwoLines рисует два ряда в одной системе координат.
//
// log включает логарифмическую шкалу. Она нужна там, где ряды разошлись на
// порядки: денежная масса выросла в шесть с лишним тысяч раз, курс — в сотню,
// и на обычной шкале вторая линия ложится на ноль, будто курс не менялся
// вовсе. Логарифм показывает не разницу величин, а разницу темпов — то, ради
// чего два ряда и кладут рядом.
func macroTwoLines(a, b []FxPoint, lang string, log bool) MacroChart {
	const w, h = 1000.0, 360.0
	if len(a) < 2 || len(b) < 2 {
		return MacroChart{}
	}
	a, b = fxThin(a, fxMaxPoints), fxThin(b, fxMaxPoints)

	lo, hi := math.Inf(1), math.Inf(-1)
	for _, set := range [][]FxPoint{a, b} {
		for _, p := range set {
			if log && p.Value <= 0 {
				continue
			}
			lo, hi = math.Min(lo, p.Value), math.Max(hi, p.Value)
		}
	}
	if math.IsInf(lo, 1) || math.IsInf(hi, -1) {
		return MacroChart{}
	}

	var ticks []FxTick
	var at func(float64) float64
	if log {
		loE, hiE := math.Floor(math.Log10(lo)), math.Ceil(math.Log10(hi))
		at = func(v float64) float64 {
			if v <= 0 {
				return h
			}
			return h - h*(math.Log10(v)-loE)/(hiE-loE)
		}
		for e := hiE; e >= loE-0.5; e-- {
			v := math.Pow(10, e)
			ticks = append(ticks, FxTick{Label: fxFormat(v, 0), Pos: (hiE - e) / (hiE - loE) * 100})
		}
	} else {
		var step float64
		lo, hi, step = fxScale(lo, hi)
		at = func(v float64) float64 { return h - h*(v-lo)/(hi-lo) }
		for v := hi; v >= lo-step/2; v -= step {
			ticks = append(ticks, FxTick{Label: fxAxisNum(v, step), Pos: (hi - v) / (hi - lo) * 100})
		}
	}

	path := func(pts []FxPoint) string {
		var s strings.Builder
		for i, p := range pts {
			if i > 0 {
				s.WriteByte(' ')
			}
			if i == 0 {
				s.WriteByte('M')
			} else {
				s.WriteByte('L')
			}
			fmt.Fprintf(&s, "%.1f %.1f", w*float64(i)/float64(len(pts)-1), at(p.Value))
		}
		return s.String()
	}

	pa := path(a)
	// Заливку под первой линией на логарифмической шкале не рисуем: площадь
	// под логарифмом ничего не означает, а глаз читает её как объём.
	area := ""
	if !log {
		area = pa + fmt.Sprintf(" L%.1f %.1f L0 %.1f Z", w, h, h)
	}
	return MacroChart{
		A: pa, B: path(b),
		Area:  area,
		Y:     ticks,
		X:     fxAxis(a, "all", lang),
		Width: int(w), Height: int(h),
		Log: log,
	}
}

// macroYears сводит инфляцию, прирост денежной массы и изменение курса в одну
// таблицу по годам. Три числа в строке отвечают на вопрос, ради которого сюда
// приходят: что происходило с деньгами и что — с ценами.
func macroYears(cpi, m3 []MacroPoint, rate []FxPoint, fund []MacroPoint) []MacroYear {
	byYear := map[int]*MacroYear{}
	get := func(y int) *MacroYear {
		if _, ok := byYear[y]; !ok {
			byYear[y] = &MacroYear{Year: y}
		}
		return byYear[y]
	}
	for _, c := range cpi {
		get(c.Period.Year()).CPI = fxPct(c.Value)
	}
	for y, d := range macroYearChange(macroPoints(m3)) {
		get(y).Money = fxPct(d)
	}
	for y, d := range macroYearChange(rate) {
		get(y).Rate = fxPct(d)
	}
	// Национальный фонд — единственная строка, где падение важнее роста:
	// уменьшение значит, что из него взяли больше, чем в него положили.
	for y, d := range macroYearChange(macroPoints(fund)) {
		r := get(y)
		r.Fund, r.FundDown = fxPct(d), d < 0
	}

	out := make([]MacroYear, 0, len(byYear))
	for _, v := range byYear {
		if v.CPI == "" && v.Money == "" && v.Rate == "" && v.Fund == "" {
			continue
		}
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Year > out[j].Year })
	return out
}

// macroYearChange считает изменение ряда от декабря к декабрю, в процентах.
func macroYearChange(pts []FxPoint) map[int]float64 {
	last := map[int]float64{}
	years := []int{}
	for _, p := range pts {
		y := p.Day.Year()
		if _, ok := last[y]; !ok {
			years = append(years, y)
		}
		last[y] = p.Value
	}
	out := map[int]float64{}
	for i := 1; i < len(years); i++ {
		prev := last[years[i-1]]
		if prev <= 0 {
			continue
		}
		out[years[i]] = (last[years[i]]/prev - 1) * 100
	}
	return out
}

// macroTenge печатает сумму в тенге, приходящую в миллионах.
func macroTenge(millions float64) string {
	switch {
	case millions >= 1e6:
		return fxFormat(millions/1e6, 1) + " трлн"
	case millions >= 1e3:
		return fxFormat(millions/1e3, 1) + " млрд"
	default:
		return fxFormat(millions, 0) + " млн"
	}
}

// macroDollars печатает сумму в долларах, приходящую в миллионах.
func macroDollars(millions float64) string {
	if millions >= 1e3 {
		return fxFormat(millions/1e3, 1) + " млрд"
	}
	return fxFormat(millions, 0) + " млн"
}

// macroMonth печатает месяц словом и годом. В подписи под графиком уместно
// сокращение, а в предложении — нет: «на июл 2026» читается как опечатка.
func macroMonth(t time.Time, lang string) string {
	if t.IsZero() {
		return ""
	}
	return fxMonthFull(t.Month(), lang) + " " + fmt.Sprintf("%d", t.Year())
}

// macroMonthIn печатает месяц в той форме, которая нужна после предлога «в».
// По-русски это предложный падеж: «в январе», а не «в январь». Остальным двум
// языкам склонение здесь не требуется.
func macroMonthIn(t time.Time, lang string) string {
	if t.IsZero() {
		return ""
	}
	if lang != LangRU {
		return macroMonth(t, lang)
	}
	prep := [12]string{"январе", "феврале", "марте", "апреле", "мае", "июне",
		"июле", "августе", "сентябре", "октябре", "ноябре", "декабре"}
	return prep[int(t.Month())-1] + " " + fmt.Sprintf("%d", t.Year())
}

// fxMonthFull — полное название месяца на языке читателя.
func fxMonthFull(m time.Month, lang string) string {
	names := map[string][12]string{
		LangKZ: {"қаңтар", "ақпан", "наурыз", "сәуір", "мамыр", "маусым",
			"шілде", "тамыз", "қыркүйек", "қазан", "қараша", "желтоқсан"},
		LangRU: {"январь", "февраль", "март", "апрель", "май", "июнь",
			"июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь"},
		LangEN: {"January", "February", "March", "April", "May", "June",
			"July", "August", "September", "October", "November", "December"},
	}
	row, ok := names[lang]
	if !ok {
		row = names[LangRU]
	}
	return row[int(m)-1]
}

// Кеш раздела.
//
// Раздел одинаков для всех сорока восьми валют и всех периодов, а его данные
// меняются раз в месяц. Считать его на каждый заход значило бы гонять шесть
// запросов и строить пять графиков ради картинки, которая всё равно та же.

// macroTTL — сколько живёт собранный раздел.
const macroTTL = time.Hour

// macroCache хранит собранный раздел по языкам.
var macroCache = struct {
	mu   sync.Mutex
	at   time.Time
	byLg map[string]MacroBlock
}{byLg: map[string]MacroBlock{}}

// macroCached отдаёт раздел, пересобирая его не чаще раза в час.
func (m *Module) macroCached(ctx context.Context, lang string) MacroBlock {
	macroCache.mu.Lock()
	fresh := time.Since(macroCache.at) < macroTTL
	got, ok := macroCache.byLg[lang]
	macroCache.mu.Unlock()
	if fresh && ok {
		return got
	}

	built := m.buildMacro(ctx, lang)
	if !built.Has {
		// Пустой раздел не запоминаем: сразу после разворачивания ряды ещё
		// загружаются, и час держать «данных нет» значит час не показывать
		// раздел, который уже готов.
		return built
	}
	macroCache.mu.Lock()
	if !fresh {
		// Устаревший набор сбрасывается целиком: держать один язык из старого
		// часа рядом с другим из нового значит показать на одной странице два
		// разных месяца.
		macroCache.byLg = map[string]MacroBlock{}
		macroCache.at = time.Now()
	}
	macroCache.byLg[lang] = built
	macroCache.mu.Unlock()
	return built
}
