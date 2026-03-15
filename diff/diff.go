package diff

import (
	"code/parsers"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

func GenDiff(path1, path2 string) (string, error) {
	content1, err := parseFileContent(path1)
	if err != nil {
		return "", err
	}
	content2, err := parseFileContent(path2)
	if err != nil {
		return "", err
	}
	return getDiff(content1, content2), nil
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

func getSortedKeys(m1, m2 map[string]any) []string {
	newMap := maps.Clone(m1)
	maps.Copy(newMap, m2)
	keys := slices.Collect(maps.Keys(newMap))
	slices.Sort(keys)
	return keys
}

func getLine(key string, val any, prefix string) string {
	return fmt.Sprintf(" %s %s: %v\n", prefix, key, val)
}

func parseFileContent(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]any{}, fmt.Errorf("error reading file %s: %w", path, err)
	}
	ext := getExtension(path)
	result, err := parsers.Parse(data, ext)
	if err != nil {
		return map[string]any{}, fmt.Errorf("error parsing content of file %s: %w", path, err)
	}
	return result, nil
}

func getExtension(path string) string {
	dotIdx := strings.LastIndex(path, ".")
	return path[dotIdx+1:]
}
