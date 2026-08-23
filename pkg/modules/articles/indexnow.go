package articles

import (
	"context"
	"time"
)

// Заявки поисковикам на постоянные страницы.
//
// Карта сайта — это приглашение зайти когда-нибудь; IndexNow говорит Bing,
// Яндексу и остальным, кто его поддерживает, прямо сейчас. Статьи уходили туда
// с самой публикации, а постоянные страницы — «Аналитика», «Курсы валют»,
// правила, тарифы — не уходили никогда: их никто не «публикует», и повода
// отправить заявку не возникало.
//
// Поэтому поводы назначены здесь. После запуска сайт объявляет весь свой
// постоянный состав: разворачивание меняет разметку и тексты, и это ровно тот
// случай, когда поисковику стоит зайти заново. Дальше раз в сутки уходят те
// страницы, содержимое которых меняется каждый день само.

const (
	// indexNowSettle — пауза после запуска. Заявка ссылается на файл ключа на
	// нашем же домене, и отправлять её раньше, чем сайт начал отвечать,
	// значит получить отказ.
	indexNowSettle = 2 * time.Minute
	// indexNowEvery — как часто объявлять ежедневно меняющиеся страницы.
	indexNowEvery = 24 * time.Hour
)

// indexNowDaily — страницы, которые меняются каждый день сами по себе: лента
// главной, курс валют и счётчики аудитории.
var indexNowDaily = []string{"/", "/rates", "/analytics"}

// runIndexNow объявляет постоянные страницы сайта до отмены контекста.
func (m *Module) runIndexNow(ctx context.Context) {
	if m.syndicate == nil {
		return
	}
	select {
	case <-ctx.Done():
		return
	case <-time.After(indexNowSettle):
	}

	all := append([]string{"/"}, publicPages...)
	m.syndicate.SubmitURLs(m.pageURLs(all), "постоянные страницы")

	t := time.NewTicker(indexNowEvery)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			m.syndicate.SubmitURLs(m.pageURLs(indexNowDaily), "ежедневные страницы")
		}
	}
}

// pageURLs разворачивает пути в полные адреса на всех трёх языках — ровно те,
// что стоят в карте сайта и в canonical. Заявка на адрес, которого сайт не
// объявляет своим, поисковику ничего не даёт.
func (m *Module) pageURLs(paths []string) []string {
	site := m.rt.Config.PublicBase()
	out := make([]string, 0, len(paths)*len(Langs))
	for _, p := range paths {
		for _, lang := range Langs {
			out = append(out, site+canonURL(p, "", lang))
		}
	}
	return out
}
