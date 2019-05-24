package flags

import (
	"testing"

	"go.xitonix.io/flags/config"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/test"
)

func TestBucket_Parse(t *testing.T) {
	testCases := []struct {
		title string
		args  []string
		flags []core.Flag
	}{
		{
			title: "unknown flags",
			args:  []string{"-flag"},
			flags: []core.Flag{newMockedFlag("unexpected", "s")},
		},
	}

	helpProvider := &config.HelpProvider{
		Writer:    &test.NullWriter{},
		Formatter: &config.TabbedHelpFormatter{},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
		})
	}
}
