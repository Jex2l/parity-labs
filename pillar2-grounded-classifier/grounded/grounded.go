// Package grounded classifies a behavioral-diff cluster by looking for
// actual evidence in the PR's diff hunks, instead of structural rules
// alone (parity's internal/classify.Heuristic) or an LLM's unverified
// say-so. If no diff line mentions a changed field, the classifier says
// so explicitly — "UNCERTAIN: no evidence found" — rather than guessing,
// the same discipline meridian-legal-intelligence applies to citations:
// an answer the system can't ground is worse than no answer.
//
// This package does not call an LLM. The grounding here is exact/
// substring field-name matching against diff text — real, deterministic,
// testable without a model endpoint, and the correct baseline to exist
// before an LLM-backed version (an LLM could later judge *semantic*
// relatedness, e.g. a renamed field, that this substring match misses;
// see the README's "not built yet" section).
package grounded

import (
	"fmt"
	"strings"
)

// ClusterView is the subset of a replay cluster this classifier needs.
// It deliberately mirrors (not imports — see the repo README for why)
// the shape of parity's internal/report.ClusterView.
type ClusterView struct {
	Endpoint     string
	ChangedPaths []string
	Severity     string // "critical" | "high" | "medium" | "low"
}

// DiffHunk is one file's added/removed lines from a PR diff.
type DiffHunk struct {
	File    string
	Added   []string // lines added by the PR, without the leading "+"
	Removed []string // lines removed by the PR, without the leading "-"
}

// Citation is one concrete piece of evidence backing a classification.
type Citation struct {
	Path string // the changed JSON path this citation supports
	File string
	Line string
}

// Label mirrors parity's classify.Label values so a future adapter can
// map 1:1 without a translation table.
type Label string

const (
	LikelyIntended   Label = "LIKELY_INTENDED"
	LikelyUnintended Label = "LIKELY_UNINTENDED"
	Uncertain        Label = "UNCERTAIN"
)

// Result is one classification, with the evidence it was grounded in.
// Citations is the key difference from classify.Heuristic.Result: every
// LIKELY_INTENDED/LIKELY_UNINTENDED verdict here must point at specific
// diff lines, or the verdict is UNCERTAIN instead.
type Result struct {
	Label      Label
	Confidence float64
	Rationale  string
	Citations  []Citation
}

// Classifier grounds its verdict in the PR's actual diff text.
type Classifier struct{}

// Classify looks for each of c.ChangedPaths' leaf field name in hunks'
// added/removed lines. A cluster is only LIKELY_INTENDED or
// LIKELY_UNINTENDED when every changed path has at least one citation;
// otherwise it's UNCERTAIN, naming exactly which paths had no evidence
// instead of silently guessing.
func (Classifier) Classify(c ClusterView, hunks []DiffHunk) Result {
	if c.Severity == "critical" {
		// Matches classify.Heuristic's rule 1: a critical-severity change
		// (e.g. a 5xx appearing) is unintended regardless of what the diff
		// says — nobody intentionally ships a crash, and grounding doesn't
		// change that judgment, only the confidence/rationale around it.
		cites := citeAll(c.ChangedPaths, hunks)
		return Result{
			Label:      LikelyUnintended,
			Confidence: 0.9,
			Rationale:  "critical severity (e.g. error response) overrides diff evidence — nobody intentionally ships a crash",
			Citations:  cites,
		}
	}

	var found, missing []string
	citations := make([]Citation, 0, len(c.ChangedPaths))
	for _, path := range c.ChangedPaths {
		leaf := leafField(path)
		if cite, ok := findCitation(path, leaf, hunks); ok {
			found = append(found, path)
			citations = append(citations, cite)
		} else {
			missing = append(missing, path)
		}
	}

	switch {
	case len(found) == 0:
		return Result{
			Label:      Uncertain,
			Confidence: 0.0,
			Rationale:  fmt.Sprintf("no PR diff line mentions any of the changed path(s) %s — cannot ground this classification in the PR's actual changes", strings.Join(c.ChangedPaths, ", ")),
			Citations:  nil,
		}
	case len(missing) > 0:
		return Result{
			Label:      Uncertain,
			Confidence: 0.3,
			Rationale:  fmt.Sprintf("found diff evidence for %s but none for %s — partial grounding isn't enough to call this intended or not", strings.Join(found, ", "), strings.Join(missing, ", ")),
			Citations:  citations,
		}
	default:
		return Result{
			Label:      LikelyIntended,
			Confidence: 0.75,
			Rationale:  fmt.Sprintf("every changed path is directly mentioned in the PR diff: %s", strings.Join(found, ", ")),
			Citations:  citations,
		}
	}
}

func citeAll(paths []string, hunks []DiffHunk) []Citation {
	out := make([]Citation, 0, len(paths))
	for _, p := range paths {
		if c, ok := findCitation(p, leafField(p), hunks); ok {
			out = append(out, c)
		}
	}
	return out
}

// findCitation searches hunks' added and removed lines for leaf as a
// whole-word-ish match (surrounded by a quote, dot, bracket, or word
// boundary on each side) so "id" doesn't spuriously match "valid".
func findCitation(path, leaf string, hunks []DiffHunk) (Citation, bool) {
	for _, h := range hunks {
		for _, line := range h.Added {
			if mentionsField(line, leaf) {
				return Citation{Path: path, File: h.File, Line: strings.TrimSpace(line)}, true
			}
		}
		for _, line := range h.Removed {
			if mentionsField(line, leaf) {
				return Citation{Path: path, File: h.File, Line: strings.TrimSpace(line)}, true
			}
		}
	}
	return Citation{}, false
}

func mentionsField(line, leaf string) bool {
	if leaf == "" {
		return false
	}
	idx := strings.Index(line, leaf)
	for idx != -1 {
		before := byte(' ')
		if idx > 0 {
			before = line[idx-1]
		}
		afterIdx := idx + len(leaf)
		after := byte(' ')
		if afterIdx < len(line) {
			after = line[afterIdx]
		}
		if !isIdentChar(before) && !isIdentChar(after) {
			return true
		}
		next := strings.Index(line[idx+1:], leaf)
		if next == -1 {
			break
		}
		idx = idx + 1 + next
	}
	return false
}

func isIdentChar(b byte) bool {
	return b == '_' || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// leafField extracts the final field name from a JSON path like
// "$.user.address.zip" or "orders[3].status", matching the path
// convention parity's internal/diff package already produces.
func leafField(path string) string {
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")
	// Take the last dot-separated segment first ("orders[3].status" ->
	// "status"), THEN strip any trailing array index from just that
	// segment ("items[2]" -> "items") — stripping at the first '[' in the
	// whole path before splitting would truncate everything after an
	// array index that isn't the final segment (caught by
	// TestLeafField_StripsPrefixAndArrayIndex: "$.orders[3].status" was
	// wrongly reduced to "orders" instead of "status").
	if i := strings.LastIndexByte(path, '.'); i != -1 {
		path = path[i+1:]
	}
	if i := strings.IndexByte(path, '['); i != -1 {
		path = path[:i]
	}
	return path
}
