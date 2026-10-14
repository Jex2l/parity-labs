# Roadmap

Status key: **shipped** (real code, tested, in production use or
validated end-to-end) / **in progress** (real code exists, incomplete)
/ **not started** (design only, no code).

## Pillar 1 — replay verification: shipped

- `parity replay` — three-way replay, diff, cluster, fidelity scoring.
- `cmd/parity-ci` + `action.yml` — packaged as a drop-in GitHub Action.
  Validated end-to-end on a real Actions run
  ([jex2l/parity#1](https://github.com/Jex2l/parity/pull/1), now
  closed); 3 real bugs found and fixed in that process (hardcoded
  branch name, a cache collision, an exit-code bug under `set -e`).
- `parity-webui` — interactive report viewer.
- `parity-mcp` — exposes the engine to coding agents over MCP.

**Next for pillar 1:** a published release binary instead of
building from source on every Action run (cuts ~10-20s per run); a
hosted multi-repo version of `parity-webui` once there's a real
customer asking for it (see `docs/limits.md` in the `parity` repo —
auth was explicitly deferred, not forgotten).

## Pillar 2 — grounded reasoning: in progress

- `pillar2-grounded-classifier/` (this repo) — a
  `classify.Classifier`-compatible implementation that grounds its
  verdict in retrieved evidence (the PR's actual diff hunks and the
  cluster's actual changed-path data) instead of `classify.Heuristic`'s
  pure structural rules.
- Status as of this commit: the interface contract and a real,
  testable evidence-retrieval + citation-check skeleton exist. No LLM
  endpoint is wired in yet (same honest gap `parity`'s own
  `internal/classify` package already documents for its heuristic
  baseline) — see that directory's README for exactly what runs today
  versus what's stubbed.

**Next for pillar 2:** wire a real model endpoint once one's available
in a dev/eval environment; build the faithfulness check (reusing the
retrieval reranker as an entailment proxy, the same technique
`meridian-legal-intelligence` already uses and has 60/60 tests
covering); evaluate against real complaint data from pillar 1's
heuristic classifier once there's a customer generating that data,
not speculatively.

## Pillar 3 — approval-gated execution: not started

Generalizing `incident-agent`'s read-unrestricted /
write-human-approval-gated pattern past incident response: a
`parity-ci` run that finds unintended diffs could let an agent propose
an auto-fix, which itself gets replayed (pillar 1, recursively) and
then requires the same human-click approval gate before it's allowed
to push a commit.

No code anywhere yet. Correctly last: highest trust, highest liability,
and the design in `jex2l/parity`'s `docs/parity-labs-architecture.md`
is explicit that this should be built against a real customer's actual
auto-fix workflow, not speculatively.

## Sequencing rationale

Each pillar only gets built once there's a concrete reason to, not on a
fixed calendar — pillar 1 shipped because it was independently useful
standalone; pillar 2 starts now because it's the natural second
`Classifier` implementation with zero changes needed to pillar 1's
code; pillar 3 waits for a customer who specifically wants
auto-applied fixes, because shipping the highest-liability feature
speculatively is the wrong order of operations.
