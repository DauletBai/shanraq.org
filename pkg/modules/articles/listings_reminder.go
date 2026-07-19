package articles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"shanraq.org/pkg/shanraq"
)

// reminderInterval is how often the expiry-reminder sweep runs.
const reminderInterval = 6 * time.Hour

// Start launches the background sweep that reminds owners about listings whose
// free window is about to end (within ~2 days).
func (m *Module) Start(ctx context.Context, _ *shanraq.Runtime) error {
	go m.reminderLoop(ctx)
	if m.infobar != nil {
		go m.infobar.Run(ctx) // background weather + exchange-rate refresher
	}
	return nil
}

func (m *Module) reminderLoop(ctx context.Context) {
	select { // let boot/migrations settle before the first sweep
	case <-ctx.Done():
		return
	case <-time.After(30 * time.Second):
	}
	m.sweepReminders(ctx)
	m.sweepExpired(ctx)
	t := time.NewTicker(reminderInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			m.sweepReminders(ctx)
			m.sweepExpired(ctx)
		}
	}
}

// sweepExpired permanently deletes listings past their 21-day window and all
// their data (owners were warned 2 days earlier by sweepReminders).
func (m *Module) sweepExpired(ctx context.Context) {
	n, err := m.listings.PurgeExpired(ctx)
	if err != nil {
		m.rt.Logger.Error("listing purge sweep", zap.Error(err))
		return
	}
	if n > 0 {
		m.rt.Logger.Info("purged expired listings", zap.Int64("count", n))
	}
}

func (m *Module) sweepReminders(ctx context.Context) {
	if m.mailer == nil {
		return
	}
	due, err := m.listings.DueReminders(ctx)
	if err != nil {
		m.rt.Logger.Error("listing reminders sweep", zap.Error(err))
		return
	}
	base := strings.TrimRight(m.rt.Config.Syndicate.BaseURL, "/")
	for _, l := range due {
		subject := "Ваше объявление скоро истекает — Shanraq"
		body := fmt.Sprintf(
			"Здравствуйте!\n\nВаше объявление «%s» истекает через %d дн. "+
				"Чтобы оно осталось в поиске, продлите его:\n"+
				"%s/listings/my\n\nЕсли не продлить, по истечении срока объявление "+
				"и все его данные будут удалены безвозвратно.\n\n— Shanraq",
			l.Title, l.DaysLeft(), base,
		)
		if err := m.mailer.Send(ctx, l.AuthorEmail, subject, body); err != nil {
			m.rt.Logger.Warn("listing reminder not sent", zap.String("to", l.AuthorEmail), zap.Error(err))
			continue // leave unmarked; retry on the next sweep
		}
		if err := m.listings.MarkReminded(ctx, l.ID); err != nil {
			m.rt.Logger.Error("mark reminded", zap.Error(err))
		}
	}
}
