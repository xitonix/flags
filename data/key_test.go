package data_test

import (
	"testing"

	"go.xitonix.io/flags/data"
)

func TestKey_SetPrefix(t *testing.T) {
	testCases := []struct {
		title          string
		inputPrefix    string
		expectedPrefix string
	}{
		{
			title:          "empty prefix",
			inputPrefix:    "",
			expectedPrefix: "",
		},
		{
			title:          "prefix casing",
			inputPrefix:    "prefix",
			expectedPrefix: "PREFIX",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			e := &data.Key{}
			e.SetPrefix(tc.inputPrefix)
			actual := e.Prefix()
			if actual != tc.expectedPrefix {
				t.Errorf("Prefix, Expected: %s, Actual:%s", tc.expectedPrefix, actual)
			}
		})
	}
}

func TestKey_Set(t *testing.T) {
	testCases := []struct {
		title         string
		inputName     string
		inputPrefix   string
		expectedName  string
		expectedIsSet bool
	}{
		{
			title:        "empty name with no prefix",
			inputName:    "",
			inputPrefix:  "",
			expectedName: "",
		},
		{
			title:        "empty name with prefix",
			inputName:    "",
			inputPrefix:  "prefix",
			expectedName: "",
		},
		{
			title:         "name with empty prefix",
			inputName:     "name",
			inputPrefix:   "",
			expectedName:  "NAME",
			expectedIsSet: true,
		},
		{
			title:         "name and prefix",
			inputName:     "name",
			inputPrefix:   "prefix",
			expectedName:  "PREFIX_NAME",
			expectedIsSet: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			e := &data.Key{}
			e.SetPrefix(tc.inputPrefix)
			e.SetID(tc.inputName)
			actual := e.String()
			if actual != tc.expectedName {
				t.Errorf("FullString(), Expected: %s, Actual:%s", tc.expectedName, actual)
			}

			if tc.expectedIsSet != e.IsSet() {
				t.Errorf("IsSet(), Expected: %s, Actual:%s", tc.expectedName, actual)
			}
		})
	}
}
