package k8seventwatcher

import (
	"errors"
	"fmt"
	"github.com/cmaster11/k8s-event-watcher/lookup"
	"gopkg.in/yaml.v2"
	"strings"
)

type EventFilter struct {
	Rules map[string]*Regexp `yaml:"rules"`
}

func (f *EventFilter) Validate() error {
	// At least one filter must exist
	if len(f.Rules) > 0 {
		return nil
	}
	return errors.New("no rules provided")
}

func (f *EventFilter) Matches(event map[string]interface{}) (map[string]interface{}, error) {
	matchedFields := make(map[string]interface{})
	for path, regex := range f.Rules {
		value, err := lookup.LookupString(event, path)
		if err != nil {
			return nil, errorf("lookup error: %s", err)
		}

		valueStr := fmt.Sprintf("%v", value.Interface())
		if !regex.MatchString(valueStr) {
			return nil, nil
		}

		matchedFields[path] = valueStr
	}

	return matchedFields, nil
}

func (f *EventFilter) String() string {
	var elements []string

	for path, regex := range f.Rules {
		elements = append(elements, fmt.Sprintf("%s=%s", path, regex.String()))
	}

	return strings.Join(elements, ",")
}

func (f *EventFilter) ToYAML() string {
	output, _ := yaml.Marshal(f)
	return string(output)
}
