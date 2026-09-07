create table articles (
    id    integer primary key,
    slug  text    not null unique,
    title text    not null,
    words integer not null default 0,
    lang  text    not null default 'kz',
    body  text    not null default ''
) strict;
