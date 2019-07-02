package flags_test

import (
	"errors"
	"reflect"
	"testing"

	"go.xitonix.io/flags"
)

func TestStringSliceMap(t *testing.T) {
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
			f := flags.StringSliceMap(tc.long, tc.usage)
			checkFlagInitialState(t, f, "[string][]string", tc.expectedUsage, tc.expectedLong, "")
			checkSliceMapFlagValues(t, map[string][]string{}, f.Get(), f.Var())
		})
	}
}

func TestStringSliceMapP(t *testing.T) {
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
			f := flags.StringSliceMapP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "[string][]string", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkSliceMapFlagValues(t, map[string][]string{}, f.Get(), f.Var())
		})
	}
}

func TestStringSliceMapFlag_WithKey(t *testing.T) {
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
			f := flags.StringSliceMap("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestStringSliceMapFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         map[string][]string
		expectedDefaultValue map[string][]string
	}{
		{
			title:                "empty default value",
			defaultValue:         map[string][]string{},
			expectedDefaultValue: map[string][]string{},
		},
		{
			title:                "non empty default value",
			defaultValue:         map[string][]string{"k": {"v"}},
			expectedDefaultValue: map[string][]string{"k": {"v"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSliceMap("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default().(map[string][]string)

			for eKey, eVal := range tc.expectedDefaultValue {
				if aVal, ok := actual[eKey]; !ok || len(aVal) != len(eVal) || !reflect.DeepEqual(aVal, eVal) {
					t.Errorf("Expected default value for '%s' key: %s Actual: %s", eKey, eVal, aVal)
				}
			}
		})
	}
}

func TestStringSliceMapFlag_Hide(t *testing.T) {
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
			f := flags.StringSliceMap("long", "usage")
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

func TestStringSliceMapFlag_IsDeprecated(t *testing.T) {
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
			f := flags.StringSliceMap("long", "usage")
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

func TestStringSliceMapFlag_IsRequired(t *testing.T) {
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
			f := flags.StringSliceMap("long", "usage")
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

func TestStringSliceMapFlag_WithDelimiter(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		delimiter     string
		expectedValue map[string][]string
	}{
		{
			title:         "empty delimiter",
			value:         `{"key" : "value1,value2"}`,
			expectedValue: map[string][]string{"key": {"value1", "value2"}},
		},
		{
			title:         "white space delimiter with white spaced input",
			value:         `{"key" : "value1 value2"}`,
			expectedValue: map[string][]string{"key": {"value1", "value2"}},
			delimiter:     " ",
		},
		{
			title:         "white space delimiter with none white spaced input",
			value:         `{"key" : "value1,value2"}`,
			expectedValue: map[string][]string{"key": {"value1,value2"}},
			delimiter:     " ",
		},
		{
			title:         "none white space delimiter",
			value:         `{"key" : "value1|value2"}`,
			expectedValue: map[string][]string{"key": {"value1", "value2"}},
			delimiter:     "|",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSliceMap("long", "usage").WithDelimiter(tc.delimiter)
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceMapFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestStringSliceMapFlag_WithTrimming(t *testing.T) {
	testCases := []struct {
		title          string
		value          string
		enableTrimming bool
		expectedValue  map[string][]string
	}{
		{
			title:          "without trimming",
			enableTrimming: false,
			value:          `{ "key": "  abc  ,  xyz  "}  `,
			expectedValue:  map[string][]string{"key": {"  abc  ", "  xyz  "}},
		},
		{
			title:          "with trimming",
			enableTrimming: true,
			value:          `{ "key": "  abc  ,  xyz  "}  `,
			expectedValue:  map[string][]string{"key": {"abc", "xyz"}},
		},
		{
			title:          "only white space input without trimming",
			enableTrimming: false,
			value:          `{ "key": "  ,   "}  `,
			expectedValue:  map[string][]string{"key": {"  ", "   "}},
		},
		{
			title:          "only white space input with trimming",
			enableTrimming: true,
			value:          `{ "key": "  ,   "}  `,
			expectedValue:  map[string][]string{"key": {"", ""}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSliceMap("long", "usage")
			if tc.enableTrimming {
				f = f.WithTrimming()
			}
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceMapFlag(t, f, err, "", tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestStringSliceMapFlag_Set(t *testing.T) {
	empty := make(map[string][]string)
	testCases := []struct {
		title         string
		value         string
		expectedValue map[string][]string
		expectedError string
	}{
		{
			title:         "empty value",
			value:         "",
			expectedValue: empty,
		},
		{
			title:         "empty map value",
			value:         "{}",
			expectedValue: empty,
		},
		{
			title:         "empty map value with white space",
			value:         " {    } ",
			expectedValue: empty,
		},
		{
			title:         "white space value",
			value:         "   ",
			expectedValue: empty,
		},
		{
			title:         "value with white space",
			value:         `  { "key" : " value1 , value2 " }  `,
			expectedValue: map[string][]string{"key": {" value1 ", " value2 "}},
		},
		{
			title:         "value without white space",
			value:         `  { "key" : "value1,value2" }  `,
			expectedValue: map[string][]string{"key": {"value1", "value2"}},
		},
		{
			title:         "invalid value",
			value:         " invalid ",
			expectedError: "'invalid' is not a valid [string][]string value for --long",
			expectedValue: empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSliceMap("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceMapFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestStringSliceMapFlag_Validation(t *testing.T) {
	empty := make(map[string][]string)
	testCases := []struct {
		title         string
		value         string
		expectedValue map[string][]string
		validationCB  func(key string, value []string) error
		expectedError string
	}{
		{
			title:         "nil validation callback",
			value:         `{"key" : "value1,value2"}`,
			expectedValue: map[string][]string{"key": {"value1", "value2"}},
			expectedError: "",
		},
		{
			title: "validation callback with no validation error",
			validationCB: func(k string, v []string) error {
				return nil
			},
			value:         `{"key" : "value1,value2"}`,
			expectedValue: map[string][]string{"key": {"value1", "value2"}},
		},
		{
			title: "validation callback with validation error",
			validationCB: func(k string, v []string) error {
				return errors.New("validation callback failed")
			},
			value:         `{"key" : "value1,value2"}`,
			expectedError: "validation callback failed",
			expectedValue: empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSliceMap("mappings", "usage")
			fVar := f.Var()
			f = f.WithValidationCallback(tc.validationCB)
			err := f.Set(tc.value)
			checkSliceMapFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestStringSliceMapFlag_ResetToDefault(t *testing.T) {
	empty := make(map[string][]string)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           map[string][]string
		defaultValue            map[string][]string
		expectedAfterResetValue map[string][]string
		expectedError           string
		setDefault              bool
		expectedIsSetAfterReset bool
	}{
		{
			title:                   "reset without defining the default value",
			value:                   `{"key" : "value1,value2"}`,
			expectedValue:           map[string][]string{"key": {"value1", "value2"}},
			expectedAfterResetValue: map[string][]string{"key": {"value1", "value2"}},
			setDefault:              false,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "reset to empty default value",
			value:                   `{"key" : "value1,value2"}`,
			expectedValue:           map[string][]string{"key": {"value1", "value2"}},
			defaultValue:            empty,
			expectedAfterResetValue: empty,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to nil default value",
			value:                   `{"key" : "value1,value2"}`,
			expectedValue:           map[string][]string{"key": {"value1", "value2"}},
			defaultValue:            nil,
			expectedAfterResetValue: nil,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "reset to non-empty default value",
			value:                   `{"key" : "value1,value2"}`,
			expectedValue:           map[string][]string{"key": {"value1", "value2"}},
			defaultValue:            map[string][]string{"default_key": {"default_value"}},
			expectedAfterResetValue: map[string][]string{"default_key": {"default_value"}},
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.StringSliceMap("long", "usage")
			if tc.setDefault {
				f = f.WithDefault(tc.defaultValue)
			}
			fVar := f.Var()
			err := f.Set(tc.value)
			checkSliceMapFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)

			f.ResetToDefault()

			if f.IsSet() != tc.expectedIsSetAfterReset {
				t.Errorf("IsSet() Expected: %v, Actual: %v", tc.expectedIsSetAfterReset, f.IsSet())
			}

			checkSliceMapFlagValues(t, tc.expectedAfterResetValue, f.Get(), fVar)
		})
	}
}
