# Go-to-market

## Wedge: pillar 1 alone, sold as "CI for AI-written code"

The cheapest possible first sale: one GitHub Action
(`uses: jex2l/parity@<tag>`), no infra to stand up beyond a capture
segment the team likely already needs for other reasons, a PR comment
and check run as the entire UX. No pillar 2 or 3 dependency.

**Target buyer:** an engineering team shipping AI-agent-authored PRs
(Copilot/Cursor/Devin-style workflows) who's already felt the pain of
reviewing agent diffs by hand, or worse, not reviewing them closely
enough.

**Proof point to lead with:** the action was validated on a real GitHub
Actions run that found and fixed 3 real bugs in the action itself —
concrete, verifiable, not a claim. Screenshot the actual PR comment
output (a real classified cluster, not a mockup) in any outbound.

## Expansion: the two proven verticals

1. **On-call/SRE agents** — `incident-agent`'s exact use case, already
   built and tested (69/69). Highest pain, highest willingness to pay,
   easiest to demo with a real MTTR number from the existing eval
   harness (`incident-agent`'s README documents a real 0%→100% fix on
   its own eval set — that's the demo).
2. **Data/ML pipeline agents** — `auto-data-pipeline-ai` +
   `NYC-Taxi-Demand-Prediction`'s MLOps lifecycle as the reference
   case study: an agent that cleans data, retrains, and redeploys a
   model, with every retrain replay-verified (pillar 1) before it goes
   live.

Both expansions reuse pillar 1's engine unchanged — they're new
*applications* of replay verification, not new core product.

## Pricing (direction, not final)

Usage-based: per verified agent-action (one replay run, one grounded
query, or one approval-gated execution), tiered by volume. Land with a
single team's coding-agent CI integration — cheapest to try, fastest to
show value — expand into SRE and MLOps once trust is established, then
an enterprise audit/compliance tier for regulated industries once
pillar 2's grounding and pillar 3's approval gate exist and there's a
customer who needs the audit trail specifically.

## What "done" looks like for the wedge, concretely

Not a revenue number — a repeatable motion: a team installs the Action,
it flags a real unintended diff within the first week that a human
review would plausibly have missed, and the team keeps it installed
past the first month without anyone championing it internally. That's
the signal pillar 1 is pulling its own weight before spending effort on
pillar 2/3's GTM.
