package types

type Node struct {
	Key      string
	OldValue any
	NewValue any
	Children []Node
	Kind     string
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
