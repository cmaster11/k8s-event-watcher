package k8seventwatcher

import (
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
		return nil, errorf("failed to initialize mapstructure decoder: %+v", err)
	}
	if err := structDecoder.Decode(evt); err != nil {
		return nil, errorf("failed to decode event struct: %+v", err)
	}

	return outMap, nil
}
