package k8seventwatcher

import (
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"testing"
)

func TestEventFilter(t *testing.T) {
	input := `
rules:
  involvedObject.kind: Job
  involvedObject.namespace: default
  involvedObject.name: test.*
  type: Warning
  reason: BackoffLimitExceeded
errorRules:
  type: "^W.*"
`

	filter := &EventFilter{}

	if err := yaml.Unmarshal([]byte(input), filter); err != nil {
		t.Fatal(err)
	}

	if err := filter.Validate(); err != nil {
		t.Fatalf("invalid rules: %s", err)
	}

	evt := v1.Event{}
	evt.InvolvedObject = v1.ObjectReference{
		Kind:      "Job",
		Namespace: "default",
		Name:      "test123",
	}
	evt.Type = "Warning"
	evt.Reason = "BackoffLimitExceeded"

	// Marshal to JSON
	obj, err := eventToMap(&evt)
	if err != nil {
		t.Fatalf("failed to cast event to map: %s", err)
	}

	matchResult, err := filter.Matches(obj)
	if err != nil {
		t.Fatalf("match error: %s", err)
	}
	if matchResult == nil {
		t.Fatalf("no match")
	}
	if matchResult.MatchedFields["involvedObject.name"] != evt.InvolvedObject.Name {
		t.Fatalf("bad matched field")
	}

	if len(matchResult.MatchedErrorFields) == 0 {
		t.Fatalf("no error matched")
	}
	if matchResult.MatchedErrorFields["type"] != evt.Type {
		t.Fatalf("bad matched error field")
	}

	output := filter.ToYAML()

	// Use marshaled output
	filter2 := &EventFilter{}

	if err := yaml.Unmarshal([]byte(output), filter2); err != nil {
		t.Fatal(err)
	}

	// Test marshal to string
	for path, regex := range filter.Rules {
		found := false
		for path2, regex2 := range filter2.Rules {
			if path == path2 {
				found = true
				if regex.String() != regex2.String() {
					t.Fatalf("wrong regex for %s", path)
				}
			}
		}
		if !found {
			t.Fatalf("path %s not found", path)
		}
	}
}
