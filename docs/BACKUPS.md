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

Repeat this after any change to the schema or the backup script. Since 19 August
it is one command rather than five remembered ones — `scripts/backup-restore-test.sh`,
run on the machine that holds the private key, because the key is deliberately
not on the server:

```sh
BACKUP_AGE_IDENTITY=~/keys/shanraq-backup.age-key \
  make restore-test ARCHIVE=~/backups/shanraq-backup-20260819-023000Z.tar.gz.age
```

It decrypts, verifies `SHA256SUMS`, restores into a scratch database with
`--exit-on-error`, prints what came back, and drops the scratch database and the
decrypted copies however it ends — including the failure paths, which is when
personal data is easiest to leave lying around.

It fails, loudly and non-zero, on a damaged archive, a wrong key, a partial
restore, a schema of fewer than twenty tables, an empty `auth_users`, no
published articles, missing migration history, or restored data holding **no
administrator** — a restore that has to be repaired by hand before the site can
be run is not a restore.

One thing it learned the hard way and now handles: on a machine with both
Postgres 14 and 16 installed, `pg_restore` from `PATH` is the older one, and it
rejects a version-16 dump with a file-header error that reads exactly like
corruption. The script picks the matching tools the way `backup.sh` does.

## The gap, and what now covers half of it

Since 18 August a second timer sends the published content to Cloudflare R2:
`shanraq-content.timer`, daily at 03:00 UTC, half an hour after the encrypted
backup so the two never compete for the database.

`scripts/content-export.sh` runs the exporter that ships inside the app image
and uploads the archive with rclone. What travels: 108 published articles with
all 324 translations, the prediction ledger, the 21 editable pages, and the 18
cover images that were uploaded rather than committed — those exist on this disk
and nowhere else, while the other ninety live under `/static/covers` in the
repository and are already off-site on GitHub.

What does not travel is a list rather than a filter, so a table added next year
cannot leak into it by being forgotten: accounts, sessions, listings (sellers'
telephone numbers), comments, votes, moderation, analytics, payments. Bylines
are included because a mirror without them is not the same publication, and the
export refuses to run when it meets an author it does not know.

The archive is not encrypted, deliberately: everything in it is already public
at shanraq.org, and an unencrypted copy is one a mirror can serve and a stranger
can restore without holding our key. It is about 6.3 MB, against R2's 10 GB free
allowance.

Two things were learned setting it up and are worth not rediscovering:

**The R2 hostname resolves to IPv6 first.** The API token is restricted to this
server's IPv4, so every request was refused with 403 until rclone was pinned
with `--bind 85.202.192.61`. The alternative is to add the v6 address to the
token's allow-list; pinning keeps the list at one address.

**Listing buckets is denied by design.** The token carries Object Read & Write
on one bucket, which does not include the account-level operation of listing
buckets. `rclone lsd r2:` failing is correct; `rclone ls r2:shanraq-content/`
is the check that means something.

Verified on 18 August by the only test that counts: the uploaded archive was
downloaded back from R2, its SHA-256 compared against the copy on the server,
extracted, and its manifest read — 108 articles, 3 predictions, 21 pages, 18
covers.

## The gap that is still open

**The full backup is still on the machine it protects.** Today's backups defend against
the likely losses — a bad migration, a mistaken `DELETE`, a corrupted table.
They do not defend against losing the machine, which is the failure the outage
was rehearsing.

Closing it is one line, once a destination exists:

```
BACKUP_UPLOAD_CMD=<command that ships {file} off the box>
```

`{file}` is replaced with the archive path; a non-zero exit fails the whole run
loudly instead of leaving a silent gap.

Until then every run says so in its log, because a gap that nothing mentions is
one you stop seeing:

```
WARNING: no BACKUP_UPLOAD_CMD — this archive stays on the machine it is
         protecting, so it survives a bad migration but not a lost host
```

And once the destination is configured, set `BACKUP_REQUIRE_OFFSITE=1` beside
it. A hand-edited `.env` is how a working off-site copy quietly turns back into
a backup sitting on the machine it protects; with that flag the run fails
instead.

**There is still no external uptime check.** `/healthz` and `/readyz` answer
correctly, and nothing outside the server asks them — so the first notice of an
outage is somebody opening the site. That wants a third-party prober and a
Telegram alert, and it wants to live somewhere the server cannot take down with
it.

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

## Deferred: Cloudflare proxying

Postponed on 14 August 2026 — the owner's call, and a defensible one: at a few
hundred readers a week the proxy buys little and costs operational complexity.
Revisit when organic growth or revenue makes the site worth attacking or
blocking. Nothing is half-applied; every record is `DNS only` and the origin is
reachable exactly as before.

Two things were established while preparing it, and are worth keeping so nobody
has to derive them again.

**The origin is ready for `Full (strict)`.** Checked 14 August: addressed by IP
with the right SNI it answers 200 with a valid Let's Encrypt chain (verify
result 0), `www` likewise, HTTP/2 on. No Origin CA certificate is needed.

**`ufw` will not do the job.** Docker publishes 80 and 443 through its own
netfilter path, so `ufw` rules never see that traffic and a "block everything
but Cloudflare" written there would appear to work and do nothing. The rules
belong in the `DOCKER-USER` chain, which is currently empty, and must be made
persistent — they do not survive a reboot on their own.

The order, when the time comes. Each step depends on the one before it:

1. SSL/TLS → **Full (strict)**. Before anything else: proxying while the mode is
   `Flexible` puts Cloudflare on HTTP to an origin that redirects to HTTPS, and
   the site dies in a redirect loop immediately.
2. Proxy **only** the apex and `www`. `mail` stays grey forever — Cloudflare
   does not proxy SMTP, and an orange cloud there points MX at HTTP addresses.
3. Verify from outside that the origin address is no longer visible and that
   client addresses still resolve correctly.
4. Apex SPF: `+a` → `ip4:85.202.192.61`. Not urgent — mail leaves through
   Resend, whose envelope sender is on `send.shanraq.org` and whose alignment
   comes from the `resend._domainkey` signature, so `+a` carries no weight
   today. But once the A record is Cloudflare's, `+a` hands every Cloudflare
   address permission to send as this domain.
5. Uncomment Cloudflare's ranges in `configs/config.prod.yaml` **and** firewall
   the origin in the same change. Either alone is worse than neither: trusting
   those ranges while the origin is directly reachable is the spoofing hole the
   firewall exists to close.

**Always Online is weaker than it sounds.** It serves from the Internet Archive,
only when Cloudflare cannot reach the origin at all, and only for pages the
archive happened to crawl — which on a young site may be a handful. It is worth
switching on because it is free, but it is not an answer to an outage. The full
archive staying readable is the static mirror's job, and that job is still open.

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
