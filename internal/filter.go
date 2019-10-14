package internal

import (
	"errors"
	v1 "k8s.io/api/core/v1"
)

type EventFilter struct {
	ObjectNamespace *Regexp `yaml:"objectNamespace,omitempty"`
	ObjectKind      *Regexp `yaml:"objectKind,omitempty"`
	ObjectName      *Regexp `yaml:"objectName,omitempty"`
	EventType       *Regexp `yaml:"eventType,omitempty"`
	EventReason     *Regexp `yaml:"eventReason,omitempty"`
}

func (f *EventFilter) Validate() error {
	// At least one filter must exist
	if f.ObjectNamespace != nil {
		return nil
	}
	if f.ObjectKind != nil {
		return nil
	}
	if f.ObjectName != nil {
		return nil
	}
	if f.EventType != nil {
		return nil
	}
	if f.EventReason != nil {
		return nil
	}
	return errors.New("no filter attributes provided")
}

func (f *EventFilter) Matches(event *v1.Event) bool {
	if f.ObjectNamespace != nil {
		if !f.ObjectNamespace.MatchString(event.InvolvedObject.Namespace) {
			return false
		}
	}
	if f.ObjectKind != nil {
		if !f.ObjectKind.MatchString(event.InvolvedObject.Kind) {
			return false
		}
	}
	if f.ObjectName != nil {
		if !f.ObjectName.MatchString(event.InvolvedObject.Name) {
			return false
		}
	}
	if f.EventType != nil {
		if !f.EventType.MatchString(event.Type) {
			return false
		}
	}
	if f.EventReason != nil {
		if !f.EventReason.MatchString(event.Reason) {
			return false
		}
	}

	return true
}
