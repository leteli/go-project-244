package code

import (
	f "code/formatters"
	"code/parsers"
	"code/types"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
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
	root := addRootNode(diff)
	return f.FormatDiff(root, format)
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
				result = append(result, types.Node{Key: key, OldValue: val1, NewValue: val1, Kind: types.Unchanged})
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

func addRootNode(nodes []types.Node) types.Node {
	return types.Node{Key: "", Kind: types.Root, Children: nodes}
}

func getSortedKeys(m1, m2 map[string]any) []string {
	keys := make([]string, 0, len(m1)+len(m2))
	for k := range m1 {
		keys = append(keys, k)
	}
	for k := range m2 {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	keys = slices.Compact(keys)
	return keys
}

func parseFileContent(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}
	ext := getExtension(path)
	config, err := parsers.Parse(data, ext)
	if err != nil {
		return nil, fmt.Errorf("error parsing content of file %s: %w", path, err)
	}
	return config, nil
}

func getExtension(path string) string {
	ext := filepath.Ext(path)
	if len(ext) <= 1 {
		return ""
	}
	return ext[1:]
}
