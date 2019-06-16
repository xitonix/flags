package flags_test

import (
	"errors"
	"testing"

	"go.xitonix.io/flags"
)

func TestString(t *testing.T) {
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
			f := flags.String(tc.long, tc.usage)
			checkFlagInitialState(t, f, "string", tc.expectedUsage, tc.expectedLong, "")
			checkFlagValues(t, "", f.Get(), f.Var())
		})
	}
}

func TestStringP(t *testing.T) {
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
			f := flags.StringP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "string", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkFlagValues(t, "", f.Get(), f.Var())
		})
	}
}

func TestStringFlag_WithKey(t *testing.T) {
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
			f := flags.String("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestStringFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         string
		expectedDefaultValue string
	}{
		{
			title:                "empty default value",
			expectedDefaultValue: "''",
		},
		{
			title:                "white space default value",
			defaultValue:         "    ",
			expectedDefaultValue: "    ",
		},
		{
			title:                "default value with white space",
			defaultValue:         "  default value  ",
			expectedDefaultValue: "  default value  ",
		},
		{
			title:                "non empty default value",
			defaultValue:         "default value",
			expectedDefaultValue: "default value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.String("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if actual != tc.expectedDefaultValue {
				t.Errorf("Expected Default Value: %s, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestStringFlag_Hide(t *testing.T) {
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
			f := flags.String("long", "usage")
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

func TestStringFlag_IsDeprecated(t *testing.T) {
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
			f := flags.String("long", "usage")
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

func TestStringFlag_Set(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		expectedError string
	}{
		{
			title: "no value",
		},
		{
			title: "white space value",
			value: "   ",
		},
		{
			title: "value with white space",
			value: "  value  ",
		},
		{
			title: "value with no white space",
			value: "value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.String("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkFlag(t, f, err, tc.expectedError, tc.value, f.Get(), fVar)
		})
	}
}

func TestStringFlag_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     string
		validationCB      func(in string) error
		setValidationCB   bool
		validationList    []string
		setValidationList bool
		ignoreCase        bool
		expectedError     string
	}{
		{
			title:           "nil validation callback",
			setValidationCB: true,
			value:           "value",
			expectedValue:   "value",
			expectedError:   "",
		},
		{
			title:             "nil validation list",
			setValidationList: true,
			value:             "value",
			expectedValue:     "value",
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             "value",
			expectedValue:     "value",
			expectedError:     "",
		},
		{
			title:             "empty validation list",
			validationList:    make([]string, 0),
			setValidationList: true,
			value:             "value",
			expectedValue:     "value",
			expectedError:     "",
		},
		{
			title:             "case sensitive validation list",
			validationList:    []string{"Green", "Red"},
			ignoreCase:        false,
			setValidationList: true,
			value:             "green",
			expectedError:     "'green' is not an acceptable value for --colours. Expected value(s): Green and Red",
		},
		{
			title:             "case insensitive validation list",
			validationList:    []string{"Green", "Red"},
			ignoreCase:        true,
			setValidationList: true,
			value:             "green",
			expectedValue:     "green",
		},
		{
			title:             "single item in the validation list",
			validationList:    []string{"Green"},
			setValidationList: true,
			value:             "blue",
			expectedError:     "'blue' is not an acceptable value for --colours. Expected value(s): Green",
		},
		{
			title:             "two items in the validation list",
			validationList:    []string{"Green", "Pink"},
			setValidationList: true,
			value:             "blue",
			expectedError:     "'blue' is not an acceptable value for --colours. Expected value(s): Green and Pink",
		},
		{
			title:             "three items in the validation list",
			validationList:    []string{"Green", "Pink", "Yellow"},
			setValidationList: true,
			value:             "blue",
			expectedError:     "'blue' is not an acceptable value for --colours. Expected value(s): Green, Pink and Yellow",
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in string) error {
				return nil
			},
			setValidationCB: true,
			value:           "blue",
			expectedValue:   "blue",
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in string) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "blue",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in string) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []string{"Green", "Pink", "Yellow"},
			setValidationList: true,
			value:             "blue",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.String("colours", "usage")
			fVar := f.Var()
			if tc.setValidationCB {
				f = f.WithValidationCallback(tc.validationCB)
			}
			if tc.setValidationList {
				f = f.WithValidRange(tc.ignoreCase, tc.validationList...)
			}
			err := f.Set(tc.value)
			checkFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestStringFlag_ResetToDefault(t *testing.T) {
	testCases := []struct {
		title                   string
		value                   string
		defaultValue            string
		expectedAfterResetValue string
		expectedError           string
		setDefault              bool
	}{
		{
			title: "no value",
		},
		{
			title:                   "reset without defining the default value",
			value:                   "value",
			expectedAfterResetValue: "value",
			setDefault:              false,
		},
		{
			title:                   "reset to empty default value",
			value:                   "value",
			defaultValue:            "",
			expectedAfterResetValue: "",
			setDefault:              true,
		},
		{
			title:                   "reset to white space default value",
			value:                   "value",
			defaultValue:            "  ",
			expectedAfterResetValue: "  ",
			setDefault:              true,
		},
		{
			title:                   "reset to non-empty default value",
			value:                   "value",
			defaultValue:            "Default   Value  ",
			expectedAfterResetValue: "Default   Value  ",
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.String("long", "usage")
			if tc.setDefault {
				f = f.WithDefault(tc.defaultValue)
			}
			fVar := f.Var()
			err := f.Set(tc.value)
			checkFlag(t, f, err, tc.expectedError, tc.value, f.Get(), fVar)

			f.ResetToDefault()

			if tc.setDefault && f.IsSet() {
				t.Error("IsSet() Expected: false, Actual: true")
			}

			checkFlagValues(t, tc.expectedAfterResetValue, f.Get(), fVar)
		})
	}
}
