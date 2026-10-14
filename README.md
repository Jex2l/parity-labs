# Parity Labs

**The trust layer for AI agents in production.**

Parity Labs is the company wrapper around an idea proven across three
working prototypes (see `docs/pitch.md` for the full case): as AI agents
get more autonomy — writing code, cleaning data, retraining models,
triaging incidents — the one unsolved problem every team independently
hits is *how do we know an agent's proposed action is safe, and its
reasoning is actually true, before it's allowed to act?*

This repo is the company's non-engine material: the pitch, the roadmap,
the go-to-market plan, and the marketing site. The actual product
engineering lives in two places, not here:

- **[`jex2l/parity`](https://github.com/Jex2l/parity)** — pillar 1
  (replay verification), fully working: `parity replay` + `cmd/parity-ci`
  + `action.yml`, validated end-to-end on a real GitHub Actions run. This
  is the product you can use today.
- **`pillar2-grounded-classifier/`** (this repo) — pillar 2 (grounded
  reasoning) is just starting here: a `classify.Classifier`
  implementation for the `parity` engine that grounds its verdict in
  retrieved evidence instead of structural heuristics alone. See that
  directory's own README for status.

Pillar 3 (approval-gated execution, generalizing `incident-agent`'s
pattern) has no code yet anywhere — it's roadmap, documented honestly as
such in `docs/roadmap.md`.

## What's in this repo

```
docs/
  pitch.md                   the one-pager: problem, product, GTM, ask
  roadmap.md                 pillar sequencing, what's built vs. not
  gtm.md                     wedge -> expansion path, pricing
  competitive-landscape.md   who else is circling this problem, and why
                              pillar 1's replay-verification wedge is
                              defensible against each
site/
  index.html                 the marketing site (static, no build step)
pillar2-grounded-classifier/
  (see its own README)       early-stage: grounding SDK for pillar 2
```

## Status

Pillar 1 is real and shipped. Everything else in this repo is either a
planning document (clearly labeled as such) or early-stage code
(`pillar2-grounded-classifier/`, which has a working skeleton and one
real capability — see its README for exactly what's tested vs. not).
Nothing here claims to be more finished than it is.
