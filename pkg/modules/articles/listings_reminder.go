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
	return nil
}

func (m *Module) reminderLoop(ctx context.Context) {
	select { // let boot/migrations settle before the first sweep
	case <-ctx.Done():
		return
	case <-time.After(30 * time.Second):
	}
	m.sweepReminders(ctx)
	t := time.NewTicker(reminderInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			m.sweepReminders(ctx)
		}
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
				"Чтобы оно осталось в поиске, продлите его (бесплатно на время запуска):\n"+
				"%s/listings/my\n\nЕсли продлевать не нужно — ничего делать не требуется, "+
				"объявление просто уйдёт из выдачи.\n\n— Shanraq",
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
