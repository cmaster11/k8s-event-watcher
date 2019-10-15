package k8seventwatcher

import (
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"testing"
)

func TestEventFilter(t *testing.T) {
	input := `
objectKind: Job
objectNamespace: default
objectName: test.*
eventType: Warning
eventReason: BackoffLimitExceeded
`

	filter := &EventFilter{}

	if err := yaml.Unmarshal([]byte(input), filter); err != nil {
		t.Fatal(err)
	}

	evt := v1.Event{}
	evt.InvolvedObject = v1.ObjectReference{
		Kind:      "Job",
		Namespace: "default",
		Name:      "test123",
	}
	evt.Type = "Warning"
	evt.Reason = "BackoffLimitExceeded"

	if !filter.Matches(&evt) {
		t.Fatal("expected match")
	}

	output := filter.ToYAML()

	// Use marshaled output
	filter2 := &EventFilter{}

	if err := yaml.Unmarshal([]byte(output), filter2); err != nil {
		t.Fatal(err)
	}

	if filter2.EventReason.String() != filter.EventReason.String() {
		t.Fatal("wrong EventReason")
	}
	if filter2.EventType.String() != filter.EventType.String() {
		t.Fatal("wrong EventType")
	}
	if filter2.ObjectKind.String() != filter.ObjectKind.String() {
		t.Fatal("wrong ObjectKind")
	}
	if filter2.ObjectName.String() != filter.ObjectName.String() {
		t.Fatal("wrong ObjectName")
	}
	if filter2.ObjectNamespace.String() != filter.ObjectNamespace.String() {
		t.Fatal("wrong ObjectNamespace")
	}
}
