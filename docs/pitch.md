# Parity Labs — the trust layer for AI agents in production

## The problem

AI agents are increasingly trusted to write code, clean data, retrain
models, and triage incidents — mostly unsupervised, because reviewing
everything by hand erases the productivity gain. Every team that adopts
this hits the same wall within months: **the agent is useful, but you
can't trust it unsupervised, and there's no general infrastructure for
verifying an agent's action before it's allowed to act.** Today, teams
solve this with ad-hoc glue code per use case, if they solve it at all.

## Why we can build this

Three independent prototypes already answer three different slices of
the trust problem, built and tested before this pitch existed:

| Slice of the problem | Prototype | Status |
|---|---|---|
| "Is this *change* safe?" | [`parity`](https://github.com/Jex2l/parity) | **Built and validated on a real CI run.** Replays real production traffic against a candidate build and reports what actually changed in behavior. |
| "Is this *claim* true?" | `meridian-legal-intelligence` | **Built, 60/60 tests passing.** Every generated answer must cite evidence actually retrieved for the query; unsupported claims are rejected, not shown. |
| "Is this *action* authorized?" | `incident-agent` | **Built, 69/69 tests passing.** Read-only investigation is unrestricted; writes require explicit human approval through an infrastructure-enforced gate. |

This isn't a pitch for something we think we could build — it's a pitch
for productizing three things that already work, proven independently,
for the same underlying problem in three different domains.

## Product: three pillars

1. **Replay verification** (pillar 1, shipped) — before any agent
   change ships, replay it against real production traffic and show
   what actually changed. Sell this alone first: "CI for AI-written
   code." Zero dependency on the other two pillars.
2. **Grounded reasoning** (pillar 2, in progress) — every claim an
   agent makes must cite the evidence it retrieved, with a faithfulness
   check, or say "I don't know." Sold as an SDK any agent builder drops
   in. First integration: a grounded implementation of pillar 1's
   `classify.Classifier` interface.
3. **Approval-gated execution** (pillar 3, roadmap) — read freely,
   write only with a human click through an infrastructure-enforced
   gate, full audit trail. The compliance story for regulated verticals.

## Go-to-market

Land with pillar 1 alone — it's useful standalone, cheapest to try
(one GitHub Action, no infra to stand up), and the fastest path to a
real customer logo. Expand into the two highest-pain verticals already
proven as prototypes: on-call/SRE (`incident-agent`'s exact use case)
and data/ML pipelines. See `docs/gtm.md` for the full wedge-to-expansion
path and pricing.

## Why now

AI agents went from toys to writing/shipping production changes in the
last 18 months. The "move fast, trust the agent" honeymoon is ending —
every engineering org is about to independently rediscover it needs
exactly this, which is why working prototypes of all three pillars
already existed before this company did.

## The ask

Not raising yet. This repo and `jex2l/parity` are the proof this is
buildable; `docs/roadmap.md` is the plan for turning pillar 1's working
wedge into pillar 2's grounded classifier, then a real first customer
conversation.
