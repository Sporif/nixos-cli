package utils_test

import (
	"testing"

	"github.com/nix-community/nixos-cli/internal/utils"
)

func TestNewNixAttrName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  utils.NixAttrName
	}{
		{
			name:  "a-z",
			input: "foobar",
			want:  utils.NixAttrName("foobar"),
		},
		{
			name:  "space and dot",
			input: "foo.bar ",
			want:  utils.NixAttrName("\"foo.bar \""),
		},
		{
			name:  "space only",
			input: " ",
			want:  utils.NixAttrName("\" \""),
		},
		{
			name:  "dot only",
			input: ".",
			want:  utils.NixAttrName("\".\""),
		},
		{
			name:  "empty attribute",
			input: "",
			want:  utils.NixAttrName(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.NewNixAttrName(tt.input)

			if got != tt.want {
				t.Errorf("Nix attribute name: got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNewNixAttrPath(t *testing.T) {
	tests := []struct {
		name  string
		input []any
		want  utils.NixAttrPath
	}{
		{
			name:  "no attributes",
			input: []any{},
			want:  utils.NixAttrPath(""),
		},
		{
			name:  "single attribute",
			input: []any{"config"},
			want:  utils.NixAttrPath("config"),
		},
		{
			name:  "multiple attributes",
			input: []any{"config", "foo", "bar"},
			want:  utils.NixAttrPath("config.foo.bar"),
		},
		{
			name:  "attribute path",
			input: []any{"config", utils.NixAttrPath("config.foo.bar"), "bar"},
			want:  utils.NixAttrPath("config.config.foo.bar.bar"),
		},
		{
			name:  "space and dot",
			input: []any{"config", "foo. ", "bar"},
			want:  utils.NixAttrPath("config.\"foo. \".bar"),
		},
		{
			name:  "space only",
			input: []any{"config", " ", "bar"},
			want:  utils.NixAttrPath("config.\" \".bar"),
		},
		{
			name:  "dot only",
			input: []any{"config", ".", "bar"},
			want:  utils.NixAttrPath("config.\".\".bar"),
		},
		{
			name:  "empty attributes",
			input: []any{"", "config", "", "foo", "", "bar", ""},
			want:  utils.NixAttrPath("config.foo.bar"),
		},
		{
			name:  "non-string types",
			input: []any{"config", 1, true},
			want:  utils.NixAttrPath("config.1.true"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.NewNixAttrPath(tt.input...)

			if got != tt.want {
				t.Errorf("Nix attribute path: got %s, want %s", got, tt.want)
			}
		})
	}
}
