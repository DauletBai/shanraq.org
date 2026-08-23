package articles

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Страница разбора курса валют.
//
// Одна валюта, один период, одна линия. Всё остальное — подписи к ней: где
// линия была выше и ниже всего, на сколько сдвинулась за период, сколько это в
// процентах. Смысл страницы в том, чтобы человек, знающий только слово «курс»,
// за один взгляд увидел форму кривой, а тот, кто пришёл считать, нашёл рядом
// числа и не считал их сам.

// fxPeriods — периоды разбора. Ключ приходит в адресе, поэтому он короткий и
// не меняется: на такие ссылки ссылаются.
var fxPeriods = []string{"month", "year", "five", "all"}

// fxDefaultCode — валюта, с которой открывается страница.
const fxDefaultCode = "USD"

// fxMaxPoints — сколько точек имеет смысл рисовать. Пять лет дневных курсов
// это около тысячи трёхсот засечек на линию шириной в тысячу единиц: соседние
// точки встают плотнее, чем толщина пера, и лишние из них не добавляют ни
// одного различимого изгиба, зато утяжеляют страницу.
const fxMaxPoints = 480

// FxTick — засечка оси значений.
type FxTick struct {
	Label string
	Pos   float64 // проценты сверху вниз
}

// FxAxis — засечка оси времени. Позиция в процентах, а не в номере точки:
// месяцы разной длины, и равномерная раскладка развела бы подпись и точку,
// к которой она относится.
type FxAxis struct {
	Label string
	Pos   float64
	Grid  bool
}

// FxChart — готовая к отрисовке линия в системе координат картинки.
type FxChart struct {
	Line   string
	Area   string
	Y      []FxTick
	X      []FxAxis
	Width  int
	Height int
	Points int
}

// FxYear — строка годовой таблицы.
type FxYear struct {
	Year  int
	Close string
	Diff  string
	Pct   string
	Up    bool
	Down  bool
}

// FxPage — данные страницы.
type FxPage struct {
	Base
	Desc       string
	Currencies []FxCurrency
	Code       string
	Name       string
	Quant      int
	Period     string
	Chart      FxChart
	HasData    bool
	Last       string
	LastDay    string
	Diff       string
	Pct        string
	Up         bool
	Down       bool
	Min        string
	MinDay     string
	Max        string
	MaxDay     string
	Avg        string
	Spread     string
	Since      string
	Years      []FxYear
	Monthly    bool
	DeepSince  string
	DailySince string
}

// handleRates отдаёт страницу разбора курса.
func (m *Module) handleRates(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	ctx := r.Context()

	page := FxPage{Period: fxPeriod(r.URL.Query().Get("p"))}
	if m.fx != nil {
		page.Currencies, _ = m.fx.Currencies(ctx)
	}
	page.Code = fxCode(r.URL.Query().Get("c"), page.Currencies)
	page.Quant = 1
	for _, c := range page.Currencies {
		if c.Code == page.Code {
			page.Name, page.Quant = c.Name, c.Quant
			page.Since = fxDate(c.Since)
			break
		}
	}
	if page.Name == "" {
		page.Name = page.Code
	}

	if m.fx != nil {
		m.fillRates(ctx, &page, lang)
	}

	page.Base = m.base(r, T(lang, "fx.title"), lang)
	page.Desc = T(lang, "fx.desc")
	m.render(w, "rates", page)
}

