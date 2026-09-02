-- +goose Up
-- What a reader has actually done, lesson by lesson.
--
-- The course page used to promise progress marks it had no way to keep. The
-- mark is deliberately not "scrolled to the bottom": a page can be flicked
-- through in a second, and a course that congratulates you for that teaches you
-- to flick. It is earned by submitting a solution to the lesson's exercise and
-- having it come back accepted.
--
-- The last review is stored beside the mark so the reader finds it again after
-- a reload, and so an attempt that failed still leaves something to work from.
CREATE TABLE IF NOT EXISTS course_progress (
    user_id    uuid        NOT NULL REFERENCES auth_users (id) ON DELETE CASCADE,
    article_id uuid        NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
    passed     boolean     NOT NULL DEFAULT false,
    attempts   int         NOT NULL DEFAULT 0,
    note       text        NOT NULL DEFAULT '',
    solution   text        NOT NULL DEFAULT '',
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, article_id)
);

-- The reader's own page asks "which lessons has this person passed"; the index
-- follows that question rather than the primary key's order.
CREATE INDEX IF NOT EXISTS idx_course_progress_user ON course_progress (user_id, passed);

-- +goose Down
DROP TABLE IF EXISTS course_progress;
