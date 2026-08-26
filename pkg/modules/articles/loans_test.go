package articles

import (
	"math"
	"testing"
)

// The figure the page is built around: at a commercial rate a flat costs its
// buyer about three times the asking price. That is what the reader sees first.
func TestLoanTotalCostIsTheHeadlineFigure(t *testing.T) {
	p := LoanCalc(LoanInput{Kind: LoanMortgage, Price: 20_000_000, Down: 4_000_000,
		RatePct: 18, Months: 240, Scheme: SchemeAnnuity})

	if p.Loan != 16_000_000 {
		t.Fatalf("loan came to %d, expected 16,000,000", p.Loan)
	}
	if p.Monthly < 246_000 || p.Monthly > 247_500 {
		t.Errorf("payment %d, expected about 246,930", p.Monthly)
	}
	if p.Multiple < 3.1 || p.Multiple > 3.3 {
		t.Errorf("purchase cost %.2f times the price, expected about 3.2", p.Multiple)
	}
	if p.TotalCost != 4_000_000+p.TotalPaid {
		t.Errorf("total %d is not the down payment plus the instalments", p.TotalCost)
	}
	if p.Overpay != p.TotalPaid-p.Loan {
		t.Errorf("overpayment %d is not the total paid less the principal", p.Overpay)
	}
}

// Equal parts of the principal cost less in interest than a level payment, and
// the first instalment is the heaviest. If that inverts, the scheme is wired
// backwards and the reader is advised into the costlier option.
func TestDifferentiatedCostsLessButStartsHeavier(t *testing.T) {
	in := LoanInput{Kind: LoanMortgage, Price: 20_000_000, Down: 4_000_000, RatePct: 18, Months: 240}
	in.Scheme = SchemeAnnuity
	ann := LoanCalc(in)
	in.Scheme = SchemeDifferentiated
	dif := LoanCalc(in)

	if dif.Interest >= ann.Interest {
		t.Errorf("differentiated charged %d in interest against the annuity's %d; it must cost less", dif.Interest, ann.Interest)
	}
	if dif.Monthly <= ann.Monthly {
		t.Errorf("the first differentiated payment %d is not above the annuity payment %d", dif.Monthly, ann.Monthly)
	}
	if dif.LastMonth >= dif.Monthly {
		t.Errorf("the last payment %d is not below the first %d; the payment must fall", dif.LastMonth, dif.Monthly)
	}
	// Whatever the scheme, the principal must be repaid exactly once.
	for _, p := range []LoanPlan{ann, dif} {
		var principal int64
		for _, r := range p.Schedule {
			principal += r.Principal
		}
		if principal != p.Loan {
			t.Errorf("%d of principal repaid on a loan of %d", principal, p.Loan)
		}
		if last := p.Schedule[len(p.Schedule)-1].Balance; last != 0 {
			t.Errorf("balance after the final payment is %d, expected zero", last)
		}
	}
}

// The point of showing an effective rate. An advertised 20 % is not 20 %:
// charged monthly it compounds to 21.94 % before a single fee is added, which
// is why the rate a lender must disclose is always above the one it prints on
// the poster. Fees then push it further.
func TestFeesPushTheEffectiveRateAboveTheAdvertisedOne(t *testing.T) {
	clean := LoanCalc(LoanInput{Kind: LoanConsumer, Amount: 2_000_000, RatePct: 20, Months: 36, Scheme: SchemeAnnuity})
	if math.Abs(clean.EffRate-21.94) > 0.15 {
		t.Errorf("with no fees the effective rate is %.2f%%, expected 21.94 -- 20%% compounded monthly", clean.EffRate)
	}
	if clean.EffRate <= 20 {
		t.Errorf("effective rate %.2f%% is not above the advertised 20%%; monthly compounding is not counted", clean.EffRate)
	}
	loaded := LoanCalc(LoanInput{Kind: LoanConsumer, Amount: 2_000_000, RatePct: 20, Months: 36,
		Scheme: SchemeAnnuity, FeeOncePct: 2, FeeMonthlyPct: 0.5, InsuranceYearPct: 1})
	if loaded.EffRate <= clean.EffRate+3 {
		t.Errorf("fees lifted the rate from %.2f only to %.2f; they are not being counted", clean.EffRate, loaded.EffRate)
	}
	if loaded.BankFees <= 0 || loaded.Insurance <= 0 {
		t.Errorf("bank fees %d and insurance %d must be reported apart", loaded.BankFees, loaded.Insurance)
	}
}

