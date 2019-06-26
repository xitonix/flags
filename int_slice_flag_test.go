package flags_test

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"

	"go.xitonix.io/flags"
)

func TestIntSlice(t *testing.T) {
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
			f := flags.IntSlice(tc.long, tc.usage)
			checkFlagInitialState(t, f, "[]int", tc.expectedUsage, tc.expectedLong, "")
			checkSliceFlagValues(t, []int{}, f.Get(), f.Var())
		})
	}
}

func TestIntSliceP(t *testing.T) {
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
			f := flags.IntSliceP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "[]int", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkSliceFlagValues(t, []int{}, f.Get(), f.Var())
		})
	}
}

func TestIntSliceFlag_WithKey(t *testing.T) {
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
			f := flags.IntSlice("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestIntSliceFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         []int
		expectedDefaultValue []int
	}{
		{
			title:                "empty default value",
			defaultValue:         []int{},
			expectedDefaultValue: []int{},
		},
		{
			title:                "non empty default value",
			defaultValue:         []int{100},
			expectedDefaultValue: []int{100},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IntSlice("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if !reflect.DeepEqual(actual.([]int), tc.expectedDefaultValue) {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestIntSliceFlag_Hide(t *testing.T) {
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
			f := flags.IntSlice("long", "usage")
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

func TestIntSliceFlag_IsDeprecated(t *testing.T) {
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
			f := flags.IntSlice("long", "usage")
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

func TestIntSliceFlag_WithDelimiter(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		delimiter     string
		expectedValue []int
	}{
		{
			title:         "empty delimiter",
			value:         "100,200",
			expectedValue: []int{100, 200},
		},
		{
			title:         "white space delimiter with white spaced input",
			value:         "100 200",
			delimiter:     " ",
			expectedValue: []int{100, 200},
		},
		{
			title:         "none white space delimiter",
			value:         "100|200",
			delimiter:     "|",
			expectedValue: []int{100, 200},
		},
		{
			title:         "no delimited input",
			value:         "100",
			delimiter:     "|",
			expectedValue: []int{100},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IntSlice("long", "usage").WithDelimiter(tc.delimiter)
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestIntSliceFlag_Set(t *testing.T) {
	empty := make([]int, 0)
	testCases := []struct {
		title         string
		value         string
		expectedValue []int
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
			expectedValue: []int{100},
		},
		{
			title:         "single value with no white space",
			value:         "100",
			expectedValue: []int{100},
		},
		{
			title:         "comma separated value with no white space",
			value:         "100,200",
			expectedValue: []int{100, 200},
		},
		{
			title:         "comma separated value with white space",
			value:         " 100 , 200 ",
			expectedValue: []int{100, 200},
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
			value:         fmt.Sprintf("%d, 0, %d", int64(math.MinInt64), int64(math.MaxInt64)),
			expectedValue: []int{math.MinInt64, 0, math.MaxInt64},
		},
		{
			title:         "invalid value",
			value:         " invalid ",
			expectedError: "'invalid' is not a valid []int value for --long",
			expectedValue: empty,
		},
		{
			title:         "partially invalid value",
			value:         "100,invalid,200",
			expectedError: "'invalid' is not a valid []int value for --long",
			expectedValue: empty,
		},
		{
			title:         "comma separated negative values",
			value:         "-100,-200,0,100,200",
			expectedValue: []int{-100, -200, 0, 100, 200},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IntSlice("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestIntSliceFlag_Validation(t *testing.T) {
	empty := make([]int, 0)
	testCases := []struct {
		title             string
		value             string
		expectedValue     []int
		validationCB      func(in int) error
		setValidationCB   bool
		validationList    []int
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback",
			setValidationCB: true,
			value:           "1,2",
			expectedValue:   []int{1, 2},
			expectedError:   "",
		},
		{
			title:             "nil validation list",
			setValidationList: true,
			value:             "1,2",
			expectedValue:     []int{1, 2},
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             "1,2",
			expectedValue:     []int{1, 2},
			expectedError:     "",
		},
		{
			title:             "empty validation list",
			validationList:    make([]int, 0),
			setValidationList: true,
			value:             "1,2",
			expectedValue:     []int{1, 2},
			expectedError:     "",
		},
		{
			title:             "none empty validation list with single item",
			validationList:    []int{100},
			setValidationList: true,
			value:             "10",
			expectedError:     "10 is not an acceptable value for --numbers. The expected value is 100.",
			expectedValue:     empty,
		},
		{
			title:             "none empty validation list with two items",
			validationList:    []int{100, 200},
			setValidationList: true,
			value:             "100,300",
			expectedError:     "300 is not an acceptable value for --numbers. The expected values are 100,200.",
			expectedValue:     empty,
		},
		{
			title:             "validation list with three items",
			validationList:    []int{100, 200, 300},
			setValidationList: true,
			value:             "7",
			expectedError:     "7 is not an acceptable value for --numbers. The expected values are 100,200,300.",
			expectedValue:     empty,
		},
		{
			title:             "duplicate items in the validation list",
			validationList:    []int{100, 100},
			setValidationList: true,
			value:             "7",
			expectedError:     "7 is not an acceptable value for --numbers. The expected value is 100.",
			expectedValue:     empty,
		},
		{
			title:             "none empty validation list",
			validationList:    []int{math.MinInt64, 0, math.MaxInt64},
			setValidationList: true,
			value:             fmt.Sprintf("%d, 0, %d", int64(math.MinInt64), int64(math.MaxInt64)),
			expectedValue:     []int{math.MinInt64, 0, math.MaxInt64},
			expectedError:     "",
		},
		{
			title:             "empty value",
			validationList:    []int{100, 200},
			setValidationList: true,
			value:             "",
			expectedError:     "",
			expectedValue:     empty,
		},
		{
			title:             "white space value",
			validationList:    []int{100, 200},
			setValidationList: true,
			value:             "  ",
			expectedValue:     empty,
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in int) error {
				return nil
			},
			setValidationCB: true,
			value:           "100",
			expectedValue:   []int{100},
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in int) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "100",
			expectedError:   "validation callback failed",
			expectedValue:   empty,
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in int) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []int{100, 200, 300},
			setValidationList: true,
			value:             "100",
			expectedError:     "validation callback failed",
			expectedValue:     empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IntSlice("numbers", "usage")
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

func TestIntSliceFlag_ResetToDefault(t *testing.T) {
	empty := make([]int, 0)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           []int
		defaultValue            []int
		expectedAfterResetValue []int
		expectedError           string
		setDefault              bool
		expectedIsSetAfterReset bool
	}{
		{
			title:                   "reset without defining the default value",
			value:                   "100",
			expectedValue:           []int{100},
			expectedAfterResetValue: []int{100},
			setDefault:              false,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "reset to empty default value",
			value:                   "100",
			expectedValue:           []int{100},
			defaultValue:            empty,
			expectedAfterResetValue: empty,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to nil default value",
			value:                   "100",
			expectedValue:           []int{100},
			defaultValue:            nil,
			expectedAfterResetValue: nil,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to non-empty default value",
			value:                   "100",
			expectedValue:           []int{100},
			defaultValue:            []int{100, 200},
			expectedAfterResetValue: []int{100, 200},
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IntSlice("long", "usage")
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
