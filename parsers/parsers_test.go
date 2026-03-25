package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		data    []byte
		f       string
		want    map[string]any
		wantErr bool
	}{
		{name: "valid json", data: []byte(`{"name":"Ann","lastname":"Smith","age":40}`), f: "json", want: map[string]any{"name": "Ann", "lastname": "Smith", "age": float64(40)}, wantErr: false},
		{name: "invalid json", data: []byte(`{"name":"Ann","lastname":`), f: "json", want: map[string]any{}, wantErr: true},
		{name: "valid yaml", data: []byte("name: Ann\nlastname: Smith\nage: 40\n"), f: "yml", want: map[string]any{"name": "Ann", "lastname": "Smith", "age": 40}, wantErr: false},
		{name: "invalid yaml", data: []byte("name: Ann\nlastname"), f: "yml", want: map[string]any{}, wantErr: true},
		{name: "unsupported format", data: []byte("name"), f: "txt", want: map[string]any{}, wantErr: true},
		{name: "string only", data: []byte("name"), f: "json", want: map[string]any{}, wantErr: true},
		{name: "json as array of invalid format", data: []byte(`["string", "one-more-string"]`), f: "json", want: map[string]any{}, wantErr: true},
		{name: "json as array with object", data: []byte(`[{"name":"Ann","lastname":"Smith","age":40}]`), f: "json", want: map[string]any{"name": "Ann", "lastname": "Smith", "age": float64(40)}, wantErr: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Parse(c.data, c.f)
			if !c.wantErr {
				require.NoError(t, err)
				assert.Equal(t, c.want, got)
			} else {
				require.Error(t, err)
			}
		})
	}
}
