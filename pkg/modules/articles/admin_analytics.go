package articles

import (
	"context"
	"fmt"
)

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
	Agents        int // registered real-estate agents

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

	_ = db.QueryRow(ctx, `SELECT COUNT(*) FROM re_agents`).Scan(&a.Agents)

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

// courseAnalyticsSince is the first day on which lesson and course openings
// use the complete audience rule (crawler, hosting-network, staff and test
// traffic removed). Older article counters remain visible as a labelled
// reference because they cannot be repaired after the request has gone.
const courseAnalyticsSince = "26.09.2026"

// CourseAnalytics is the course section of the admin dashboard. The figures
// intentionally stop at aggregate cohorts: guest/signed-in and reading
// language. Shanraq does not keep a visitor's identity or lesson path.
type CourseAnalytics struct {
	Since string
	Day   Audience
	Week  Audience
	Month Audience
	All   Audience
	Hubs  Audience

	Learners int64
	Attempts int64
	Passed   int64

	Courses []CourseAnalyticsRow
	Lessons []CourseLessonAnalytics
}

type CourseAnalyticsRow struct {
	Slug  string
	Title string

	Lessons int64
	Views   Audience
	Hub     Audience
	KZ      int64
	RU      int64
	EN      int64

	// LegacyViews predates the complete audience filter and can contain live
	// browser checks made while a course was being published.
	LegacyViews int64
	Started     int64
	Finished    int64
	Learners    int64
	Attempts    int64
	Passed      int64
}

type CourseLessonAnalytics struct {
	CourseSlug  string
	CourseTitle string
	Slug        string
	Title       string
	Position    int
	Views       Audience
	LegacyViews int64
	Started     int64
	Finished    int64
}

