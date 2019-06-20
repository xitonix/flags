package flags_test

import (
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
			title:         "invalid value",
			value:         " invalid ",
			expectedError: "is not a valid []int value",
			expectedValue: empty,
		},
		{
			title:         "partially invalid value",
			value:         "100,invalid,200",
			expectedError: "is not a valid []int value",
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
			expectedAfterResetValue: []int{100},
			setDefault:              true,
			expectedIsSetAfterReset: true,
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
