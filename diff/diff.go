package diff

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

func ParseFileContent(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]any{}, err
	}
	ext := getExtension(path)
	switch ext {
	case "json":
		var parsed map[string]any
		if err := json.Unmarshal(data, &parsed); err != nil {
			return map[string]any{}, err
		}
		return parsed, nil
	default:
		return map[string]any{}, fmt.Errorf("unsupported file extension %s", ext)
	}
}

func getExtension(path string) string {
	dotIdx := strings.LastIndex(path, ".")
	return path[dotIdx+1:]
}

func GenDiff(path1, path2 string) (string, error) {
	content1, err := ParseFileContent(path1)
	if err != nil {
		return "", err
	}
	content2, err := ParseFileContent(path2)
	if err != nil {
		return "", err
	}
	return getDiff(content1, content2), nil
}

func getSortedKeys(m1 map[string]any, m2 map[string]any) []string {
	newMap := maps.Clone(m1)
	maps.Copy(newMap, m2)
	keys := slices.Collect(maps.Keys(newMap))
	slices.Sort(keys)
	return keys
}

func getLine(key string, val any, prefix string) string {
	return fmt.Sprintf(" %s %s: %v\n", prefix, key, val)
}

func getDiff(content1, content2 map[string]any) string {
	sortedKeys := getSortedKeys(content1, content2)
	var sb strings.Builder
	sb.WriteString("{\n")
	for _, key := range sortedKeys {
		val1, ok1 := content1[key]
		val2, ok2 := content2[key]
		if ok1 && ok2 {
			if val1 == val2 {
				sb.WriteString(getLine(key, val1, " "))
			} else {
				sb.WriteString(getLine(key, val1, "-"))
				sb.WriteString(getLine(key, val2, "+"))
			}
		} else if ok1 {
			sb.WriteString(getLine(key, val1, "-"))
		} else {
			sb.WriteString(getLine(key, val2, "+"))
		}
	}
	sb.WriteString("}")
	return sb.String()
}
