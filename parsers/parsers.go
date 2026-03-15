package parsers

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte, format string) (map[string]any, error) {
	var parsed map[string]any
	switch format {
	case "json":
		if err := json.Unmarshal(data, &parsed); err != nil {
			return map[string]any{}, err
		}
		return parsed, nil
	case "yml", "yaml":
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return map[string]any{}, err
		}
		return parsed, nil
	default:
		return map[string]any{}, fmt.Errorf("unsupported format %s", format)
	}
}
