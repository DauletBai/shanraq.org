# Admin access

There is **no standing administrator account** in the database and **no way to
become one from the web**. Nothing on the public surface can mint an
administrator — no bootstrap page, no "first user becomes admin" rule, no
secret URL. An endpoint that can create an admin is an endpoint an attacker can
try, so the capability lives outside the web application entirely.

Access is provisioned with `adminctl`, run by an operator on the server, where
the database already trusts them. This is the same shape as Django's
`createsuperuser` or `rails console` — the standard pattern for this problem.

## Creating the first administrator

Over SSH on the server (or `docker compose exec app ...`):

```sh
DATABASE_URL="postgres://…" adminctl create \
  -email you@example.com -first Имя -last Фамилия [-middle Отчество]
```

`adminctl` prompts for the password twice, without echo. It is never passed as
a flag: flags land in shell history and are visible to any user who can run
`ps`. For unattended provisioning set `ADMIN_PASSWORD` in the environment
instead, and unset it afterwards.

Every field is validated with the same rules the site uses: a real name
(letters only, capitalized) and a password of at least 8 characters containing
both letters and digits.

## Other commands

```sh
adminctl promote -email someone@example.com [-role admin|director]   # existing account
adminctl list                                                        # who has staff access
```

`promote` is the right tool when the person already registered normally — it
raises an existing account rather than creating a second one.

## Roles

| Role | Can |
|---|---|
| `admin`, `director` | everything, including managing users and service switches |
| `manager` | dashboard + finance |
| `editor` | moderation |

`admin` is the superuser role and is intentionally not assignable from the
admin panel's role form — only from `adminctl`. That keeps the highest
privilege off the web surface.

## Recommended hardening before the public launch

These are not implemented yet; they are the next steps, in order of value:

1. **TOTP for staff.** The `auth_mfa_totp` table already exists. Requiring a
   second factor for `admin`/`director` is the single biggest win — a leaked
   admin password stops being enough.
2. **Restrict `/admin` at the edge.** An IP allowlist or a VPN/Tailscale-only
   route in the reverse proxy, so the panel is not reachable from the open
   internet at all.
3. **Separate the session.** A shorter session lifetime for staff than for
   readers, and re-authentication before destructive actions.
4. **Audit trail.** Moderation decisions are already logged; extend it to role
   changes and service-switch flips (who, what, when).

Provisioning an admin is an operational event: it should appear in the server's
shell history and deploy log, and be done deliberately — not be something the
application can do to itself.
