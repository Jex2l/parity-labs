package grounded

import (
	"strings"
	"testing"
)

func TestClassify_FieldMentionedInDiff_LikelyIntended(t *testing.T) {
	// This is the real diff shape from jex2l/parity PR #1's self-test:
	// a currency field added to /stats/revenue's response map.
	hunks := []DiffHunk{
		{
			File: "eval/demo-service/main.go",
			Added: []string{
				`		writeJSON(w, 200, map[string]any{"total_revenue": total, "order_count": len(orders), "generated_at": time.Now().UTC().Format(time.RFC3339), "currency": "USD"})`,
			},
			Removed: []string{
				`		writeJSON(w, 200, map[string]any{"total_revenue": total, "order_count": len(orders), "generated_at": time.Now().UTC().Format(time.RFC3339)})`,
			},
		},
	}
	c := ClusterView{
		Endpoint:     "/stats/revenue",
		ChangedPaths: []string{"$.currency"},
		Severity:     "low",
	}

	got := Classifier{}.Classify(c, hunks)

	if got.Label != LikelyIntended {
		t.Fatalf("Label = %s, want %s", got.Label, LikelyIntended)
	}
	if len(got.Citations) != 1 {
		t.Fatalf("Citations = %v, want exactly 1", got.Citations)
	}
	if !strings.Contains(got.Citations[0].Line, "currency") {
		t.Fatalf("citation line doesn't mention currency: %q", got.Citations[0].Line)
	}
	if got.Citations[0].File != "eval/demo-service/main.go" {
		t.Fatalf("citation file = %q, want eval/demo-service/main.go", got.Citations[0].File)
	}
}

func TestClassify_NoMatchingDiffLine_Uncertain(t *testing.T) {
	hunks := []DiffHunk{
		{File: "unrelated.go", Added: []string{`		fmt.Println("hello")`}},
	}
	c := ClusterView{Endpoint: "/users/{id}", ChangedPaths: []string{"$.email"}, Severity: "medium"}

	got := Classifier{}.Classify(c, hunks)

	if got.Label != Uncertain {
		t.Fatalf("Label = %s, want %s", got.Label, Uncertain)
	}
	if got.Confidence != 0.0 {
		t.Fatalf("Confidence = %v, want 0", got.Confidence)
	}
	if len(got.Citations) != 0 {
		t.Fatalf("expected no citations, got %v", got.Citations)
	}
	if !strings.Contains(got.Rationale, "$.email") {
		t.Fatalf("rationale should name the ungrounded path, got: %s", got.Rationale)
	}
}

func TestClassify_PartialEvidence_UncertainWithPartialCitations(t *testing.T) {
	hunks := []DiffHunk{
		{File: "f.go", Added: []string{`"currency": "USD",`}},
	}
	c := ClusterView{
		Endpoint:     "/stats/revenue",
		ChangedPaths: []string{"$.currency", "$.exchange_rate"},
		Severity:     "low",
	}

	got := Classifier{}.Classify(c, hunks)

	if got.Label != Uncertain {
		t.Fatalf("Label = %s, want %s", got.Label, Uncertain)
	}
	if len(got.Citations) != 1 {
		t.Fatalf("expected exactly 1 partial citation, got %v", got.Citations)
	}
	if !strings.Contains(got.Rationale, "exchange_rate") {
		t.Fatalf("rationale should name the missing path, got: %s", got.Rationale)
	}
}

func TestClassify_CriticalSeverity_AlwaysUnintendedRegardlessOfDiff(t *testing.T) {
	hunks := []DiffHunk{{File: "f.go", Added: []string{`"status_code": 500,`}}}
	c := ClusterView{Endpoint: "/orders", ChangedPaths: []string{"$.status_code"}, Severity: "critical"}

	got := Classifier{}.Classify(c, hunks)

	if got.Label != LikelyUnintended {
		t.Fatalf("Label = %s, want %s", got.Label, LikelyUnintended)
	}
	if got.Confidence != 0.9 {
		t.Fatalf("Confidence = %v, want 0.9", got.Confidence)
	}
}

func TestMentionsField_WholeWordOnly_NoSubstringFalsePositive(t *testing.T) {
	// "id" must not match inside "valid" or "guid_string".
	if mentionsField(`"is_valid": true`, "id") {
		t.Fatalf("mentionsField matched 'id' inside 'valid' — should require a word boundary")
	}
	if !mentionsField(`"id": 42`, "id") {
		t.Fatalf("mentionsField failed to match a real whole-word 'id'")
	}
	if !mentionsField(`  "user_id": 7,`, "user_id") {
		t.Fatalf("mentionsField failed to match 'user_id' as a full field name")
	}
}

func TestLeafField_StripsPrefixAndArrayIndex(t *testing.T) {
	cases := map[string]string{
		"$.currency":         "currency",
		"$.user.address.zip": "zip",
		"items[2]":           "items",
		"$.orders[3].status": "status",
		"plain":              "plain",
	}
	for path, want := range cases {
		if got := leafField(path); got != want {
			t.Errorf("leafField(%q) = %q, want %q", path, got, want)
		}
	}
}
