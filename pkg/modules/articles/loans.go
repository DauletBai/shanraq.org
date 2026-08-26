package articles

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The loan calculator: mortgage, car, consumer, instalment, business.
//
// Two things make it worth building rather than linking to a bank's own.
//
// First, what it puts first. Every lender's calculator leads with the monthly
// instalment, because an instalment is a small number and a small number sells.
// This one leads with what the thing ends up costing -- the down payment plus
// everything handed over -- and that total as a multiple of the price.
//
// Second, the effective rate. An advertised rate is not what a loan costs: fees
// and compulsory insurance sit outside it. Kazakh law knows this and requires
// lenders to disclose an annual effective rate, capped at 56 % unsecured and
// 40 % secured. We compute the same kind of figure from the fees the reader
// enters, and say plainly that it is our arithmetic and not the bank's official
// disclosure -- a bank folds in charges we cannot see from outside.
//
// The arithmetic is an identity, not a forecast, so it is a pure function and
// tested as one.

// Loan kinds. Stable identifiers: they key the programme table, the form and
// the translations.
const (
	LoanMortgage    = "mortgage"
	LoanAuto        = "auto"
	LoanConsumer    = "consumer"
	LoanInstallment = "installment"
	LoanBusiness    = "business"
)

// LoanKinds is the display order of the loan types.
var LoanKinds = []string{LoanMortgage, LoanAuto, LoanConsumer, LoanInstallment, LoanBusiness}

// Repayment schemes. Annuity keeps the payment level; differentiated repays
// equal parts of the principal, so the first payment is the largest and the
// total interest is lower. Kazakh banks offer both and the choice is worth real
// money, so the calculator has to model it.
const (
	SchemeAnnuity        = "annuity"
	SchemeDifferentiated = "differentiated"
)

// LoanInput is everything the form collects.
type LoanInput struct {
	Kind    string
	Price   int64 // price of the home or car; 0 for a plain cash loan
	Down    int64 // down payment
	Amount  int64 // borrowed directly, when there is no Price
	RatePct float64
	Months  int
	Scheme  string

	FeeOncePct       float64 // one-off fee, % of the loan
	FeeOnceFixed     int64   // one-off fee, tenge
	FeeMonthlyPct    float64 // monthly servicing fee, % of the loan
	InsuranceYearPct float64 // compulsory insurance, % of the loan a year

	// MarketRatePct is what the same loan would cost without state help. When
	// it is above RatePct the gap is somebody's bill, and that somebody is the
	// budget: a subsidised rate is not a cheaper loan, it is the same loan with
	// part of it paid by taxpayers. Left at zero, no subsidy is shown.
	MarketRatePct float64
}

// LoanRow is one month of the schedule.
type LoanRow struct {
	N         int
	Payment   int64
	Principal int64
	Interest  int64
	Fees      int64
	Balance   int64
}

// LoanShares is where the borrower's money ends up.
//
// A loan has more parties than the two on the contract, and only one line here
// buys anything that exists: what goes to the seller. The rest is the price of
// the money. Naming the recipients is the whole point -- a total overpayment is
// abstract, "this much to the bank, this much to the insurer" is not.
type LoanShares struct {
	Seller       int64 // the price of the thing itself
	BankInterest int64
	BankFees     int64
	Insurer      int64
	Subsidy      int64 // what the budget adds on a subsidised rate, estimated
}

// LoanPlan is the outcome of one calculation.
type LoanPlan struct {
	Loan      int64
	Monthly   int64 // annuity: the payment; differentiated: the first one
	LastMonth int64 // differentiated: the final, smallest payment
	TotalPaid int64 // everything handed over, fees and insurance included
	Interest  int64
	Fees      int64
	Overpay   int64   // TotalPaid minus the loan: the price of the money
	TotalCost int64   // down payment plus TotalPaid
	Multiple  float64 // TotalCost over the price, when there is a price
	EffRate   float64 // annual effective rate, our own calculation
	BankFees  int64   // fees to the lender, one-off and monthly
	Insurance int64   // premiums to the insurer
	Shares    LoanShares
	Schedule  []LoanRow
}

