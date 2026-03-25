package parsers

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte, format string) (map[string]any, error) {
	var parsed any
	switch format {
	case "json":
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
	case "yml", "yaml":
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported format %s", format)
	}
	switch val := parsed.(type) {
	case map[string]any:
		return val, nil
	case []any:
		if len(val) == 1 {
			if m, isMap := val[0].(map[string]any); isMap {
				return m, nil
			}
		}
		return nil, fmt.Errorf("expected object, got an array of invalid format")

	default:
		return nil, fmt.Errorf("expected object, got type %T", val)
	}
}
