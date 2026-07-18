package articles

// Ad is one creative in the sidebar ad carousel. Demo ads are self-made,
// clearly-labeled placeholders — no real brand, logo, or photo — so nothing
// implies a commercial relationship that does not exist. Real ads would later
// come from config or a paid-placement store; the template contract stays the
// same.
type Ad struct {
	Image string // /static/... illustration
	Title string
	Price string
	Desc  string
	URL   string // click target ("#" for demo)
}

// demoAds returns three localized placeholder car ads for the sidebar slot,
// used until real paid placements are wired up.
func demoAds(lang string) []Ad {
	return []Ad{
		{Image: "/static/demo/ads/car-sedan.svg", URL: "#",
			Title: T(lang, "ad.sedan_title"), Desc: T(lang, "ad.sedan_desc"), Price: "11 990 000 ₸"},
		{Image: "/static/demo/ads/car-suv.svg", URL: "#",
			Title: T(lang, "ad.suv_title"), Desc: T(lang, "ad.suv_desc"), Price: "15 490 000 ₸"},
		{Image: "/static/demo/ads/car-hatch.svg", URL: "#",
			Title: T(lang, "ad.hatch_title"), Desc: T(lang, "ad.hatch_desc"), Price: "9 890 000 ₸"},
	}
}
