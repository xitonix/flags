package core_test

import (
	"strings"
	"testing"

	"go.xitonix.io/flags/core"
)

func TestErrInvalidFlag_Error(t *testing.T) {
	testCases := []struct {
		title            string
		long, short, key string
		msg              string
	}{
		{
			title: "flag with long name",
			long:  "--long",
			msg:   "error message",
		},
		{
			title: "flag with short name only_this should never happen",
			short: "-short",
			msg:   "error message",
		},
		{
			title: "flag with long and short names",
			long:  "--long",
			short: "-short",
			msg:   "error message",
		},
		{
			title: "flag with long and short names along with a key",
			long:  "--long",
			short: "-short",
			key:   "FLAG_KEY",
			msg:   "error message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := core.NewInvalidFlagErr(tc.long, tc.short, tc.key, tc.msg)
			actual := err.Error()
			if !strings.HasSuffix(actual, tc.msg) {
				t.Errorf("Expected suffix: %v, Actual: %s", tc.msg, actual)
			}

			if tc.key != "" {
				if !strings.Contains(actual, tc.key) {
					t.Errorf("Expected %v, Actual: %s", tc.key, actual)
				}
				if strings.Contains(actual, tc.long) {
					t.Errorf("Did not expect to see --%v, Actual: %s", tc.long, actual)
				}
				if strings.Contains(actual, tc.short) {
					t.Errorf("Did not expect to see -%v, Actual: %s", tc.short, actual)
				}
				return
			}

			if tc.long != "" && !strings.Contains(actual, tc.long) {
				t.Errorf("Expected --%v, Actual: %s", tc.long, actual)
			}

			if tc.short != "" && !strings.Contains(actual, tc.short) {
				t.Errorf("Expected -%v, Actual: %s", tc.short, actual)
			}
		})
	}
}
