package flags_test

import (
	"testing"
	"time"

	"go.xitonix.io/flags"
	"go.xitonix.io/flags/test"
)

func TestTime(t *testing.T) {
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
			f := flags.Time(tc.long, tc.usage)
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

			if f.Type() != "time" {
				t.Errorf("The flag type was expected to be 'time.Time', but it was %s", f.Type())
			}

			if !f.Get().IsZero() {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
		})
	}
}

func TestTimeP(t *testing.T) {
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
			f := flags.TimeP(tc.long, tc.usage, tc.short)
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

			if f.Type() != "time" {
				t.Errorf("The flag type was expected to be 'time.Time', but it was %s", f.Type())
			}

			if !f.Get().IsZero() {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
		})
	}
}

func TestTimeFlag_WithKey(t *testing.T) {
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
			f := flags.Time("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestTimeFlag_WithDefault(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		title                string
		defaultValue         time.Time
		expectedDefaultValue time.Time
	}{
		{
			title:                "zero default value",
			defaultValue:         time.Time{},
			expectedDefaultValue: time.Time{},
		},
		{
			title:                "non zero default value",
			defaultValue:         now,
			expectedDefaultValue: now,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if actual != tc.expectedDefaultValue {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestTimeFlag_Hide(t *testing.T) {
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
			f := flags.Time("long", "usage")
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

func TestTimeFlag_IsDeprecated(t *testing.T) {
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
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Set(t *testing.T) {
	zero := time.Time{}
	testCases := []struct {
		title         string
		value         string
		expectedValue time.Time
		expectedError string
	}{
		{
			title:         "no value",
			expectedValue: zero,
		},
		{
			title:         "whitespace value",
			value:         "   ",
			expectedValue: zero,
		},
		{
			title:         "value with whitespace",
			value:         "  2020-01-15T00:00:00.000000000Z00:00  ",
			expectedValue: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			title:         "value with no whitespaces",
			value:         "2020-01-15",
			expectedValue: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			title:         "invalid value",
			value:         "abc",
			expectedError: "is not a valid time value",
			expectedValue: zero,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_ResetToDefault(t *testing.T) {
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           time.Time
		defaultValue            time.Time
		expectedAfterResetValue time.Time
		expectedError           string
		setDefault              bool
	}{
		{
			title: "no value",
		},
		{
			title:                   "reset without defining the default value",
			value:                   "2020-01-15",
			expectedValue:           time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			expectedAfterResetValue: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			setDefault:              false,
		},
		{
			title:                   "reset to zero default value",
			value:                   "2020-01-15",
			expectedValue:           time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			defaultValue:            time.Time{},
			expectedAfterResetValue: time.Time{},
			setDefault:              true,
		},
		{
			title:                   "reset to non-zero default value",
			value:                   "2020-01-15",
			expectedValue:           time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			defaultValue:            time.Date(2020, 10, 17, 0, 0, 0, 0, time.UTC),
			expectedAfterResetValue: time.Date(2020, 10, 17, 0, 0, 0, 0, time.UTC),
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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
