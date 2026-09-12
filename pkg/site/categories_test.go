package site

import "testing"

func TestCategories(t *testing.T) {
	if len(Categories) == 0 {
		t.Fatal("no categories defined")
	}
	c := Categories[0]
	if !IsCategory(c) {
		t.Errorf("%q should be a category", c)
	}
	if IsCategory("no-such-category") {
		t.Error("unknown category must be false")
	}
	if NormalizeCategory(c) != c {
		t.Error("valid category should pass through NormalizeCategory")
	}
	if NormalizeCategory("no-such-category") != CategoryGeneral {
		t.Errorf("invalid category should become %q", CategoryGeneral)
	}
	// Find any subcategory and verify the parent round-trip.
	var sub, parent string
	for cat, subs := range Subcategories {
		if len(subs) > 0 {
			sub, parent = subs[0], cat
			break
		}
	}
	if sub != "" {
		if !IsSubcategory(sub) {
			t.Errorf("%q should be a subcategory", sub)
		}
		if SubcategoryParent(sub) != parent {
			t.Errorf("parent of %q = %q, want %q", sub, SubcategoryParent(sub), parent)
		}
		if NormalizeSubcategory(parent, sub) != sub {
			t.Error("sub under its parent should pass NormalizeSubcategory")
		}
		if NormalizeSubcategory("wrong-parent", sub) != "" {
			t.Error("sub under a wrong parent should be dropped")
		}
	}
}
