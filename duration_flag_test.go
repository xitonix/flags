package flags_test

import (
	"errors"
	"testing"
	"time"

	"go.xitonix.io/flags"
)

func TestDuration(t *testing.T) {
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
			f := flags.Duration(tc.long, tc.usage)
			checkFlagInitialState(t, f, "duration", tc.expectedUsage, tc.expectedLong, "")
			checkFlagValues(t, time.Duration(0), f.Get(), f.Var())
		})
	}
}

func TestDurationP(t *testing.T) {
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
			f := flags.DurationP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "duration", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkFlagValues(t, time.Duration(0), f.Get(), f.Var())
		})
	}
}

func TestDurationFlag_WithKey(t *testing.T) {
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
			f := flags.Duration("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestDurationFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         time.Duration
		expectedDefaultValue time.Duration
	}{
		{
			title:                "zero default value",
			defaultValue:         0,
			expectedDefaultValue: 0,
		},
		{
			title:                "non zero default value",
			defaultValue:         1 * time.Second,
			expectedDefaultValue: 1 * time.Second,
		},
		{
			title:                "negative default value",
			defaultValue:         -1 * time.Second,
			expectedDefaultValue: -1 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Duration("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if actual != tc.expectedDefaultValue {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestDurationFlag_Hide(t *testing.T) {
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
			f := flags.Duration("long", "usage")
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

func TestDurationFlag_IsDeprecated(t *testing.T) {
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
			f := flags.Duration("long", "usage")
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

func TestDurationFlag_IsRequired(t *testing.T) {
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
			f := flags.Duration("long", "usage")
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

func TestDurationFlag_Set(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		expectedValue time.Duration
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
			value:         "  10s  ",
			expectedValue: 10 * time.Second,
		},
		{
			title:         "numeric value",
			value:         "100",
			expectedError: "is not a valid duration value",
			expectedValue: 0,
		},
		{
			title:         "value with no white space",
			value:         "2s",
			expectedValue: 2 * time.Second,
		},
		{
			title:         "invalid value",
			value:         "abc",
			expectedError: "is not a valid duration value",
			expectedValue: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Duration("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestDurationFlag_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Duration
		validationCB      func(in time.Duration) error
		setValidationCB   bool
		validationList    []time.Duration
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback",
			setValidationCB: true,
			value:           "100s",
			expectedValue:   100 * time.Second,
			expectedError:   "",
		},
		{
			title:             "nil validation list",
			setValidationList: true,
			value:             "100s",
			expectedValue:     100 * time.Second,
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             "100s",
			expectedValue:     100 * time.Second,
			expectedError:     "",
		},
		{
			title:             "empty validation list",
			validationList:    make([]time.Duration, 0),
			setValidationList: true,
			value:             "100s",
			expectedValue:     100 * time.Second,
			expectedError:     "",
		},
		{
			title:             "single item in the validation list",
			validationList:    []time.Duration{1 * time.Second},
			setValidationList: true,
			value:             "200s",
			expectedError:     "200s is not an acceptable value for --long. The expected value is 1s.",
		},
		{
			title:             "two items in the validation list",
			validationList:    []time.Duration{1 * time.Second, 2 * time.Second},
			setValidationList: true,
			value:             "300s",
			expectedError:     "300s is not an acceptable value for --long. The expected values are 1s,2s.",
		},
		{
			title:             "three items in the validation list",
			validationList:    []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second},
			setValidationList: true,
			value:             "400s",
			expectedError:     "400s is not an acceptable value for --long. The expected values are 1s,2s,3s.",
		},
		{
			title:             "duplicate items in the validation list",
			validationList:    []time.Duration{1 * time.Second, 1 * time.Second},
			setValidationList: true,
			value:             "400s",
			expectedError:     "400s is not an acceptable value for --long. The expected value is 1s.",
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in time.Duration) error {
				return nil
			},
			setValidationCB: true,
			value:           "100s",
			expectedValue:   100 * time.Second,
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in time.Duration) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "100s",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in time.Duration) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []time.Duration{100 * time.Second, 200 * time.Second, 300 * time.Second},
			setValidationList: true,
			value:             "100s",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Duration("long", "usage")
			fVar := f.Var()
			if tc.setValidationCB {
				f = f.WithValidationCallback(tc.validationCB)
			}
			if tc.setValidationList {
				f = f.WithValidRange(tc.validationList...)
			}
			err := f.Set(tc.value)
			checkFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestDurationFlag_ResetToDefault(t *testing.T) {
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           time.Duration
		defaultValue            time.Duration
		expectedAfterResetValue time.Duration
		expectedError           string
		setDefault              bool
	}{
		{
			title: "no value",
		},
		{
			title:                   "reset without defining the default value",
			value:                   "10s",
			expectedValue:           10 * time.Second,
			expectedAfterResetValue: 10 * time.Second,
			setDefault:              false,
		},
		{
			title:                   "reset to zero default value",
			value:                   "10s",
			expectedValue:           10 * time.Second,
			defaultValue:            0,
			expectedAfterResetValue: 0,
			setDefault:              true,
		},
		{
			title:                   "reset to non-zero default value",
			value:                   "10s",
			expectedValue:           10 * time.Second,
			defaultValue:            50 * time.Second,
			expectedAfterResetValue: 50 * time.Second,
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Duration("long", "usage")
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