// fillRates собирает ряд и все числа к нему.
func (m *Module) fillRates(ctx context.Context, page *FxPage, lang string) {
	var (
		pts []FxPoint
		err error
	)
	now := time.Now().UTC()
	switch page.Period {
	case "all":
		page.Monthly = true
		pts, err = m.fx.SeriesMonthly(ctx, page.Code, time.Time{})
	case "five":
		page.Monthly = true
		pts, err = m.fx.SeriesMonthly(ctx, page.Code, monthStart(now.AddDate(-5, 0, 0)))
	case "year":
		pts, err = m.fx.Series(ctx, page.Code, now.AddDate(-1, 0, 0))
	default:
		pts, err = m.fx.Series(ctx, page.Code, now.AddDate(0, 0, -30))
	}
	if err != nil || len(pts) < 2 {
		return
	}

	// Приводим к привычной кратности: банк объявляет иены сотнями, и человек
	// ищет глазами именно то число, которое видел в строке курсов.
	q := float64(page.Quant)
	for i := range pts {
		pts[i].Value *= q
	}

	page.HasData = true
	page.Chart = fxBuildChart(pts, page.Period, lang)

	first, last := pts[0], pts[len(pts)-1]
	lo, hi, sum := pts[0], pts[0], 0.0
	for _, p := range pts {
		if p.Value < lo.Value {
			lo = p
		}
		if p.Value > hi.Value {
			hi = p
		}
		sum += p.Value
	}
	diff := last.Value - first.Value

	page.Last, page.LastDay = fxNum(last.Value), fxDate(last.Day)
	page.Diff = fxDelta(diff, last.Value)
	if first.Value > 0 {
		page.Pct = fxPct(diff / first.Value * 100)
	}
	page.Up, page.Down = diff > 0, diff < 0
	page.Min, page.MinDay = fxNum(lo.Value), fxDate(lo.Day)
	page.Max, page.MaxDay = fxNum(hi.Value), fxDate(hi.Day)
	page.Avg = fxNum(sum / float64(len(pts)))
	page.Spread = fxNum(hi.Value - lo.Value)
	page.Years = fxYears(pts)
}

// fxBuildChart переводит ряд в координаты картинки.
func fxBuildChart(pts []FxPoint, period, lang string) FxChart {
	const w, h = 1000.0, 360.0
	pts = fxThin(pts, fxMaxPoints)

	lo, hi := pts[0].Value, pts[0].Value
	for _, p := range pts {
		lo = math.Min(lo, p.Value)
		hi = math.Max(hi, p.Value)
	}
	lo, hi, step := fxScale(lo, hi)

	x := func(i int) float64 { return w * float64(i) / float64(len(pts)-1) }
	y := func(v float64) float64 { return h - h*(v-lo)/(hi-lo) }

	var line strings.Builder
	for i, p := range pts {
		if i > 0 {
			line.WriteByte(' ')
		}
		if i == 0 {
			line.WriteByte('M')
		} else {
			line.WriteByte('L')
		}
		fmt.Fprintf(&line, "%.1f %.1f", x(i), y(p.Value))
	}
	area := line.String() + fmt.Sprintf(" L%.1f %.1f L0 %.1f Z", w, h, h)

	ticks := []FxTick{}
	for v := hi; v >= lo-step/2; v -= step {
		ticks = append(ticks, FxTick{Label: fxAxisNum(v, step), Pos: (hi - v) / (hi - lo) * 100})
	}

	return FxChart{
		Line: line.String(), Area: area, Y: ticks, X: fxAxis(pts, period, lang),
		Width: int(w), Height: int(h), Points: len(pts),
	}
}

// fxAxis размечает ось времени по естественным границам выбранного периода:
// месяц — по дням, год — по месяцам, пять лет и вся история — по годам.
// Равномерные засечки через каждые N точек давали бы подписи вроде «14.03» и
// «27.07», по которым нельзя ни найти нужное место, ни сравнить два года.
func fxAxis(pts []FxPoint, period, lang string) []FxAxis {
	if len(pts) < 2 {
		return nil
	}
	pos := func(i int) float64 { return 100 * float64(i) / float64(len(pts)-1) }
	out := []FxAxis{}

	switch period {
	case "month":
		// Каждый день своим числом: месяц читают по числам, а не по датам.
		for i, p := range pts {
			out = append(out, FxAxis{Label: strconv.Itoa(p.Day.Day()), Pos: pos(i),
				Grid: p.Day.Day() == 1})
		}
	case "year":
		for i, p := range pts {
			if i > 0 && p.Day.Month() == pts[i-1].Day.Month() {
				continue
			}
			lab := fxMonthShort(p.Day.Month(), lang)
			if p.Day.Month() == time.January {
				lab = strconv.Itoa(p.Day.Year())
			}
			out = append(out, FxAxis{Label: lab, Pos: pos(i), Grid: p.Day.Month() == time.January})
		}
	default:
		// Пять лет и вся история: подписываем годы, но так, чтобы их не стало
		// три десятка — на тридцати трёх годах подписи слиплись бы в полосу.
		step := 1
		if years := pts[len(pts)-1].Day.Year() - pts[0].Day.Year(); years > 12 {
			step = 5
		} else if years > 6 {
			step = 2
		}
		for i, p := range pts {
			if i > 0 && p.Day.Year() == pts[i-1].Day.Year() {
				continue
			}
			if p.Day.Year()%step != 0 {
				continue
			}
			out = append(out, FxAxis{Label: strconv.Itoa(p.Day.Year()), Pos: pos(i), Grid: true})
		}
	}
	return fxSpaceOut(out)
}

