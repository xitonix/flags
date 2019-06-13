package flags_test

import (
	"testing"

	"go.xitonix.io/flags"
	"go.xitonix.io/flags/test"
)

func TestStringSlice(t *testing.T) {
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
			title:         "whitespace usage",
			long:          "long",
			expectedLong:  "long",
			usage:         "     ",
			expectedUsage: "     ",
		},
		{
			title:         "long name with whitespace",
			long:          "   long  ",
			expectedLong:  "long",
			usage:         "     ",
			expectedUsage: "     ",
		},
		{
			title:         "whitespace long name will be validated at parse time",
			long:          "   ",
			expectedLong:  "",
			usage:         "",
			expectedUsage: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSlice(tc.long, tc.usage)
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

			if f.Type() != "[]string" {
				t.Errorf("The flag type was expected to be '[]string', but it was %s", f.Type())
			}

			if !test.StringsEqual(f.Get(), []string{}) {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
		})
	}
}

func TestStringSliceP(t *testing.T) {
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
			title:         "whitespace usage",
			long:          "long",
			expectedLong:  "long",
			usage:         "     ",
			expectedUsage: "     ",
		},
		{
			title:         "long name with whitespace",
			long:          "   long  ",
			expectedLong:  "long",
			usage:         "     ",
			expectedUsage: "     ",
		},
		{
			title:         "whitespace long name will be validated at parse time",
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
			title:         "long and short names with whitespace",
			long:          " Long ",
			expectedLong:  "long",
			short:         " Short ",
			expectedShort: "Short",
		},
		{
			title:         "whitespace long and short names will be validated at parse time",
			long:          "  ",
			expectedLong:  "",
			short:         "    ",
			expectedShort: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSliceP(tc.long, tc.usage, tc.short)
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

			if f.Type() != "[]string" {
				t.Errorf("The flag type was expected to be '[]string', but it was %s", f.Type())
			}

			if !test.StringsEqual(f.Get(), []string{}) {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
		})
	}
}

