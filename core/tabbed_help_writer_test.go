package core

import (
	"bytes"
	"testing"
)

func TestTabbedHelpWriter_Write(t *testing.T) {
	testCases := []struct {
		title    string
		input    []byte
		expected string
	}{
		{
			title: "nil input",
			input: nil,
		},
		{
			title: "empty input",
			input: make([]byte, 0),
		},
		{
			title:    "non empty input",
			input:    []byte("input"),
			expected: "input",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := NewTabbedHelpWriter(buf)
			w.w = w.w.Init(buf, 1, 1, 1, 0, 1)
			n, err := w.Write(tc.input)
			if err != nil {
				t.Errorf("Write: Expected a nil error, Actual: %v", err)
			}

			err = w.Close()
			if err != nil {
				t.Errorf("Close: Expected a nil error, Actual: %v", err)
			}

			if n != len(tc.expected) {
				t.Errorf("Expected number of written bytes: %d, Actual: %d", len(tc.expected), n)
			}
			if n != len(tc.expected) {
				t.Errorf("Expected number of written bytes: %d, Actual: %d", len(tc.expected), n)
			}

			actual := buf.String()
			if actual != tc.expected {
				t.Errorf("Expected written string: %s, Actual: %s", tc.expected, actual)
			}
		})
	}
}