// fxMinGap — наименьший просвет между подписями оси, в процентах ширины.
const fxMinGap = 2.6

// fxSpaceOut убирает подписи, налезающие на соседнюю справа. Ряд почти всегда
// начинается посреди месяца, и «авг» первой точки садится вплотную к «сен»
// первого сентября, склеиваясь с ним в нечитаемое пятно. Идём с конца, потому
// что последняя подпись — сегодняшняя дата, и она нужнее любой предыдущей.
func fxSpaceOut(a []FxAxis) []FxAxis {
	if len(a) < 2 {
		return a
	}
	last := 101.0
	for i := len(a) - 1; i >= 0; i-- {
		if last-a[i].Pos < fxMinGap {
			// Гасим только подпись: засечка на сетке размечает график и без
			// неё, а убери мы её вместе с текстом — пропала бы и линия месяца.
			a[i].Label = ""
			continue
		}
		last = a[i].Pos
	}
	return a
}

// fxMonthShort — короткое название месяца на языке читателя.
func fxMonthShort(m time.Month, lang string) string {
	names := map[string][12]string{
		LangKZ: {"қаң", "ақп", "нау", "сәу", "мам", "мау", "шіл", "там", "қыр", "қаз", "қар", "жел"},
		LangRU: {"янв", "фев", "мар", "апр", "май", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"},
		LangEN: {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	}
	row, ok := names[lang]
	if !ok {
		row = names[LangRU]
	}
	return row[int(m)-1]
}

// fxScale подбирает границы и шаг оси значений так, чтобы подписи были
// круглыми числами. Ось, размеченная делениями вроде «472,91» и «478,79»,
// читается хуже, чем никакая: по ней нельзя ни прикинуть значение точки, ни
// сравнить два графика между собой.
func fxScale(lo, hi float64) (float64, float64, float64) {
	if hi <= lo {
		hi = lo + math.Max(math.Abs(lo)*0.01, 0.01)
	}
	step := fxNiceStep((hi - lo) / 4)
	lo = math.Floor(lo/step) * step
	hi = math.Ceil(hi/step) * step
	// Курс не бывает отрицательным, и место под минус на оси — потерянная
	// высота графика.
	if lo < 0 {
		lo = 0
	}
	return lo, hi, step
}

// fxNiceStep округляет шаг вверх до единицы, двойки, четвертушки или половины
// от ближайшей степени десяти.
func fxNiceStep(raw float64) float64 {
	if raw <= 0 {
		return 1
	}
	mag := math.Pow(10, math.Floor(math.Log10(raw)))
	for _, m := range []float64{1, 2, 2.5, 5} {
		if raw <= m*mag {
			return m * mag
		}
	}
	return 10 * mag
}

// fxThin прореживает ряд до max точек, сохраняя первую и последнюю.
func fxThin(pts []FxPoint, max int) []FxPoint {
	if len(pts) <= max {
		return pts
	}
	out := make([]FxPoint, 0, max)
	for i := 0; i < max-1; i++ {
		out = append(out, pts[i*(len(pts)-1)/(max-1)])
	}
	return append(out, pts[len(pts)-1])
}

// fxYears сводит ряд по годам: чем закончился год и на сколько он сдвинул
// курс. Длинную линию глазом на годы не разложить, а вопрос «какой год был
// плохим» задают именно так.
func fxYears(pts []FxPoint) []FxYear {
	last := map[int]FxPoint{}
	years := []int{}
	for _, p := range pts {
		y := p.Day.Year()
		if _, ok := last[y]; !ok {
			years = append(years, y)
		}
		last[y] = p
	}
	if len(years) < 2 {
		return nil
	}
	out := make([]FxYear, 0, len(years))
	for i := len(years) - 1; i >= 0; i-- {
		y := years[i]
		row := FxYear{Year: y, Close: fxNum(last[y].Value)}
		if i > 0 {
			prev := last[years[i-1]].Value
			d := last[y].Value - prev
			row.Diff, row.Up, row.Down = fxDelta(d, last[y].Value), d > 0, d < 0
			if prev > 0 {
				row.Pct = fxPct(d / prev * 100)
			}
		}
		out = append(out, row)
	}
	return out
}

// fxPeriod приводит период из адреса к известному.
func fxPeriod(s string) string {
	for _, p := range fxPeriods {
		if s == p {
			return p
		}
	}
	return "month"
}

// fxCode приводит валюту из адреса к той, что у нас есть.
func fxCode(s string, list []FxCurrency) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	for _, c := range list {
		if c.Code == s {
			return s
		}
	}
	// Валюты в списке идут от ходовых к прочим, и первая — это доллар, если он
	// у нас есть. Опечатка в адресе не должна оставлять человека с пустотой.
	if len(list) > 0 {
		return list[0].Code
	}
	return fxDefaultCode
}

// fxNum печатает курс: два знака после запятой для обычных чисел и четыре для
// мелких, где два знака превратили бы весь ряд в одинаковые нули.
func fxNum(v float64) string { return fxFormat(v, fxDigits(v)) }

// fxAxisNum печатает деление оси с той точностью, которую задаёт её шаг. Шаг в
// десять тенге не нуждается в копейках: «480,00» вместо «480» — это два знака
// шума в каждой подписи.
func fxAxisNum(v, step float64) string {
	digits := 0
	if step < 1 {
		digits = int(math.Ceil(-math.Log10(step)))
	}
	return fxFormat(v, digits)
}

// fxFormat печатает число с запятой в дробной части и пробелами в разрядах.
func fxFormat(v float64, digits int) string {
	s := strconv.FormatFloat(v, 'f', digits, 64)
	whole, frac, _ := strings.Cut(s, ".")
	neg := strings.HasPrefix(whole, "-")
	whole = strings.TrimPrefix(whole, "-")
	// Разряды разделяет узкий неразрывный пробел: он не даёт числу
	// переноситься на новую строку посреди тысяч и не расталкивает цифры
	// так широко, как обычный пробел.
	var b strings.Builder
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteRune('\u202f')
		}
		b.WriteRune(r)
	}
	out := b.String()
	if frac != "" {
		out += "," + frac
	}
	if neg {
		out = "−" + out
	}
	return out
}

