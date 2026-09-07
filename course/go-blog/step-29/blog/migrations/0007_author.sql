-- Every article gets an owner. The column is nullable on purpose: a migration
-- has to decide what happens to the rows that are already there, and there is
-- no honest owner to invent for them. An article with no owner belongs to
-- nobody, and nobody can edit it -- which is the safe answer.
alter table articles add column author_id integer references users (id);

create index articles_author on articles (author_id);
