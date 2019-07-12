package core_test

import (
	"errors"
	"testing"

	"github.com/xitonix/flags"
	"github.com/xitonix/flags/core"
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
			checkFlagInitialState(t, f, "bool", tc.expectedUsage, tc.expectedLong, "")
			checkFlagValues(t, false, f.Get(), f.Var())
		})
	}
}

func TestBoolFlag_WithShort(t *testing.T) {
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
			short:         "S",
			expectedShort: "S",
		},
		{
			title:         "long and short names with white space",
			long:          " Long ",
			expectedLong:  "long",
			short:         " S ",
			expectedShort: "S",
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
			f := flags.Bool(tc.long, tc.usage).WithShort(tc.short)
			checkFlagInitialState(t, f, "bool", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkFlagValues(t, false, f.Get(), f.Var())
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

func TestBoolFlag_IsRequired(t *testing.T) {
	testCases := []struct {
		title      string
		isRequired bool
	}{
		{
			title: "not required by default",
		},
		{
			title:      "required flag",
			isRequired: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Bool("long", "usage")
			if tc.isRequired {
				f = f.Required()
			}
			actual := f.IsRequired()
			if actual != tc.isRequired {
				t.Errorf("Expected IsRequired: %v, Actual: %v", tc.isRequired, actual)
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
			checkFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestBoolFlag_Validation(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		expectedValue bool
		validationCB  func(in bool) error
		expectedError string
	}{
		{
			title:         "nil validation callback",
			value:         "true",
			expectedValue: true,
			expectedError: "",
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in bool) error {
				return nil
			},
			value:         "true",
			expectedValue: true,
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in bool) error {
				return errors.New("validation callback failed")
			},
			value:         "true",
			expectedError: "validation callback failed",
			expectedValue: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Bool("long", "usage")
			fVar := f.Var()
			f = f.WithValidationCallback(tc.validationCB)
			err := f.Set(tc.value)
			checkFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
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
			checkFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)

			f.ResetToDefault()

			if tc.setDefault && f.IsSet() {
				t.Error("IsSet() Expected: false, Actual: true")
			}

			checkFlagValues(t, tc.expectedAfterResetValue, f.Get(), fVar)
		})
	}
}

func TestBoolFlag_EmptyValue(t *testing.T) {
	var f core.Flag = core.NewBool("bool", "usage")
	fb, ok := f.(core.EmptyValueProvider)
	if !ok {
		t.Error("a boolean flag must implement core.EmptyValueProvider interface")
	}
	actual := fb.EmptyValue()
	if actual != "true" {
		t.Errorf("Expected Empty Value: true, Actual: %s", actual)
	}
}
