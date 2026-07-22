package articles

import "context"

// AdminAnalytics is the growth dashboard: the numbers a founder actually needs
// to see whether the platform is picking up — organized by the funnel it
// describes, from audience through authors and listings to the referral loop.
type AdminAnalytics struct {
	// Audience.
	Subscribers int

	// People.
	Users         int
	VerifiedEmail int
	VerifiedPhone int
	Authors       int // users who can publish (real name + verified phone)

	// Content.
	Published int
	InReview  int
	NeedsWork int
	Drafts    int
	AIColumns int // first-party AI Dake columns

	// Listings.
	Listings      int
	ActiveListing int
	Promoted      int

	// Engagement.
	Comments int
	Hidden   int

	// Referral loop.
	Invited      int // people who signed up via someone's link
	Qualified    int // …who then posted a real listing (the reward trigger)
	TopReferrers []ReferrerRow
	CreditGiven  int // promotion-days granted
	CreditSpent  int // promotion-days used
}

// ConversionPct is the share of invited users who became qualified — the single
// number that says whether the referral loop actually works.
func (a AdminAnalytics) ConversionPct() int {
	if a.Invited == 0 {
		return 0
	}
	return a.Qualified * 100 / a.Invited
}

// ReferrerRow is one entry in the top-referrers table.
type ReferrerRow struct {
	Email     string
	Invited   int
	Qualified int
}

// adminAnalytics gathers the dashboard in a handful of grouped queries. Each is
// a cheap aggregate; at launch scale this is well under a millisecond.
func (m *Module) adminAnalytics(ctx context.Context) (AdminAnalytics, error) {
	var a AdminAnalytics
	db := m.rt.DB

	_ = db.QueryRow(ctx, `SELECT COUNT(*) FROM subscribers`).Scan(&a.Subscribers)

	_ = db.QueryRow(ctx, `
		SELECT COUNT(*),
		       COUNT(*) FILTER (WHERE email_verified_at IS NOT NULL),
		       COUNT(*) FILTER (WHERE phone_verified_at IS NOT NULL),
		       COUNT(*) FILTER (WHERE phone_verified_at IS NOT NULL AND first_name <> '' AND last_name <> '')
		  FROM auth_users`).Scan(&a.Users, &a.VerifiedEmail, &a.VerifiedPhone, &a.Authors)

	_ = db.QueryRow(ctx, `
		SELECT COUNT(*) FILTER (WHERE status = 'published'),
		       COUNT(*) FILTER (WHERE status = 'review'),
		       COUNT(*) FILTER (WHERE status = 'needs_work'),
		       COUNT(*) FILTER (WHERE status = 'draft'),
		       COUNT(*) FILTER (WHERE status = 'published' AND author_id = $1)
		  FROM articles`, SanaAuthorID).Scan(&a.Published, &a.InReview, &a.NeedsWork, &a.Drafts, &a.AIColumns)

	_ = db.QueryRow(ctx, `
		SELECT COUNT(*),
		       COUNT(*) FILTER (WHERE status = 'published' AND expires_at > NOW()),
		       COUNT(*) FILTER (WHERE promoted_until > NOW())
		  FROM listings`).Scan(&a.Listings, &a.ActiveListing, &a.Promoted)

	_ = db.QueryRow(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'hidden')
		  FROM comments`).Scan(&a.Comments, &a.Hidden)

	_ = db.QueryRow(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'qualified')
		  FROM referrals`).Scan(&a.Invited, &a.Qualified)

	// Promotion-day credit split into granted (positive deltas) and spent
	// (absolute of the negatives), so the two read naturally in the UI.
	_ = db.QueryRow(ctx, `
		SELECT COALESCE(SUM(delta_days) FILTER (WHERE delta_days > 0), 0),
		       COALESCE(-SUM(delta_days) FILTER (WHERE delta_days < 0), 0)
		  FROM promo_credit_ledger`).Scan(&a.CreditGiven, &a.CreditSpent)

	rows, err := db.Query(ctx, `
		SELECT COALESCE(u.email, ''), COUNT(r.*),
		       COUNT(r.*) FILTER (WHERE r.status = 'qualified')
		  FROM referrals r
		  LEFT JOIN auth_users u ON u.id = r.referrer_id
		 GROUP BY u.email
		 ORDER BY COUNT(r.*) FILTER (WHERE r.status = 'qualified') DESC, COUNT(r.*) DESC
		 LIMIT 10`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rr ReferrerRow
			if err := rows.Scan(&rr.Email, &rr.Invited, &rr.Qualified); err == nil {
				a.TopReferrers = append(a.TopReferrers, rr)
			}
		}
	}
	return a, nil
}
