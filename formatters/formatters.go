package formatters

import (
	"code/types"
	"fmt"
	"maps"
	"slices"
	"strings"
)

func FormatDiff(diff []types.Node, format string) string {
	switch format {
	case types.Plain:
		return FormatPlainDiff(diff)
	default:
		return FormatStylishDiff(diff)
	}
}

func FormatStylishDiff(diff []types.Node) string {
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
				val := inner(node.Children, level+1)
				sb.WriteString(getLine(node.Key, val, "  "))
			case types.Changed:
				oldVal := formatValue(node.OldValue, level)
				sb.WriteString(getLine(node.Key, oldVal, "- "))
				sb.WriteString(lineOffset)
				newVal := formatValue(node.NewValue, level)
				sb.WriteString(getLine(node.Key, newVal, "+ "))
			case types.Added:
				newVal := formatValue(node.NewValue, level)
				sb.WriteString(getLine(node.Key, newVal, "+ "))
			case types.Removed:
				oldVal := formatValue(node.OldValue, level)
				sb.WriteString(getLine(node.Key, oldVal, "- "))
			case types.Unchanged:
				val := formatValue(node.NewValue, level)
				sb.WriteString(getLine(node.Key, val, "  "))
			}
		}
		sb.WriteString(getOffset(level, types.BracketOffset))
		sb.WriteString("}")
		return sb.String()
	}
	return inner(diff, level)
}

func FormatPlainDiff(diff []types.Node) string {
	var parentKey string
	var inner func([]types.Node, string) string

	inner = func(diff []types.Node, parentKey string) string {
		var sb strings.Builder
		for _, node := range diff {
			key := getPlainKey(parentKey, node.Key)
			switch node.Kind {
			case types.Nested:
				res := inner(node.Children, key)
				sb.WriteString(res)
			case types.Changed:
				oldVal := formatPlainValue(node.OldValue)
				newVal := formatPlainValue(node.NewValue)
				line := fmt.Sprintf("Property '%s' was updated. From %s to %s\n", key, oldVal, newVal)
				sb.WriteString(line)
			case types.Added:
				val := formatPlainValue(node.NewValue)
				line := fmt.Sprintf("Property '%s' was added with value: %s\n", key, val)
				sb.WriteString(line)
			case types.Removed:
				line := fmt.Sprintf("Property '%s' was removed\n", key)
				sb.WriteString(line)
			}
		}
		return sb.String()
	}
	result := inner(diff, parentKey)
	return strings.TrimSpace(result)
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
func getPlainKey(parentKey, key string) string {
	if parentKey == "" {
		// return "'" + key + "'"
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
		return fmt.Sprintf("%v", val)
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
		v := formatValue(m[k], level)
		sb.WriteString(getOffset(level, types.MapOffset))
		sb.WriteString(getLine(k, v, ""))
	}
	sb.WriteString(getOffset(level, types.BracketOffset))
	sb.WriteString("}")
	return sb.String()
}
