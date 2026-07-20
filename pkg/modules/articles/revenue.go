package articles

// Revenue policy — the single source of truth, so the public texts, the
// advertiser cabinet, and any future payout code cannot drift apart.
//
// Two kinds of authorship, two different answers:
//
//   - Platform AI agents (AI Dake, and any agent added to aiAgentAuthors).
//     These are our own models running on our own infrastructure. They are not
//     counterparties and cannot hold rights or receive payment. All revenue
//     earned around their materials belongs to the project in full.
//
//   - Independent human authors. A revenue-share program is planned at
//     RevenueAuthorPct / RevenuePlatformPct, and will be fixed in a written
//     author agreement — not by a sentence on a web page.
//
// There is deliberately no "community fund" share. An earlier draft promised
// 10% to one; no such entity was ever constituted, had a regulator, or had
// distribution rules, so the promise was removed rather than left dangling.
//
// RevenueProgramLive gates every public statement of the split. It stays false
// until money can actually move: an author agreement, withholding of ИПН and
// social contributions by the TOO as tax agent, and a payout ledger. Publicly
// naming a split before then is a public offer under art. 395 of the Civil
// Code of the RK that we would be unable to honour.
const (
	// RevenueAuthorPct is the independent author's share of revenue
	// attributable to their material, once the program launches.
	RevenueAuthorPct = 50
	// RevenuePlatformPct is the platform's share of the same revenue.
	RevenuePlatformPct = 100 - RevenueAuthorPct
	// RevenueProgramLive reports whether author payouts actually run. While
	// false, public pages describe the program as planned, never as terms.
	RevenueProgramLive = false
)

// aiAgentAuthors are the platform's own AI author accounts. Material published
// under these IDs is work of the platform: revenue from it goes to the project
// in full and no author share is ever accrued against it.
var aiAgentAuthors = map[string]string{
	SanaAuthorID: SanaName, // AI Dake — the AI columnist
	// Future agents (e.g. AI Make) are registered here and nowhere else.
}

// IsAIAgentAuthor reports whether an author ID is one of the platform's AI
// agents. This is authorship, not tooling: it is a different question from
// whether a translation was machine-generated (Source == "ai"), which is also
// true for a human author who used the AI assistant.
func IsAIAgentAuthor(authorID string) bool {
	_, ok := aiAgentAuthors[authorID]
	return ok
}

// RevenueShare returns the author/platform split for material by the given
// author. Platform AI agents keep no share — the project takes all of it.
func RevenueShare(authorID string) (authorPct, platformPct int) {
	if IsAIAgentAuthor(authorID) {
		return 0, 100
	}
	return RevenueAuthorPct, RevenuePlatformPct
}
