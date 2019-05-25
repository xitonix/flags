package core_test

import (
	"testing"

	"go.xitonix.io/flags/core"
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

func TestKey_Auto(t *testing.T) {
	testCases := []struct {
		title        string
		inputName    string
		inputPrefix  string
		expectedName string
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
			title:        "name with empty prefix",
			inputName:    "name",
			inputPrefix:  "",
			expectedName: "NAME",
		},
		{
			title:        "name and prefix",
			inputName:    "name",
			inputPrefix:  "prefix",
			expectedName: "PREFIX_NAME",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			e := &core.Key{}
			e.SetID(tc.inputName, true)
			e.SetPrefix(tc.inputPrefix)
			actual := e.Value()
			if actual != tc.expectedName {
				t.Errorf("Value, Expected: %s, Actual:%s", tc.expectedName, actual)
			}
		})
	}
}

func TestKey_Set(t *testing.T) {
	testCases := []struct {
		title        string
		inputName    string
		inputPrefix  string
		expectedName string
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
			title:        "name with empty prefix",
			inputName:    "name",
			inputPrefix:  "",
			expectedName: "NAME",
		},
		{
			title:        "name and prefix",
			inputName:    "name",
			inputPrefix:  "prefix",
			expectedName: "PREFIX_NAME",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			e := &core.Key{}
			e.SetID(tc.inputName, false)
			e.SetPrefix(tc.inputPrefix)
			actual := e.Value()
			if actual != tc.expectedName {
				t.Errorf("Value, Expected: %s, Actual:%s", tc.expectedName, actual)
			}
		})
	}
}

func TestKey_Set_Overrides_Auto(t *testing.T) {
	testCases := []struct {
		title         string
		inputAutoName string
		inputName     string
		inputPrefix   string
		expectedName  string
	}{
		{
			title:        "empty names and prefix",
			inputName:    "",
			inputPrefix:  "",
			expectedName: "",
		},
		{
			title:        "empty names with prefix",
			inputName:    "",
			inputPrefix:  "prefix",
			expectedName: "",
		},
		{
			title:         "explicit name with empty auto name and empty prefix",
			inputName:     "name",
			inputAutoName: "",
			inputPrefix:   "",
			expectedName:  "NAME",
		},
		{
			title:         "explicit name with auto name and empty prefix",
			inputName:     "name",
			inputAutoName: "auto",
			inputPrefix:   "",
			expectedName:  "NAME",
		},
		{
			title:         "explicit name with auto name and prefix",
			inputName:     "name",
			inputAutoName: "auto",
			inputPrefix:   "prefix",
			expectedName:  "PREFIX_NAME",
		},
		{
			title:         "empty explicit name with auto name and empty prefix",
			inputName:     "",
			inputAutoName: "auto",
			inputPrefix:   "",
			expectedName:  "",
		},
		{
			title:         "empty explicit name with auto name and prefix",
			inputName:     "",
			inputAutoName: "auto",
			inputPrefix:   "prefix",
			expectedName:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			e := &core.Key{}
			e.SetPrefix(tc.inputPrefix)
			e.SetID(tc.inputAutoName, true)
			e.SetID(tc.inputName, false)
			actual := e.Value()
			if actual != tc.expectedName {
				t.Errorf("Value, Expected: %s, Actual:%s", tc.expectedName, actual)
			}
		})
	}
}

func TestKey_Auto_Does_Not_Override_Set(t *testing.T) {
	testCases := []struct {
		title         string
		inputAutoName string
		inputName     string
		inputPrefix   string
		expectedName  string
	}{
		{
			title:         "empty explicit name with auto name and empty prefix",
			inputName:     "",
			inputAutoName: "auto",
			inputPrefix:   "",
			expectedName:  "",
		},
		{
			title:         "empty explicit name with auto name and prefix",
			inputName:     "",
			inputAutoName: "auto",
			inputPrefix:   "prefix",
			expectedName:  "",
		},
		{
			title:         "explicit name with auto name and empty prefix",
			inputName:     "name",
			inputAutoName: "auto",
			inputPrefix:   "",
			expectedName:  "NAME",
		},
		{
			title:         "explicit name with auto name and prefix",
			inputName:     "name",
			inputAutoName: "auto",
			inputPrefix:   "prefix",
			expectedName:  "PREFIX_NAME",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			e := &core.Key{}
			e.SetPrefix(tc.inputPrefix)
			e.SetID(tc.inputName, false)
			e.SetID(tc.inputAutoName, true)
			actual := e.Value()
			if actual != tc.expectedName {
				t.Errorf("Value, Expected: %s, Actual:%s", tc.expectedName, actual)
			}
		})
	}
}
