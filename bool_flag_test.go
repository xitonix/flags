package flags_test

import (
	"testing"

	"go.xitonix.io/flags"
	"go.xitonix.io/flags/test"
)

func TestBool(t *testing.T) {
	testCases := []struct {
		title         string
		long          string
		expectedLong  string
		usage         string
		expectedUsage string
	}{
		{
			title:         "lowercase long name with usage",
			long:          "long",
			expectedLong:  "long",
			usage:         "usage",
			expectedUsage: "usage",
		},
		{
			title:         "uppercase long name with usage",
			long:          "LONG",
			expectedLong:  "long",
			usage:         " I must Stay Unchanged   ",
			expectedUsage: " I must Stay Unchanged   ",
		},
		{
			title:         "white space usage",
			long:          "long",
			expectedLong:  "long",
			usage:         "     ",
			expectedUsage: "     ",
		},
		{
			title:         "long name with white space",
			long:          "   long  ",
			expectedLong:  "long",
			usage:         "     ",
			expectedUsage: "     ",
		},
		{
			title:         "white space long name will be validated at parse time",
			long:          "   ",
			expectedLong:  "",
			usage:         "",
			expectedUsage: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Bool(tc.long, tc.usage)
			if f.LongName() != tc.expectedLong {
				t.Errorf("Expected Long Name: %s, Actual: %s", tc.expectedLong, f.LongName())
			}
			if f.Usage() != tc.expectedUsage {
				t.Errorf("Expected Usage: %s, Actual: %s", tc.expectedUsage, f.Usage())
			}

			if f.IsDeprecated() {
				t.Error("The flag must not be marked as deprecated by default")
			}

			if f.IsHidden() {
				t.Error("The flag must not be marked as hidden by default")
			}

			if f.IsSet() {
				t.Error("The flag value must not be set initially")
			}

			if f.ShortName() != "" {
				t.Errorf("The short name was expected to be empty but it was %s", f.ShortName())
			}

			if f.Default() != nil {
				t.Errorf("The initial default value was expected to be nil, but it was %v", f.Default())
			}

			if f.Type() != "bool" {
				t.Errorf("The flag type was expected to be 'bool', but it was %s", f.Type())
			}

			if f.Get() {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
		})
	}
}

func TestBoolP(t *testing.T) {
	testCases := []struct {
		title         string
		long, short   string
		expectedLong  string
		expectedShort string
		usage         string
		expectedUsage string
	}{
		{
			title: "empty long and short names",
		},
		{
			title:         "lowercase long name with usage",
			long:          "long",
			expectedLong:  "long",
			usage:         "usage",
			expectedUsage: "usage",
		},
		{
			title:         "uppercase long name with usage",
			long:          "LONG",
			expectedLong:  "long",
			usage:         " I must Stay Unchanged   ",
			expectedUsage: " I must Stay Unchanged   ",
		},
		{
			title:         "white space usage",
			long:          "long",
			expectedLong:  "long",
			usage:         "     ",
			expectedUsage: "     ",
		},
		{
			title:         "long name with white space",
			long:          "   long  ",
			expectedLong:  "long",
			usage:         "     ",
			expectedUsage: "     ",
		},
		{
			title:         "white space long name will be validated at parse time",
			long:          "   ",
			expectedLong:  "",
			usage:         "",
			expectedUsage: "",
		},
		{
			title:         "lowercase long and short names",
			long:          "long",
			expectedLong:  "long",
			short:         "s",
			expectedShort: "s",
		},
		{
			title:         "uppercase long and short names",
			long:          "Long",
			expectedLong:  "long",
			short:         "Short",
			expectedShort: "Short",
		},
		{
			title:         "long and short names with white space",
			long:          " Long ",
			expectedLong:  "long",
			short:         " Short ",
			expectedShort: "Short",
		},
		{
			title:         "white space long and short names will be validated at parse time",
			long:          "  ",
			expectedLong:  "",
			short:         "    ",
			expectedShort: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.BoolP(tc.long, tc.usage, tc.short)
			if f.LongName() != tc.expectedLong {
				t.Errorf("Expected Long Name: %s, Actual: %s", tc.expectedLong, f.LongName())
			}
			if f.Usage() != tc.expectedUsage {
				t.Errorf("Expected Usage: %s, Actual: %s", tc.expectedUsage, f.Usage())
			}

			if f.IsDeprecated() {
				t.Error("The flag must not be marked as deprecated by default")
			}

			if f.IsHidden() {
				t.Error("The flag must not be marked as hidden by default")
			}

			if f.IsSet() {
				t.Error("The flag value must not be set initially")
			}

			if f.ShortName() != tc.expectedShort {
				t.Errorf("The short name was expected to be %s but it was %s", tc.expectedShort, f.ShortName())
			}

			if f.Default() != nil {
				t.Errorf("The initial default value was expected to be nil, but it was %v", f.Default())
			}

			if f.Type() != "bool" {
				t.Errorf("The flag type was expected to be 'bool', but it was %s", f.Type())
			}

			if f.Get() {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
		})
	}
}

