package articles

import "testing"

// A cloud address must never reach the audience, however ordinary its browser
// string. This is the defect the panel was built on for months: the network was
// resolved one branch too late and kept as something to slice by, so automation
// sending a normal Chrome header was counted as a reader in views, visitors,
// sources, devices and OS at once. Moving one line breaks it again, so the rule
// is nailed down here.
func TestAudienceBucket(t *testing.T) {
	const chrome = "" // a browser string leaves botLabel empty

	cases := []struct {
		name, bot, country, want string
	}{
		{"облако с браузерным UA — не аудитория", chrome, datacenterLabel, bucketDatacenter},
		{"настоящая страна — аудитория", chrome, "KZ", bucketAudience},
		{"страна не определилась — аудитория", chrome, "", bucketAudience},
		{"известный краулер — в панель ботов", "google", "US", bucketBot},
		{"краулер из облака — считается краулером, не датацентром", "ai", datacenterLabel, bucketBot},
		{"seo-сканер — выбрасывается совсем", "seo", "NL", bucketDrop},
		{"seo-сканер из облака — тоже выбрасывается", "seo", datacenterLabel, bucketDrop},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := audienceBucket(c.bot, c.country); got != c.want {
				t.Errorf("audienceBucket(%q, %q) = %q, хотели %q", c.bot, c.country, got, c.want)
			}
		})
	}
}

// Only one bucket may be the audience: if a second one ever starts counting as
// a reader, the panel goes back to being unreadable and nobody notices.
func TestOnlyOneAudienceBucket(t *testing.T) {
	seen := 0
	for _, bot := range []string{"", "google", "ai", "seo", "other"} {
		for _, country := range []string{"KZ", "US", "", datacenterLabel} {
			if audienceBucket(bot, country) == bucketAudience {
				seen++
				if bot != "" {
					t.Errorf("бот %q попал в аудиторию (страна %q)", bot, country)
				}
				if country == datacenterLabel {
					t.Errorf("датацентр попал в аудиторию (бот %q)", bot)
				}
			}
		}
	}
	if seen == 0 {
		t.Fatal("в аудиторию не попадает никто — фильтр перекрыл всё")
	}
}
