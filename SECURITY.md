# Security Policy

We take the security of Shanraq seriously — resilience and user safety are core
to the project.

## Reporting a vulnerability

**Please do not open a public GitHub issue for security problems.**

Report vulnerabilities privately by either:

- Emailing **shanirak.org@gmail.com** with the details, or
- Using GitHub's **[private vulnerability reporting](https://github.com/DauletBai/shanraq.org/security/advisories/new)**
  (Security → Report a vulnerability).

Please include:

- a description of the issue and its impact,
- steps to reproduce (proof-of-concept if possible),
- affected version/commit and environment.

## What to expect

- **Acknowledgement** within 72 hours.
- An initial assessment and, where confirmed, a remediation plan.
- Credit for responsible disclosure if you would like it.

We ask that you give us reasonable time to fix an issue before any public
disclosure, and that you avoid privacy violations, data destruction, or service
disruption while researching.

## Scope

In scope: this repository and the production service at `shanraq.org`.
Out of scope: third-party services (hosting, DNS, email provider), volumetric
denial-of-service, and social-engineering of staff.

## Good to know

- Secrets are never committed — configuration is env-first and `.env` is
  git-ignored. If you believe a secret was exposed in git history, report it
  privately as above.
