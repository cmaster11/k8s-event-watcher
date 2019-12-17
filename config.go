package k8seventwatcher

import (
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Filters []*EventFilter `yaml:"filters"`

	// If true, accepts only events generated since the app has been launched
	SinceNow bool `yaml:"sinceNow"`
}

func (c *Config) Validate() error {
	// Technically, we may want to process ALL k8s events
	// if len(c.Filters) == 0 {
	// 	return errors.New("zero filters provided")
	// }

	for _, filter := range c.Filters {
		if err := filter.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) MatchingEventFilter(event map[string]interface{}) (*EventFilter, map[string]interface{}, error) {
	for _, filter := range c.Filters {
		matchedFields, err := filter.Matches(event)
		if err != nil {
			return nil, nil, errorf("error matching filter: %s", err)
		}
		if matchedFields != nil {
			return filter, matchedFields, nil
		}
	}

	return nil, nil, nil
}

func (c *Config) Dump() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("failed to dump config: %+v", c)
	}

	return string(data)
}
