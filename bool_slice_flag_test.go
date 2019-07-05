package flags_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/xitonix/flags"
)

func TestBoolSlice(t *testing.T) {
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
			f := flags.BoolSlice(tc.long, tc.usage)
			checkFlagInitialState(t, f, "[]bool", tc.expectedUsage, tc.expectedLong, "")
			checkSliceFlagValues(t, []bool{}, f.Get(), f.Var())
		})
	}
}

func TestBoolSliceFlag_WithShort(t *testing.T) {
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
			f := flags.BoolSlice(tc.long, tc.usage).WithShort(tc.short)
			checkFlagInitialState(t, f, "[]bool", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkSliceFlagValues(t, []bool{}, f.Get(), f.Var())
		})
	}
}

func TestBoolSliceFlag_WithKey(t *testing.T) {
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
			f := flags.BoolSlice("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestBoolSliceFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         []bool
		expectedDefaultValue []bool
	}{
		{
			title:                "empty default value",
			defaultValue:         []bool{},
			expectedDefaultValue: []bool{},
		},
		{
			title:                "non empty default value",
			defaultValue:         []bool{true},
			expectedDefaultValue: []bool{true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.BoolSlice("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if !reflect.DeepEqual(actual.([]bool), tc.expectedDefaultValue) {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestBoolSliceFlag_Hide(t *testing.T) {
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
			f := flags.BoolSlice("long", "usage")
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

func TestBoolSliceFlag_IsDeprecated(t *testing.T) {
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
			f := flags.BoolSlice("long", "usage")
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

func TestBoolSliceFlag_IsRequired(t *testing.T) {
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
			f := flags.BoolSlice("long", "usage")
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

func TestBoolSliceFlag_WithDelimiter(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		delimiter     string
		expectedValue []bool
	}{
		{
			title:         "empty delimiter",
			value:         "true,false",
			expectedValue: []bool{true, false},
		},
		{
			title:         "empty delimiter with mixed input",
			value:         "true,false,0,1",
			expectedValue: []bool{true, false, false, true},
		},
		{
			title:         "white space delimiter with white spaced input",
			value:         "true false",
			delimiter:     " ",
			expectedValue: []bool{true, false},
		},
		{
			title:         "none white space delimiter",
			value:         "true|false",
			delimiter:     "|",
			expectedValue: []bool{true, false},
		},
		{
			title:         "none white space delimiter with mixed input",
			value:         "1|true|false|1",
			delimiter:     "|",
			expectedValue: []bool{true, true, false, true},
		},
		{
			title:         "no delimited input",
			value:         "true",
			delimiter:     "|",
			expectedValue: []bool{true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.BoolSlice("long", "usage").WithDelimiter(tc.delimiter)
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestBoolSliceFlag_Set(t *testing.T) {
	empty := make([]bool, 0)
	testCases := []struct {
		title         string
		value         string
		expectedValue []bool
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
			value:         "  true  ",
			expectedValue: []bool{true},
		},
		{
			title:         "numeric single value with white space",
			value:         "  1  ",
			expectedValue: []bool{true},
		},
		{
			title:         "single value with no white space",
			value:         "true",
			expectedValue: []bool{true},
		},
		{
			title:         "numeric single value with no white space",
			value:         "1",
			expectedValue: []bool{true},
		},
		{
			title:         "comma separated value with no white space",
			value:         "true,false,true",
			expectedValue: []bool{true, false, true},
		},
		{
			title:         "mixed comma separated value with no white space",
			value:         "1,true,false,1,true,0",
			expectedValue: []bool{true, true, false, true, true, false},
		},
		{
			title:         "comma separated value with white space",
			value:         " true, false , true ",
			expectedValue: []bool{true, false, true},
		},
		{
			title:         "mixed comma separated value with white space",
			value:         " 1 , true, false , true , 1 ",
			expectedValue: []bool{true, true, false, true, true},
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
			expectedError: "'invalid' is not a valid []bool value for --long",
			expectedValue: empty,
		},
		{
			title:         "partially invalid value",
			value:         "1,invalid,0",
			expectedError: "'invalid' is not a valid []bool value for --long",
			expectedValue: empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.BoolSlice("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestBoolSliceFlag_Validation(t *testing.T) {
	empty := make([]bool, 0)
	testCases := []struct {
		title           string
		value           string
		expectedValue   []bool
		validationCB    func(in bool) error
		setValidationCB bool
		expectedError   string
	}{
		{
			title:           "nil validation callback",
			setValidationCB: true,
			value:           "0,1",
			expectedValue:   []bool{false, true},
			expectedError:   "",
		},
		{
			title:           "mixed input with nil validation callback",
			setValidationCB: true,
			value:           "true,0,1,false",
			expectedValue:   []bool{true, false, true, false},
			expectedError:   "",
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in bool) error {
				return nil
			},
			setValidationCB: true,
			value:           "1",
			expectedValue:   []bool{true},
		},
		{
			title: "mixed input with validation callback and no validation error",
			validationCB: func(in bool) error {
				return nil
			},
			setValidationCB: true,
			value:           "true,0,1,false",
			expectedValue:   []bool{true, false, true, false},
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in bool) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "true",
			expectedError:   "validation callback failed",
			expectedValue:   empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.BoolSlice("numbers", "usage")
			fVar := f.Var()
			if tc.setValidationCB {
				f = f.WithValidationCallback(tc.validationCB)
			}
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestBoolSliceFlag_ResetToDefault(t *testing.T) {
	empty := make([]bool, 0)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           []bool
		defaultValue            []bool
		expectedAfterResetValue []bool
		expectedError           string
		setDefault              bool
		expectedIsSetAfterReset bool
	}{
		{
			title:                   "reset without defining the default value",
			value:                   "true",
			expectedValue:           []bool{true},
			expectedAfterResetValue: []bool{true},
			setDefault:              false,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "reset to empty default value",
			value:                   "true",
			expectedValue:           []bool{true},
			defaultValue:            empty,
			expectedAfterResetValue: empty,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to nil default value",
			value:                   "true",
			expectedValue:           []bool{true},
			defaultValue:            nil,
			expectedAfterResetValue: nil,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to non-empty default value",
			value:                   "true,0",
			expectedValue:           []bool{true, false},
			defaultValue:            []bool{false, false},
			expectedAfterResetValue: []bool{false, false},
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.BoolSlice("long", "usage")
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
