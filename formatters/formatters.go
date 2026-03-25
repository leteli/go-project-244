package formatters

import (
	"code/types"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
)

func FormatDiff(diff types.Node, format string) (string, error) {
	switch format {
	case types.Plain:
		return FormatDiffPlain(diff.Children), nil
	case types.JSON:
		return FormatDiffJSON(diff)
	default:
		return FormatDiffStylish(diff.Children), nil
	}
}

func FormatDiffStylish(diff []types.Node) string {
	var level = 1
	var inner func([]types.Node, int) string
	inner = func(diff []types.Node, level int) string {
		var sb strings.Builder
		lineOffset := getOffset(level, types.LineOffset)
		sb.WriteString("{\n")
		for _, node := range diff {
			sb.WriteString(lineOffset)
			switch node.Kind {
			case types.Nested:
				sb.WriteString(getLine(node.Key, inner(node.Children, level+1), "  "))
			case types.Changed:
				sb.WriteString(getLine(node.Key, formatValue(node.OldValue, level), "- "))
				sb.WriteString(lineOffset)
				sb.WriteString(getLine(node.Key, formatValue(node.NewValue, level), "+ "))
			case types.Added:
				sb.WriteString(getLine(node.Key, formatValue(node.NewValue, level), "+ "))
			case types.Removed:
				sb.WriteString(getLine(node.Key, formatValue(node.OldValue, level), "- "))
			case types.Unchanged:
				sb.WriteString(getLine(node.Key, formatValue(node.NewValue, level), "  "))
			}
		}
		sb.WriteString(getOffset(level, types.BracketOffset))
		sb.WriteString("}")
		return sb.String()
	}
	return inner(diff, level)
}

func FormatDiffPlain(diff []types.Node) string {
	var parentKey string
	var inner func([]types.Node, string) string

	inner = func(diff []types.Node, parentKey string) string {
		var sb strings.Builder
		for _, node := range diff {
			key := getPlainKey(parentKey, node.Key)
			switch node.Kind {
			case types.Nested:
				sb.WriteString(inner(node.Children, key))
			case types.Changed:
				sb.WriteString(fmt.Sprintf("Property '%s' was updated. From %s to %s\n", key, formatPlainValue(node.OldValue), formatPlainValue(node.NewValue)))
			case types.Added:
				sb.WriteString(fmt.Sprintf("Property '%s' was added with value: %s\n", key, formatPlainValue(node.NewValue)))
			case types.Removed:
				sb.WriteString(fmt.Sprintf("Property '%s' was removed\n", key))
			}
		}
		return sb.String()
	}
	return strings.TrimSpace(inner(diff, parentKey))
}

func formatPlainValue(val any) string {
	switch v := val.(type) {
	case []any, map[string]any:
		return "[complex value]"
	case nil:
		return "null"
	case string:
		return fmt.Sprintf("'%s'", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func FormatDiffJSON(diff types.Node) (string, error) {
	raw, err := json.Marshal(diff)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func getPlainKey(parentKey, key string) string {
	if parentKey == "" {
		return key
	}
	return fmt.Sprintf("%s.%s", parentKey, key)
}

func formatValue(val any, level int) string {
	switch v := val.(type) {
	case []any:
		return fmt.Sprintf("%v", v)
	case map[string]any:
		return stringifyMap(v, level+1)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getLine(key string, val string, prefix string) string {
	return fmt.Sprintf("%s%s: %s\n", prefix, key, val)
}

func getOffset(level int, kind string) string {
	if level <= 0 {
		return ""
	}
	var repeat int
	switch kind {
	case types.LineOffset:
		repeat = level*4 - 2
	case types.BracketOffset:
		repeat = (level - 1) * 4
	case types.MapOffset:
		repeat = level * 4
	}
	return strings.Repeat(" ", repeat)
}

func stringifyMap(m map[string]any, level int) string {
	var sb strings.Builder
	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)
	sb.WriteString("{\n")
	for _, k := range keys {
		sb.WriteString(getOffset(level, types.MapOffset))
		sb.WriteString(getLine(k, formatValue(m[k], level), ""))
	}
	sb.WriteString(getOffset(level, types.BracketOffset))
	sb.WriteString("}")
	return sb.String()
}
