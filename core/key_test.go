package core_test

import (
	"testing"

	"github.com/xitonix/flags/core"
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
			e := &core.Key{}
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
			title:         "dash with no prefix",
			inputName:     "-",
			inputPrefix:   "",
			expectedName:  "",
			expectedIsSet: true,
		},
		{
			title:         "dash with prefix",
			inputName:     "-",
			inputPrefix:   "PREFIX",
			expectedName:  "",
			expectedIsSet: true,
		},
		{
			title:         "dash with white space prefix",
			inputName:     "-",
			inputPrefix:   "  ",
			expectedName:  "",
			expectedIsSet: true,
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
			e := &core.Key{}
			e.SetPrefix(tc.inputPrefix)
			e.SetID(tc.inputName)
			actual := e.String()
			if actual != tc.expectedName {
				t.Errorf("String(), Expected: %s, Actual:%s", tc.expectedName, actual)
			}

			actualIsSet := e.IsSet()
			if tc.expectedIsSet != actualIsSet {
				t.Errorf("IsSet(), Expected: %v, Actual:%v", tc.expectedIsSet, actualIsSet)
			}
		})
	}
}
