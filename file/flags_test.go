package file

import (
	"reflect"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "valid flags",
			args: []string{
				"add",
				"--description", "Coffee",
				"--amount", "5",
			},
			want: map[string]string{
				"--description": "Coffee",
				"--amount":      "5",
			},
		},
		{
			name:    "missing value",
			args:    []string{"add", "--amount"},
			wantErr: true,
		},
		{
			name:    "value is another flag",
			args:    []string{"add", "--amount", "--id"},
			wantErr: true,
		},
		{
			name:    "argument without flag",
			args:    []string{"add", "Coffee"},
			wantErr: true,
		},
		{
			name: "duplicate flag",
			args: []string{
				"add",
				"--amount", "5",
				"--amount", "10",
			},
			wantErr: true,
		},
		{
			name: "no flags",
			args: []string{"list"},
			want: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFlags(tt.args)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFlags() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFlags() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
