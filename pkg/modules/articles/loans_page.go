package articles

import (
	"encoding/json"
	"net/http"
)

// CalcPage is the loan calculator page.
//
// Programmes reach the browser as JSON in a data attribute, the way the rate
// charts hand over their series, so picking a type and a programme refills the
// form without a round trip and the totals follow every keystroke.
//
// The page still works with scripts off: the worked example is computed on the
// server and the table of programmes -- rates, terms, sources, dates -- is
// plain HTML. That table is also the part worth indexing.
type CalcPage struct {
	Base
	Desc     string
	Kinds    []string
	Programs []LoanProgram
	DataJSON string
	// MarketJSON maps a loan kind to the dearest rate on offer for it. That is
	// the yardstick a subsidised rate is measured against: the gap between them
	// is what the budget covers, and naming it is the point.
	MarketJSON string
	// BannerAds fills the shared real-estate sidebar. The template ranges over
	// it, and a page that omits the field aborts the render half way -- which is
	// how this was found.
	BannerAds []*Listing
	// LabelsJSON carries the handful of strings the script writes into the page.
	// Passing them in an attribute rather than interpolating them into inline
	// script keeps the template out of a JavaScript context it cannot parse.
	LabelsJSON string

	// Worked example: the choice a Kazakh buyer actually faces, a state
	// programme against a commercial one, computed server-side so the page
	// answers something before anyone touches a field.
	ExPrice, ExDown int64
	ExRate          float64
	ExMonths        int
	ExPlan          LoanPlan
	ExStateRate     float64
	ExStateMonths   int
	ExState         LoanPlan
}

// calcJSON is one programme as the form needs it: rate in percent rather than
// basis points, and nothing the form does not fill in.
type calcJSON struct {
	Code   string  `json:"code"`
	Kind   string  `json:"kind"`
	Name   string  `json:"name"`
	Lender string  `json:"lender"`
	Rate   float64 `json:"rate"`
	Months int     `json:"months"`
	Down   int     `json:"down"`
	Note   string  `json:"note"`
}

// The worked example: a flat near the middle of the regional market, the
// smallest down payment a programme allows, and the two rates that bracket the
// real choice.
const (
	calcExPrice       = 20_000_000
	calcExDown        = 4_000_000
	calcExRate        = 18
	calcExMonths      = 240
	calcExStateRate   = 7
	calcExStateMonths = 300
)

func (m *Module) handleCalculator(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)

	page := CalcPage{
		Kinds:         LoanKinds,
		ExPrice:       calcExPrice,
		ExDown:        calcExDown,
		ExRate:        calcExRate,
		ExMonths:      calcExMonths,
		ExStateRate:   calcExStateRate,
		ExStateMonths: calcExStateMonths,
	}
	page.ExPlan = LoanCalc(LoanInput{
		Kind: LoanMortgage, Price: page.ExPrice, Down: page.ExDown,
		RatePct: page.ExRate, Months: page.ExMonths, Scheme: SchemeAnnuity,
	})
	page.ExState = LoanCalc(LoanInput{
		Kind: LoanMortgage, Price: page.ExPrice, Down: page.ExDown,
		RatePct: page.ExStateRate, Months: page.ExStateMonths, Scheme: SchemeAnnuity,
	})

	if m.loans != nil {
		page.Programs, _ = m.loans.List(r.Context(), lang)
	}
	if m.listings != nil {
		page.BannerAds, _ = m.listings.BannerListings(r.Context(), 2)
	}

	rows := make([]calcJSON, 0, len(page.Programs))
	for _, p := range page.Programs {
		rows = append(rows, calcJSON{
			Code: p.Code, Kind: p.Kind, Name: p.Name, Lender: p.Lender,
			Rate: p.Rate(), Months: p.MaxMonths, Down: p.MinDownPct, Note: p.Note,
		})
	}
	if b, err := json.Marshal(rows); err == nil {
		page.DataJSON = string(b)
	} else {
		page.DataJSON = "[]"
	}

	market := map[string]float64{}
	for _, p := range page.Programs {
		if r := p.Rate(); r > market[p.Kind] {
			market[p.Kind] = r
		}
	}
	if b, err := json.Marshal(market); err == nil {
		page.MarketJSON = string(b)
	} else {
		page.MarketJSON = "{}"
	}

	labels := map[string]string{
		"multiple": T(lang, "calc.r_multiple"),
		"monthly":  T(lang, "calc.r_monthly"),
		"first":    T(lang, "calc.r_first"),
		"last":     T(lang, "calc.r_last"),
		"overpay":  T(lang, "calc.r_overpay_short"),
		"interest": T(lang, "calc.r_interest"),
		"fees":     T(lang, "calc.r_fees"),
		"seller":   T(lang, "calc.w_seller"),
		"wint":     T(lang, "calc.w_int"),
		"wfees":    T(lang, "calc.w_fees"),
		"wins":     T(lang, "calc.w_ins"),
		"subsidy":  T(lang, "calc.w_subsidy"),
		"you":      T(lang, "calc.w_you"),
	}
	if b, err := json.Marshal(labels); err == nil {
		page.LabelsJSON = string(b)
	} else {
		page.LabelsJSON = "{}"
	}

	page.Base = m.base(r, T(lang, "calc.title"), lang)
	// The shared sidebar links to the calculator; on the calculator itself that
	// link points at the page you are already reading.
	page.Base.Active = "calculator"
	page.Desc = T(lang, "calc.desc")
	m.render(w, "calculator", page)
}

// ProgramsByKind groups the programmes for the plain-HTML table, so a reader
// without script still sees them sorted by what they are.
func (p CalcPage) ProgramsByKind(kind string) []LoanProgram {
	var out []LoanProgram
	for _, x := range p.Programs {
		if x.Kind == kind {
			out = append(out, x)
		}
	}
	return out
}
