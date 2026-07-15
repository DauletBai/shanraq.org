package articles

// CategoryGeneral is the default rubric for unclassified articles.
const CategoryGeneral = "general"

// Categories is the ordered list of rubrics shown in the portal menu
// ("general" is the fallback bucket and is not shown as a menu item).
var Categories = []string{"sport", "society", "politics", "economy", "culture", "technology", "it", "opinion", "world"}

// Subcategories maps each category to its subcategory slugs. Categories in the
// top menu are ordered by editorial importance; subcategories within a menu are
// ordered alphabetically (by their Russian label) so a long list is easy to
// scan. Slugs are globally unique.
var Subcategories = map[string][]string{
	"sport":      {"basketball", "boxing", "wrestling", "cycling", "athletics", "mma", "tennis", "football", "hockey", "chess"},
	"society":    {"charity", "health", "crime", "migration", "youth", "education", "holidays", "religion", "family", "ecology"},
	"politics":   {"elections", "diplomacy", "law", "corruption", "defense", "parliament", "parties", "government", "regions"},
	"economy":    {"banks", "business", "industry", "labor", "agriculture", "startups", "trade", "finance", "prices", "energy"},
	"culture":    {"architecture", "art", "history", "cinema", "literature", "fashion", "music", "theatre", "traditions", "language"},
	"technology": {"aviation", "autotech", "agrotech", "biotech", "space", "machinery", "nanotech", "science", "robotics", "shipbuilding", "telecom", "electronics"},
	"it":         {"databases", "backend", "webdev", "devops", "ai", "internet", "cybersecurity", "mobile", "cloud", "gamedev", "frontend"},
	"opinion":    {"analytics", "blogs", "debate", "interview", "column", "review", "letters", "editorial", "satire"},
	"world":      {"asia", "africa", "middle_east", "europe", "oceania", "russia", "north_america", "cis", "central_asia", "south_america"},
}

var validCategories = func() map[string]bool {
	m := map[string]bool{CategoryGeneral: true}
	for _, c := range Categories {
		m[c] = true
	}
	return m
}()

var subToCategory = func() map[string]string {
	m := map[string]string{}
	for cat, subs := range Subcategories {
		for _, s := range subs {
			m[s] = cat
		}
	}
	return m
}()

// IsCategory reports whether s is a known category slug.
func IsCategory(s string) bool { return validCategories[s] }

// NormalizeCategory returns s if valid, else the general bucket.
func NormalizeCategory(s string) string {
	if IsCategory(s) {
		return s
	}
	return CategoryGeneral
}

// Subcats returns the subcategory slugs of a category (nil if none).
func Subcats(cat string) []string { return Subcategories[cat] }

// IsSubcategory reports whether s is a known subcategory slug.
func IsSubcategory(s string) bool { _, ok := subToCategory[s]; return ok }

// SubcategoryParent returns the category a subcategory belongs to.
func SubcategoryParent(s string) string { return subToCategory[s] }

// NormalizeSubcategory keeps sub only if it belongs to cat; otherwise "".
func NormalizeSubcategory(cat, sub string) string {
	if sub != "" && subToCategory[sub] == cat {
		return sub
	}
	return ""
}