// fxDelta печатает изменение со знаком, с той же точностью, что и сам курс:
// в колонке, где рядом стоят «+2,25» и «+0,5900», глаз спотыкается о разную
// длину дробной части и перестаёт сравнивать числа. Точность задаёт ref —
// величина курса, а не величина самого изменения.
//
// Минус здесь настоящий, а не дефис: в столбце чисел дефис по высоте не
// совпадает с плюсом, и ряд рябит.
func fxDelta(v, ref float64) string {
	return fxSign(v, fxFormat(v, fxDigits(ref)))
}

// fxPct печатает проценты — всегда с двумя знаками, каким бы ни был курс.
func fxPct(v float64) string {
	return fxSign(v, fxFormat(v, 2)) + "%"
}

// fxSign приписывает плюс: минус уже стоит в самом числе.
func fxSign(v float64, s string) string {
	if v > 0 {
		return "+" + s
	}
	return s
}

// fxDigits — сколько знаков после запятой уместно для курса такого размера.
// У валюты дешевле тенге два знака превратили бы весь ряд в одинаковые нули.
func fxDigits(v float64) int {
	if math.Abs(v) < 1 {
		return 4
	}
	return 2
}

// fxDate печатает дату коротко. Формат один на все три языка: числовая дата
// читается одинаково и не требует названий месяцев в трёх падежах.
func fxDate(d time.Time) string {
	if d.IsZero() {
		return ""
	}
	return d.Format("02.01.2006")
}