func (m *Module) courseAnalytics(ctx context.Context, lang string) (CourseAnalytics, error) {
	out := CourseAnalytics{Since: courseAnalyticsSince}
	db := m.rt.DB

	// The four windows apply to filtered lesson openings only. Hub openings are
	// shown separately: opening a course map is intent, opening a lesson is use.
	if err := db.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(n) FILTER (WHERE kind=$1 AND day=(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Almaty')::date AND is_guest),0),
			COALESCE(SUM(n) FILTER (WHERE kind=$1 AND day=(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Almaty')::date AND NOT is_guest),0),
			COALESCE(SUM(n) FILTER (WHERE kind=$1 AND day>=(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Almaty')::date-6 AND is_guest),0),
			COALESCE(SUM(n) FILTER (WHERE kind=$1 AND day>=(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Almaty')::date-6 AND NOT is_guest),0),
			COALESCE(SUM(n) FILTER (WHERE kind=$1 AND day>=(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Almaty')::date-29 AND is_guest),0),
			COALESCE(SUM(n) FILTER (WHERE kind=$1 AND day>=(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Almaty')::date-29 AND NOT is_guest),0),
			COALESCE(SUM(n) FILTER (WHERE kind=$1 AND is_guest),0),
			COALESCE(SUM(n) FILTER (WHERE kind=$1 AND NOT is_guest),0),
			COALESCE(SUM(n) FILTER (WHERE kind=$2 AND is_guest),0),
			COALESCE(SUM(n) FILTER (WHERE kind=$2 AND NOT is_guest),0)
		FROM analytics_daily WHERE kind IN ($1,$2)`, metricCourseLesson, metricCourseHub).
		Scan(&out.Day.Guest, &out.Day.Registered, &out.Week.Guest, &out.Week.Registered,
			&out.Month.Guest, &out.Month.Registered, &out.All.Guest, &out.All.Registered,
			&out.Hubs.Guest, &out.Hubs.Registered); err != nil {
		return out, fmt.Errorf("course analytics totals: %w", err)
	}

	rows, err := db.Query(ctx, `
		WITH legacy AS (
			SELECT article_id, SUM(views) AS views FROM article_views_daily GROUP BY article_id
		), reads AS (
			SELECT article_id, finished FROM article_reads
		), depth AS (
			SELECT article_id,
			       SUM(count) FILTER (WHERE depth=25) AS started
			FROM reading_depth GROUP BY article_id
		), item_stats AS (
			SELECT i.series_id,
			       COUNT(*) FILTER (WHERE a.status='published' AND i.position>=10) AS lessons,
			       COALESCE(SUM(l.views) FILTER (WHERE a.status='published'),0) AS legacy_views,
			       COALESCE(SUM(d.started) FILTER (WHERE a.status='published'),0) AS started,
			       COALESCE(SUM(r.finished) FILTER (WHERE a.status='published'),0) AS finished
			FROM article_series_items i
			JOIN articles a ON a.id=i.article_id
			LEFT JOIN legacy l ON l.article_id=a.id
			LEFT JOIN reads r ON r.article_id=a.id
			LEFT JOIN depth d ON d.article_id=a.id
			GROUP BY i.series_id
		), progress AS (
			SELECT i.series_id, COUNT(DISTINCT p.user_id) AS learners,
			       COALESCE(SUM(p.attempts),0) AS attempts,
			       COUNT(*) FILTER (WHERE p.passed) AS passed
			FROM article_series_items i
			JOIN course_progress p ON p.article_id=i.article_id
			GROUP BY i.series_id
		), clean AS (
			SELECT split_part(label,'|',1) AS slug,
			       COALESCE(SUM(n) FILTER (WHERE is_guest),0) AS guest,
			       COALESCE(SUM(n) FILTER (WHERE NOT is_guest),0) AS registered,
			       COALESCE(SUM(n) FILTER (WHERE split_part(label,'|',3)='kz'),0) AS kz,
			       COALESCE(SUM(n) FILTER (WHERE split_part(label,'|',3)='ru'),0) AS ru,
			       COALESCE(SUM(n) FILTER (WHERE split_part(label,'|',3)='en'),0) AS en
			FROM analytics_daily WHERE kind=$1 GROUP BY 1
		), hubs AS (
			SELECT split_part(label,'|',1) AS slug,
			       COALESCE(SUM(n) FILTER (WHERE is_guest),0) AS guest,
			       COALESCE(SUM(n) FILTER (WHERE NOT is_guest),0) AS registered
			FROM analytics_daily WHERE kind=$2 GROUP BY 1
		)
		SELECT s.slug, COALESCE(NULLIF(t.title,''),s.slug),
		       COALESCE(i.lessons,0), COALESCE(c.guest,0), COALESCE(c.registered,0),
		       COALESCE(h.guest,0), COALESCE(h.registered,0),
		       COALESCE(c.kz,0), COALESCE(c.ru,0), COALESCE(c.en,0),
		       COALESCE(i.legacy_views,0), COALESCE(i.started,0), COALESCE(i.finished,0),
		       COALESCE(p.learners,0), COALESCE(p.attempts,0), COALESCE(p.passed,0)
		FROM article_series s
		LEFT JOIN article_series_i18n t ON t.series_id=s.id AND t.lang=$3
		LEFT JOIN item_stats i ON i.series_id=s.id
		LEFT JOIN progress p ON p.series_id=s.id
		LEFT JOIN clean c ON c.slug=s.slug
		LEFT JOIN hubs h ON h.slug=s.slug
		WHERE s.status='published'
		ORDER BY COALESCE(c.guest,0)+COALESCE(c.registered,0) DESC,
		         COALESCE(i.legacy_views,0) DESC`, metricCourseLesson, metricCourseHub, lang)
	if err != nil {
		return out, fmt.Errorf("course analytics rows: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var r CourseAnalyticsRow
		if err := rows.Scan(&r.Slug, &r.Title, &r.Lessons,
			&r.Views.Guest, &r.Views.Registered, &r.Hub.Guest, &r.Hub.Registered,
			&r.KZ, &r.RU, &r.EN, &r.LegacyViews, &r.Started, &r.Finished,
			&r.Learners, &r.Attempts, &r.Passed); err != nil {
			return out, fmt.Errorf("scan course analytics: %w", err)
		}
		out.Courses = append(out.Courses, r)
	}
	if err := rows.Err(); err != nil {
		return out, err
	}
	if err := db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT p.user_id), COALESCE(SUM(p.attempts),0),
		       COUNT(*) FILTER (WHERE p.passed)
		FROM course_progress p
		WHERE EXISTS (
			SELECT 1 FROM article_series_items i
			JOIN article_series s ON s.id=i.series_id AND s.status='published'
			WHERE i.article_id=p.article_id
		)`).
		Scan(&out.Learners, &out.Attempts, &out.Passed); err != nil {
		return out, fmt.Errorf("course progress totals: %w", err)
	}

	lessons, err := db.Query(ctx, `
		WITH clean AS (
			SELECT split_part(label,'|',1) AS course_slug,
			       split_part(label,'|',2) AS article_slug,
			       COALESCE(SUM(n) FILTER (WHERE is_guest),0) AS guest,
			       COALESCE(SUM(n) FILTER (WHERE NOT is_guest),0) AS registered
			FROM analytics_daily WHERE kind=$1 GROUP BY 1,2
		), legacy AS (
			SELECT article_id, SUM(views) AS views FROM article_views_daily GROUP BY article_id
		), depth AS (
			SELECT article_id, SUM(count) FILTER (WHERE depth=25) AS started
			FROM reading_depth GROUP BY article_id
		)
		SELECT s.slug, COALESCE(NULLIF(st.title,''),s.slug), a.slug,
		       COALESCE(NULLIF(at.title,''),NULLIF(orig.title,''),a.slug), i.position,
		       COALESCE(c.guest,0), COALESCE(c.registered,0), COALESCE(l.views,0),
		       COALESCE(d.started,0), COALESCE(ar.finished,0)
		FROM article_series s
		JOIN article_series_items i ON i.series_id=s.id
		JOIN articles a ON a.id=i.article_id AND a.status='published'
		LEFT JOIN article_series_i18n st ON st.series_id=s.id AND st.lang=$2
		LEFT JOIN article_translations at ON at.article_id=a.id AND at.lang=$2
		LEFT JOIN article_translations orig ON orig.article_id=a.id AND orig.lang=a.original_lang
		LEFT JOIN clean c ON c.course_slug=s.slug AND c.article_slug=a.slug
		LEFT JOIN legacy l ON l.article_id=a.id
		LEFT JOIN depth d ON d.article_id=a.id
		LEFT JOIN article_reads ar ON ar.article_id=a.id
		WHERE s.status='published'
		ORDER BY COALESCE(c.guest,0)+COALESCE(c.registered,0) DESC,
		         COALESCE(l.views,0) DESC, s.slug, i.position`, metricCourseLesson, lang)
	if err != nil {
		return out, fmt.Errorf("course lesson analytics: %w", err)
	}
	defer lessons.Close()
	for lessons.Next() {
		var r CourseLessonAnalytics
		if err := lessons.Scan(&r.CourseSlug, &r.CourseTitle, &r.Slug, &r.Title,
			&r.Position, &r.Views.Guest, &r.Views.Registered, &r.LegacyViews,
			&r.Started, &r.Finished); err != nil {
			return out, fmt.Errorf("scan lesson analytics: %w", err)
		}
		out.Lessons = append(out.Lessons, r)
	}
	return out, lessons.Err()
}
