create table sessions (
    token_hash text primary key,
    user_id    integer not null references users (id) on delete cascade,
    expires_at text not null
) strict;

create index sessions_user on sessions (user_id);
