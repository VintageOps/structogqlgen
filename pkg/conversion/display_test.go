package conversion

import (
	"strings"
	"testing"
)

// TestGqlPrettyPrint is a test function for the GqlPrettyPrint function.
func TestGqlPrettyPrint(t *testing.T) {
	tests := []struct {
		name    string
		input   []GqlTypeDefinition
		opts    *PrettyPrintOptions
		want    string
		wantErr bool
	}{
		{
			name:    "Empty GQL Type Definition",
			input:   []GqlTypeDefinition{},
			opts:    &PrettyPrintOptions{},
			want:    "\n",
			wantErr: false,
		},
		// Will add more real test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GqlPrettyPrint(tt.input, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GqlPrettyPrint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.ReplaceAll(got, " ", "") != strings.ReplaceAll(tt.want, " ", "") {
				t.Errorf("GqlPrettyPrint() = %v, want %v", got, tt.want)
			}
		})
	}
}
