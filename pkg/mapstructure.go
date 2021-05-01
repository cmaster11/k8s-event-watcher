package pkg

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	v1 "k8s.io/api/core/v1"
)

func eventToMap(evt *v1.Event) (map[string]interface{}, error) {
	// Convert the event to a map
	outMap := make(map[string]interface{})
	structDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &outMap,
		TagName: "json",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mapstructure decoder: %w", err)
	}
	if err := structDecoder.Decode(evt); err != nil {
		return nil, fmt.Errorf("failed to decode event struct: %w", err)
	}

	return outMap, nil
}
