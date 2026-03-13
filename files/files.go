package files

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func ParseFileContent(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ext := getExtension(path)
	switch ext {
	case "json":
		var parsed any // TODO: maybe use map[string]any if json objects only
		if err := json.Unmarshal(data, &parsed); err != nil {
			return err
		}
		fmt.Println(parsed)
		return nil
	default:
		return fmt.Errorf("unsupported file extension %s", ext)
	}
}

func getExtension(path string) string {
	dotIdx := strings.LastIndex(path, ".")
	return path[dotIdx+1:]
}
