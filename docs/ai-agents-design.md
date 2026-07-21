# AI agents: night desks and service agents

Decision of 21 July 2026. This document is the contract; the code follows it.

## The rule that shapes everything

**Names are for functions, not for opinions.**

One editorial voice — AI Dake — carries every column and every judgment. Three
further agents exist, but none of them has a byline, because none of them
makes claims about the world. They perform procedures a user can appeal.

The reason is not aesthetics. Four personas would be one model wearing four
names, and a reader seeing four columnists infers four minds: if three of them
agree, that reads as independent corroboration. There would be no
corroboration. Manufacturing the appearance of plurality is the one thing a
platform whose whole premise is honesty about AI cannot afford.

The secondary reason is practical: the same model under different briefs
contradicts itself, and there is no editorial mind to reconcile the versions.
Readers notice before the publisher does.

## Night pipeline — four desks, one byline

Four parallel sweeps, then verification, then synthesis:

| Desk | Scope |
|---|---|
| 1 | Politics and economy — they move each other, so one desk holds both |
| 2 | World and security — conflicts, energy routes, sanctions |
| 3 | Technology and society — AI, regulation, how people live and work |
| 4 | Kazakhstan's regions — what the capital-based press misses |

**The verification stage is not optional and is not a formality.** Its job is
to attack the sweep's output: find the figure that comes from an interested
party, the number that only one aggregator carries, the "agreement" that is a
memorandum. This session produced the evidence for insisting on it — a
settlement whose population Wikidata gave as 12,228 and district totals as
3,603; "24 intergovernmental documents" that belonged to a different summit in
a different year; fourteen economic figures that could not be verified and were
therefore not printed.

A second verifier is worth more than a second byline.

## What may publish unattended

**Only fact digests with sources.** What happened, when, the numbers, the
links. These go out in the morning under a neutral desk byline.

**Every judgment, comparison and conclusion goes to drafts** and waits for a
human. One fabricated figure published at six in the morning under an AI byline
costs more than a month of night work, because it destroys the only thing the
AI column has: that its numbers can be checked.

The digest must carry a "could not verify" block. Publishing the gaps is what
makes the rest credible.

## Service agents — the other three names

Дәке, Бәке, Мәке, Жәке are respectful forms of address. The set is coherent, so
keep it — but for services, where a name genuinely helps a user understand who
acted:

| Agent | Function |
|---|---|
| AI Bake | Moderation — comments, articles, listings. Every decision logged and appealable. |
| AI Make | Support — questions about rules, publication, listings. Escalates to a human at the edge of its competence. |
| AI Jake | Listing assistant — helps format an ad, checks photos and completeness. Never invents facts about a property. |

These are **not** registered in `aiAgentAuthors`: they publish nothing, so they
need no byline and no revenue treatment. They need logs, appeal routes, and a
visible statement that a machine acted.

Priority note: the production-readiness review flagged moderation as an
unclosed hole and a launch blocker. The moderation agent is worth more than a
fourth columnist.

## Cost architecture

Running four research sweeps a night at frontier-model quality costs roughly
$110/month for collection alone — over half the $200 budget, before any writing.

Split the work by model instead of by agent:

- **Sweep and triage**: a cheap fast model. High volume, low judgment.
- **Verification**: mid-tier. This stage earns its cost.
- **Synthesis of one or two pieces**: the strong model, and only here.

That lands near $40–60/month. The saving comes from matching the model to the
task, not from cutting the number of desks.

Hard requirements: a token ceiling per night, and fail-closed behaviour — if the
budget or a source is unavailable, publish nothing rather than publish thin.

## Open question

Whether readers follow bylines or rubrics is unknown for this audience. The
argument for four columnists is that people subscribe to people. I believe a
rubric subscription serves the same need without the honesty cost, but that is
a prediction, not a measurement. Revisit once there is real subscription data.
