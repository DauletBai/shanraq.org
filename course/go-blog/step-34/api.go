package main

import (
	"encoding/json"
	"net/http"
	"time"

	"myblog/blog"
)

// apiArticle is what leaves the blog. It is a type of its own rather than
// blog.Article with tags on it: the shape of the answer is a promise to
// whoever reads the API, and it should not change every time a column is
// added to the database.
type apiArticle struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Lang        string    `json:"lang"`
	ReadingTime int       `json:"reading_time"`
	Tags        []string  `json:"tags,omitempty"`
	Cover       string    `json:"cover,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitzero"`
}

func newAPIArticle(a blog.Article, tags []string) apiArticle {
	return apiArticle{
		Slug: a.Slug, Title: a.Title, Lang: a.Lang,
		ReadingTime: a.ReadingTime(), Tags: tags,
		Cover: a.Cover, UpdatedAt: a.UpdatedAt,
	}
}

// writeJSON sets the header before the status: WriteHeader sends the status
// and every header at once, and after it a header cannot be changed.
func writeJSON(w http.ResponseWriter, r *http.Request, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// The status is already on its way; all that is left is the log.
		logger.Error("json", "id", reqID(w), "жол", r.URL.Path, "қате", err)
	}
}

// apiError keeps one shape for every failure. Whoever reads the API parses
// the answer, so an HTML page in place of JSON breaks them at the first byte.
func apiError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	writeJSON(w, r, code, map[string]string{"error": msg})
}
