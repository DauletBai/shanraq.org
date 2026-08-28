package articles

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Turning an article into something worth hearing.
//
// Two jobs live here, and they have to be done in one place because the page
// depends on both agreeing. The first is deciding what the blocks are: the
// player highlights block N while cue N is sounding, so the sequence the
// synthesiser is given must be the sequence the browser walks. The second is
// deciding what the words are, because a page is written for eyes and a good
// deal of it means nothing to an ear.

// narrationBlocks are the tags read aloud, matching the selector in listen.js.
// A block nested inside another is read once, at the outer one.
var narrationBlocks = map[string]bool{
	"p": true, "h2": true, "h3": true, "h4": true, "li": true,
	"blockquote": true, "figcaption": true, "th": true, "td": true,
}

var narrationSkip = map[string]bool{"pre": true, "code": true}

// Characters that carry meaning to the eye and none to the ear. Read out, a
// backslash or a hash is noise dropped into the middle of a sentence.
var muteChars = regexp.MustCompile(`[\\|*#~` + "`" + `^<>{}\[\]_]+`)

// Brackets and quotation marks are removed while their contents stay. The marks
// are punctuation for the eye: they group and attribute, and a voice conveys
// neither by pronouncing them. What is inside them is ordinary prose and has to
// survive -- dropping the aside along with its brackets would lose sentences.
var bracketsQuotes = regexp.MustCompile(`[()«»„“”"']`)

// A bare address is unlistenable -- "h t t p s colon slash slash" for twenty
// characters. A link keeps its visible label, which is prose; the address goes.
var urlPattern = regexp.MustCompile(`(?i)\bhttps?://\S+|\bwww\.\S+`)

// The separator between sources reads as "middle dot". It is a pause.
var dotSeparators = regexp.MustCompile(`\s*[·•]\s*`)

// A thousands separator is a space to the eye and a full stop to the voice:
// "16 700" comes out as "sixteen, seven hundred". Closed up it is read as the
// one number it is. Only groups of exactly three digits are joined, so ordinary
// prose is left alone.
var digitGroups = regexp.MustCompile(`(\d)[ \x{00a0}\x{202f}](\d{3})\b`)

var manySpaces = regexp.MustCompile(`\s+`)

// speechText rewrites one block into what should actually be said.
func speechText(s string) string {
	s = urlPattern.ReplaceAllString(s, " ")
	s = muteChars.ReplaceAllString(s, " ")
	s = bracketsQuotes.ReplaceAllString(s, " ")
	// Twice: 55 800 000 carries two separators, and one pass closes one gap.
	s = digitGroups.ReplaceAllString(s, "$1$2")
	s = digitGroups.ReplaceAllString(s, "$1$2")
	s = dotSeparators.ReplaceAllString(s, ", ")
	// A dash between words is a pause; a dash against a word is a hyphen and
	// must stay, or compound words come apart.
	s = strings.ReplaceAll(s, " — ", ", ")
	s = strings.ReplaceAll(s, " – ", ", ")
	s = manySpaces.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// NarrationBlocks walks rendered article HTML and returns what to read, in the
// order the page shows it.
//
// It parses the same HTML the browser receives rather than the Markdown behind
// it. The cue map is only useful if block N here is block N in the DOM, and the
// renderer is free to add, merge or wrap things on the way; reading the source
// instead would be a second guess at what it did.
func NarrationBlocks(rendered string) []string {
	doc, err := html.Parse(strings.NewReader(rendered))
	if err != nil {
		return nil
	}
	var out []string
	var walk func(n *html.Node, inBlock, skipping bool)
	walk = func(n *html.Node, inBlock, skipping bool) {
		if n.Type == html.ElementNode {
			if narrationSkip[n.Data] || hasHiddenAttr(n) {
				skipping = true
			} else if narrationBlocks[n.Data] && !inBlock && !skipping {
				if t := speechText(textOf(n)); len([]rune(t)) >= 2 {
					out = append(out, t)
				}
				// Everything below has been taken with the outer block.
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c, inBlock, skipping)
		}
	}
	walk(doc, false, false)
	return out
}

func hasHiddenAttr(n *html.Node) bool {
	for _, a := range n.Attr {
		if a.Key == "aria-hidden" && a.Val == "true" {
			return true
		}
		if a.Key == "hidden" {
			return true
		}
	}
	return false
}

// textOf collects the visible text of a node, skipping what is not read.
func textOf(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
			return
		}
		if n.Type == html.ElementNode && (narrationSkip[n.Data] || hasHiddenAttr(n)) {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return b.String()
}
