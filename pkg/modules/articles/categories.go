package articles

// CategoryGeneral is the default rubric for unclassified articles.
const CategoryGeneral = "general"

// Categories is the ordered list of real rubrics shown in the portal nav
// ("general" is the fallback bucket and is not shown as a nav chip).
var Categories = []string{"society", "politics", "economy", "culture", "technology", "opinion", "world"}

var validCategories = func() map[string]bool {
	m := map[string]bool{CategoryGeneral: true}
	for _, c := range Categories {
		m[c] = true
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
