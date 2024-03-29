package pkg

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Filters []*EventFilter `mapstructure:"filters" yaml:"filters"`

	// If true, accepts only events generated since the app has been launched
	SinceNow bool `mapstructure:"sinceNow" yaml:"sinceNow"`
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

func (c *Config) MatchingEventFilter(event map[string]interface{}) (*EventFilter, *MatchResult, error) {
	for _, filter := range c.Filters {
		matchResult, err := filter.Matches(event)
		if err != nil {
			return nil, nil, fmt.Errorf("error matching filter: %w", err)
		}
		if matchResult != nil {
			return filter, matchResult, nil
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
