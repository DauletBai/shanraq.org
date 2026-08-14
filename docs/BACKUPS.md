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
