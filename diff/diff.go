package diff

import (
	f "code/formatters"
	"code/parsers"
	"code/types"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

func GenDiff(path1, path2, format string) (string, error) {
	content1, err := parseFileContent(path1)
	if err != nil {
		return "", err
	}
	content2, err := parseFileContent(path2)
	if err != nil {
		return "", err
	}
	diff := BuildDiff(content1, content2)
	return f.FormatDiff(diff, format), nil
}

func BuildDiff(content1, content2 map[string]any) []types.Node {
	sortedKeys := getSortedKeys(content1, content2)
	result := make([]types.Node, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		val1, ok1 := content1[key]
		val2, ok2 := content2[key]
		if ok1 && ok2 {
			map1, isMap1 := val1.(map[string]any)
			map2, isMap2 := val2.(map[string]any)
			if isMap1 && isMap2 {
				nested := BuildDiff(map1, map2)
				result = append(result, types.Node{Key: key, Children: nested, Kind: types.Nested})
			} else if reflect.DeepEqual(val1, val2) {
				result = append(result, types.Node{Key: key, NewValue: val1, Kind: types.Unchanged})
			} else {
				result = append(result, types.Node{Key: key, OldValue: val1, NewValue: val2, Kind: types.Changed})
			}
		} else if ok1 {
			result = append(result, types.Node{Key: key, OldValue: val1, Kind: types.Removed})
		} else {
			result = append(result, types.Node{Key: key, NewValue: val2, Kind: types.Added})
		}
	}
	return result
}

func getSortedKeys(m1, m2 map[string]any) []string {
	newMap := maps.Clone(m1)
	maps.Copy(newMap, m2)
	keys := slices.Collect(maps.Keys(newMap))
	slices.Sort(keys)
	return keys
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
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}
	return strings.Replace(ext, ".", "", 1)
}
