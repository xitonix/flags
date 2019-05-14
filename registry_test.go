package flags

import (
	"strings"
	"testing"

	"go.xitonix.io/flags/core"
)

func TestRegistry_Add(t *testing.T) {
	testCases := []struct {
		title       string
		f           core.Flag
		expectedErr string
	}{
		{
			title: "",
			f:     newGeneric("flag", "f"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			reg := newRegistry()
			err := reg.add(tc.f)
			if tc.expectedErr == "" && err == nil {
				return
			}

			if err != nil {
				switch {
				case tc.expectedErr == "":
					t.Errorf("Did not expect an error, but received %v", err)
				case !strings.Contains(err.Error(), tc.expectedErr):
					t.Errorf("Expected to get an error containing '%s' string, but received %v", tc.expectedErr, err)
				}
				return
			}
			if tc.expectedErr != "" {
				t.Errorf("Expected to get an error containing '%s' string, but received nil", tc.expectedErr)
			}
		})
	}
}
