package types

type Node struct {
	Key      string `json:"key"`
	OldValue any    `json:"old_value"` // NB: maybe omitempty
	NewValue any    `json:"new_value"` // NB: maybe omitempty
	Children []Node `json:"children,omitempty"`
	Kind     string `json:"kind"`
}

const (
	Added     = "added"
	Removed   = "removed"
	Changed   = "changed"
	Unchanged = "unchanged"
	Nested    = "nested"
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
