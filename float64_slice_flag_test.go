package flags_test

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"

	"go.xitonix.io/flags"
)

func TestFloat64Slice(t *testing.T) {
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
			f := flags.Float64Slice(tc.long, tc.usage)
			checkFlagInitialState(t, f, "[]float64", tc.expectedUsage, tc.expectedLong, "")
			checkSliceFlagValues(t, []float64{}, f.Get(), f.Var())
		})
	}
}

func TestFloat64SliceP(t *testing.T) {
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
			f := flags.Float64SliceP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "[]float64", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkSliceFlagValues(t, []float64{}, f.Get(), f.Var())
		})
	}
}

func TestFloat64SliceFlag_WithKey(t *testing.T) {
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
			f := flags.Float64Slice("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestFloat64SliceFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         []float64
		expectedDefaultValue []float64
	}{
		{
			title:                "empty default value",
			defaultValue:         []float64{},
			expectedDefaultValue: []float64{},
		},
		{
			title:                "non empty default value",
			defaultValue:         []float64{100.87},
			expectedDefaultValue: []float64{100.87},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Float64Slice("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if !reflect.DeepEqual(actual.([]float64), tc.expectedDefaultValue) {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestFloat64SliceFlag_Hide(t *testing.T) {
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
			f := flags.Float64Slice("long", "usage")
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

func TestFloat64SliceFlag_IsDeprecated(t *testing.T) {
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
			f := flags.Float64Slice("long", "usage")
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

func TestFloat64SliceFlag_IsRequired(t *testing.T) {
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
			f := flags.Float64Slice("long", "usage")
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

func TestFloat64SliceFlag_WithDelimiter(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		delimiter     string
		expectedValue []float64
	}{
		{
			title:         "empty delimiter",
			value:         "100.87,200.90",
			expectedValue: []float64{100.87, 200.90},
		},
		{
			title:         "white space delimiter with white spaced input",
			value:         "100.87 200.90",
			delimiter:     " ",
			expectedValue: []float64{100.87, 200.90},
		},
		{
			title:         "none white space delimiter",
			value:         "100.87|200.90",
			delimiter:     "|",
			expectedValue: []float64{100.87, 200.90},
		},
		{
			title:         "no delimited input",
			value:         "100.87",
			delimiter:     "|",
			expectedValue: []float64{100.87},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Float64Slice("long", "usage").WithDelimiter(tc.delimiter)
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestFloat64SliceFlag_Set(t *testing.T) {
	empty := make([]float64, 0)
	testCases := []struct {
		title         string
		value         string
		expectedValue []float64
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
			value:         "  100.87  ",
			expectedValue: []float64{100.87},
		},
		{
			title:         "single value with no white space",
			value:         "100.87",
			expectedValue: []float64{100.87},
		},
		{
			title:         "comma separated value with no white space",
			value:         "0,100.87,200.90",
			expectedValue: []float64{0, 100.87, 200.90},
		},
		{
			title:         "comma separated value with white space",
			value:         " 0, 100.87 , 200.90 ",
			expectedValue: []float64{0, 100.87, 200.90},
		},
		{
			title:         "comma separated empty string",
			value:         ",,",
			expectedValue: empty,
		},
		{
			title:         "comma separated range",
			value:         fmt.Sprintf("0,100.87,200.90,%f", math.MaxFloat64),
			expectedValue: []float64{0, 100.87, 200.90, math.MaxFloat64},
		},
		{
			title:         "comma separated white space string",
			value:         " , , ",
			expectedValue: empty,
		},
		{
			title:         "invalid value",
			value:         " invalid ",
			expectedError: "'invalid' is not a valid []float64 value for --long",
			expectedValue: empty,
		},
		{
			title:         "partially invalid value",
			value:         "100.87,invalid,200.90",
			expectedError: "'invalid' is not a valid []float64 value for --long",
			expectedValue: empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Float64Slice("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestFloat64SliceFlag_Validation(t *testing.T) {
	empty := make([]float64, 0)
	testCases := []struct {
		title             string
		value             string
		expectedValue     []float64
		validationCB      func(in float64) error
		setValidationCB   bool
		validationList    []float64
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback",
			setValidationCB: true,
			value:           "2.5, 3.14",
			expectedValue:   []float64{2.5, 3.14},
			expectedError:   "",
		},
		{
			title:             "nil validation list",
			setValidationList: true,
			value:             "2.5, 3.14",
			expectedValue:     []float64{2.5, 3.14},
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             "2.5, 3.14",
			expectedValue:     []float64{2.5, 3.14},
			expectedError:     "",
		},
		{
			title:             "empty validation list",
			validationList:    make([]float64, 0),
			setValidationList: true,
			value:             "2.5, 3.14",
			expectedValue:     []float64{2.5, 3.14},
			expectedError:     "",
		},
		{
			title:             "none empty validation list with single item",
			validationList:    []float64{100.87},
			setValidationList: true,
			value:             "10",
			expectedError:     "10 is not an acceptable value for --numbers. The expected value is 100.87.",
			expectedValue:     empty,
		},
		{
			title:             "none empty validation list with two items",
			validationList:    []float64{100.87, 200.90},
			setValidationList: true,
			value:             "100.87,314.32",
			expectedError:     "314.32 is not an acceptable value for --numbers. The expected values are 100.87,200.9.",
			expectedValue:     empty,
		},
		{
			title:             "duplicate items in the validation list",
			validationList:    []float64{100.87, 100.87},
			setValidationList: true,
			value:             "100.87,314.32",
			expectedError:     "314.32 is not an acceptable value for --numbers. The expected value is 100.87.",
			expectedValue:     empty,
		},
		{
			title:             "validation list with three items",
			validationList:    []float64{100.87, 200.90, 314.32},
			setValidationList: true,
			value:             "7.5",
			expectedError:     "7.5 is not an acceptable value for --numbers. The expected values are 100.87,200.9,314.32.",
			expectedValue:     empty,
		},
		{
			title:             "none empty validation list",
			validationList:    []float64{math.MinInt64, 0, 100.87, 200.90, math.MaxFloat64},
			setValidationList: true,
			value:             fmt.Sprintf("%f, 0, 100.87 ,200.90 , %f", float64(math.MinInt64), math.MaxFloat64),
			expectedValue:     []float64{math.MinInt64, 0, 100.87, 200.90, math.MaxFloat64},
			expectedError:     "",
		},
		{
			title:             "empty value",
			validationList:    []float64{100.87, 200.90},
			setValidationList: true,
			value:             "",
			expectedError:     "",
			expectedValue:     empty,
		},
		{
			title:             "white space value",
			validationList:    []float64{100.87, 200.90},
			setValidationList: true,
			value:             "  ",
			expectedValue:     empty,
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in float64) error {
				return nil
			},
			setValidationCB: true,
			value:           "100.87",
			expectedValue:   []float64{100.87},
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in float64) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "100.87",
			expectedError:   "validation callback failed",
			expectedValue:   empty,
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in float64) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []float64{100.87, 200.90, 314.32},
			setValidationList: true,
			value:             "100.87",
			expectedError:     "validation callback failed",
			expectedValue:     empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Float64Slice("numbers", "usage")
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

func TestFloat64SliceFlag_ResetToDefault(t *testing.T) {
	empty := make([]float64, 0)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           []float64
		defaultValue            []float64
		expectedAfterResetValue []float64
		expectedError           string
		setDefault              bool
		expectedIsSetAfterReset bool
	}{
		{
			title:                   "reset without defining the default value",
			value:                   "100.87",
			expectedValue:           []float64{100.87},
			expectedAfterResetValue: []float64{100.87},
			setDefault:              false,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "reset to empty default value",
			value:                   "100.87",
			expectedValue:           []float64{100.87},
			defaultValue:            empty,
			expectedAfterResetValue: empty,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to nil default value",
			value:                   "100.87",
			expectedValue:           []float64{100.87},
			defaultValue:            nil,
			expectedAfterResetValue: nil,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to non-empty default value",
			value:                   "100.87",
			expectedValue:           []float64{100.87},
			defaultValue:            []float64{100.87, 200.90},
			expectedAfterResetValue: []float64{100.87, 200.90},
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Float64Slice("long", "usage")
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
