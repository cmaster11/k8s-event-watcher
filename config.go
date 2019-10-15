package k8seventwatcher

import (
	"errors"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
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

func (c *Config) MatchingEventFilter(event *v1.Event) *EventFilter {
	for _, filter := range c.Filters {
		if filter.Matches(event) {
			return filter
		}
	}

	return nil
}

func (c *Config) Dump() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("failed to dump Config: %+v", c)
	}

	return string(data)
}
