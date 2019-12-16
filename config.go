package k8seventwatcher

import (
	"errors"
	"gopkg.in/yaml.v2"
	"log"
)

type Config struct {
	Filters []*EventFilter `yaml:"filters"`

	// If true, accepts only events generated since the app has been launched
	SinceNow bool `yaml:"sinceNow"`
}

func (c *Config) Validate() error {
	if len(c.Filters) == 0 {
		return errors.New("zero filters provided")
	}

	for _, filter := range c.Filters {
		if err := filter.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) MatchingEventFilter(event map[string]interface{}) (*EventFilter, error) {
	for _, filter := range c.Filters {
		matches, err := filter.Matches(event)
		if err != nil {
			return nil, errorf("error matching filter: %s", err)
		}
		if matches {
			return filter, nil
		}
	}

	return nil, nil
}

func (c *Config) Dump() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("failed to dump Config: %+v", c)
	}

	return string(data)
}