// loanPrincipal is what is actually borrowed.
func loanPrincipal(in LoanInput) (loan, down, price int64) {
	price, down = in.Price, in.Down
	if price > 0 {
		if down < 0 {
			down = 0
		}
		if down > price {
			down = price
		}
		return price - down, down, price
	}
	if in.Amount > 0 {
		return in.Amount, 0, 0
	}
	return 0, 0, 0
}

// LoanCalc builds the whole schedule and the totals.
//
// Guards return an empty plan rather than an error: the page recalculates on
// every keystroke, and a half-typed sum is an ordinary state, not a fault.
func LoanCalc(in LoanInput) LoanPlan {
	loan, down, price := loanPrincipal(in)
	if loan <= 0 || in.Months <= 0 {
		return LoanPlan{}
	}
	n := in.Months
	if n > 600 {
		n = 600 // fifty years: past this the schedule is noise, not a plan
	}
	rate := in.RatePct
	if rate < 0 {
		rate = 0
	}
	i := rate / 100 / 12

	feeMonthly := int64(math.Round(float64(loan) * in.FeeMonthlyPct / 100))
	insMonthly := int64(math.Round(float64(loan) * in.InsuranceYearPct / 100 / 12))
	perMonthFees := feeMonthly + insMonthly

	feeOnce := in.FeeOnceFixed + int64(math.Round(float64(loan)*in.FeeOncePct/100))
	if feeOnce < 0 {
		feeOnce = 0
	}

	p := LoanPlan{Loan: loan, Schedule: make([]LoanRow, 0, n)}

	balance := loan
	var basePay float64
	if in.Scheme != SchemeDifferentiated {
		if i <= 0 {
			basePay = float64(loan) / float64(n)
		} else {
			basePay = float64(loan) * i / (1 - math.Pow(1+i, -float64(n)))
		}
	}
	flatPrincipal := float64(loan) / float64(n)

	for k := 1; k <= n; k++ {
		interest := int64(math.Round(float64(balance) * i))
		var principal int64
		if in.Scheme == SchemeDifferentiated {
			principal = int64(math.Round(flatPrincipal))
		} else {
			principal = int64(math.Round(basePay)) - interest
		}
		if k == n || principal > balance {
			principal = balance // the last month clears whatever rounding left
		}
		if principal < 0 {
			principal = 0
		}
		balance -= principal

		row := LoanRow{N: k, Principal: principal, Interest: interest, Fees: perMonthFees, Balance: balance}
		row.Payment = principal + interest + perMonthFees
		p.Schedule = append(p.Schedule, row)

		p.Interest += interest
		p.Fees += perMonthFees
		p.BankFees += feeMonthly
		p.Insurance += insMonthly
		p.TotalPaid += row.Payment
	}

	p.Fees += feeOnce
	p.BankFees += feeOnce
	p.TotalPaid += feeOnce
	if len(p.Schedule) > 0 {
		p.Monthly = p.Schedule[0].Payment
		p.LastMonth = p.Schedule[len(p.Schedule)-1].Payment
	}
	p.Overpay = p.TotalPaid - loan
	p.TotalCost = down + p.TotalPaid
	if price > 0 {
		p.Multiple = float64(p.TotalCost) / float64(price)
	}
	p.EffRate = loanEffectiveRate(loan, feeOnce, p.Schedule)

	// Who ends up with the money. The seller's line is the only one that buys
	// something that exists; everything after it is the price of the money.
	p.Shares = LoanShares{
		Seller:       price,
		BankInterest: p.Interest,
		BankFees:     p.BankFees,
		Insurer:      p.Insurance,
	}
	if price == 0 {
		p.Shares.Seller = loan // a cash loan hands the borrower the sum itself
	}
	// A subsidised rate is not a cheaper loan: it is the same loan with part of
	// the interest paid from the budget. Estimated as the interest the same
	// schedule would carry at the market rate, less the interest actually
	// charged. MarketRatePct is cleared on the inner call so this cannot recur.
	if in.MarketRatePct > rate {
		at := in
		at.RatePct, at.MarketRatePct = in.MarketRatePct, 0
		if full := LoanCalc(at); full.Interest > p.Interest {
			p.Shares.Subsidy = full.Interest - p.Interest
		}
	}
	return p
}

