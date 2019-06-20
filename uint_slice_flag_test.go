package flags_test

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"

	"go.xitonix.io/flags"
)

func TestUIntSlice(t *testing.T) {
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
			f := flags.UIntSlice(tc.long, tc.usage)
			checkFlagInitialState(t, f, "[]uint", tc.expectedUsage, tc.expectedLong, "")
			checkSliceFlagValues(t, []uint{}, f.Get(), f.Var())
		})
	}
}

func TestUIntSliceP(t *testing.T) {
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
			f := flags.UIntSliceP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "[]uint", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkSliceFlagValues(t, []uint{}, f.Get(), f.Var())
		})
	}
}

func TestUIntSliceFlag_WithKey(t *testing.T) {
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
			f := flags.UIntSlice("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestUIntSliceFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         []uint
		expectedDefaultValue []uint
	}{
		{
			title:                "empty default value",
			defaultValue:         []uint{},
			expectedDefaultValue: []uint{},
		},
		{
			title:                "non empty default value",
			defaultValue:         []uint{100},
			expectedDefaultValue: []uint{100},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.UIntSlice("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if !reflect.DeepEqual(actual.([]uint), tc.expectedDefaultValue) {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestUIntSliceFlag_Hide(t *testing.T) {
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
			f := flags.UIntSlice("long", "usage")
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

func TestUIntSliceFlag_IsDeprecated(t *testing.T) {
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
			f := flags.UIntSlice("long", "usage")
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

func TestUIntSliceFlag_WithDelimiter(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		delimiter     string
		expectedValue []uint
	}{
		{
			title:         "empty delimiter",
			value:         "100,200",
			expectedValue: []uint{100, 200},
		},
		{
			title:         "white space delimiter with white spaced input",
			value:         "100 200",
			delimiter:     " ",
			expectedValue: []uint{100, 200},
		},
		{
			title:         "none white space delimiter",
			value:         "100|200",
			delimiter:     "|",
			expectedValue: []uint{100, 200},
		},
		{
			title:         "no delimited input",
			value:         "100",
			delimiter:     "|",
			expectedValue: []uint{100},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.UIntSlice("long", "usage").WithDelimiter(tc.delimiter)
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestUIntSliceFlag_Set(t *testing.T) {
	empty := make([]uint, 0)
	testCases := []struct {
		title         string
		value         string
		expectedValue []uint
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
			value:         "  100  ",
			expectedValue: []uint{100},
		},
		{
			title:         "single value with no white space",
			value:         "100",
			expectedValue: []uint{100},
		},
		{
			title:         "comma separated value with no white space",
			value:         "0,100,200",
			expectedValue: []uint{0, 100, 200},
		},
		{
			title:         "comma separated value with white space",
			value:         " 0, 100 , 200 ",
			expectedValue: []uint{0, 100, 200},
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
			title:         "comma separated range",
			value:         fmt.Sprintf("0,100,200,%d", uint64(math.MaxUint64)),
			expectedValue: []uint{0, 100, 200, math.MaxUint64},
		},
		{
			title:         "invalid value",
			value:         " invalid ",
			expectedError: "is not a valid []uint value",
			expectedValue: empty,
		},
		{
			title:         "partially invalid value",
			value:         "100,invalid,200",
			expectedError: "is not a valid []uint value",
			expectedValue: empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.UIntSlice("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestUIntSliceFlag_Validation(t *testing.T) {
	empty := make([]uint, 0)
	testCases := []struct {
		title             string
		value             string
		expectedValue     []uint
		validationCB      func(in uint) error
		setValidationCB   bool
		validationList    []uint
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback",
			setValidationCB: true,
			value:           "1,2",
			expectedValue:   []uint{1, 2},
			expectedError:   "",
		},
		{
			title:             "nil validation list",
			setValidationList: true,
			value:             "1,2",
			expectedValue:     []uint{1, 2},
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             "1,2",
			expectedValue:     []uint{1, 2},
			expectedError:     "",
		},
		{
			title:             "empty validation list",
			validationList:    make([]uint, 0),
			setValidationList: true,
			value:             "1,2",
			expectedValue:     []uint{1, 2},
			expectedError:     "",
		},
		{
			title:             "none empty validation list with single item",
			validationList:    []uint{100, 200},
			setValidationList: true,
			value:             "10",
			expectedError:     "10 is not an acceptable value for --numbers. The expected values are 100 and 200.",
			expectedValue:     empty,
		},
		{
			title:             "none empty validation list with multiple items",
			validationList:    []uint{100, 200},
			setValidationList: true,
			value:             "100,300",
			expectedError:     "300 is not an acceptable value for --numbers. The expected values are 100 and 200.",
			expectedValue:     empty,
		},
		{
			title:             "validation list with three entries",
			validationList:    []uint{100, 200, 300},
			setValidationList: true,
			value:             "7",
			expectedError:     "7 is not an acceptable value for --numbers. The expected values are 100, 200 and 300.",
			expectedValue:     empty,
		},
		{
			title:             "none empty validation list",
			validationList:    []uint{0, 100, 200, math.MaxUint64},
			setValidationList: true,
			value:             fmt.Sprintf("0, 100,200, %v", uint64(math.MaxUint64)),
			expectedError:     "",
			expectedValue:     []uint{0, 100, 200, math.MaxUint64},
		},
		{
			title:             "empty value",
			validationList:    []uint{100, 200},
			setValidationList: true,
			value:             "",
			expectedError:     "",
			expectedValue:     empty,
		},
		{
			title:             "white space value",
			validationList:    []uint{100, 200},
			setValidationList: true,
			value:             "  ",
			expectedValue:     empty,
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in uint) error {
				return nil
			},
			setValidationCB: true,
			value:           "100",
			expectedValue:   []uint{100},
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in uint) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "100",
			expectedError:   "validation callback failed",
			expectedValue:   empty,
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in uint) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []uint{100, 200, 300},
			setValidationList: true,
			value:             "100",
			expectedError:     "validation callback failed",
			expectedValue:     empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.UIntSlice("numbers", "usage")
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

func TestUIntSliceFlag_ResetToDefault(t *testing.T) {
	empty := make([]uint, 0)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           []uint
		defaultValue            []uint
		expectedAfterResetValue []uint
		expectedError           string
		setDefault              bool
		expectedIsSetAfterReset bool
	}{
		{
			title:                   "reset without defining the default value",
			value:                   "100",
			expectedValue:           []uint{100},
			expectedAfterResetValue: []uint{100},
			setDefault:              false,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "reset to empty default value",
			value:                   "100",
			expectedValue:           []uint{100},
			defaultValue:            empty,
			expectedAfterResetValue: empty,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to nil default value",
			value:                   "100",
			expectedValue:           []uint{100},
			defaultValue:            nil,
			expectedAfterResetValue: []uint{100},
			setDefault:              true,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "reset to non-empty default value",
			value:                   "100",
			expectedValue:           []uint{100},
			defaultValue:            []uint{100, 200},
			expectedAfterResetValue: []uint{100, 200},
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.UIntSlice("long", "usage")
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
