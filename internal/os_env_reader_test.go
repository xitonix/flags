package internal_test

import (
	"go.xitonix.io/flags/internal"
	"os"
	"testing"
)

func TestOSEnvReader_Get(t *testing.T) {
	testCases := []struct {
		title      string
		key, value string
	}{
		{
			title: "empty value",
			key:   "FLAGS_TEST",
		},
		{
			title: "none empty value",
			value: "ENV VALUE",
			key:   "FLAGS_TEST",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			reader := internal.OSEnvReader{}
			_ = os.Setenv(tc.key, tc.value)
			defer func() {
				_ = os.Unsetenv(tc.key)
			}()

			actual, ok := reader.Get(tc.key)
			if !ok {
				t.Error("Expected the environment variable to exist")
			}
			if actual != tc.value {
				t.Errorf("Expected Env Variable: %s, Actual: %s", tc.value, actual)
			}
		})
	}
}
