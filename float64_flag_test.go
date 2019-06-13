package flags_test

import (
	"fmt"
	"math"
	"testing"

	"go.xitonix.io/flags"
	"go.xitonix.io/flags/test"
)

func TestFloat64(t *testing.T) {
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
			f := flags.Float64(tc.long, tc.usage)
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

			if f.Type() != "float64" {
				t.Errorf("The flag type was expected to be 'float64', but it was %s", f.Type())
			}

			if f.Get() != 0.0 {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
		})
	}
}

func TestFloat64P(t *testing.T) {
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
			f := flags.Float64P(tc.long, tc.usage, tc.short)
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

			if f.Type() != "float64" {
				t.Errorf("The flag type was expected to be 'float64', but it was %s", f.Type())
			}

			if f.Get() != 0.0 {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
		})
	}
}

func TestFloat64Flag_WithKey(t *testing.T) {
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
			f := flags.Float64("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestFloat64Flag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         float64
		expectedDefaultValue float64
	}{
		{
			title:                "zero default value",
			defaultValue:         0.0,
			expectedDefaultValue: 0.0,
		},
		{
			title:                "non zero default value",
			defaultValue:         100.5,
			expectedDefaultValue: 100.5,
		},
		{
			title:                "negative default value",
			defaultValue:         -100.5,
			expectedDefaultValue: -100.5,
		},
		{
			title:                "max float64",
			defaultValue:         math.MaxFloat64,
			expectedDefaultValue: math.MaxFloat64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Float64("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if actual != tc.expectedDefaultValue {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestFloat64Flag_Hide(t *testing.T) {
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
			f := flags.Float64("long", "usage")
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

func TestFloat64Flag_IsDeprecated(t *testing.T) {
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
			f := flags.Float64("long", "usage")
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

func TestFloat64Flag_Set(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		expectedValue float64
		expectedError string
	}{
		{
			title:         "no value",
			expectedValue: 0.0,
		},
		{
			title:         "white space value",
			value:         "   ",
			expectedValue: 0.0,
		},
		{
			title:         "value with white space",
			value:         "  100.5  ",
			expectedValue: 100.5,
		},
		{
			title:         "negative value",
			value:         "-100.5",
			expectedValue: -100.5,
		},
		{
			title:         "value with no white space",
			value:         "100.5",
			expectedValue: 100.5,
		},
		{
			title:         "invalid value",
			value:         "abc",
			expectedError: "is not a valid float64 value",
			expectedValue: 0.0,
		},
		{
			title:         "max float64",
			value:         fmt.Sprintf("%f", math.MaxFloat64),
			expectedValue: math.MaxFloat64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Float64("long", "usage")
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

func TestFloat64Flag_ResetToDefault(t *testing.T) {
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           float64
		defaultValue            float64
		expectedAfterResetValue float64
		expectedError           string
		setDefault              bool
	}{
		{
			title: "no value",
		},
		{
			title:                   "reset without defining the default value",
			value:                   "100.5",
			expectedValue:           100.5,
			expectedAfterResetValue: 100.5,
			setDefault:              false,
		},
		{
			title:                   "reset to zero default value",
			value:                   "100.5",
			expectedValue:           100.5,
			defaultValue:            0.0,
			expectedAfterResetValue: 0.0,
			setDefault:              true,
		},
		{
			title:                   "reset to non-zero default value",
			value:                   "100.5",
			expectedValue:           100.5,
			defaultValue:            200.06,
			expectedAfterResetValue: 200.06,
			setDefault:              true,
		},
		{
			title:                   "reset to negative default value",
			value:                   "100.5",
			expectedValue:           100.5,
			defaultValue:            -200.06,
			expectedAfterResetValue: -200.06,
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Float64("long", "usage")
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
