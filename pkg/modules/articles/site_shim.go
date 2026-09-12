package articles

import (
	"html/template"

	"shanraq.org/pkg/site"
)

// The page frame — the languages, the UI dictionary, and the structs the shared
// header, sidebar and footer render — lives in pkg/site, so that a module that
// is not "articles" (the shop, the weather, the payments desk) can render a
// page of this site without importing this package.
//
// The names below are aliases, not copies: articles.Base and site.Base are the
// same type. They exist so that moving the definitions did not have to touch a
// thousand lines that only mention them. Each one disappears when the last file
// that uses it leaves this package for a module of its own.
const (
	LangKZ = site.LangKZ
	LangRU = site.LangRU
	LangEN = site.LangEN

	langCookieName = site.LangCookie
)

var (
	Langs      = site.Langs
	LangLabels = site.LangLabels
	LangNames  = site.LangNames
)

type (
	Base        = site.Base
	ServiceView = site.ServiceView
	FeedItem    = site.FeedItem
	Ad          = site.Ad
	Rate        = site.Rate
	SocialLink  = site.SocialLink
	InfoBarData = site.InfoBarData
)

// The rubric list is the site's own menu rather than one module's taxonomy:
// articles are filed under it, listings sit beside it, and the header draws it.
const CategoryGeneral = site.CategoryGeneral

var (
	Categories    = site.Categories
	Subcategories = site.Subcategories
)

// payKindAdOrder is what an advertising order is called in the payments ledger.
// The payments module treats it as an opaque label; this package is the one
// that knows what it means.
const payKindAdOrder = "ad_order"

// TOCItem is the frame's, since the frame renders the Markdown.
type TOCItem = site.TOCItem

// RenderMarkdown is re-exported: the studio preview and the tools that render
// a body outside a page call it by this name.
func RenderMarkdown(source string) template.HTML { return site.RenderMarkdown(source) }
