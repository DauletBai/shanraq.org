create table users (
    id            integer primary key,
    email         text not null unique,
    password_hash text not null,
    created_at    text not null default (datetime('now'))
) strict;