// loanEffectiveRate is the annual rate at which the payments actually made
// discount back to the money actually received.
//
// This is the arithmetic behind the disclosure Kazakh law requires: the sum a
// borrower gets is the loan less any fee withheld at signing, and what they
// give back is every instalment including servicing charges and insurance. If
// the nominal rate were the whole story the two would agree; they do not, and
// the gap is the point of showing it.
//
// Solved by bisection rather than Newton: a monotone function on a bracket that
// cannot be escaped, so it converges from any input a form can produce, and no
// derivative can send it off to infinity.
func loanEffectiveRate(loan, feeOnce int64, sched []LoanRow) float64 {
	net := float64(loan - feeOnce)
	if net <= 0 || len(sched) == 0 {
		return 0
	}
	npv := func(monthly float64) float64 {
		var sum float64
		for _, r := range sched {
			sum += float64(r.Payment) / math.Pow(1+monthly, float64(r.N))
		}
		return sum - net
	}
	// At a zero rate the payments cannot discount to less than the money taken,
	// unless the loan is interest-free and fee-free -- in which case the answer
	// is zero and there is nothing to search for.
	if npv(0) <= 0 {
		return 0
	}
	lo, hi := 0.0, 1.0 // 100 % a month is far beyond any legal Kazakh loan
	for npv(hi) > 0 && hi < 10 {
		hi *= 2
	}
	for k := 0; k < 200; k++ {
		mid := (lo + hi) / 2
		if npv(mid) > 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	monthly := (lo + hi) / 2
	return (math.Pow(1+monthly, 12) - 1) * 100
}

// LoanProgram is one lending programme as the operator entered it.
type LoanProgram struct {
	Code       string
	Kind       string
	Lender     string
	Name       string
	Note       string
	RateBP     int // basis points: 700 = 7.00 %
	MaxMonths  int
	MinDownPct int
	SourceURL  string
	CheckedOn  string // dd.mm.yyyy, empty when never checked
}

// Rate returns the programme's annual rate as a percentage.
func (p LoanProgram) Rate() float64 { return float64(p.RateBP) / 100 }

// LoanStore reads the programme table.
type LoanStore struct{ db *pgxpool.Pool }

// NewLoanStore returns a store over the pool.
func NewLoanStore(db *pgxpool.Pool) *LoanStore { return &LoanStore{db: db} }

// List returns the active programmes in display order, with names and notes in
// one language. A row left untranslated falls back to Russian, so a partly
// filled table still renders something a reader can use.
func (s *LoanStore) List(ctx context.Context, lang string) ([]LoanProgram, error) {
	if s == nil || s.db == nil {
		return nil, nil
	}
	rows, err := s.db.Query(ctx, `
		SELECT code, kind, lender,
		       COALESCE(NULLIF(CASE $1 WHEN 'kz' THEN name_kz WHEN 'en' THEN name_en ELSE name_ru END, ''), name_ru),
		       COALESCE(NULLIF(CASE $1 WHEN 'kz' THEN note_kz WHEN 'en' THEN note_en ELSE note_ru END, ''), note_ru),
		       rate_bp, max_months, min_down_pct, source_url,
		       COALESCE(to_char(checked_on, 'DD.MM.YYYY'), '')
		FROM loan_programs
		WHERE active
		ORDER BY sort, code`, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []LoanProgram
	for rows.Next() {
		var p LoanProgram
		if err := rows.Scan(&p.Code, &p.Kind, &p.Lender, &p.Name, &p.Note,
			&p.RateBP, &p.MaxMonths, &p.MinDownPct, &p.SourceURL, &p.CheckedOn); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Save writes one programme back, for the admin panel.
func (s *LoanStore) Save(ctx context.Context, p LoanProgram, active bool) error {
	if s == nil || s.db == nil {
		return nil
	}
	_, err := s.db.Exec(ctx, `
		UPDATE loan_programs
		   SET active = $2, rate_bp = $3, max_months = $4, min_down_pct = $5,
		       source_url = $6, checked_on = NULLIF($7, '')::date, updated_at = now()
		 WHERE code = $1`,
		p.Code, active, p.RateBP, p.MaxMonths, p.MinDownPct, p.SourceURL, p.CheckedOn)
	return err
}
