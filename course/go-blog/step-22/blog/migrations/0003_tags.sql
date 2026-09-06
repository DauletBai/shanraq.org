create table tags (
    id   integer primary key,
    name text not null unique
) strict;

create table article_tags (
    article_id integer not null references articles (id) on delete cascade,
    tag_id     integer not null references tags (id),
    primary key (article_id, tag_id)
) strict;
