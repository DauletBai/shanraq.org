package ratings

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// These run against a real database because what is worth checking here is the
// arithmetic the SQL does -- the weighted sums, and what a vote is allowed to
// touch. A mock would only assert that the queries we wrote are the queries we
// wrote.

// fixture is one author, one article, and however many voters a test asks for,
// all removed afterwards.
type fixture struct {
	pool    *pgxpool.Pool
	store   *Store
	author  uuid.UUID
	article uuid.UUID
	ctx     context.Context
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the ratings integration tests")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	f := &fixture{pool: pool, store: &Store{db: pool}, ctx: ctx,
		author: uuid.New(), article: uuid.New()}
	f.user(t, f.author)
	if _, err := pool.Exec(ctx,
		`INSERT INTO articles (id, author_id, slug, original_lang, status, published_at)
		 VALUES ($1,$2,$3,'ru','published',NOW())`,
		f.article, f.author, "rt-"+f.article.String()[:12]); err != nil {
		t.Fatalf("insert article: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM article_votes WHERE article_id=$1`, f.article)
		_, _ = pool.Exec(ctx, `DELETE FROM articles WHERE id=$1`, f.article)
		pool.Close()
	})
	return f
}

// user makes an account and schedules its removal, along with anything hanging
// off it that would hold the row down.
func (f *fixture) user(t *testing.T, id uuid.UUID) uuid.UUID {
	t.Helper()
	if _, err := f.pool.Exec(f.ctx,
		`INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'x','user')`,
		id, "rt-"+id.String()+"@t.test"); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM comment_votes WHERE user_id=$1`, id)
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM article_votes WHERE user_id=$1`, id)
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM author_reputation WHERE user_id=$1`, id)
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM comments WHERE user_id=$1`, id)
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM auth_users WHERE id=$1`, id)
	})
	return id
}

// karma sets a voter's standing directly, which is how a test buys the weight
// it wants to check without casting a hundred votes to earn it.
func (f *fixture) karma(t *testing.T, id uuid.UUID, k int) {
	t.Helper()
	if _, err := f.pool.Exec(f.ctx,
		`INSERT INTO author_reputation (user_id, karma) VALUES ($1,$2)
		 ON CONFLICT (user_id) DO UPDATE SET karma = EXCLUDED.karma`, id, k); err != nil {
		t.Fatalf("set karma: %v", err)
	}
}

func TestAnAuthorCannotVoteOnTheirOwnArticle(t *testing.T) {
	f := newFixture(t)
	if _, err := f.store.Vote(f.ctx, f.article, f.author, f.author, VoteUp); !errors.Is(err, ErrSelfVote) {
		t.Errorf("voting on one's own article returned %v, want ErrSelfVote", err)
	}
}

func TestAVoteMustBeUpDownOrNothing(t *testing.T) {
	f := newFixture(t)
	voter := f.user(t, uuid.New())
	for _, bad := range []int{2, -2, 100} {
		if _, err := f.store.Vote(f.ctx, f.article, voter, f.author, bad); err == nil {
			t.Errorf("value %d was accepted as a vote", bad)
		}
	}
}

// Up, then changed to down, then withdrawn. The score follows each time, and
// changing a vote replaces it rather than adding a second one.
func TestChangingAndWithdrawingAVote(t *testing.T) {
	f := newFixture(t)
	voter := f.user(t, uuid.New())

	if score, err := f.store.Vote(f.ctx, f.article, voter, f.author, VoteUp); err != nil || score != 1 {
		t.Fatalf("after an upvote: score %d, err %v; want 1", score, err)
	}
	if score, err := f.store.Vote(f.ctx, f.article, voter, f.author, VoteDown); err != nil || score != -1 {
		t.Fatalf("after changing to a downvote: score %d, err %v; want -1", score, err)
	}
	if score, err := f.store.Vote(f.ctx, f.article, voter, f.author, VoteNone); err != nil || score != 0 {
		t.Fatalf("after withdrawing: score %d, err %v; want 0", score, err)
	}
}

// The point of the weighting: a reader with standing moves a score further than
// a fresh account, and the cap stops one voter from deciding it alone.
func TestAVoteCarriesTheWeightOfItsVoter(t *testing.T) {
	f := newFixture(t)
	fresh := f.user(t, uuid.New())
	seasoned := f.user(t, uuid.New())
	capped := f.user(t, uuid.New())
	f.karma(t, seasoned, 250)  // weight 3
	f.karma(t, capped, 100000) // weight 5, the ceiling

	if score, _ := f.store.Vote(f.ctx, f.article, fresh, f.author, VoteUp); score != 1 {
		t.Errorf("a fresh account's vote counted %d, want 1", score)
	}
	if score, _ := f.store.Vote(f.ctx, f.article, seasoned, f.author, VoteUp); score != 4 {
		t.Errorf("250 karma should add 3, giving 4; got %d", score)
	}
	if score, _ := f.store.Vote(f.ctx, f.article, capped, f.author, VoteUp); score != 9 {
		t.Errorf("weight is capped at %d, so the total should be 9; got %d", maxWeight, score)
	}
}

// Karma is the weighted sum across everything the author has written, so a vote
// on one article shows up in the reputation the next one is judged by.
func TestKarmaFollowsTheVotesOnAnAuthorsArticles(t *testing.T) {
	f := newFixture(t)
	voter := f.user(t, uuid.New())

	if k, err := f.store.AuthorKarma(f.ctx, f.author); err != nil || k != 0 {
		t.Fatalf("an author with no votes has karma %d, err %v; want 0", k, err)
	}
	if _, err := f.store.Vote(f.ctx, f.article, voter, f.author, VoteUp); err != nil {
		t.Fatal(err)
	}
	if k, _ := f.store.AuthorKarma(f.ctx, f.author); k != 1 {
		t.Errorf("one upvote should leave karma 1, got %d", k)
	}
	if _, err := f.store.Vote(f.ctx, f.article, voter, f.author, VoteDown); err != nil {
		t.Fatal(err)
	}
	if k, _ := f.store.AuthorKarma(f.ctx, f.author); k != -1 {
		t.Errorf("changing that vote to a downvote should leave karma -1, got %d", k)
	}
}

func TestForArticleReportsTheViewersOwnVote(t *testing.T) {
	f := newFixture(t)
	voter := f.user(t, uuid.New())
	other := f.user(t, uuid.New())
	if _, err := f.store.Vote(f.ctx, f.article, voter, f.author, VoteUp); err != nil {
		t.Fatal(err)
	}

	if r, err := f.store.ForArticle(f.ctx, f.article, voter); err != nil || r.Score != 1 || r.UserVote != VoteUp {
		t.Errorf("the voter sees %+v, err %v; want score 1 and their own +1", r, err)
	}
	if r, _ := f.store.ForArticle(f.ctx, f.article, other); r.Score != 1 || r.UserVote != VoteNone {
		t.Errorf("a reader who has not voted sees %+v; want score 1 and no vote of their own", r)
	}
	// An article that is not there is a zero, not an error: the page asks about
	// whatever id is in the address.
	if r, err := f.store.ForArticle(f.ctx, uuid.New(), voter); err != nil || r.Score != 0 {
		t.Errorf("an unknown article gave %+v, err %v; want a plain zero", r, err)
	}
}

// The invariant the package is explicit about: an author's reputation is built
// from their articles, and letting comment votes into it would turn every
// argument in a thread into a standing on the author's record.
func TestCommentVotesDoNotTouchKarma(t *testing.T) {
	f := newFixture(t)
	commenter := f.user(t, uuid.New())
	voter := f.user(t, uuid.New())

	var commentID uuid.UUID
	if err := f.pool.QueryRow(f.ctx,
		`INSERT INTO comments (article_id, user_id, body) VALUES ($1,$2,'x') RETURNING id`,
		f.article, commenter).Scan(&commentID); err != nil {
		t.Fatalf("insert comment: %v", err)
	}

	before, _ := f.store.AuthorKarma(f.ctx, commenter)
	if score, err := f.store.VoteComment(f.ctx, commentID, voter, VoteUp); err != nil || score != 1 {
		t.Fatalf("upvoting a comment: score %d, err %v; want 1", score, err)
	}
	after, _ := f.store.AuthorKarma(f.ctx, commenter)
	if before != after {
		t.Errorf("a comment vote moved its author's karma from %d to %d; it must not", before, after)
	}
	if v, err := f.store.CommentVote(f.ctx, commentID, voter); err != nil || v != VoteUp {
		t.Errorf("the voter's own comment vote reads %d, err %v; want +1", v, err)
	}
	if v, _ := f.store.CommentVote(f.ctx, commentID, uuid.New()); v != VoteNone {
		t.Errorf("someone who has not voted reads %d; want 0", v)
	}
}

func TestVotingOnYourOwnCommentAndOnNothing(t *testing.T) {
	f := newFixture(t)
	commenter := f.user(t, uuid.New())
	var commentID uuid.UUID
	if err := f.pool.QueryRow(f.ctx,
		`INSERT INTO comments (article_id, user_id, body) VALUES ($1,$2,'x') RETURNING id`,
		f.article, commenter).Scan(&commentID); err != nil {
		t.Fatalf("insert comment: %v", err)
	}

	if _, err := f.store.VoteComment(f.ctx, commentID, commenter, VoteUp); !errors.Is(err, ErrOwnComment) {
		t.Errorf("voting on one's own comment returned %v, want ErrOwnComment", err)
	}
	if _, err := f.store.VoteComment(f.ctx, uuid.New(), commenter, VoteUp); !errors.Is(err, ErrNotFound) {
		t.Errorf("voting on a comment that does not exist returned %v, want ErrNotFound", err)
	}
}
