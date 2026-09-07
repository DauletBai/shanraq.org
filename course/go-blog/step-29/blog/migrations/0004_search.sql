-- The search table holds no text of its own: content= points it at articles,
-- so an article is stored once and indexed once.
create virtual table articles_fts using fts5(
    title, body, content='articles', content_rowid='id'
);

-- Whatever is already written has to get into the index too: a migration
-- that only starts indexing from today leaves the old articles unfindable.
insert into articles_fts (rowid, title, body)
select id, title, body from articles;

create trigger articles_ai after insert on articles begin
    insert into articles_fts (rowid, title, body)
    values (new.id, new.title, new.body);
end;

-- A delete is written as an insert of the 'delete' command together with the
-- OLD values: the index needs them to find what to take out.
create trigger articles_ad after delete on articles begin
    insert into articles_fts (articles_fts, rowid, title, body)
    values ('delete', old.id, old.title, old.body);
end;

create trigger articles_au after update on articles begin
    insert into articles_fts (articles_fts, rowid, title, body)
    values ('delete', old.id, old.title, old.body);
    insert into articles_fts (rowid, title, body)
    values (new.id, new.title, new.body);
end;
