package code

import (
	"code/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	basePath := filepath.Join("testdata", "fixture")
	flatPath := filepath.Join(basePath, "flat")
	nestedPath := filepath.Join(basePath, "nested")
	invalidPath := filepath.Join(basePath, "invalid")
	cases := []struct {
		name     string
		path1    string
		path2    string
		format   string
		wantPath string
		wantErr  bool
	}{
		{
			name:     "2 valid json format stylish",
			path1:    filepath.Join(flatPath, "file1.json"),
			path2:    filepath.Join(flatPath, "file2.json"),
			format:   types.Stylish,
			wantPath: filepath.Join(flatPath, "file1_file2_result.txt"),
			wantErr:  false,
		},
		{
			name:    "1 json as string only",
			path1:   filepath.Join(flatPath, "file1.json"),
			path2:   filepath.Join(invalidPath, "string_only.json"),
			wantErr: true,
		},
		{
			name:    "1 json as array of invalid format",
			path1:   filepath.Join(flatPath, "file1.json"),
			path2:   filepath.Join(invalidPath, "array_invalid.json"),
			wantErr: true,
		},
		// {
		// 	name:    "1 json as array with object",
		// 	path1:   filepath.Join(flatPath, "file1.json"),
		// 	path2:   filepath.Join(flatPath, "file2_in_array.json"),
		// 	wantErr: false,
		// },
		{
			name:    "1 invalid json",
			path1:   filepath.Join(flatPath, "file1.json"),
			path2:   filepath.Join(invalidPath, "invalid.json"),
			wantErr: true,
		},
		{
			name:    "unsupported file extension",
			path1:   filepath.Join(flatPath, "file1.json"),
			path2:   filepath.Join(invalidPath, "wrong_ext.txt"),
			wantErr: true,
		},
		{
			name:     "valid yaml format stylish",
			path1:    filepath.Join(flatPath, "file1.yml"),
			path2:    filepath.Join(flatPath, "file2.yaml"),
			format:   types.Stylish,
			wantPath: filepath.Join(flatPath, "file1_file2_result.txt"),
			wantErr:  false,
		},
		{
			name:     "valid json and yaml format stylish",
			path1:    filepath.Join(flatPath, "file1.json"),
			path2:    filepath.Join(flatPath, "file2.yaml"),
			format:   types.Stylish,
			wantPath: filepath.Join(flatPath, "file1_file2_result.txt"),
			wantErr:  false,
		},
		{
			name:     "invalid yaml",
			path1:    filepath.Join(flatPath, "file1.yml"),
			path2:    filepath.Join(invalidPath, "invalid.yml"),
			wantPath: "",
			wantErr:  true,
		},
		{
			name:     "valid nested json format stylish",
			path1:    filepath.Join(nestedPath, "file1.json"),
			path2:    filepath.Join(nestedPath, "file2.json"),
			format:   types.Stylish,
			wantPath: filepath.Join(nestedPath, "file1_file2_result.txt"),
			wantErr:  false,
		},
		{
			name:     "valid nested yaml format stylish",
			path1:    filepath.Join(nestedPath, "file1.yml"),
			path2:    filepath.Join(nestedPath, "file2.yml"),
			format:   types.Stylish,
			wantPath: filepath.Join(nestedPath, "file1_file2_result.txt"),
			wantErr:  false,
		},
		{
			name:     "valid nested json and yaml format stylish",
			path1:    filepath.Join(nestedPath, "file1.yml"),
			path2:    filepath.Join(nestedPath, "file2.yml"),
			format:   types.Stylish,
			wantPath: filepath.Join(nestedPath, "file1_file2_result.txt"),
			wantErr:  false,
		},
		{
			name:     "valid json format plain",
			path1:    filepath.Join(flatPath, "file1.json"),
			path2:    filepath.Join(flatPath, "file2.json"),
			format:   types.Plain,
			wantPath: filepath.Join(flatPath, "file1_file2_result_plain.txt"),
			wantErr:  false,
		},
		{
			name:     "valid nested yaml format plain",
			path1:    filepath.Join(nestedPath, "file1.yml"),
			path2:    filepath.Join(nestedPath, "file2.yml"),
			format:   types.Plain,
			wantPath: filepath.Join(nestedPath, "file1_file2_result_plain.txt"),
			wantErr:  false,
		},
		{
			name:     "valid yaml format json",
			path1:    filepath.Join(flatPath, "file1.yml"),
			path2:    filepath.Join(flatPath, "file2.yaml"),
			format:   types.JSON,
			wantPath: filepath.Join(flatPath, "file1_file2_result_json.json"),
			wantErr:  false,
		},
		{
			name:     "valid nested json and yaml format json",
			path1:    filepath.Join(nestedPath, "file1.json"),
			path2:    filepath.Join(nestedPath, "file2.yml"),
			format:   types.JSON,
			wantPath: filepath.Join(nestedPath, "file1_file2_result_json.json"),
			wantErr:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := GenDiff(c.path1, c.path2, c.format)
			if !c.wantErr {
				require.NoError(t, err)
				want := getExpectedDiffContent(t, c.wantPath)
				assert.Equal(t, want, got)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestBuildDiff(t *testing.T) {
	cases := []struct {
		name string
		m1   map[string]any
		m2   map[string]any
		want []types.Node
	}{{
		name: "flat diff",
		m1: map[string]any{
			"follow":  false,
			"host":    "hexlet.io",
			"proxy":   "123.234.53.22",
			"timeout": 50.0,
		},
		m2: map[string]any{
			"host":    "hexlet.io",
			"timeout": 20.0,
			"verbose": true,
		}, want: []types.Node{
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
	},
		{
			name: "nested diff",
			m1: map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting2": 200.0,
					"setting3": true,
					"setting5": map[string]any{},
				},
			},
			m2: map[string]any{
				"common": map[string]any{
					"setting1": "Value 1",
					"setting3": nil,
					"setting4": "blah blah",
					"setting5": map[string]any{
						"key5": "value5",
					},
				},
			},
			want: []types.Node{
				{
					Key:  "common",
					Kind: types.Nested,
					Children: []types.Node{
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
							Kind: types.Nested,
							Children: []types.Node{
								{
									Key:      "key5",
									Kind:     types.Added,
									NewValue: "value5",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := BuildDiff(c.m1, c.m2)
			assert.Equal(t, c.want, got)
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
