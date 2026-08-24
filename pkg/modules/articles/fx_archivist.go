package articles

import (
	"context"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Backfilling the rates.
//
// The National Bank serves one specific day's rate and holds about five years
// back; for earlier dates an empty document arrives. So we build the archive
// ourselves: one day at a time, unhurried, until we hit the edge of somebody
// else's window, and after that we fetch today's every hour.
//
// Slowly, on purpose. The full history is about two thousand requests, and firing
// them in a volley would mean staging a small assault on someone else's server for
// the sake of a page that is in no hurry at all.

const (
	// fxStep is the pause between requests while backfilling.
	fxStep = 2 * time.Second
	// fxIdle is how often to check whether a new day has appeared once there is
	// nothing left to backfill.
	fxIdle = time.Hour
	// fxEmptyRun is how many consecutive empty days count as the source's window
	// edge. A holiday week gives up to nine, so the threshold sits well above.
	fxEmptyRun = 21
)

// fxFloor is the depth below which we do not ask. The bank falls silent around
// the middle of 2021; a year of margin costs one extra day of probing.
var fxFloor = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

// RunFxArchivist keeps the rate archive until the context is cancelled.
func (m *Module) RunFxArchivist(ctx context.Context) {
	if m.fx == nil {
		return
	}
	// The deep archive is pulled right after the first day — that day tells us
	// which currencies we need — and once an hour after that. Waiting for the
	// daily backfill to finish is not an option: it runs for over an hour, and
	// "all time" would be empty for the whole of it.
	var nextDeep time.Time
	for {
		done, err := m.fxOnce(ctx)
		if err != nil {
			m.rt.Logger.Warn("архив курсов", zap.Error(err))
		}
		if !time.Now().Before(nextDeep) {
			if err := m.fxHistory(ctx); err != nil {
				m.rt.Logger.Warn("месячный архив курсов", zap.Error(err))
			}
			nextDeep = time.Now().Add(fxIdle)
		}
		wait := fxStep
		if done {
			wait = fxIdle
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
	}
}

// fxHistory pulls the monthly depth when it is missing or has fallen behind.
func (m *Module) fxHistory(ctx context.Context) error {
	fresh, err := m.fx.MonthlyFresh(ctx)
	if err != nil || fresh {
		return err
	}
	rows, err := m.fetchBISMonthly(ctx)
	if err != nil || len(rows) == 0 {
		return err
	}
	if err := m.fx.SaveMonthly(ctx, rows); err != nil {
		return err
	}
	m.rt.Logger.Info("месячный архив курсов обновлён", zap.Int("точек", len(rows)))
	return nil
}

// fxOnce handles one day. It returns true once there is nothing left to
// backfill.
func (m *Module) fxOnce(ctx context.Context) (bool, error) {
	day, ok, err := m.fx.NextToProbe(ctx, fxFloor)
	if err != nil || !ok {
		return true, err
	}
	// We have reached the edge: below this the source is silent, and there is no
	// point going further.
	if run, err := m.fx.EmptyRunBelow(ctx, day); err == nil && run >= fxEmptyRun {
		return true, nil
	}

	rates, err := m.fetchFxDay(ctx, day)
	if err != nil {
		// A network error is not the source's answer. The day is left unmarked so
		// it will be tried again: marking it empty would lose the day for good.
		return false, err
	}
	if err := m.fx.Save(ctx, rates); err != nil {
		return false, err
	}
	return false, m.fx.MarkProbed(ctx, day, len(rates))
}

// fetchFxDay fetches and parses one day's rates.
func (m *Module) fetchFxDay(ctx context.Context, day time.Time) ([]FxRate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fxRatesURL(day), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "shanraq.org/1.0")
	resp, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	return parseFxRates(body, day)
}
