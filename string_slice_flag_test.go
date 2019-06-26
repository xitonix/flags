package flags_test

import (
	"errors"
	"reflect"
	"testing"

	"go.xitonix.io/flags"
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
			f := flags.StringSlice(tc.long, tc.usage)
			checkFlagInitialState(t, f, "[]string", tc.expectedUsage, tc.expectedLong, "")
			checkSliceFlagValues(t, []string{}, f.Get(), f.Var())
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
			f := flags.StringSliceP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "[]string", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkSliceFlagValues(t, []string{}, f.Get(), f.Var())
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
			if !reflect.DeepEqual(actual.([]string), tc.expectedDefaultValue) {
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
			title:         "white space delimiter with white spaced input",
			value:         "abc xyz",
			delimiter:     " ",
			expectedValue: []string{"abc", "xyz"},
		},
		{
			title:         "white space delimiter with none white spaced input",
			value:         "abc,xyz",
			delimiter:     " ",
			expectedValue: []string{"abc,xyz"},
		},
		{
			title:         "none white space delimiter",
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
			checkSliceFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestStringSliceFlag_WithTrimming(t *testing.T) {
	testCases := []struct {
		title          string
		value          string
		enableTrimming bool
		expectedValue  []string
	}{
		{
			title:          "without trimming",
			enableTrimming: false,
			value:          "  abc  ,  xyz  ",
			expectedValue:  []string{"  abc  ", "  xyz  "},
		},
		{
			title:          "with trimming",
			enableTrimming: true,
			value:          "  abc  ,  xyz  ",
			expectedValue:  []string{"abc", "xyz"},
		},
		{
			title:          "only white space input without trimming",
			enableTrimming: false,
			value:          "   ",
			expectedValue:  []string{"   "},
		},
		{
			title:          "only white space input with trimming",
			enableTrimming: true,
			value:          "   ,   ",
			expectedValue:  []string{"", ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSlice("long", "usage")
			if tc.enableTrimming {
				f = f.WithTrimming()
			}
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
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
			title:         "white space value",
			value:         "   ",
			expectedValue: []string{"   "},
		},
		{
			title:         "value with white space",
			value:         "  abc  ",
			expectedValue: []string{"  abc  "},
		},
		{
			title:         "value with no white space",
			value:         "abc",
			expectedValue: []string{"abc"},
		},
		{
			title:         "comma separated value with no white space",
			value:         "abc,efg",
			expectedValue: []string{"abc", "efg"},
		},
		{
			title:         "comma separated value with white space",
			value:         " abc , efg ",
			expectedValue: []string{" abc ", " efg "},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSlice("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestStringSliceFlag_Validation(t *testing.T) {
	empty := make([]string, 0)
	testCases := []struct {
		title             string
		value             string
		expectedValue     []string
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
			expectedValue:   []string{"value"},
			expectedError:   "",
		},
		{
			title:             "nil validation list",
			setValidationList: true,
			value:             "value",
			expectedValue:     []string{"value"},
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             "value",
			expectedValue:     []string{"value"},
			expectedError:     "",
		},
		{
			title:             "empty validation list",
			validationList:    make([]string, 0),
			setValidationList: true,
			value:             "value",
			expectedValue:     []string{"value"},
			expectedError:     "",
		},
		{
			title:             "invalid case sensitive validation list with single item",
			validationList:    []string{"Green", "Red"},
			ignoreCase:        false,
			setValidationList: true,
			value:             "green",
			expectedError:     "green is not an acceptable value for --colours. The expected values are Green,Red.",
			expectedValue:     empty,
		},
		{
			title:             "duplicate case sensitive validation list",
			validationList:    []string{"Green", "Green"},
			ignoreCase:        false,
			setValidationList: true,
			value:             "green",
			expectedError:     "green is not an acceptable value for --colours. The expected value is Green.",
			expectedValue:     empty,
		},
		{
			title:             "duplicate case insensitive validation list",
			validationList:    []string{"Green", "Green"},
			ignoreCase:        false,
			setValidationList: true,
			value:             "Pink",
			expectedError:     "Pink is not an acceptable value for --colours. The expected value is Green.",
			expectedValue:     empty,
		},
		{
			title:             "invalid case sensitive validation list with multiple items",
			validationList:    []string{"Green", "Red", "Blue"},
			ignoreCase:        false,
			setValidationList: true,
			value:             "Green,red,Blue",
			expectedError:     "red is not an acceptable value for --colours. The expected values are Green,Red,Blue.",
			expectedValue:     empty,
		},
		{
			title:             "valid case sensitive validation list with multiple items",
			validationList:    []string{"Green", "Red"},
			ignoreCase:        false,
			setValidationList: true,
			value:             "Green,Red",
			expectedError:     "",
			expectedValue:     []string{"Green", "Red"},
		},
		{
			title:             "empty value",
			validationList:    []string{"Green", "Red"},
			ignoreCase:        false,
			setValidationList: true,
			value:             "",
			expectedError:     "",
			expectedValue:     empty,
		},
		{
			title:             "white space value",
			validationList:    []string{"Green", "Red"},
			ignoreCase:        false,
			setValidationList: true,
			value:             "  ",
			expectedError:     "'  ' is not an acceptable value for --colours. The expected values are Green,Red.",
			expectedValue:     empty,
		},
		{
			title:             "acceptable white space value",
			validationList:    []string{"Green", "Red", "  "},
			ignoreCase:        false,
			setValidationList: true,
			value:             "  ",
			expectedError:     "",
			expectedValue:     []string{"  "},
		},
		{
			title:             "invalid case insensitive validation list with multiple items",
			validationList:    []string{"Green", "Red"},
			ignoreCase:        true,
			setValidationList: true,
			value:             "Green,Red,Pink",
			expectedError:     "Pink is not an acceptable value for --colours. The expected values are Green,Red.",
			expectedValue:     empty,
		},
		{
			title:             "valid case insensitive validation list with multiple items",
			validationList:    []string{"Green", "Red"},
			ignoreCase:        true,
			setValidationList: true,
			value:             "green,red",
			expectedValue:     []string{"green", "red"},
		},
		{
			title:             "three items in the validation list",
			validationList:    []string{"Green", "Pink", "Yellow"},
			setValidationList: true,
			value:             "blue",
			expectedError:     "blue is not an acceptable value for --colours. The expected values are Green,Pink,Yellow.",
			expectedValue:     empty,
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in string) error {
				return nil
			},
			setValidationCB: true,
			value:           "blue",
			expectedValue:   []string{"blue"},
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in string) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "blue",
			expectedError:   "validation callback failed",
			expectedValue:   empty,
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
			expectedValue:     empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSlice("colours", "usage")
			fVar := f.Var()
			if tc.setValidationCB {
				f = f.WithValidationCallback(tc.validationCB)
			}
			if tc.setValidationList {
				f = f.WithValidRange(tc.ignoreCase, tc.validationList...)
			}
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
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
			expectedAfterResetValue: nil,
			setDefault:              true,
			expectedIsSetAfterReset: false,
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
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)

			f.ResetToDefault()

			if f.IsSet() != tc.expectedIsSetAfterReset {
				t.Errorf("IsSet() Expected: %v, Actual: %v", tc.expectedIsSetAfterReset, f.IsSet())
			}

			checkSliceFlagValues(t, tc.expectedAfterResetValue, f.Get(), fVar)
		})
	}
}
