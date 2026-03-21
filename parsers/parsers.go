package parsers

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte, format string) (any, error) {
	var parsed any
	switch format {
	case "json":
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
		return parsed, nil
	case "yml", "yaml":
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
		return parsed, nil
	default:
		return nil, fmt.Errorf("unsupported format %s", format)
	}
}
