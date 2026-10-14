# pillar2-grounded-classifier

Pillar 2 of the Parity Labs platform (see `../docs/pitch.md` and
`../docs/roadmap.md`): "is this claim true?" applied to pillar 1's
cluster classifications.

**Status update: wired into production, but as a port, not this code
directly.** This package's logic now runs live in `jex2l/parity` as
`internal/classify/grounded` — confirmed on a real PR
(`jex2l/parity-client-test#6`): 80% confidence with a cited diff line
for a genuine field addition, vs. the structural heuristic's 55% with
no citation. It's a reimplementation, not an import of this code (see
"Not wired into `parity` at all" below, which this update doesn't
delete — the reasoning there is still exactly why a straight import was
never possible). The production version also fixed a real gap this
prototype had: its original field matching was exact-substring, which
would have missed most real Go APIs (a JSON tag like `in_stock` rarely
appears literally in Go source using the field name `InStock`) — the
production version matches identifier tokens case/underscore-
insensitively instead. That fix has not been backported to this
standalone prototype; this file (below) still describes the original,
narrower matching this specific code does.

## What's built

`grounded.Classifier` — a real, tested classification engine that
grounds its verdict in the PR's actual diff text instead of structural
heuristics alone. For each cluster's changed JSON paths, it searches the
PR's diff hunks for a line that actually mentions the field:

- Every changed path has a citation → `LIKELY_INTENDED`, with the
  citing diff line attached as evidence.
- No changed path has a citation → `UNCERTAIN`, explicitly naming which
  paths it couldn't ground — not a guess dressed up as a verdict.
- Some but not all → `UNCERTAIN`, naming both what was found and what's
  still missing.
- Critical severity still always wins (matches
  `classify.Heuristic`'s rule 1: nobody intentionally ships a crash,
  regardless of what the diff says) — but even that path attaches
  whatever citations it can find, for a human to double-check.

**This does not call an LLM.** The grounding is deterministic
substring/word-boundary matching against diff text — real and testable
without a model endpoint, same discipline `internal/classify.Heuristic`
already applies in the `parity` repo: a structural/deterministic
baseline has to exist and be tested before an LLM-backed version is
worth building on top of it.

6 tests, including one built directly from the real diff in
[jex2l/parity PR #1](https://github.com/Jex2l/parity/pull/1)'s
self-test (the `currency` field added to `/stats/revenue`) — and one
real bug the tests caught and fixed before this was committed:
`leafField` was stripping everything after the *first* `[` in a path,
so `$.orders[3].status` wrongly reduced to `orders` instead of
`status`. Fixed by extracting the last dot-separated segment first,
then stripping any array index from just that segment.

## What's explicitly not built

- **No LLM-backed semantic matching.** A renamed field (`currency` →
  `currency_code`) won't be grounded by this version even if the diff
  clearly shows the rename — only exact/substring field-name matches
  count as evidence today. This is the actual next increment, not a
  hidden limitation: wire a model call that can judge *semantic*
  relatedness between a changed path and a diff hunk when the
  deterministic matcher finds nothing, falling back to `UNCERTAIN`
  exactly as today whenever the model also can't find support (same
  "say I don't know" discipline, just with a smarter first pass).
- **No faithfulness check yet.** `meridian-legal-intelligence` reuses
  its retrieval reranker as a cheap entailment proxy to catch citations
  that exist but don't actually support the claim attached to them.
  Nothing equivalent exists here yet — today, *any* line containing the
  leaf field name counts as a citation, even one that doesn't actually
  explain the behavioral change (e.g. a comment mentioning the field
  name, not the code that changes it).
- **Not wired into `parity` at all.** This is the integration gap worth
  naming precisely: `parity`'s `classify.Classifier` interface,
  `report.ClusterView`, and `classify.PRContext`/`Result` types all live
  under that repo's `internal/` tree, which Go's package-visibility
  rules mean no external module — including this one — can import.
  `grounded.ClusterView`/`Result` above deliberately mirror those types'
  shape so an adapter is mechanical once this gap is closed, but closing
  it requires one of: (a) moving `parity`'s relevant types to a public
  `pkg/` path, (b) vendoring this package's logic directly inside the
  `parity` module as a new `internal/classify/grounded` implementation
  of the existing `Classifier` interface, or (c) running this as a
  standalone HTTP service `parity-ci` calls out to. None of these three
  has been decided or started.

## Running the tests

```bash
go test ./... -v
```

No external dependencies, no network calls, no model endpoint required
— matches the project-wide pattern of not needing infrastructure that
isn't actually available to run the test suite.
