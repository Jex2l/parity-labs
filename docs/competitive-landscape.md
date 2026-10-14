# Competitive landscape

**Caveat up front:** this is directional analysis from general category
knowledge, not primary research (no customer interviews, no vendor
comparison calls done). Treat it as a starting map to refine, not a
researched conclusion — specific vendor names and feature claims below
should be verified before being used in anything customer-facing.

## Adjacent categories, and why pillar 1 is a different bet

- **Contract-testing tools** (the Pact-style category: consumer-driven
  API contracts) — verify a predefined contract, not real production
  behavior. Pillar 1's bet is that AI-authored changes break things
  nobody thought to write a contract for; replaying *actual* traffic
  catches that class of regression by construction, contract tools
  structurally can't.
- **LLM-as-judge code review tools** (various AI PR-review bots) —
  evaluate a diff's *code*, not the service's *runtime behavior*. A
  change can look fine to a code reviewer (human or LLM) and still
  behave differently under real traffic; pillar 1 is behavior-first,
  code-review-agnostic.
- **General APM/observability vendors** — catch a regression *after*
  it's in production, via dashboards and alerts a human has to notice.
  Pillar 1 catches it *before merge*, as a CI gate, with no human
  needing to notice anything for the common case.
- **Traditional snapshot/golden-file testing** — requires someone to
  have written the test and the golden file in the first place. Pillar
  1's traffic is real and ambient; no one has to anticipate the case.

## Where real competition is most likely to emerge

The most plausible near-term competitor isn't a dedicated "replay
verification" company — it's the AI coding-agent platforms themselves
(Copilot, Cursor, Devin-style tools) building a thinner version of this
directly into their own PR workflow, as a feature rather than a
separate product. That's the actual strategic risk to watch, not a
point-solution startup in the same category.

## Why pillar 1 is still worth building despite that risk

A platform-native version would be scoped to that platform's own
agent's changes. Pillar 1 is platform-agnostic by design (it replays
*any* candidate build against *any* captured traffic, regardless of
which agent — or human — wrote the change) — the defensible position is
being the neutral verification layer across agent vendors, not
competing to be a better agent vendor.

## Unknowns this doc doesn't resolve

- Whether any well-funded company has already announced something in
  this exact space — not checked as part of this pass.
- Real pricing data for adjacent categories (contract testing, AI code
  review) — the pricing direction in `docs/gtm.md` is a starting
  hypothesis, not benchmarked against real comparables yet.
