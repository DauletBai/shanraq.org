// Package ratings implements article voting and author reputation ("karma").
//
// A reader casts a +1 or -1 vote on an article. The vote counts for a weight
// derived from the voter's own karma, so established members influence scores
// more than fresh (possibly sock-puppet) accounts. An article's cached score is
// the weighted sum of its votes; an author's karma is the weighted sum across
// all of their articles.
package ratings

// VoteUp / VoteDown / VoteNone are the accepted vote values.
const (
	VoteDown = -1
	VoteNone = 0
	VoteUp   = 1
)

// maxWeight caps how much a single highly-reputed voter can swing a score.
const maxWeight = 5

// Weight maps a voter's karma to their per-vote weight. New or negative-karma
// voters count for 1; weight rises one point per 100 karma, capped at maxWeight.
func Weight(karma int) int {
	w := 1 + karma/100
	if w < 1 {
		w = 1
	}
	if w > maxWeight {
		w = maxWeight
	}
	return w
}

// Rating captures an article's cached score plus the viewing user's own vote.
type Rating struct {
	Score    int
	UserVote int // -1, 0, or +1
}
