package blog

import (
	"errors"
	"fmt"
)

// ErrNotFound is returned when nothing is filed under that address. The caller
// compares against it with errors.Is rather than reading the message.
var ErrNotFound = errors.New("мақала табылмады")

// Store keeps the articles. The fields are lower case on purpose: what leaves
// this package is the ability, not the arrangement, so the map could become a
// database tomorrow without anything outside noticing.
type Store struct {
	order []string
	items map[string]Article
}

func NewStore() *Store {
	return &Store{items: map[string]Article{}}
}

// Add files an article and remembers where it goes in the running order.
func (s *Store) Add(a Article) {
	if _, seen := s.items[a.Slug]; !seen {
		s.order = append(s.order, a.Slug)
	}
	s.items[a.Slug] = a
}

func (s *Store) Get(slug string) (Article, error) {
	a, ok := s.items[slug]
	if !ok {
		return Article{}, fmt.Errorf("get %q: %w", slug, ErrNotFound)
	}
	return a, nil
}

// All hands back the articles in the order they were added.
func (s *Store) All() []Article {
	out := make([]Article, 0, len(s.order))
	for _, slug := range s.order {
		out = append(out, s.items[slug])
	}
	return out
}

func (s *Store) Len() int { return len(s.items) }
