package articles

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// Hosts, visitors and visits.
//
// analytics_daily counts events and cannot tell one reader from another, which
// is why it reports views and nothing else. These three figures need people
// told apart, so this file does that -- under two rules that keep it from
// becoming a profile.
//
// The first is the salt. An identifier here is a truncated HMAC of the address
// and the browser string under a key that is generated fresh every day. The
// same reader on two days produces two unrelated values, so a day's distinct
// visitors can be counted and nobody can be followed past midnight. The
// addresses themselves are never written down.
//
// The second is the slot. A row is one visitor inside one half-hour, which is
// also the definition of a visit: click through five pages and you stay one
// visit, come back after a break and you are a new one. Every figure is then a
// plain query over stored rows -- no session table, no timers, nothing to
// rebuild after a restart.

// idLen is how much of the HMAC is kept. Sixteen bytes is far past the point
// where two readers collide, and stopping there means the full digest -- the
// only value an attacker holding the salt could test addresses against -- is
// never stored.
const idLen = 16

// slotKey is one visitor in one half-hour, with the audience switches they
// belong to.
type slotKey struct {
	slot     time.Time
	vid      [idLen]byte
	host     [idLen]byte
	isKZ     bool
	isMobile bool
}

// slotOf truncates a moment to the half-hour it falls in.
func slotOf(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), (t.Minute()/30)*30, 0, 0, time.UTC)
}

// saltCache holds the day's key so the hot path does not query for it. It is
// refreshed when the day turns; the stored copy exists only so a restart does
// not split one day into two populations of unrelated hashes.
type saltCache struct {
	mu  sync.Mutex
	day string
	key []byte
}

// salt returns today's key, creating and storing one if this is its first use.
// On any database trouble it falls back to a process-local random key: losing
// continuity across a restart is a worse report, never a leak.
func (s *saltCache) salt(ctx context.Context, mt *Metrics) []byte {
	day := time.Now().UTC().Format("2006-01-02")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.day == day && len(s.key) > 0 {
		return s.key
	}
	fresh := make([]byte, 32)
	if _, err := rand.Read(fresh); err != nil {
		return nil
	}
	var stored []byte
	if mt != nil && mt.db != nil {
		err := mt.db.QueryRow(ctx, `
			INSERT INTO analytics_salt (day, salt) VALUES ($1, $2)
			ON CONFLICT (day) DO UPDATE SET salt = analytics_salt.salt
			RETURNING salt`, day, fresh).Scan(&stored)
		if err != nil && mt.log != nil {
			mt.log.Warn("analytics salt", zap.Error(err))
		}
	}
	if len(stored) == 0 {
		stored = fresh
	}
	s.day, s.key = day, stored
	return stored
}

// ident derives the pair of identifiers for one hit. host covers the address
// alone, so two people behind one office router are one host and two visitors,
// which is what those words mean everywhere else.
func ident(salt []byte, ip, ua string) (vid, host [idLen]byte) {
	h := hmac.New(sha256.New, salt)
	h.Write([]byte(ip))
	copy(host[:], h.Sum(nil))
	h.Reset()
	h.Write([]byte(ip))
	h.Write([]byte{0})
	h.Write([]byte(ua))
	copy(vid[:], h.Sum(nil))
	return vid, host
}

// noteSlot buffers one page view against its visitor-slot. A missing address or
// an unusable salt drops the hit rather than counting an unattributable one.
func (mt *Metrics) noteSlot(ctx context.Context, r *http.Request, isKZ, isMobile bool) {
	if mt == nil {
		return
	}
	ip := clientIP(r)
	if len(ip) == 0 {
		return
	}
	salt := mt.salts.salt(ctx, mt)
	if len(salt) == 0 {
		return
	}
	vid, host := ident(salt, ip.String(), r.Header.Get("User-Agent"))
	k := slotKey{slot: slotOf(time.Now()), vid: vid, host: host, isKZ: isKZ, isMobile: isMobile}
	mt.mu.Lock()
	if mt.slots == nil {
		mt.slots = map[slotKey]int64{}
	}
	mt.slots[k]++
	mt.mu.Unlock()
}

// flushSlots writes the buffered visitor-slots. Views accumulate onto an
// existing row; the audience flags are set once, by the hit that opened it.
func (mt *Metrics) flushSlots(ctx context.Context) {
	if mt == nil || mt.db == nil {
		return
	}
	mt.mu.Lock()
	if len(mt.slots) == 0 {
		mt.mu.Unlock()
		return
	}
	batch := mt.slots
	mt.slots = map[slotKey]int64{}
	mt.mu.Unlock()

	b := &pgx.Batch{}
	for k, n := range batch {
		b.Queue(`
			INSERT INTO analytics_slots (slot, vid, host, is_kz, is_mobile, views)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (slot, vid)
			DO UPDATE SET views = analytics_slots.views + EXCLUDED.views`,
			k.slot, k.vid[:], k.host[:], k.isKZ, k.isMobile, n)
	}
	res := mt.db.SendBatch(ctx, b)
	defer res.Close()
	for range batch {
		if _, err := res.Exec(); err != nil {
			if mt.log != nil {
				mt.log.Warn("analytics slots flush", zap.Error(err))
			}
			return
		}
	}
}

// purge enforces what the privacy policy promises.
//
// The policy says yesterday's key cannot be recovered, and that is only true if
// yesterday's key is gone. Salts are dropped after two days -- one day of grace
// so a flush straddling midnight still finds the key it hashed with -- and from
// then on no stored identifier can be traced back to an address by anyone,
// ourselves included.
//
// The rows themselves outlive their salt on purpose: unlinkable counts are what
// the chart is made of. They go at two years, which is the longest window the
// month view asks for.
func (mt *Metrics) purge(ctx context.Context) {
	if mt == nil || mt.db == nil {
		return
	}
	if _, err := mt.db.Exec(ctx,
		`DELETE FROM analytics_salt WHERE day < CURRENT_DATE - INTERVAL '2 days'`); err != nil {
		if mt.log != nil {
			mt.log.Warn("analytics salt purge", zap.Error(err))
		}
	}
	if _, err := mt.db.Exec(ctx,
		`DELETE FROM analytics_slots WHERE slot < NOW() - INTERVAL '2 years'`); err != nil {
		if mt.log != nil {
			mt.log.Warn("analytics slots purge", zap.Error(err))
		}
	}
}
