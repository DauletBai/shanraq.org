package articles

// CategoryGeneral is the default rubric for unclassified articles.
const CategoryGeneral = "general"

// Categories is the ordered list of rubrics shown in the portal menu
// ("general" is the fallback bucket and is not shown as a menu item).
var Categories = []string{"sport", "society", "politics", "economy", "culture", "technology", "opinion", "world"}

// Subcategories maps each category to its ordered list of subcategory slugs
// (the submenu). Slugs are globally unique.
var Subcategories = map[string][]string{
	"sport":      {"football", "boxing", "hockey", "athletics", "basketball", "tennis", "wrestling", "mma", "chess", "cycling"},
	"society":    {"education", "health", "family", "religion", "ecology", "migration", "charity", "youth", "crime", "holidays"},
	"politics":   {"elections", "parliament", "regions", "government", "diplomacy", "law", "defense", "parties", "corruption"},
	"economy":    {"prices", "business", "energy", "finance", "agriculture", "industry", "trade", "startups", "labor", "banks"},
	"culture":    {"literature", "cinema", "music", "theatre", "art", "history", "language", "fashion", "architecture", "traditions"},
	"technology": {"internet", "ai", "science", "gadgets", "space", "software", "cybersecurity", "telecom", "gaming", "autotech"},
	"opinion":    {"column", "debate", "interview", "analytics", "editorial", "letters", "blogs", "review", "satire"},
	"world":      {"central_asia", "europe", "asia", "north_america", "south_america", "africa", "oceania", "middle_east", "cis", "russia"},
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
