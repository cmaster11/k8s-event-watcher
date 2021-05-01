package pkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mcuadros/go-lookup"
	"gopkg.in/yaml.v2"
)

type EventFilter struct {

	// Rules used to match the event
	Rules map[string]*Regexp `mapstructure:"rules" yaml:"rules" json:"rules"`

	// If all these rules match, the event is considered an error
	ErrorRules map[string]*Regexp `mapstructure:"errorRules" yaml:"errorRules" json:"errorRules"`
}

func (f *EventFilter) Validate() error {
	// At least one rule must exist, once a filter is defined
	if len(f.Rules) == 0 {
		return errors.New("no rules provided")
	}

	return nil
}

type MatchResult struct {
	MatchedFields      map[string]interface{} `yaml:"matchedFields" json:"matchedFields"`
	MatchedErrorFields map[string]interface{} `yaml:"matchedErrorFields" json:"matchedErrorFields"`
}

func (f *EventFilter) Matches(event map[string]interface{}) (*MatchResult, error) {
	matchedFields := make(map[string]interface{})
	for path, regex := range f.Rules {
		value, err := lookup.LookupString(event, path)
		if err != nil {
			return nil, fmt.Errorf("lookup error: %s", err)
		}

		valueStr := fmt.Sprintf("%v", value.Interface())
		if !regex.MatchString(valueStr) {
			return nil, nil
		}

		matchedFields[path] = valueStr
	}

	matchedErrorFields := make(map[string]interface{})
	for path, regex := range f.ErrorRules {
		value, err := lookup.LookupString(event, path)
		if err != nil {
			return nil, fmt.Errorf("lookup error: %s", err)
		}

		valueStr := fmt.Sprintf("%v", value.Interface())
		if !regex.MatchString(valueStr) {
			matchedErrorFields = nil
			break
		}

		matchedErrorFields[path] = valueStr
	}

	return &MatchResult{
		MatchedFields:      matchedFields,
		MatchedErrorFields: matchedErrorFields,
	}, nil
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
