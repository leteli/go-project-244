package diff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	basePath := filepath.Join("..", "testdata", "fixture")
	cases := []struct {
		name     string
		path1    string
		path2    string
		wantPath string
		wantErr  bool
	}{
		{name: "2 valid json configs", path1: filepath.Join(basePath, "file1.json"), path2: filepath.Join(basePath, "file2.json"), wantPath: filepath.Join(basePath, "file1_file2_json_result.txt"), wantErr: false},
		{name: "1 json with string only", path1: filepath.Join(basePath, "file1.json"), path2: filepath.Join(basePath, "string_only.json"), wantPath: "", wantErr: true},
		{name: "1 invalid json", path1: filepath.Join(basePath, "file1.json"), path2: filepath.Join(basePath, "invalid.json"), wantPath: "", wantErr: true},
		{name: "unsupported file extension", path1: filepath.Join(basePath, "file1.json"), path2: filepath.Join(basePath, "wrong_ext.txt"), wantPath: "", wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := GenDiff(c.path1, c.path2)
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

func getExpectedDiffContent(t testing.TB, path string) string {
	t.Helper()
	want, err := os.ReadFile(path)
	if err != nil {
		require.NoError(t, err)
	}
	return string(want)
}
