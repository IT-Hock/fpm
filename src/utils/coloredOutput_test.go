package utils

import (
	"fmt"
	"testing"
)

func TestColorizeHtml(t *testing.T) {
	var tests = []struct {
		input, expected string
		args            []interface{}
	}{
		{"", "", nil},
		{"<", "<", nil},
		{"< ", "< ", nil},
		{"<red>", "<red>", nil},
		{"<red> ", "<red> ", nil},
		{"<red>foo", "<red>foo", nil},
		{"<red>foo</red>", "\x1b[31mfoo\x1b[0m", nil},
		{"<red>foo</red> ", "\x1b[31mfoo\x1b[0m ", nil},
		{"<red>foo</red> bar", "\x1b[31mfoo\x1b[0m bar", nil},
		{"<red>foo</red> <blue>bar</blue>", "\x1b[31mfoo\x1b[0m \x1b[34mbar\x1b[0m", nil},
		{"<red>%s</red>", "\x1b[31m%s\x1b[0m", []interface{}{"foo"}},
		{"<red>%s</red> <blue>%s</blue>", "\x1b[31m%s\x1b[0m \x1b[34m%s\x1b[0m", []interface{}{"foo", "bar"}},
		{"<red>%s</red> <blue>%s</blue>", "\x1b[31m%s\x1b[0m \x1b[34m%s\x1b[0m", []interface{}{"foo"}},
		{"<red>%s</red> <blue>%s</blue>", "\x1b[31m%s\x1b[0m \x1b[34m%s\x1b[0m", []interface{}{"foo", "<red>bar</red>"}},
	}

	for _, test := range tests {
		actual := ColorizeHtml(test.input, test.args...)
		expected := fmt.Sprintf(test.expected, test.args...)
		if actual != expected {
			t.Errorf("ColorizeHtml(%q) = %q, expected %q", test.input, actual, test.expected)
		}
	}
}
