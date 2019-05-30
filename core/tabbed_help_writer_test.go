package core_test

import (
	"go.xitonix.io/flags/core"
	"testing"
)

func TestTabbedHelpWriter_Write_Empty(t *testing.T) {
	testCases := []struct {
		title string
		input []byte
	}{
		{
			title: "nil input",
			input: nil,
		},
		{
			title: "empty input",
			input: make([]byte, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := core.NewTabbedHelpWriter()
			n, err := w.Write(tc.input)
			if err != nil {
				t.Errorf("Expected a nil error, Actual: %v", err)
			}
			if n > 0 {
				t.Errorf("Expected number of written bytes: 0, Actual: %d", n)
			}
		})
	}
}
