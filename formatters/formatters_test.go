package formatters

import (
	"code/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatDiff(t *testing.T) {
	basePath := filepath.Join("..", "testdata", "fixture")
	flatDiff := types.Node{
		Kind: types.Root,
		Children: []types.Node{
			{
				Key:      "follow",
				Kind:     types.Removed,
				OldValue: false,
			},
			{
				Key:      "host",
				Kind:     types.Unchanged,
				OldValue: "hexlet.io",
				NewValue: "hexlet.io",
			},
			{
				Key:      "proxy",
				Kind:     types.Removed,
				OldValue: "123.234.53.22",
			},
			{
				Key:      "timeout",
				Kind:     types.Changed,
				OldValue: 50.0,
				NewValue: 20.0,
			},
			{
				Key:      "verbose",
				Kind:     types.Added,
				NewValue: true,
			},
		},
	}

	nestedDiff := types.Node{
		Kind: types.Root,
		Children: []types.Node{
			{
				Key:  "common",
				Kind: types.Nested,
				Children: []types.Node{
					{
						Key:      "follow",
						Kind:     types.Added,
						NewValue: false,
					},
					{
						Key:      "setting1",
						Kind:     types.Unchanged,
						OldValue: "Value 1",
						NewValue: "Value 1",
					},
					{
						Key:      "setting2",
						Kind:     types.Removed,
						OldValue: 200.0,
					},
					{
						Key:      "setting3",
						Kind:     types.Changed,
						OldValue: true,
						NewValue: nil,
					},
					{
						Key:      "setting4",
						Kind:     types.Added,
						NewValue: "blah blah",
					},
					{
						Key:  "setting5",
						Kind: types.Added,
						NewValue: map[string]any{
							"key5": "value5",
						},
					},
					{
						Key:  "setting6",
						Kind: types.Nested,
						Children: []types.Node{
							{
								Key:  "doge",
								Kind: types.Nested,
								Children: []types.Node{
									{
										Key:      "wow",
										Kind:     types.Changed,
										OldValue: "",
										NewValue: "so much",
									},
								},
							},
							{
								Key:      "key",
								Kind:     types.Unchanged,
								OldValue: "value",
								NewValue: "value",
							},
							{
								Key:      "ops",
								Kind:     types.Added,
								NewValue: "vops",
							},
						},
					},
				},
			},
			{
				Key:  "group1",
				Kind: types.Nested,
				Children: []types.Node{
					{
						Key:      "baz",
						Kind:     types.Changed,
						OldValue: "bas",
						NewValue: "bars",
					},
					{
						Key:      "foo",
						Kind:     types.Unchanged,
						OldValue: "bar",
						NewValue: "bar",
					},
					{
						Key:  "nest",
						Kind: types.Changed,
						OldValue: map[string]any{
							"key": "value",
						},
						NewValue: "str",
					},
				},
			},
			{
				Key:  "group2",
				Kind: types.Removed,
				OldValue: map[string]any{
					"abc": 12345.0,
					"deep": map[string]any{
						"id": 45.0,
					},
				},
			},
			{
				Key:  "group3",
				Kind: types.Added,
				NewValue: map[string]any{
					"deep": map[string]any{
						"id": map[string]any{
							"number": 45.0,
						},
					},
					"fee": 100500,
				},
			}}}

	cases := []struct {
		name     string
		diff     types.Node
		format   string
		wantPath string
	}{
		{
			name:     "flat diff - stylish",
			diff:     flatDiff,
			format:   types.Stylish,
			wantPath: filepath.Join(basePath, "flat", "file1_file2_result.txt"),
		},
		{
			name:     "flat diff - plain",
			diff:     flatDiff,
			format:   types.Plain,
			wantPath: filepath.Join(basePath, "flat", "file1_file2_result_plain.txt"),
		},
		{
			name:     "flat diff - json",
			diff:     flatDiff,
			format:   types.JSON,
			wantPath: filepath.Join(basePath, "flat", "file1_file2_result_json.json"),
		},
		{
			name:     "nested diff stylish",
			diff:     nestedDiff,
			format:   types.Stylish,
			wantPath: filepath.Join(basePath, "nested", "file1_file2_result.txt"),
		},
		{
			name:     "nested diff plain",
			diff:     nestedDiff,
			format:   types.Plain,
			wantPath: filepath.Join(basePath, "nested", "file1_file2_result_plain.txt"),
		},
		{
			name:     "nested diff json",
			diff:     nestedDiff,
			format:   types.JSON,
			wantPath: filepath.Join(basePath, "nested", "file1_file2_result_json.json"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := FormatDiff(c.diff, c.format)
			want := getExpectedDiffContent(t, c.wantPath)
			assert.Equal(t, want, got)
			require.NoError(t, err)
		})
	}
}

func getExpectedDiffContent(t testing.TB, path string) string {
	t.Helper()
	want, err := os.ReadFile(path)
	if err != nil {
		require.NoError(t, err)
	}
	return string(want)
}
