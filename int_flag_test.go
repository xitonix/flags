package flags_test

import (
	"testing"

	"go.xitonix.io/flags"
)

func TestInt(t *testing.T) {
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
			f := flags.Int(tc.long, tc.usage)
			checkFlagInitialState(t, f, "int", tc.expectedUsage, tc.expectedLong, "")
			checkFlagValues(t, 0, f.Get(), f.Var())
		})
	}
}

func TestIntP(t *testing.T) {
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
			f := flags.IntP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "int", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkFlagValues(t, 0, f.Get(), f.Var())
		})
	}
}

func TestIntFlag_WithKey(t *testing.T) {
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
			f := flags.Int("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestIntFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         int
		expectedDefaultValue int
	}{
		{
			title:                "zero default value",
			defaultValue:         0,
			expectedDefaultValue: 0,
		},
		{
			title:                "non zero default value",
			defaultValue:         100,
			expectedDefaultValue: 100,
		},
		{
			title:                "negative default value",
			defaultValue:         -100,
			expectedDefaultValue: -100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Int("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if actual != tc.expectedDefaultValue {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestIntFlag_Hide(t *testing.T) {
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
			f := flags.Int("long", "usage")
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

func TestIntFlag_IsDeprecated(t *testing.T) {
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
			f := flags.Int("long", "usage")
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

func TestIntFlag_Set(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		expectedValue int
		expectedError string
	}{
		{
			title:         "no value",
			expectedValue: 0,
		},
		{
			title:         "white space value",
			value:         "   ",
			expectedValue: 0,
		},
		{
			title:         "value with white space",
			value:         "  100  ",
			expectedValue: 100,
		},
		{
			title:         "value with no white space",
			value:         "100",
			expectedValue: 100,
		},
		{
			title:         "negative value",
			value:         "-100",
			expectedValue: -100,
		},
		{
			title:         "invalid value",
			value:         "abc",
			expectedError: "is not a valid int value",
			expectedValue: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Int("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestIntFlag_ResetToDefault(t *testing.T) {
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           int
		defaultValue            int
		expectedAfterResetValue int
		expectedError           string
		setDefault              bool
	}{
		{
			title: "no value",
		},
		{
			title:                   "reset without defining the default value",
			value:                   "100",
			expectedValue:           100,
			expectedAfterResetValue: 100,
			setDefault:              false,
		},
		{
			title:                   "reset to zero default value",
			value:                   "100",
			expectedValue:           100,
			defaultValue:            0,
			expectedAfterResetValue: 0,
			setDefault:              true,
		},
		{
			title:                   "reset to non-zero default value",
			value:                   "100",
			expectedValue:           100,
			defaultValue:            200,
			expectedAfterResetValue: 200,
			setDefault:              true,
		},
		{
			title:                   "reset to negative default value",
			value:                   "100",
			expectedValue:           100,
			defaultValue:            -200,
			expectedAfterResetValue: -200,
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Int("long", "usage")
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
