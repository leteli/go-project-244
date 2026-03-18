package formatters

import (
	"code/types"
	"fmt"
	"maps"
	"slices"
	"strings"
)

func FormatDiff(diff []types.Node, level int) string { // TODO handle stylish arg
	var sb strings.Builder
	lineOffset := getOffset(level, types.LineOffset)
	sb.WriteString("{\n")
	for _, node := range diff {
		sb.WriteString(lineOffset)
		switch node.Kind {
		case types.Nested:
			val := FormatDiff(node.Children, level+1)
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

func formatValue(val any, level int) string {
	if s, ok := val.([]any); ok {
		return fmt.Sprintf("%v", s) // TODO: stringify correctly
	}
	if m, ok := val.(map[string]any); ok {
		return stringifyMap(m, level+1)
	}
	return fmt.Sprintf("%v", val)
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
