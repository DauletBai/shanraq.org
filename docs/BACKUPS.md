# Backups and resilience

Written 14 August 2026, the morning after a 4-hour 29-minute outage in which
this project came within one failed disk of losing everything it had, because
no backup existed at all.

## What runs now

Two systemd timers on the production host. The unit files live in
[`deploy/systemd/`](../deploy/systemd/) and are copied to
`/etc/systemd/system/`, so a rebuilt machine can be brought back to this state
from the repository rather than from memory.

| Timer | When | What |
|---|---|---|
| `shanraq-backup.timer` | daily, 02:30 UTC (07:30 Astana) | encrypted dump of the database and the uploaded media |
| `docker-cache-prune.timer` | Sundays, 03:30 UTC | drops Docker build cache older than 72 hours |

Both are `Persistent=true`: a machine that was down at the scheduled minute
runs the job when it comes back. On 13 August that was not hypothetical.

### The backup

`scripts/backup.sh`, driven by `BACKUP_*` variables in `/opt/shanraq/.env`:

```
BACKUP_COMPOSE_FILE=/opt/shanraq/docker-compose.prod.yml   # pg_dump inside the db container
BACKUP_MEDIA_VOLUME=shanraq_media-data                     # uploaded images and avatars
BACKUP_DIR=/opt/shanraq/backups
BACKUP_RETENTION=14
BACKUP_AGE_RECIPIENT=age1…                                 # public half only
```

One archive per run: `db.dump` (custom format, restores with `pg_restore`),
`media.tar.gz`, and `SHA256SUMS` over both — then the whole thing encrypted
with `age`. About 12 MB today.

### The key

`age` is asymmetric, and only the **public** half is on the server. That is
deliberate: whoever takes the machine would otherwise take the backups and the
key to read them in the same breath.

The private key is the owner's, kept off the server. **Without it the backups
are unreadable — there is no recovery path, no reset, nobody to ask.** It
belongs in a password manager and in one offline place, not in this repository
and not on the server.

## Verified, not assumed

A backup nobody has restored is a hope, not a backup. On 14 August the first
archive was taken through the whole chain:

1. Written and encrypted on the server.
2. Downloaded, decrypted with the private key held on the owner's machine.
3. `SHA256SUMS` verified — both members clean.
4. `pg_restore` into a scratch database on a real Postgres 16.

What came back: 47 tables, 108 published articles, **324 translations** —
exactly 108 × 3 languages — 6 users, 4 listings, 21 comments, goose at
`20251108001300`. Identical to production. The scratch database was dropped and
the decrypted copies deleted; they hold personal data and must not linger.

Repeat this after any change to the schema or the backup script. The command
that matters:

```sh
age -d -i <private-key> -o backup.tar.gz shanraq-backup-*.age
tar xzf backup.tar.gz && shasum -a 256 -c SHA256SUMS
pg_restore -d <scratch-db> --no-owner db.dump
```

## The gap that is still open

**Every copy is on the machine it protects.** Today's backups defend against
the likely losses — a bad migration, a mistaken `DELETE`, a corrupted table.
They do not defend against losing the machine, which is the failure the outage
was rehearsing.

Closing it is one line, once a destination exists:

```
BACKUP_UPLOAD_CMD=<command that ships {file} off the box>
```

`{file}` is replaced with the archive path; a non-zero exit fails the whole run
loudly instead of leaving a silent gap.

### Where the copies may legally live

The Law on Personal Data, Article 12(2), requires personal data to be stored
in a database located in Kazakhstan. The archive holds accounts, e-mail
addresses, sellers' phone numbers and avatars, so the split is not a matter of
taste:

| Copy | Contains personal data | May live |
|---|---|---|
| Full archive | yes | second location **inside Kazakhstan** — different provider, different city |
| Content-only export | no | anywhere, including abroad |

A content-only export — articles, translations, predictions, editable pages,
cover images — carries nothing personal and can go to free object storage
outside the country. It is also exactly the material a read-only mirror would
need, so the two jobs are one job. Not built yet.

## Where the single points of failure now stand

| Layer | Before 13 August | Now |
|---|---|---|
| DNS | all three nameservers at the hosting provider | Cloudflare, independent |
| Domain registration | PS.KZ | PS.KZ — unchanged |
| Server | one VPS, Kosshy | one VPS, Kosshy — unchanged |
| Backups | none | daily, encrypted, verified — but on the same machine |
| Monitoring | none — the owner learned from readers | none |

The DNS move is the one that changed the shape of the problem: repointing the
domain at another host is now minutes of work and needs nobody's permission.

## What is worth doing next, in order

1. **A copy that leaves the machine.** Everything else is second.
2. **Uptime monitoring** to an address that is not on this server, so an outage
   is reported by an instrument rather than by a reader.
3. **A content-only mirror abroad**, which also answers the blocking scenario a
   second Kazakh location cannot.
4. **A standby host** at another Kazakh provider, cold, restored from the
   backup.
5. **Move the registrar** away from the hosting provider, so one company does
   not hold both handles.