// Every tenge the borrower parts with has a recipient, and the sum of the
// recipients has to be the sum paid. A breakdown that does not add up is worse
// than no breakdown.
func TestSharesAccountForEveryTenge(t *testing.T) {
	p := LoanCalc(LoanInput{Kind: LoanMortgage, Price: 20_000_000, Down: 4_000_000,
		RatePct: 18, Months: 240, Scheme: SchemeAnnuity, FeeOncePct: 1, InsuranceYearPct: 0.5})

	sum := p.Shares.Seller + p.Shares.BankInterest + p.Shares.BankFees + p.Shares.Insurer
	if sum != p.TotalCost {
		t.Errorf("the recipients total %d while the buyer parts with %d", sum, p.TotalCost)
	}
	if p.Shares.Seller != 20_000_000 {
		t.Errorf("the seller received %d against a price of 20,000,000", p.Shares.Seller)
	}
}

// A subsidised rate is the same loan with part of the interest paid from the
// budget. The calculator has to name that number, not hide it in a lower rate.
func TestSubsidyIsWhatTheBudgetAdds(t *testing.T) {
	p := LoanCalc(LoanInput{Kind: LoanMortgage, Price: 20_000_000, Down: 4_000_000,
		RatePct: 7, Months: 300, Scheme: SchemeAnnuity, MarketRatePct: 18})
	if p.Shares.Subsidy <= 0 {
		t.Fatal("no budget subsidy computed for a subsidised rate")
	}
	market := LoanCalc(LoanInput{Kind: LoanMortgage, Price: 20_000_000, Down: 4_000_000,
		RatePct: 18, Months: 300, Scheme: SchemeAnnuity})
	if want := market.Interest - p.Interest; p.Shares.Subsidy != want {
		t.Errorf("subsidy %d, expected %d: the interest gap between the market and the subsidised rate", p.Shares.Subsidy, want)
	}
	// Without a market rate to compare against there is nothing to claim.
	if q := LoanCalc(LoanInput{Kind: LoanMortgage, Price: 20_000_000, Down: 4_000_000,
		RatePct: 7, Months: 300, Scheme: SchemeAnnuity}); q.Shares.Subsidy != 0 {
		t.Errorf("a subsidy of %d was computed with no market rate to compare against", q.Shares.Subsidy)
	}
}

// The page recalculates on every keystroke, so half-entered and absurd input is
// an ordinary state. None of it may panic or print a nonsense figure.
func TestHalfEnteredInputGivesAnEmptyPlan(t *testing.T) {
	cases := []struct {
		name string
		in   LoanInput
	}{
		{"price not entered", LoanInput{Kind: LoanMortgage, Months: 240, RatePct: 18}},
		{"term not entered", LoanInput{Kind: LoanMortgage, Price: 20_000_000, RatePct: 18}},
		{"negative price", LoanInput{Kind: LoanMortgage, Price: -5, Months: 240, RatePct: 18}},
		{"down payment above the price", LoanInput{Kind: LoanMortgage, Price: 20_000_000, Down: 30_000_000, Months: 240, RatePct: 18}},
		{"negative rate", LoanInput{Kind: LoanConsumer, Amount: 1_000_000, Months: 12, RatePct: -5}},
		{"term beyond the cap", LoanInput{Kind: LoanConsumer, Amount: 1_000_000, Months: 99999, RatePct: 10}},
	}
	for _, c := range cases {
		p := LoanCalc(c.in)
		if p.Monthly < 0 || p.TotalPaid < 0 || p.Loan < 0 ||
			math.IsNaN(p.Multiple) || math.IsInf(p.Multiple, 0) ||
			math.IsNaN(p.EffRate) || math.IsInf(p.EffRate, 0) {
			t.Errorf("%s: nonsensical plan %+v", c.name, p)
		}
		if len(p.Schedule) > 600 {
			t.Errorf("%s: schedule runs to %d months", c.name, len(p.Schedule))
		}
	}
}

// An interest-free instalment plan divides by zero in the annuity formula, and
// it is the ordinary case for the "рассрочка" tab, not an edge one.
func TestZeroRateRepaysInEqualParts(t *testing.T) {
	p := LoanCalc(LoanInput{Kind: LoanInstallment, Price: 1_200_000, Down: 0, RatePct: 0, Months: 12, Scheme: SchemeAnnuity})
	if p.Monthly != 100_000 {
		t.Errorf("payment %d, expected 100,000", p.Monthly)
	}
	if p.Overpay != 0 || p.EffRate != 0 {
		t.Errorf("at a zero rate the overpayment is %d and the rate %.2f, both expected to be zero", p.Overpay, p.EffRate)
	}
}

// Basis points keep the table in integers; the form wants percent.
func TestProgramRateReadsBasisPoints(t *testing.T) {
	if got := (LoanProgram{RateBP: 700}).Rate(); got != 7 {
		t.Errorf("700 basis points read as %v%%, expected 7", got)
	}
	if got := (LoanProgram{RateBP: 1875}).Rate(); got != 18.75 {
		t.Errorf("1875 basis points read as %v%%, expected 18.75", got)
	}
}