func TestBoolFlag_WithKey(t *testing.T) {
	testCases := []struct {
		title       string
		key         string
		expectedKey string
	}{
		{
			title: "empty key",
		},
		{
			title: "white space key",
			key:   "      ",
		},
		{
			title:       "lowercase key",
			key:         "key",
			expectedKey: "KEY",
		},
		{
			title:       "key with white space",
			key:         "   key   ",
			expectedKey: "KEY",
		},
		{
			title:       "key with white space in the middle",
			key:         "   key with white space  ",
			expectedKey: "KEY_WITH_WHITE_SPACE",
		},
		{
			title:       "key with hyphens",
			key:         "------key-------with-----hyphen----",
			expectedKey: "_KEY_WITH_HYPHEN_",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Bool("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestBoolFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         bool
		expectedDefaultValue bool
	}{
		{
			title:                "false default value",
			defaultValue:         false,
			expectedDefaultValue: false,
		},
		{
			title:                "true default value",
			defaultValue:         true,
			expectedDefaultValue: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Bool("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if actual != tc.expectedDefaultValue {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestBoolFlag_Hide(t *testing.T) {
	testCases := []struct {
		title    string
		isHidden bool
	}{
		{
			title: "visible by default",
		},
		{
			title:    "hidden flag",
			isHidden: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Bool("long", "usage")
			if tc.isHidden {
				f = f.Hide()
			}
			actual := f.IsHidden()
			if actual != tc.isHidden {
				t.Errorf("Expected IsHidden: %v, Actual: %v", tc.isHidden, actual)
			}
		})
	}
}

func TestBoolFlag_IsDeprecated(t *testing.T) {
	testCases := []struct {
		title        string
		isDeprecated bool
	}{
		{
			title: "not deprecated by default",
		},
		{
			title:        "deprecated flag",
			isDeprecated: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Bool("long", "usage")
			if tc.isDeprecated {
				f = f.MarkAsDeprecated()
			}
			actual := f.IsDeprecated()
			if actual != tc.isDeprecated {
				t.Errorf("Expected IsDeprecated: %v, Actual: %v", tc.isDeprecated, actual)
			}
		})
	}
}

func TestBoolFlag_Set(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		expectedValue bool
		expectedError string
	}{
		{
			title:         "no value",
			expectedValue: false,
		},
		{
			title:         "white space value",
			value:         "   ",
			expectedValue: false,
		},
		{
			title:         "true value with white space",
			value:         "  true  ",
			expectedValue: true,
		},
		{
			title:         "false value with white space",
			value:         "  false  ",
			expectedValue: false,
		},
		{
			title:         "false with no white space",
			value:         "false",
			expectedValue: false,
		},
		{
			title:         "true with no white space",
			value:         "true",
			expectedValue: true,
		},
		{
			title:         "one",
			value:         "1",
			expectedValue: true,
		},
		{
			title:         "zero",
			value:         "0",
			expectedValue: false,
		},
		{
			title:         "greater than one",
			value:         "2",
			expectedError: "is not a valid bool value",
		},
		{
			title:         "negative value",
			value:         "-1",
			expectedError: "is not a valid bool value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Bool("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			if !test.ErrorContains(err, tc.expectedError) {
				t.Errorf("Expected to receive an error with '%s', but received %s", tc.expectedError, err)
			}
			actual := f.Get()
			if actual != tc.expectedValue {
				t.Errorf("Expected value: %v, Actual: %v", tc.expectedValue, actual)
			}

			if *fVar != tc.expectedValue {
				t.Errorf("Expected flag variable: %v, Actual: %v", tc.expectedValue, *fVar)
			}
		})
	}
}

func TestBoolFlag_ResetToDefault(t *testing.T) {
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           bool
		defaultValue            bool
		expectedAfterResetValue bool
		expectedError           string
		setDefault              bool
	}{
		{
			title: "no value",
		},
		{
			title:                   "reset without defining the default value",
			value:                   "true",
			expectedValue:           true,
			expectedAfterResetValue: true,
			setDefault:              false,
		},
		{
			title:                   "reset to false default value",
			value:                   "true",
			expectedValue:           true,
			defaultValue:            false,
			expectedAfterResetValue: false,
			setDefault:              true,
		},
		{
			title:                   "reset to true default value",
			value:                   "true",
			expectedValue:           true,
			defaultValue:            true,
			expectedAfterResetValue: true,
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Bool("long", "usage")
			if tc.setDefault {
				f = f.WithDefault(tc.defaultValue)
			}
			fVar := f.Var()
			err := f.Set(tc.value)
			if !test.ErrorContains(err, tc.expectedError) {
				t.Errorf("Expected to receive an error with '%s', but received %s", tc.expectedError, err)
			}
			actual := f.Get()
			if actual != tc.expectedValue {
				t.Errorf("Expected value: %v, Actual: %v", tc.expectedValue, actual)
			}

			if *fVar != tc.expectedValue {
				t.Errorf("Expected flag variable: %v, Actual: %v", tc.expectedValue, *fVar)
			}

			f.ResetToDefault()

			if tc.setDefault && f.IsSet() {
				t.Error("IsSet() Expected: false, Actual: true")
			}

			actual = f.Get()
			if actual != tc.expectedAfterResetValue {
				t.Errorf("Expected value after reset: %v, Actual: %v", tc.expectedAfterResetValue, actual)
			}

			if *fVar != tc.expectedAfterResetValue {
				t.Errorf("Expected flag variable after reset: %v, Actual: %v", tc.expectedAfterResetValue, *fVar)
			}
		})
	}
}
