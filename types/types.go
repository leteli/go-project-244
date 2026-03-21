package types

type Node struct {
	Key      string `json:"key"`
	OldValue any    `json:"old_value,omitempty"`
	NewValue any    `json:"new_value,omitempty"`
	Children []Node `json:"children,omitempty"`
	Kind     string `json:"kind"`
}

const (
	Added     = "added"
	Removed   = "removed"
	Changed   = "changed"
	Unchanged = "unchanged"
	Nested    = "nested"
	Root      = "root"
)

const (
	LineOffset    = "lineOffset"
	BracketOffset = "bracketOffset"
	MapOffset     = "mapOffset"
)

const (
	Stylish = "stylish"
	Plain   = "plain"
	JSON    = "json"
)
