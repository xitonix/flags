package flags_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"go.xitonix.io/flags"
)

func TestDurationSlice(t *testing.T) {
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
			f := flags.DurationSlice(tc.long, tc.usage)
			checkFlagInitialState(t, f, "[]duration", tc.expectedUsage, tc.expectedLong, "")
			checkSliceFlagValues(t, []time.Duration{}, f.Get(), f.Var())
		})
	}
}

func TestDurationSliceP(t *testing.T) {
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
			f := flags.DurationSliceP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "[]duration", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkSliceFlagValues(t, []time.Duration{}, f.Get(), f.Var())
		})
	}
}

func TestDurationSliceFlag_WithKey(t *testing.T) {
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
			f := flags.DurationSlice("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestDurationSliceFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         []time.Duration
		expectedDefaultValue []time.Duration
	}{
		{
			title:                "empty default value",
			defaultValue:         []time.Duration{},
			expectedDefaultValue: []time.Duration{},
		},
		{
			title:                "non empty default value",
			defaultValue:         []time.Duration{2 * time.Second},
			expectedDefaultValue: []time.Duration{2 * time.Second},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.DurationSlice("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if !reflect.DeepEqual(actual.([]time.Duration), tc.expectedDefaultValue) {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestDurationSliceFlag_Hide(t *testing.T) {
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
			f := flags.DurationSlice("long", "usage")
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

func TestDurationSliceFlag_IsDeprecated(t *testing.T) {
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
			f := flags.DurationSlice("long", "usage")
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

func TestDurationSliceFlag_IsRequired(t *testing.T) {
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
			f := flags.DurationSlice("long", "usage")
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

func TestDurationSliceFlag_WithDelimiter(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		delimiter     string
		expectedValue []time.Duration
	}{
		{
			title:         "empty delimiter",
			value:         "2s,3s",
			expectedValue: []time.Duration{2 * time.Second, 3 * time.Second},
		},
		{
			title:         "white space delimiter with white spaced input",
			value:         "2s 3s",
			delimiter:     " ",
			expectedValue: []time.Duration{2 * time.Second, 3 * time.Second},
		},
		{
			title:         "none white space delimiter",
			value:         "2s|3s",
			delimiter:     "|",
			expectedValue: []time.Duration{2 * time.Second, 3 * time.Second},
		},
		{
			title:         "no delimited input",
			value:         "2s",
			delimiter:     "|",
			expectedValue: []time.Duration{2 * time.Second},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.DurationSlice("long", "usage").WithDelimiter(tc.delimiter)
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestDurationSliceFlag_Set(t *testing.T) {
	empty := make([]time.Duration, 0)
	testCases := []struct {
		title         string
		value         string
		expectedValue []time.Duration
		expectedError string
	}{
		{
			title:         "empty value",
			value:         "",
			expectedValue: empty,
		},
		{
			title:         "white space value",
			value:         "   ",
			expectedValue: empty,
		},
		{
			title:         "single value with white space",
			value:         "  2s  ",
			expectedValue: []time.Duration{2 * time.Second},
		},
		{
			title:         "single value with no white space",
			value:         "2s",
			expectedValue: []time.Duration{2 * time.Second},
		},
		{
			title:         "comma separated value with no white space",
			value:         "0s,2s,3s",
			expectedValue: []time.Duration{0, 2 * time.Second, 3 * time.Second},
		},
		{
			title:         "comma separated value with white space",
			value:         " 0, 2s   ,   3s ",
			expectedValue: []time.Duration{0, 2 * time.Second, 3 * time.Second},
		},
		{
			title:         "comma separated empty string",
			value:         ",,",
			expectedValue: empty,
		},
		{
			title:         "comma separated white space string",
			value:         " , , ",
			expectedValue: empty,
		},
		{
			title:         "invalid value",
			value:         " invalid ",
			expectedError: "'invalid' is not a valid []duration value for --long",
			expectedValue: empty,
		},
		{
			title:         "partially invalid value",
			value:         "2s,invalid,3s",
			expectedError: "'invalid' is not a valid []duration value for --long",
			expectedValue: empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.DurationSlice("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestDurationSliceFlag_Validation(t *testing.T) {
	empty := make([]time.Duration, 0)
	testCases := []struct {
		title             string
		value             string
		expectedValue     []time.Duration
		validationCB      func(in time.Duration) error
		setValidationCB   bool
		validationList    []time.Duration
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback",
			setValidationCB: true,
			value:           "2s, 1s",
			expectedValue:   []time.Duration{2 * time.Second, 1 * time.Second},
			expectedError:   "",
		},
		{
			title:             "nil validation list",
			setValidationList: true,
			value:             "2s, 1s",
			expectedValue:     []time.Duration{2 * time.Second, 1 * time.Second},
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             "2s, 1s",
			expectedValue:     []time.Duration{2 * time.Second, 1 * time.Second},
			expectedError:     "",
		},
		{
			title:             "empty validation list",
			validationList:    make([]time.Duration, 0),
			setValidationList: true,
			value:             "2s, 1s",
			expectedValue:     []time.Duration{2 * time.Second, 1 * time.Second},
			expectedError:     "",
		},
		{
			title:             "none empty validation list with single item",
			validationList:    []time.Duration{2 * time.Second},
			setValidationList: true,
			value:             "1s",
			expectedError:     "1s is not an acceptable value for --durations. The expected value is 2s.",
			expectedValue:     empty,
		},
		{
			title:             "none empty validation list with two items",
			validationList:    []time.Duration{2 * time.Second, 3 * time.Second},
			setValidationList: true,
			value:             "1s",
			expectedError:     "1s is not an acceptable value for --durations. The expected values are 2s,3s.",
			expectedValue:     empty,
		},
		{
			title:             "duplicate items in the validation list",
			validationList:    []time.Duration{3 * time.Second, 3 * time.Second},
			setValidationList: true,
			value:             "1s",
			expectedError:     "1s is not an acceptable value for --durations. The expected value is 3s.",
			expectedValue:     empty,
		},
		{
			title:             "validation list with three entries",
			validationList:    []time.Duration{3 * time.Second, 2 * time.Second, 1 * time.Second},
			setValidationList: true,
			value:             "4s",
			expectedError:     "4s is not an acceptable value for --durations. The expected values are 3s,2s,1s.",
			expectedValue:     empty,
		},
		{
			title:             "none empty validation list with valid value",
			validationList:    []time.Duration{2 * time.Second, 3 * time.Second},
			setValidationList: true,
			value:             "2s",
			expectedError:     "",
			expectedValue:     []time.Duration{2 * time.Second},
		},
		{
			title:             "empty value",
			validationList:    []time.Duration{2 * time.Second},
			setValidationList: true,
			value:             "",
			expectedError:     "",
			expectedValue:     empty,
		},
		{
			title:             "white space value",
			validationList:    []time.Duration{2 * time.Second},
			setValidationList: true,
			value:             "  ",
			expectedValue:     empty,
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in time.Duration) error {
				return nil
			},
			setValidationCB: true,
			value:           "2s",
			expectedValue:   []time.Duration{2 * time.Second},
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in time.Duration) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "2s",
			expectedError:   "validation callback failed",
			expectedValue:   empty,
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in time.Duration) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []time.Duration{2 * time.Second, 3 * time.Second},
			setValidationList: true,
			value:             "2s",
			expectedError:     "validation callback failed",
			expectedValue:     empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.DurationSlice("durations", "usage")
			fVar := f.Var()
			if tc.setValidationCB {
				f = f.WithValidationCallback(tc.validationCB)
			}
			if tc.setValidationList {
				f = f.WithValidRange(tc.validationList...)
			}
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestDurationSliceFlag_ResetToDefault(t *testing.T) {
	empty := make([]time.Duration, 0)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           []time.Duration
		defaultValue            []time.Duration
		expectedAfterResetValue []time.Duration
		expectedError           string
		setDefault              bool
		expectedIsSetAfterReset bool
	}{
		{
			title:                   "reset without defining the default value",
			value:                   "2s",
			expectedValue:           []time.Duration{2 * time.Second},
			expectedAfterResetValue: []time.Duration{2 * time.Second},
			setDefault:              false,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "reset to empty default value",
			value:                   "2s",
			expectedValue:           []time.Duration{2 * time.Second},
			defaultValue:            empty,
			expectedAfterResetValue: empty,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to nil default value",
			value:                   "2s",
			expectedValue:           []time.Duration{2 * time.Second},
			defaultValue:            nil,
			expectedAfterResetValue: nil,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to non-empty default value",
			value:                   "2s",
			expectedValue:           []time.Duration{2 * time.Second},
			defaultValue:            []time.Duration{2 * time.Second, 3 * time.Second},
			expectedAfterResetValue: []time.Duration{2 * time.Second, 3 * time.Second},
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.DurationSlice("long", "usage")
			if tc.setDefault {
				f = f.WithDefault(tc.defaultValue)
			}
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)

			f.ResetToDefault()

			if f.IsSet() != tc.expectedIsSetAfterReset {
				t.Errorf("IsSet() Expected: %v, Actual: %v", tc.expectedIsSetAfterReset, f.IsSet())
			}

			checkSliceFlagValues(t, tc.expectedAfterResetValue, f.Get(), fVar)
		})
	}
}
