package core_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/xitonix/flags"
)

func TestUInt16(t *testing.T) {
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
			f := flags.UInt16(tc.long, tc.usage)
			checkFlagInitialState(t, f, "uint16", tc.expectedUsage, tc.expectedLong, "")
			checkFlagValues(t, uint16(0), f.Get(), f.Var())
		})
	}
}

func TestUInt16Flag_WithShort(t *testing.T) {
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
			f := flags.UInt16(tc.long, tc.usage).WithShort(tc.short)
			checkFlagInitialState(t, f, "uint16", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkFlagValues(t, uint16(0), f.Get(), f.Var())
		})
	}
}

func TestUInt16Flag_WithKey(t *testing.T) {
	testCases := []struct {
		title       string
		key         string
		expectedKey string
	}{
		{
			title: "empty key",
		},
		{
			title:       "dash key",
			key:         "-",
			expectedKey: "",
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
			f := flags.UInt16("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestUInt16Flag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         uint16
		expectedDefaultValue uint16
	}{
		{
			title:                "zero default value",
			defaultValue:         0,
			expectedDefaultValue: 0,
		},
		{
			title:                "non zero default value",
			defaultValue:         100,
			expectedDefaultValue: 100,
		},
		{
			title:                "max uint16",
			defaultValue:         math.MaxUint16,
			expectedDefaultValue: math.MaxUint16,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.UInt16("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if actual != tc.expectedDefaultValue {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestUInt16Flag_Hide(t *testing.T) {
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
			f := flags.UInt16("long", "usage")
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

func TestUInt16Flag_IsDeprecated(t *testing.T) {
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
			f := flags.UInt16("long", "usage")
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

func TestUInt16Flag_IsRequired(t *testing.T) {
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
			f := flags.UInt16("long", "usage")
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

func TestUInt16Flag_Set(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		expectedValue uint16
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
			value:         "  100  ",
			expectedValue: 100,
		},
		{
			title:         "value with no white space",
			value:         "100",
			expectedValue: 100,
		},
		{
			title:         "invalid value",
			value:         "abc",
			expectedError: "is not a valid uint16 value",
			expectedValue: 0,
		},
		{
			title:         "max uint16",
			value:         fmt.Sprintf("%d", uint16(math.MaxUint16)),
			expectedValue: math.MaxUint16,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.UInt16("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestUInt16Flag_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     uint16
		validationCB      func(in uint16) error
		setValidationCB   bool
		validationList    []uint16
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback",
			setValidationCB: true,
			value:           "100",
			expectedValue:   100,
			expectedError:   "",
		},
		{
			title:             "nil validation list",
			setValidationList: true,
			value:             "100",
			expectedValue:     100,
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             "100",
			expectedValue:     100,
			expectedError:     "",
		},
		{
			title:             "empty validation list",
			validationList:    make([]uint16, 0),
			setValidationList: true,
			value:             "100",
			expectedValue:     100,
			expectedError:     "",
		},
		{
			title:             "single item in the validation list",
			validationList:    []uint16{100},
			setValidationList: true,
			value:             "101",
			expectedError:     "101 is not an acceptable value for --long. The expected value is 100.",
		},
		{
			title:             "two items in the validation list",
			validationList:    []uint16{100, 101},
			setValidationList: true,
			value:             "102",
			expectedError:     "102 is not an acceptable value for --long. The expected values are 100,101.",
		},
		{
			title:             "three items in the validation list",
			validationList:    []uint16{100, 101, 102},
			setValidationList: true,
			value:             "104",
			expectedError:     "104 is not an acceptable value for --long. The expected values are 100,101,102.",
		},
		{
			title:             "min and max",
			validationList:    []uint16{0, math.MaxUint16},
			setValidationList: true,
			value:             "104",
			expectedError:     fmt.Sprintf("104 is not an acceptable value for --long. The expected values are %v,%v.", 0, math.MaxUint16),
		},
		{
			title:             "duplicate items in the validation list",
			validationList:    []uint16{100, 100},
			setValidationList: true,
			value:             "101",
			expectedError:     "101 is not an acceptable value for --long. The expected value is 100.",
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(in uint16) error {
				return nil
			},
			setValidationCB: true,
			value:           "100",
			expectedValue:   100,
		},
		{
			title: "validation callback with validation error",
			validationCB: func(in uint16) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "100",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in uint16) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []uint16{100, 101, 102},
			setValidationList: true,
			value:             "100",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.UInt16("long", "usage")
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

func TestUInt16Flag_ResetToDefault(t *testing.T) {
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           uint16
		defaultValue            uint16
		expectedAfterResetValue uint16
		expectedError           string
		setDefault              bool
	}{
		{
			title: "no value",
		},
		{
			title:                   "reset without defining the default value",
			value:                   "100",
			expectedValue:           100,
			expectedAfterResetValue: 100,
			setDefault:              false,
		},
		{
			title:                   "reset to zero default value",
			value:                   "100",
			expectedValue:           100,
			defaultValue:            0,
			expectedAfterResetValue: 0,
			setDefault:              true,
		},
		{
			title:                   "reset to non-zero default value",
			value:                   "100",
			expectedValue:           100,
			defaultValue:            200,
			expectedAfterResetValue: 200,
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.UInt16("long", "usage")
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
