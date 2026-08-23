package articles

import (
	"context"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Догрузка курсов.
//
// Нацбанк отдаёт курс за конкретный день и держит около пяти лет назад; за
// более ранние даты приходит пустой документ. Поэтому архив собираем сами:
// один день за раз, не спеша, пока не упрёмся в край чужого окна, а дальше
// каждый час забираем сегодняшний.
//
// Медленно — намеренно. Полная история это около двух тысяч запросов, и
// выпустить их залпом значило бы устроить чужому серверу маленький штурм ради
// страницы, которая никуда не торопится.

const (
	// fxStep — пауза между запросами при догрузке.
	fxStep = 2 * time.Second
	// fxIdle — как часто проверять, не появился ли новый день, когда догружать
	// уже нечего.
	fxIdle = time.Hour
	// fxEmptyRun — сколько подряд пустых дней считать краем окна источника.
	// Праздничная неделя даёт до девяти, поэтому порог заметно выше.
	fxEmptyRun = 21
)

// fxFloor — глубина, ниже которой мы не спрашиваем. Банк начинает молчать
// примерно на середине 2021 года; год запаса стоит одного лишнего дня опроса.
var fxFloor = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

// RunFxArchivist ведёт архив курсов до отмены контекста.
func (m *Module) RunFxArchivist(ctx context.Context) {
	if m.fx == nil {
		return
	}
	// Глубину подтягиваем сразу после первого же дня — он подсказывает, какие
	// валюты нам нужны, — и дальше раз в час. Ждать конца дневной догрузки
	// нельзя: она идёт больше часа, и всё это время «весь период» был бы пуст.
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

// fxHistory подтягивает месячную глубину, если её ещё нет или она отстала.
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

// fxOnce обрабатывает один день. Возвращает true, когда догружать больше
// нечего.
func (m *Module) fxOnce(ctx context.Context) (bool, error) {
	day, ok, err := m.fx.NextToProbe(ctx, fxFloor)
	if err != nil || !ok {
		return true, err
	}
	// Дошли до края: ниже источник молчит, и дальше идти незачем.
	if run, err := m.fx.EmptyRunBelow(ctx, day); err == nil && run >= fxEmptyRun {
		return true, nil
	}

	rates, err := m.fetchFxDay(ctx, day)
	if err != nil {
		// Сетевая ошибка — не ответ источника. День не помечаем, попробуем
		// снова: пометить его пустым значило бы потерять день навсегда.
		return false, err
	}
	if err := m.fx.Save(ctx, rates); err != nil {
		return false, err
	}
	return false, m.fx.MarkProbed(ctx, day, len(rates))
}

// fetchFxDay забирает и разбирает курсы за день.
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
