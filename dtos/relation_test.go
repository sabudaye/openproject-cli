package dtos_test

import (
	"encoding/json"
	"testing"

	"github.com/sabudaye/openproject-cli/dtos"
)

func TestRelationDtoConvertsLinksWithTitles(t *testing.T) {
	var dto dtos.RelationDto
	payload := `{
		"id": 650,
		"name": "duplicated by",
		"type": "duplicated",
		"reverseType": "duplicates",
		"lag": 2,
		"description": "Those are the same...",
		"_links": {
			"from": {"href": "/api/v3/work_packages/108", "title": "Destroy Rebel Base"},
			"to": {"href": "/api/v3/work_packages/78", "title": "Annihilate Yavin 4"}
		}
	}`
	if err := json.Unmarshal([]byte(payload), &dto); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	rel := dto.Convert()
	if rel.Id != 650 || rel.Name != "duplicated by" || rel.Type != "duplicated" {
		t.Errorf("basic fields not converted: %+v", rel)
	}
	if rel.Lag != 2 || !rel.HasLag {
		t.Errorf("Lag = %d (HasLag=%v), want 2 (true)", rel.Lag, rel.HasLag)
	}
	if rel.FromId != 108 || rel.FromTitle != "Destroy Rebel Base" {
		t.Errorf("From = #%d %q, want #108 %q", rel.FromId, rel.FromTitle, "Destroy Rebel Base")
	}
	if rel.ToId != 78 || rel.ToTitle != "Annihilate Yavin 4" {
		t.Errorf("To = #%d %q, want #78 %q", rel.ToId, rel.ToTitle, "Annihilate Yavin 4")
	}
}

func TestRelationDtoOmittedLagStaysAbsent(t *testing.T) {
	var dto dtos.RelationDto
	payload := `{"id": 651, "name": "relates to", "type": "relates", "lag": null}`
	if err := json.Unmarshal([]byte(payload), &dto); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if dto.Lag != nil {
		t.Errorf("Lag = %v, want nil for an omitted/null lag", dto.Lag)
	}
	if got := dto.Convert().HasLag; got {
		t.Error("HasLag = true, want false for a relation without lag")
	}
}

// Instances may return the work packages as _embedded objects without
// titles; the ids must still be picked up.
func TestRelationDtoConvertsEmbeddedWorkPackages(t *testing.T) {
	var dto dtos.RelationDto
	payload := `{
		"id": 651,
		"name": "relates to",
		"type": "relates",
		"_embedded": {"from": {"id": 42}, "to": {"id": 77}}
	}`
	if err := json.Unmarshal([]byte(payload), &dto); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	rel := dto.Convert()
	if rel.FromId != 42 || rel.ToId != 77 {
		t.Errorf("FromId/ToId = %d/%d, want 42/77", rel.FromId, rel.ToId)
	}
	if rel.FromTitle != "" || rel.ToTitle != "" {
		t.Errorf("titles = %q/%q, want empty", rel.FromTitle, rel.ToTitle)
	}
}

func TestRelationCollectionDtoConvertsElements(t *testing.T) {
	var dto dtos.RelationCollectionDto
	payload := `{
		"_type": "Collection",
		"total": 2,
		"count": 2,
		"_embedded": {"elements": [{"id": 1}, {"id": 2}]}
	}`
	if err := json.Unmarshal([]byte(payload), &dto); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	rels := dto.Convert()
	if len(rels) != 2 || rels[0].Id != 1 || rels[1].Id != 2 {
		t.Errorf("Convert() = %+v, want 2 relations with ids 1 and 2", rels)
	}
}