func TestStringSliceFlag_WithKey(t *testing.T) {
	testCases := []struct {
		title       string
		key         string
		expectedKey string
	}{
		{
			title: "empty key",
		},
		{
			title: "whitespace key",
			key:   "      ",
		},
		{
			title:       "lowercase key",
			key:         "key",
			expectedKey: "KEY",
		},
		{
			title:       "key with whitespace",
			key:         "   key   ",
			expectedKey: "KEY",
		},
		{
			title:       "key with whitespace in the middle",
			key:         "   key with whitespace  ",
			expectedKey: "KEY_WITH_WHITESPACE",
		},
		{
			title:       "key with hyphens",
			key:         "------key-------with-----hyphen----",
			expectedKey: "_KEY_WITH_HYPHEN_",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSlice("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestStringSliceFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         []string
		expectedDefaultValue []string
	}{
		{
			title:                "empty default value",
			defaultValue:         []string{},
			expectedDefaultValue: []string{},
		},
		{
			title:                "non empty default value",
			defaultValue:         []string{"abc"},
			expectedDefaultValue: []string{"abc"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSlice("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if !test.StringsEqual(actual.([]string), tc.expectedDefaultValue) {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestStringSliceFlag_Hide(t *testing.T) {
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
			f := flags.StringSlice("long", "usage")
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

func TestStringSliceFlag_IsDeprecated(t *testing.T) {
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
			f := flags.StringSlice("long", "usage")
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

func TestStringSliceFlag_WithDelimiter(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		delimiter     string
		expectedValue []string
	}{
		{
			title:         "empty delimiter",
			value:         "abc,xyz",
			expectedValue: []string{"abc", "xyz"},
		},
		{
			title:         "whitespace delimiter with white spaced input",
			value:         "abc xyz",
			delimiter:     " ",
			expectedValue: []string{"abc", "xyz"},
		},
		{
			title:         "whitespace delimiter with none white spaced input",
			value:         "abc,xyz",
			delimiter:     " ",
			expectedValue: []string{"abc,xyz"},
		},
		{
			title:         "none whitespace delimiter",
			value:         "abc|xyz",
			delimiter:     "|",
			expectedValue: []string{"abc", "xyz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSlice("long", "usage").WithDelimiter(tc.delimiter)
			fVar := f.Var()
			err := f.Set(tc.value)
			if err != nil {
				t.Errorf("Did not expect to receive an error, but received %s", err)
			}
			actual := f.Get()
			if !test.StringsEqual(actual, tc.expectedValue) {
				t.Errorf("Expected value: %v, Actual: %v", tc.expectedValue, actual)
			}

			if !test.StringsEqual(actual, tc.expectedValue) {
				t.Errorf("Expected flag variable: %v, Actual: %v", tc.expectedValue, *fVar)
			}
		})
	}
}

func TestStringSliceFlag_Set(t *testing.T) {
	empty := make([]string, 0)
	testCases := []struct {
		title         string
		value         string
		expectedValue []string
	}{
		{
			title:         "empty value",
			value:         "",
			expectedValue: empty,
		},
		{
			title:         "whitespace value",
			value:         "   ",
			expectedValue: []string{"   "},
		},
		{
			title:         "value with whitespace",
			value:         "  abc  ",
			expectedValue: []string{"  abc  "},
		},
		{
			title:         "value with no whitespaces",
			value:         "abc",
			expectedValue: []string{"abc"},
		},
		{
			title:         "comma separated value with no whitespace",
			value:         "abc,efg",
			expectedValue: []string{"abc", "efg"},
		},
		{
			title:         "comma separated value with whitespaces",
			value:         " abc , efg ",
			expectedValue: []string{" abc ", " efg "},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSlice("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			if err != nil {
				t.Errorf("Did not expect to receive an error, but received %s", err)
			}
			actual := f.Get()
			if !test.StringsEqual(actual, tc.expectedValue) {
				t.Errorf("Expected value: %v, Actual: %v", tc.expectedValue, actual)
			}

			if !test.StringsEqual(actual, tc.expectedValue) {
				t.Errorf("Expected flag variable: %v, Actual: %v", tc.expectedValue, *fVar)
			}
		})
	}
}

func TestStringSliceFlag_ResetToDefault(t *testing.T) {
	empty := make([]string, 0)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           []string
		defaultValue            []string
		expectedAfterResetValue []string
		expectedError           string
		setDefault              bool
		expectedIsSetAfterReset bool
	}{
		{
			title:                   "reset without defining the default value",
			value:                   "abc",
			expectedValue:           []string{"abc"},
			expectedAfterResetValue: []string{"abc"},
			setDefault:              false,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "reset to empty default value",
			value:                   "abc",
			expectedValue:           []string{"abc"},
			defaultValue:            empty,
			expectedAfterResetValue: empty,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to nil default value",
			value:                   "abc",
			expectedValue:           []string{"abc"},
			defaultValue:            nil,
			expectedAfterResetValue: []string{"abc"},
			setDefault:              true,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "reset to non-empty default value",
			value:                   "abc",
			expectedValue:           []string{"abc"},
			defaultValue:            []string{"abc", "efg"},
			expectedAfterResetValue: []string{"abc", "efg"},
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSlice("long", "usage")
			if tc.setDefault {
				f = f.WithDefault(tc.defaultValue)
			}
			fVar := f.Var()
			err := f.Set(tc.value)
			if !test.ErrorContains(err, tc.expectedError) {
				t.Errorf("Expected to receive an error with '%s', but received %s", tc.expectedError, err)
			}
			actual := f.Get()
			if !test.StringsEqual(actual, tc.expectedValue) {
				t.Errorf("Expected value: %v, Actual: %v", tc.expectedValue, actual)
			}

			if !test.StringsEqual(*fVar, tc.expectedValue) {
				t.Errorf("Expected flag variable: %v, Actual: %v", tc.expectedValue, *fVar)
			}

			f.ResetToDefault()

			if f.IsSet() != tc.expectedIsSetAfterReset {
				t.Errorf("IsSet() Expected: %v, Actual: %v", tc.expectedIsSetAfterReset, f.IsSet())
			}

			actual = f.Get()
			if !test.StringsEqual(actual, tc.expectedAfterResetValue) {
				t.Errorf("Expected value after reset: %v, Actual: %v", tc.expectedAfterResetValue, actual)
			}

			if !test.StringsEqual(*fVar, tc.expectedAfterResetValue) {
				t.Errorf("Expected flag variable after reset: %v, Actual: %v", tc.expectedAfterResetValue, *fVar)
			}
		})
	}
}
