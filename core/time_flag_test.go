package core_test

import (
	"errors"
	"testing"
	"time"

	"github.com/xitonix/flags"
)

func TestTime(t *testing.T) {
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
			f := flags.Time(tc.long, tc.usage)
			checkFlagInitialState(t, f, "time", tc.expectedUsage, tc.expectedLong, "")
			checkFlagValues(t, time.Time{}, f.Get(), f.Var())
		})
	}
}

func TestTimeFlag_WithShort(t *testing.T) {
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
			f := flags.Time(tc.long, tc.usage).WithShort(tc.short)
			checkFlagInitialState(t, f, "time", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkFlagValues(t, time.Time{}, f.Get(), f.Var())
		})
	}
}

func TestTimeFlag_WithKey(t *testing.T) {
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
			f := flags.Time("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestTimeFlag_WithDefault(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		title                string
		defaultValue         time.Time
		expectedDefaultValue time.Time
	}{
		{
			title:                "zero default value",
			defaultValue:         time.Time{},
			expectedDefaultValue: time.Time{},
		},
		{
			title:                "non zero default value",
			defaultValue:         now,
			expectedDefaultValue: now,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if actual != tc.expectedDefaultValue {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestTimeFlag_Hide(t *testing.T) {
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
			f := flags.Time("long", "usage")
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

func TestTimeFlag_IsDeprecated(t *testing.T) {
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
			f := flags.Time("long", "usage")
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

func TestTimeFlag_IsRequired(t *testing.T) {
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
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Set(t *testing.T) {
	zero := time.Time{}
	testCases := []struct {
		title         string
		value         string
		expectedValue time.Time
		expectedError string
	}{
		{
			title:         "no value",
			expectedValue: zero,
		},
		{
			title:         "white space value",
			value:         "   ",
			expectedValue: zero,
		},
		{
			title:         "value with white space",
			value:         "  27/08/1980  ",
			expectedValue: time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
		},
		{
			title:         "invalid value",
			value:         "abc",
			expectedError: "is not a valid time value",
			expectedValue: zero,
		},
		{
			title:         "12hr time format with 24hr time only value",
			value:         "14:22:20AM",
			expectedError: "is not a valid time value",
			expectedValue: zero,
		},
		{
			title:         "12hr time format with 24hr full date value",
			value:         "27-08-1980T14:22:20 PM",
			expectedError: "is not a valid time value",
			expectedValue: zero,
		},
		{
			title:         "12hr time format with 24hr timestamp value",
			value:         "Aug 27 14:22:20 PM",
			expectedError: "is not a valid time value",
			expectedValue: zero,
		},
		// Full date and time with dash
		{
			title:         "full dash separated date and time with 24 hrs format and time indicator",
			value:         "27-08-1980T14:22:20",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and time with 24 hrs format and no time indicator",
			value:         "27-08-1980 14:22:20",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 24 hrs format and time indicator",
			value:         "27-08-1980T14:22:20.027081980",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 24 hrs format and no time indicator",
			value:         "27-08-1980 14:22:20.027081980",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},

		{
			title:         "full dash separated date and time with 12 hrs upper case spaced and time indicator",
			value:         "27-08-1980T02:22:20 PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and time with 12 hrs upper case spaced and no time indicator",
			value:         "27-08-1980 02:22:20 PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 12 hrs upper case spaced and time indicator",
			value:         "27-08-1980T02:22:20.027081980 PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 12 hrs upper case spaced and no time indicator",
			value:         "27-08-1980 02:22:20.027081980 PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},

		{
			title:         "full dash separated date and time with 12 hrs upper case no space and time indicator",
			value:         "27-08-1980T02:22:20PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and time with 12 hrs upper case no space and no time indicator",
			value:         "27-08-1980 02:22:20PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 12 hrs upper case no space and time indicator",
			value:         "27-08-1980T02:22:20.027081980PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 12 hrs upper case no space and no time indicator",
			value:         "27-08-1980 02:22:20.027081980PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},

		{
			title:         "full dash separated date and time with 12 hrs lower case spaced and time indicator",
			value:         "27-08-1980T02:22:20 pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and time with 12 hrs lower case spaced and no time indicator",
			value:         "27-08-1980 02:22:20 pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 12 hrs lower case spaced and time indicator",
			value:         "27-08-1980T02:22:20.027081980 pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 12 hrs lower case spaced and no time indicator",
			value:         "27-08-1980 02:22:20.027081980 pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},

		{
			title:         "full dash separated date and time with 12 hrs lower case no space and time indicator",
			value:         "27-08-1980T02:22:20pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and time with 12 hrs lower case no space and no time indicator",
			value:         "27-08-1980 02:22:20pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 12 hrs lower case no space and time indicator",
			value:         "27-08-1980T02:22:20.027081980pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full dash separated date and nano time with 12 hrs lower case no space and no time indicator",
			value:         "27-08-1980 02:22:20.027081980pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		// Full date and time with forward slash
		{
			title:         "full forward slash separated date and time with 24 hrs format and time indicator",
			value:         "27/08/1980T14:22:20",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and time with 24 hrs format and no time indicator",
			value:         "27/08/1980 14:22:20",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 24 hrs format and time indicator",
			value:         "27/08/1980T14:22:20.027081980",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 24 hrs format and no time indicator",
			value:         "27/08/1980 14:22:20.027081980",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},

		{
			title:         "full forward slash separated date and time with 12 hrs upper case spaced and time indicator",
			value:         "27/08/1980T02:22:20 PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and time with 12 hrs upper case spaced and no time indicator",
			value:         "27/08/1980 02:22:20 PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 12 hrs upper case spaced and time indicator",
			value:         "27/08/1980T02:22:20.027081980 PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 12 hrs upper case spaced and no time indicator",
			value:         "27/08/1980 02:22:20.027081980 PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full forward slash separated date and time with 12 hrs upper case no space and time indicator",
			value:         "27/08/1980T02:22:20PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and time with 12 hrs upper case no space and no time indicator",
			value:         "27/08/1980 02:22:20PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 12 hrs upper case no space and time indicator",
			value:         "27/08/1980T02:22:20.027081980PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 12 hrs upper case no space and no time indicator",
			value:         "27/08/1980 02:22:20.027081980PM",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},

		{
			title:         "full forward slash separated date and time with 12 hrs lower case spaced and time indicator",
			value:         "27/08/1980T02:22:20 pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and time with 12 hrs lower case spaced and no time indicator",
			value:         "27/08/1980 02:22:20 pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 12 hrs lower case spaced and time indicator",
			value:         "27/08/1980T02:22:20.027081980 pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 12 hrs lower case spaced and no time indicator",
			value:         "27/08/1980 02:22:20.027081980 pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},

		{
			title:         "full forward slash separated date and time with 12 hrs lower case no space and time indicator",
			value:         "27/08/1980T02:22:20pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and time with 12 hrs lower case no space and no time indicator",
			value:         "27/08/1980 02:22:20pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 12 hrs lower case no space and time indicator",
			value:         "27/08/1980T02:22:20.027081980pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		{
			title:         "full forward slash separated date and nano time with 12 hrs lower case no space and no time indicator",
			value:         "27/08/1980 02:22:20.027081980pm",
			expectedValue: time.Date(1980, 8, 27, 14, 22, 20, 27081980, time.UTC),
		},
		// Date Only
		{
			title:         "full forward slash separated date only",
			value:         "27/08/1980",
			expectedValue: time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
		},
		{
			title:         "full dash separated date only",
			value:         "27-08-1980",
			expectedValue: time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
		},
		// Timestamp
		{
			title:         "time stamp with 24 hrs time format and no nano seconds",
			value:         "Aug 27 14:22:20",
			expectedValue: time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "lowercase time stamp with 24 hrs time format and no nano seconds",
			value:         "aug 27 14:22:20",
			expectedValue: time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "time stamp with 24 hrs time format and nano seconds",
			value:         "Aug 27 14:22:20.000000876",
			expectedValue: time.Date(0, 8, 27, 14, 22, 20, 876, time.UTC),
		},
		{
			title:         "time stamp with uppercase 12 hrs spaced time format and nano seconds",
			value:         "Aug 27 02:22:20.000000876 PM",
			expectedValue: time.Date(0, 8, 27, 14, 22, 20, 876, time.UTC),
		},
		{
			title:         "time stamp with uppercase 12 hrs time format and nano seconds",
			value:         "Aug 27 02:22:20.000000876PM",
			expectedValue: time.Date(0, 8, 27, 14, 22, 20, 876, time.UTC),
		},
		{
			title:         "time stamp with lowercase 12 hrs spaced time format and nano seconds",
			value:         "Aug 27 02:22:20.000000876 pm",
			expectedValue: time.Date(0, 8, 27, 14, 22, 20, 876, time.UTC),
		},
		{
			title:         "time stamp with lowercase 12 hrs time format and nano seconds",
			value:         "Aug 27 02:22:20.000000876pm",
			expectedValue: time.Date(0, 8, 27, 14, 22, 20, 876, time.UTC),
		},
		// Time only
		{
			title:         "24 hrs time format and no nano seconds",
			value:         "14:22:20",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "24 hrs time format and nano seconds",
			value:         "14:22:20.999999999",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title:         "spaced 12 hrs uppercase time format and no nano seconds",
			value:         "02:22:20 PM",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "spaced 12 hrs uppercase time format and nano seconds",
			value:         "02:22:20.999999999 PM",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title:         "12 hrs uppercase time format and no nano seconds",
			value:         "02:22:20PM",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "12 hrs uppercase time format and nano seconds",
			value:         "02:22:20.999999999PM",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title:         "spaced 12 hrs lowercase time format and no nano seconds",
			value:         "02:22:20 pm",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "spaced 12 hrs lowercase time format and nano seconds",
			value:         "02:22:20.999999999 pm",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title:         "12 hrs lowercase time format and no nano seconds",
			value:         "02:22:20pm",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title:         "12 hrs lowercase time format and nano seconds",
			value:         "02:22:20.999999999pm",
			expectedValue: time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestTimeFlag_Slashed_24Hrs_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T14:22:20.999999999",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980T14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T14:22:20.999999999",
			expectedError:     "25/12/2019T14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27/08/1980T14:22:20.999999999.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T14:22:20.999999999",
			expectedError:     "25/12/2019T14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27/08/1980T14:22:20.999999999,27/08/1981T14:22:20.999999999.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T14:22:20.999999999",
			expectedError:     "25/12/2019T14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27/08/1980T14:22:20.999999999,27/08/1981T14:22:20.999999999,27/08/1982T14:22:20.999999999.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T14:22:20.999999999",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T14:22:20.999999999",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T14:22:20.999999999",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T14:22:20",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980T14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T14:22:20",
			expectedError:     "25/12/2019T14:22:20 is not an acceptable value for --long. You must pick a value from 27/08/1980T14:22:20.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T14:22:20",
			expectedError:     "25/12/2019T14:22:20 is not an acceptable value for --long. You must pick a value from 27/08/1980T14:22:20,27/08/1981T14:22:20.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T14:22:20",
			expectedError:     "25/12/2019T14:22:20 is not an acceptable value for --long. You must pick a value from 27/08/1980T14:22:20,27/08/1981T14:22:20,27/08/1982T14:22:20.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T14:22:20",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T14:22:20",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T14:22:20",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_24Hrs_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 14:22:20.999999999",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980 14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 14:22:20.999999999",
			expectedError:     "25/12/2019 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27/08/1980 14:22:20.999999999.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 14:22:20.999999999",
			expectedError:     "25/12/2019 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27/08/1980 14:22:20.999999999,27/08/1981 14:22:20.999999999.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 14:22:20.999999999",
			expectedError:     "25/12/2019 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27/08/1980 14:22:20.999999999,27/08/1981 14:22:20.999999999,27/08/1982 14:22:20.999999999.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 14:22:20.999999999",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 14:22:20.999999999",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 14:22:20.999999999",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 14:22:20",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980 14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 14:22:20",
			expectedError:     "25/12/2019 14:22:20 is not an acceptable value for --long. You must pick a value from 27/08/1980 14:22:20.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 14:22:20",
			expectedError:     "25/12/2019 14:22:20 is not an acceptable value for --long. You must pick a value from 27/08/1980 14:22:20,27/08/1981 14:22:20.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 14:22:20",
			expectedError:     "25/12/2019 14:22:20 is not an acceptable value for --long. You must pick a value from 27/08/1980 14:22:20,27/08/1981 14:22:20,27/08/1982 14:22:20.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 14:22:20",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 14:22:20",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 14:22:20",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_Uppercase_12Hrs_Spaced_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T02:22:20.999999999 PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980T02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999 PM",
			expectedError:     "25/12/2019T02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999 PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999 PM",
			expectedError:     "25/12/2019T02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999 PM,27/08/1981T02:22:20.999999999 PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999 PM",
			expectedError:     "25/12/2019T02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999 PM,27/08/1981T02:22:20.999999999 PM,27/08/1982T02:22:20.999999999 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20.999999999 PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20.999999999 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999 PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T02:22:20 PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980T02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T02:22:20 PM",
			expectedError:     "25/12/2019T02:22:20 PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20 PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20 PM",
			expectedError:     "25/12/2019T02:22:20 PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20 PM,27/08/1981T02:22:20 PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20 PM",
			expectedError:     "25/12/2019T02:22:20 PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20 PM,27/08/1981T02:22:20 PM,27/08/1982T02:22:20 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20 PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20 PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_Uppercase_12Hrs_Spaced_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 02:22:20.999999999 PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980 02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999 PM",
			expectedError:     "25/12/2019 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999 PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999 PM",
			expectedError:     "25/12/2019 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999 PM,27/08/1981 02:22:20.999999999 PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999 PM",
			expectedError:     "25/12/2019 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999 PM,27/08/1981 02:22:20.999999999 PM,27/08/1982 02:22:20.999999999 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20.999999999 PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20.999999999 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999 PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 02:22:20 PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980 02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 02:22:20 PM",
			expectedError:     "25/12/2019 02:22:20 PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20 PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20 PM",
			expectedError:     "25/12/2019 02:22:20 PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20 PM,27/08/1981 02:22:20 PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20 PM",
			expectedError:     "25/12/2019 02:22:20 PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20 PM,27/08/1981 02:22:20 PM,27/08/1982 02:22:20 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20 PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20 PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_Lowercase_12Hrs_Spaced_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T02:22:20.999999999 pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980T02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999 pm",
			expectedError:     "25/12/2019T02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999 pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999 pm",
			expectedError:     "25/12/2019T02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999 pm,27/08/1981T02:22:20.999999999 pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999 pm",
			expectedError:     "25/12/2019T02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999 pm,27/08/1981T02:22:20.999999999 pm,27/08/1982T02:22:20.999999999 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20.999999999 pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20.999999999 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999 pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T02:22:20 pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980T02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T02:22:20 pm",
			expectedError:     "25/12/2019T02:22:20 pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20 pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20 pm",
			expectedError:     "25/12/2019T02:22:20 pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20 pm,27/08/1981T02:22:20 pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20 pm",
			expectedError:     "25/12/2019T02:22:20 pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20 pm,27/08/1981T02:22:20 pm,27/08/1982T02:22:20 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20 pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20 pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_Lowercase_12Hrs_Spaced_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 02:22:20.999999999 pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980 02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999 pm",
			expectedError:     "25/12/2019 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999 pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999 pm",
			expectedError:     "25/12/2019 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999 pm,27/08/1981 02:22:20.999999999 pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999 pm",
			expectedError:     "25/12/2019 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999 pm,27/08/1981 02:22:20.999999999 pm,27/08/1982 02:22:20.999999999 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20.999999999 pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20.999999999 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999 pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 02:22:20 pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980 02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 02:22:20 pm",
			expectedError:     "25/12/2019 02:22:20 pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20 pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20 pm",
			expectedError:     "25/12/2019 02:22:20 pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20 pm,27/08/1981 02:22:20 pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20 pm",
			expectedError:     "25/12/2019 02:22:20 pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20 pm,27/08/1981 02:22:20 pm,27/08/1982 02:22:20 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20 pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20 pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_Uppercase_12Hrs_Attached_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T02:22:20.999999999PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980T02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999PM",
			expectedError:     "25/12/2019T02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999PM",
			expectedError:     "25/12/2019T02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999PM,27/08/1981T02:22:20.999999999PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999PM",
			expectedError:     "25/12/2019T02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999PM,27/08/1981T02:22:20.999999999PM,27/08/1982T02:22:20.999999999PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20.999999999PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20.999999999PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T02:22:20PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980T02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T02:22:20PM",
			expectedError:     "25/12/2019T02:22:20PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20PM",
			expectedError:     "25/12/2019T02:22:20PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20PM,27/08/1981T02:22:20PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20PM",
			expectedError:     "25/12/2019T02:22:20PM is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20PM,27/08/1981T02:22:20PM,27/08/1982T02:22:20PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_Uppercase_12Hrs_Attached_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 02:22:20.999999999PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980 02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999PM",
			expectedError:     "25/12/2019 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999PM",
			expectedError:     "25/12/2019 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999PM,27/08/1981 02:22:20.999999999PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999PM",
			expectedError:     "25/12/2019 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999PM,27/08/1981 02:22:20.999999999PM,27/08/1982 02:22:20.999999999PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20.999999999PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20.999999999PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 02:22:20PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980 02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 02:22:20PM",
			expectedError:     "25/12/2019 02:22:20PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20PM",
			expectedError:     "25/12/2019 02:22:20PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20PM,27/08/1981 02:22:20PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20PM",
			expectedError:     "25/12/2019 02:22:20PM is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20PM,27/08/1981 02:22:20PM,27/08/1982 02:22:20PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_Lowercase_12Hrs_Attached_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T02:22:20.999999999pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980T02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999pm",
			expectedError:     "25/12/2019T02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999pm",
			expectedError:     "25/12/2019T02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999pm,27/08/1981T02:22:20.999999999pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999pm",
			expectedError:     "25/12/2019T02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20.999999999pm,27/08/1981T02:22:20.999999999pm,27/08/1982T02:22:20.999999999pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20.999999999pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20.999999999pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20.999999999pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980T02:22:20pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980T02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980T02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980T02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019T02:22:20pm",
			expectedError:     "25/12/2019T02:22:20pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20pm",
			expectedError:     "25/12/2019T02:22:20pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20pm,27/08/1981T02:22:20pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20pm",
			expectedError:     "25/12/2019T02:22:20pm is not an acceptable value for --long. You must pick a value from 27/08/1980T02:22:20pm,27/08/1981T02:22:20pm,27/08/1982T02:22:20pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019T02:22:20pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019T02:22:20pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_Lowercase_12Hrs_Attached_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 02:22:20.999999999pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27/08/1980 02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999pm",
			expectedError:     "25/12/2019 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999pm",
			expectedError:     "25/12/2019 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999pm,27/08/1981 02:22:20.999999999pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999pm",
			expectedError:     "25/12/2019 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20.999999999pm,27/08/1981 02:22:20.999999999pm,27/08/1982 02:22:20.999999999pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20.999999999pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20.999999999pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20.999999999pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980 02:22:20pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980 02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980 02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980 02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019 02:22:20pm",
			expectedError:     "25/12/2019 02:22:20pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20pm",
			expectedError:     "25/12/2019 02:22:20pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20pm,27/08/1981 02:22:20pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20pm",
			expectedError:     "25/12/2019 02:22:20pm is not an acceptable value for --long. You must pick a value from 27/08/1980 02:22:20pm,27/08/1981 02:22:20pm,27/08/1982 02:22:20pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019 02:22:20pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019 02:22:20pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_24Hrs_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T14:22:20.999999999",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980T14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T14:22:20.999999999",
			expectedError:     "25-12-2019T14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27-08-1980T14:22:20.999999999.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T14:22:20.999999999",
			expectedError:     "25-12-2019T14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27-08-1980T14:22:20.999999999,27-08-1981T14:22:20.999999999.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T14:22:20.999999999",
			expectedError:     "25-12-2019T14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27-08-1980T14:22:20.999999999,27-08-1981T14:22:20.999999999,27-08-1982T14:22:20.999999999.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T14:22:20.999999999",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T14:22:20.999999999",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T14:22:20.999999999",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T14:22:20",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980T14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T14:22:20",
			expectedError:     "25-12-2019T14:22:20 is not an acceptable value for --long. You must pick a value from 27-08-1980T14:22:20.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T14:22:20",
			expectedError:     "25-12-2019T14:22:20 is not an acceptable value for --long. You must pick a value from 27-08-1980T14:22:20,27-08-1981T14:22:20.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T14:22:20",
			expectedError:     "25-12-2019T14:22:20 is not an acceptable value for --long. You must pick a value from 27-08-1980T14:22:20,27-08-1981T14:22:20,27-08-1982T14:22:20.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T14:22:20",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T14:22:20",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T14:22:20",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_24Hrs_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 14:22:20.999999999",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980 14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 14:22:20.999999999",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 14:22:20.999999999",
			expectedError:     "25-12-2019 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27-08-1980 14:22:20.999999999.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 14:22:20.999999999",
			expectedError:     "25-12-2019 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27-08-1980 14:22:20.999999999,27-08-1981 14:22:20.999999999.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 14:22:20.999999999",
			expectedError:     "25-12-2019 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 27-08-1980 14:22:20.999999999,27-08-1981 14:22:20.999999999,27-08-1982 14:22:20.999999999.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 14:22:20.999999999",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 14:22:20.999999999",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 14:22:20.999999999",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 14:22:20",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980 14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 14:22:20",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 14:22:20",
			expectedError:     "25-12-2019 14:22:20 is not an acceptable value for --long. You must pick a value from 27-08-1980 14:22:20.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 14:22:20",
			expectedError:     "25-12-2019 14:22:20 is not an acceptable value for --long. You must pick a value from 27-08-1980 14:22:20,27-08-1981 14:22:20.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 14:22:20",
			expectedError:     "25-12-2019 14:22:20 is not an acceptable value for --long. You must pick a value from 27-08-1980 14:22:20,27-08-1981 14:22:20,27-08-1982 14:22:20.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 14:22:20",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 14:22:20",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 14:22:20",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_Uppercase_12Hrs_Spaced_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T02:22:20.999999999 PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980T02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999 PM",
			expectedError:     "25-12-2019T02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999 PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999 PM",
			expectedError:     "25-12-2019T02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999 PM,27-08-1981T02:22:20.999999999 PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999 PM",
			expectedError:     "25-12-2019T02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999 PM,27-08-1981T02:22:20.999999999 PM,27-08-1982T02:22:20.999999999 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20.999999999 PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20.999999999 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999 PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T02:22:20 PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980T02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T02:22:20 PM",
			expectedError:     "25-12-2019T02:22:20 PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20 PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20 PM",
			expectedError:     "25-12-2019T02:22:20 PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20 PM,27-08-1981T02:22:20 PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20 PM",
			expectedError:     "25-12-2019T02:22:20 PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20 PM,27-08-1981T02:22:20 PM,27-08-1982T02:22:20 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20 PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20 PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_Uppercase_12Hrs_Spaced_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 02:22:20.999999999 PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980 02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 02:22:20.999999999 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999 PM",
			expectedError:     "25-12-2019 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999 PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999 PM",
			expectedError:     "25-12-2019 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999 PM,27-08-1981 02:22:20.999999999 PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999 PM",
			expectedError:     "25-12-2019 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999 PM,27-08-1981 02:22:20.999999999 PM,27-08-1982 02:22:20.999999999 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20.999999999 PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20.999999999 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999 PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 02:22:20 PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980 02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 02:22:20 PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 02:22:20 PM",
			expectedError:     "25-12-2019 02:22:20 PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20 PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20 PM",
			expectedError:     "25-12-2019 02:22:20 PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20 PM,27-08-1981 02:22:20 PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20 PM",
			expectedError:     "25-12-2019 02:22:20 PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20 PM,27-08-1981 02:22:20 PM,27-08-1982 02:22:20 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20 PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20 PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_Lowercase_12Hrs_Spaced_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T02:22:20.999999999 pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980T02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999 pm",
			expectedError:     "25-12-2019T02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999 pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999 pm",
			expectedError:     "25-12-2019T02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999 pm,27-08-1981T02:22:20.999999999 pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999 pm",
			expectedError:     "25-12-2019T02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999 pm,27-08-1981T02:22:20.999999999 pm,27-08-1982T02:22:20.999999999 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20.999999999 pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20.999999999 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999 pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T02:22:20 pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980T02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T02:22:20 pm",
			expectedError:     "25-12-2019T02:22:20 pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20 pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20 pm",
			expectedError:     "25-12-2019T02:22:20 pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20 pm,27-08-1981T02:22:20 pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20 pm",
			expectedError:     "25-12-2019T02:22:20 pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20 pm,27-08-1981T02:22:20 pm,27-08-1982T02:22:20 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20 pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20 pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_Lowercase_12Hrs_Spaced_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 02:22:20.999999999 pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980 02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 02:22:20.999999999 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999 pm",
			expectedError:     "25-12-2019 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999 pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999 pm",
			expectedError:     "25-12-2019 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999 pm,27-08-1981 02:22:20.999999999 pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999 pm",
			expectedError:     "25-12-2019 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999 pm,27-08-1981 02:22:20.999999999 pm,27-08-1982 02:22:20.999999999 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20.999999999 pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20.999999999 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999 pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 02:22:20 pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980 02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 02:22:20 pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 02:22:20 pm",
			expectedError:     "25-12-2019 02:22:20 pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20 pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20 pm",
			expectedError:     "25-12-2019 02:22:20 pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20 pm,27-08-1981 02:22:20 pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20 pm",
			expectedError:     "25-12-2019 02:22:20 pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20 pm,27-08-1981 02:22:20 pm,27-08-1982 02:22:20 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20 pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20 pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_Uppercase_12Hrs_Attached_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T02:22:20.999999999PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980T02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999PM",
			expectedError:     "25-12-2019T02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999PM",
			expectedError:     "25-12-2019T02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999PM,27-08-1981T02:22:20.999999999PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999PM",
			expectedError:     "25-12-2019T02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999PM,27-08-1981T02:22:20.999999999PM,27-08-1982T02:22:20.999999999PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20.999999999PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20.999999999PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T02:22:20PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980T02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T02:22:20PM",
			expectedError:     "25-12-2019T02:22:20PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20PM",
			expectedError:     "25-12-2019T02:22:20PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20PM,27-08-1981T02:22:20PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20PM",
			expectedError:     "25-12-2019T02:22:20PM is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20PM,27-08-1981T02:22:20PM,27-08-1982T02:22:20PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_Uppercase_12Hrs_Attached_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 02:22:20.999999999PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980 02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 02:22:20.999999999PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999PM",
			expectedError:     "25-12-2019 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999PM",
			expectedError:     "25-12-2019 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999PM,27-08-1981 02:22:20.999999999PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999PM",
			expectedError:     "25-12-2019 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999PM,27-08-1981 02:22:20.999999999PM,27-08-1982 02:22:20.999999999PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20.999999999PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20.999999999PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 02:22:20PM",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980 02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 02:22:20PM",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 02:22:20PM",
			expectedError:     "25-12-2019 02:22:20PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20PM",
			expectedError:     "25-12-2019 02:22:20PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20PM,27-08-1981 02:22:20PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20PM",
			expectedError:     "25-12-2019 02:22:20PM is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20PM,27-08-1981 02:22:20PM,27-08-1982 02:22:20PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20PM",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_Lowercase_12Hrs_Attached_With_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T02:22:20.999999999pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980T02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999pm",
			expectedError:     "25-12-2019T02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999pm",
			expectedError:     "25-12-2019T02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999pm,27-08-1981T02:22:20.999999999pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999pm",
			expectedError:     "25-12-2019T02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20.999999999pm,27-08-1981T02:22:20.999999999pm,27-08-1982T02:22:20.999999999pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20.999999999pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20.999999999pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20.999999999pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980T02:22:20pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980T02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980T02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980T02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019T02:22:20pm",
			expectedError:     "25-12-2019T02:22:20pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20pm",
			expectedError:     "25-12-2019T02:22:20pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20pm,27-08-1981T02:22:20pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20pm",
			expectedError:     "25-12-2019T02:22:20pm is not an acceptable value for --long. You must pick a value from 27-08-1980T02:22:20pm,27-08-1981T02:22:20pm,27-08-1982T02:22:20pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019T02:22:20pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019T02:22:20pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_Lowercase_12Hrs_Attached_Without_T_Full_Date_Time_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 02:22:20.999999999pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "27-08-1980 02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 02:22:20.999999999pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999pm",
			expectedError:     "25-12-2019 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999pm",
			expectedError:     "25-12-2019 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999pm,27-08-1981 02:22:20.999999999pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999pm",
			expectedError:     "25-12-2019 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20.999999999pm,27-08-1981 02:22:20.999999999pm,27-08-1982 02:22:20.999999999pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20.999999999pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20.999999999pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20.999999999pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980 02:22:20pm",
			expectedValue:   time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980 02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980 02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980 02:22:20pm",
			expectedValue:     time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019 02:22:20pm",
			expectedError:     "25-12-2019 02:22:20pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20pm",
			expectedError:     "25-12-2019 02:22:20pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20pm,27-08-1981 02:22:20pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1981, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(1982, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20pm",
			expectedError:     "25-12-2019 02:22:20pm is not an acceptable value for --long. You must pick a value from 27-08-1980 02:22:20pm,27-08-1981 02:22:20pm,27-08-1982 02:22:20pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20pm",
			expectedValue:   time.Date(2019, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019 02:22:20pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019 02:22:20pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Slashed_Date_Only_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27/08/1980",
			expectedValue:   time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27/08/1980",
			expectedValue:     time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27/08/1980",
			expectedValue:     time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27/08/1980",
			expectedValue:     time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC)},
			setValidationList: true,
			value:             "25/12/2019",
			expectedError:     "25/12/2019 is not an acceptable value for --long. You must pick a value from 27/08/1980.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
				time.Date(1981, 8, 27, 0, 0, 0, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019",
			expectedError:     "25/12/2019 is not an acceptable value for --long. You must pick a value from 27/08/1980,27/08/1981.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
				time.Date(1981, 8, 27, 0, 0, 0, 0, time.UTC),
				time.Date(1982, 8, 27, 0, 0, 0, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019",
			expectedError:     "25/12/2019 is not an acceptable value for --long. You must pick a value from 27/08/1980,27/08/1981,27/08/1982.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25/12/2019",
			expectedValue:   time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25/12/2019",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25/12/2019",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Dashed_Date_Only_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "27-08-1980",
			expectedValue:   time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "27-08-1980",
			expectedValue:     time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "27-08-1980",
			expectedValue:     time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "27-08-1980",
			expectedValue:     time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC)},
			setValidationList: true,
			value:             "25-12-2019",
			expectedError:     "25-12-2019 is not an acceptable value for --long. You must pick a value from 27-08-1980.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
				time.Date(1981, 8, 27, 0, 0, 0, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019",
			expectedError:     "25-12-2019 is not an acceptable value for --long. You must pick a value from 27-08-1980,27-08-1981.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
				time.Date(1981, 8, 27, 0, 0, 0, 0, time.UTC),
				time.Date(1982, 8, 27, 0, 0, 0, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019",
			expectedError:     "25-12-2019 is not an acceptable value for --long. You must pick a value from 27-08-1980,27-08-1981,27-08-1982.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "25-12-2019",
			expectedValue:   time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "25-12-2019",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(1980, 8, 27, 0, 0, 0, 0, time.UTC),
			},
			setValidationList: true,
			value:             "25-12-2019",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_24Hrs_With_Timestamp_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "Aug 27 14:22:20.999999999",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "Aug 27 14:22:20.999999999",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 14:22:20.999999999",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 14:22:20.999999999",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 14:22:20.999999999",
			expectedError:     "Dec 25 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from Aug 27 14:22:20.999999999.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 14:22:20.999999999",
			expectedError:     "Dec 25 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from Aug 27 14:22:20.999999999,Aug 28 14:22:20.999999999.",
		},
		{
			title: "duplicate items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 14:22:20.999999999",
			expectedError:     "Dec 25 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from Aug 27 14:22:20.999999999.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 14:22:20.999999999",
			expectedError:     "Dec 25 14:22:20.999999999 is not an acceptable value for --long. You must pick a value from Aug 27 14:22:20.999999999,Aug 28 14:22:20.999999999,Aug 29 14:22:20.999999999.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 14:22:20.999999999",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 14:22:20.999999999",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 14:22:20.999999999",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "Aug 27 14:22:20",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "Aug 27 14:22:20",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 14:22:20",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 14:22:20",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 14:22:20",
			expectedError:     "Dec 25 14:22:20 is not an acceptable value for --long. You must pick a value from Aug 27 14:22:20.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 14:22:20",
			expectedError:     "Dec 25 14:22:20 is not an acceptable value for --long. You must pick a value from Aug 27 14:22:20,Aug 28 14:22:20.",
		},
		{
			title: "duplicate items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 14:22:20",
			expectedError:     "Dec 25 14:22:20 is not an acceptable value for --long. You must pick a value from Aug 27 14:22:20.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 14:22:20",
			expectedError:     "Dec 25 14:22:20 is not an acceptable value for --long. You must pick a value from Aug 27 14:22:20,Aug 28 14:22:20,Aug 29 14:22:20.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 14:22:20",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 14:22:20",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 14:22:20",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Uppercase_12Hrs_Spaced_With_Timestamp_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "Aug 27 02:22:20.999999999 PM",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "Aug 27 02:22:20.999999999 PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 02:22:20.999999999 PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 02:22:20.999999999 PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 PM",
			expectedError:     "Dec 25 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999 PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 PM",
			expectedError:     "Dec 25 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999 PM,Aug 28 02:22:20.999999999 PM.",
		},
		{
			title: "duplicate items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 PM",
			expectedError:     "Dec 25 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999 PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 PM",
			expectedError:     "Dec 25 02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999 PM,Aug 28 02:22:20.999999999 PM,Aug 29 02:22:20.999999999 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20.999999999 PM",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20.999999999 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "Aug 27 02:22:20 PM",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "Aug 27 02:22:20 PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 02:22:20 PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 02:22:20 PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 02:22:20 PM",
			expectedError:     "Dec 25 02:22:20 PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20 PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20 PM",
			expectedError:     "Dec 25 02:22:20 PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20 PM,Aug 28 02:22:20 PM.",
		},
		{
			title: "duplicate items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20 PM",
			expectedError:     "Dec 25 02:22:20 PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20 PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20 PM",
			expectedError:     "Dec 25 02:22:20 PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20 PM,Aug 28 02:22:20 PM,Aug 29 02:22:20 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20 PM",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20 PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Lowercase_12Hrs_Spaced_With_Timestamp_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "Aug 27 02:22:20.999999999 pm",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "Aug 27 02:22:20.999999999 pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 02:22:20.999999999 pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 02:22:20.999999999 pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 pm",
			expectedError:     "Dec 25 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999 pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 pm",
			expectedError:     "Dec 25 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999 pm,Aug 28 02:22:20.999999999 pm.",
		},
		{
			title: "duplicate items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 pm",
			expectedError:     "Dec 25 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999 pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 pm",
			expectedError:     "Dec 25 02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999 pm,Aug 28 02:22:20.999999999 pm,Aug 29 02:22:20.999999999 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20.999999999 pm",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20.999999999 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999 pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "Aug 27 02:22:20 pm",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "Aug 27 02:22:20 pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 02:22:20 pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 02:22:20 pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 02:22:20 pm",
			expectedError:     "Dec 25 02:22:20 pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20 pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20 pm",
			expectedError:     "Dec 25 02:22:20 pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20 pm,Aug 28 02:22:20 pm.",
		},
		{
			title: "duplicate items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20 pm",
			expectedError:     "Dec 25 02:22:20 pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20 pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20 pm",
			expectedError:     "Dec 25 02:22:20 pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20 pm,Aug 28 02:22:20 pm,Aug 29 02:22:20 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20 pm",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20 pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Uppercase_12Hrs_Attached_With_Timestamp_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "Aug 27 02:22:20.999999999PM",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "Aug 27 02:22:20.999999999PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 02:22:20.999999999PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 02:22:20.999999999PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999PM",
			expectedError:     "Dec 25 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999PM",
			expectedError:     "Dec 25 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999PM,Aug 28 02:22:20.999999999PM.",
		},
		{
			title: "duplicate items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999PM",
			expectedError:     "Dec 25 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999PM",
			expectedError:     "Dec 25 02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999PM,Aug 28 02:22:20.999999999PM,Aug 29 02:22:20.999999999PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20.999999999PM",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20.999999999PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "Aug 27 02:22:20PM",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "Aug 27 02:22:20PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 02:22:20PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 02:22:20PM",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 02:22:20PM",
			expectedError:     "Dec 25 02:22:20PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20PM",
			expectedError:     "Dec 25 02:22:20PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20PM,Aug 28 02:22:20PM.",
		},
		{
			title: "duplicate items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20PM",
			expectedError:     "Dec 25 02:22:20PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20PM",
			expectedError:     "Dec 25 02:22:20PM is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20PM,Aug 28 02:22:20PM,Aug 29 02:22:20PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20PM",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Lowercase_12Hrs_Attached_With_Timestamp_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "Aug 27 02:22:20.999999999pm",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "Aug 27 02:22:20.999999999pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 02:22:20.999999999pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 02:22:20.999999999pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999pm",
			expectedError:     "Dec 25 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999pm",
			expectedError:     "Dec 25 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999pm,Aug 28 02:22:20.999999999pm.",
		},
		{
			title: "duplicate items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999pm",
			expectedError:     "Dec 25 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 999999999, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999pm",
			expectedError:     "Dec 25 02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20.999999999pm,Aug 28 02:22:20.999999999pm,Aug 29 02:22:20.999999999pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20.999999999pm",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20.999999999pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20.999999999pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "Aug 27 02:22:20pm",
			expectedValue:   time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "Aug 27 02:22:20pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "Aug 27 02:22:20pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "Aug 27 02:22:20pm",
			expectedValue:     time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "Dec 25 02:22:20pm",
			expectedError:     "Dec 25 02:22:20pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20pm",
			expectedError:     "Dec 25 02:22:20pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20pm,Aug 28 02:22:20pm.",
		},
		{
			title: "duplicate items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20pm",
			expectedError:     "Dec 25 02:22:20pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 28, 14, 22, 20, 0, time.UTC),
				time.Date(0, 8, 29, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20pm",
			expectedError:     "Dec 25 02:22:20pm is not an acceptable value for --long. You must pick a value from Aug 27 02:22:20pm,Aug 28 02:22:20pm,Aug 29 02:22:20pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20pm",
			expectedValue:   time.Date(0, 12, 25, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "Dec 25 02:22:20pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 8, 27, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "Dec 25 02:22:20pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_24Hrs_With_Time_Only_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "14:22:20.999999999",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "14:22:20.999999999",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "14:22:20.999999999",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "14:22:20.999999999",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "14:22:20.999999999",
			expectedError:     "14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 15:22:20.999999999.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "14:22:20.999999999",
			expectedError:     "14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 15:22:20.999999999,16:22:20.999999999.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "14:22:20.999999999",
			expectedError:     "14:22:20.999999999 is not an acceptable value for --long. You must pick a value from 15:22:20.999999999,16:22:20.999999999,17:22:20.999999999.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "14:22:20.999999999",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "14:22:20.999999999",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "14:22:20.999999999",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "14:22:20",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "14:22:20",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "14:22:20",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "14:22:20",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "14:22:20",
			expectedError:     "14:22:20 is not an acceptable value for --long. You must pick a value from 15:22:20.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "14:22:20",
			expectedError:     "14:22:20 is not an acceptable value for --long. You must pick a value from 15:22:20,16:22:20.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "14:22:20",
			expectedError:     "14:22:20 is not an acceptable value for --long. You must pick a value from 15:22:20,16:22:20,17:22:20.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "14:22:20",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "14:22:20",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "14:22:20",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Uppercase_12Hrs_Spaced_With_Time_Only_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "02:22:20.999999999 PM",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "02:22:20.999999999 PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "02:22:20.999999999 PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "02:22:20.999999999 PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "02:22:20.999999999 PM",
			expectedError:     "02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 03:22:20.999999999 PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999 PM",
			expectedError:     "02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 03:22:20.999999999 PM,04:22:20.999999999 PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999 PM",
			expectedError:     "02:22:20.999999999 PM is not an acceptable value for --long. You must pick a value from 03:22:20.999999999 PM,04:22:20.999999999 PM,05:22:20.999999999 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "02:22:20.999999999 PM",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "02:22:20.999999999 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999 PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "02:22:20 PM",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "02:22:20 PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "02:22:20 PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "02:22:20 PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "02:22:20 PM",
			expectedError:     "02:22:20 PM is not an acceptable value for --long. You must pick a value from 03:22:20 PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20 PM",
			expectedError:     "02:22:20 PM is not an acceptable value for --long. You must pick a value from 03:22:20 PM,04:22:20 PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20 PM",
			expectedError:     "02:22:20 PM is not an acceptable value for --long. You must pick a value from 03:22:20 PM,04:22:20 PM,05:22:20 PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "02:22:20 PM",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "02:22:20 PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20 PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Lowercase_12Hrs_Spaced_With_Time_Only_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "02:22:20.999999999 pm",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "02:22:20.999999999 pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "02:22:20.999999999 pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "02:22:20.999999999 pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "02:22:20.999999999 pm",
			expectedError:     "02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 03:22:20.999999999 pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999 pm",
			expectedError:     "02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 03:22:20.999999999 pm,04:22:20.999999999 pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999 pm",
			expectedError:     "02:22:20.999999999 pm is not an acceptable value for --long. You must pick a value from 03:22:20.999999999 pm,04:22:20.999999999 pm,05:22:20.999999999 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "02:22:20.999999999 pm",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "02:22:20.999999999 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999 pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "02:22:20 pm",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "02:22:20 pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "02:22:20 pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "02:22:20 pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "02:22:20 pm",
			expectedError:     "02:22:20 pm is not an acceptable value for --long. You must pick a value from 03:22:20 pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20 pm",
			expectedError:     "02:22:20 pm is not an acceptable value for --long. You must pick a value from 03:22:20 pm,04:22:20 pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20 pm",
			expectedError:     "02:22:20 pm is not an acceptable value for --long. You must pick a value from 03:22:20 pm,04:22:20 pm,05:22:20 pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "02:22:20 pm",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "02:22:20 pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20 pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Uppercase_12Hrs_Attached_With_Time_Only_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "02:22:20.999999999PM",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "02:22:20.999999999PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "02:22:20.999999999PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "02:22:20.999999999PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "02:22:20.999999999PM",
			expectedError:     "02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 03:22:20.999999999PM.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999PM",
			expectedError:     "02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 03:22:20.999999999PM,04:22:20.999999999PM.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999PM",
			expectedError:     "02:22:20.999999999PM is not an acceptable value for --long. You must pick a value from 03:22:20.999999999PM,04:22:20.999999999PM,05:22:20.999999999PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "02:22:20.999999999PM",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "02:22:20.999999999PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999PM",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "02:22:20PM",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "02:22:20PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "02:22:20PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "02:22:20PM",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "02:22:20PM",
			expectedError:     "02:22:20PM is not an acceptable value for --long. You must pick a value from 03:22:20PM.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20PM",
			expectedError:     "02:22:20PM is not an acceptable value for --long. You must pick a value from 03:22:20PM,04:22:20PM.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20PM",
			expectedError:     "02:22:20PM is not an acceptable value for --long. You must pick a value from 03:22:20PM,04:22:20PM,05:22:20PM.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "02:22:20PM",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "02:22:20PM",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20PM",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_Lowercase_12Hrs_Attached_With_Time_Only_Validation(t *testing.T) {
	testCases := []struct {
		title             string
		value             string
		expectedValue     time.Time
		validationCB      func(in time.Time) error
		setValidationCB   bool
		validationList    []time.Time
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with nano seconds",
			setValidationCB: true,
			value:           "02:22:20.999999999pm",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list with nano seconds",
			setValidationList: true,
			value:             "02:22:20.999999999pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "02:22:20.999999999pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list with nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "02:22:20.999999999pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list with nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC)},
			setValidationList: true,
			value:             "02:22:20.999999999pm",
			expectedError:     "02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 03:22:20.999999999pm.",
		},
		{
			title: "two items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999pm",
			expectedError:     "02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 03:22:20.999999999pm,04:22:20.999999999pm.",
		},
		{
			title: "three items in the validation list with nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 999999999, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999pm",
			expectedError:     "02:22:20.999999999pm is not an acceptable value for --long. You must pick a value from 03:22:20.999999999pm,04:22:20.999999999pm,05:22:20.999999999pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "02:22:20.999999999pm",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
		},
		{
			title: "validation callback with validation error and with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "02:22:20.999999999pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list with nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 999999999, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20.999999999pm",
			expectedError:     "validation callback failed",
		},

		{
			title:           "nil validation callback without nano seconds",
			setValidationCB: true,
			value:           "02:22:20pm",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:   "",
		},
		{
			title:             "nil validation list without nano seconds",
			setValidationList: true,
			value:             "02:22:20pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback without nano seconds",
			setValidationList: true,
			setValidationCB:   true,
			value:             "02:22:20pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "empty validation list without nano seconds",
			validationList:    make([]time.Time, 0),
			setValidationList: true,
			value:             "02:22:20pm",
			expectedValue:     time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			expectedError:     "",
		},
		{
			title:             "single item in the validation list without nano seconds",
			validationList:    []time.Time{time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC)},
			setValidationList: true,
			value:             "02:22:20pm",
			expectedError:     "02:22:20pm is not an acceptable value for --long. You must pick a value from 03:22:20pm.",
		},
		{
			title: "two items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20pm",
			expectedError:     "02:22:20pm is not an acceptable value for --long. You must pick a value from 03:22:20pm,04:22:20pm.",
		},
		{
			title: "three items in the validation list without nano seconds",
			validationList: []time.Time{
				time.Date(0, 1, 1, 15, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 16, 22, 20, 0, time.UTC),
				time.Date(0, 1, 1, 17, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20pm",
			expectedError:     "02:22:20pm is not an acceptable value for --long. You must pick a value from 03:22:20pm,04:22:20pm,05:22:20pm.",
		},
		{
			title: "validation callback with no validation error and nano seconds",
			validationCB: func(in time.Time) error {
				return nil
			},
			setValidationCB: true,
			value:           "02:22:20pm",
			expectedValue:   time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
		},
		{
			title: "validation callback with validation error and without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           "02:22:20pm",
			expectedError:   "validation callback failed",
		},
		{
			title: "validation callback takes priority over validation list without nano seconds",
			validationCB: func(in time.Time) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			validationList: []time.Time{
				time.Date(0, 1, 1, 14, 22, 20, 0, time.UTC),
			},
			setValidationList: true,
			value:             "02:22:20pm",
			expectedError:     "validation callback failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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

func TestTimeFlag_ResetToDefault(t *testing.T) {
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           time.Time
		defaultValue            time.Time
		expectedAfterResetValue time.Time
		expectedError           string
		setDefault              bool
	}{
		{
			title: "no value",
		},
		{
			title:                   "reset without defining the default value",
			value:                   "27/08/1980T14:22:20.999999999",
			expectedValue:           time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			expectedAfterResetValue: time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			setDefault:              false,
		},
		{
			title:                   "reset to zero default value",
			value:                   "27/08/1980T14:22:20.999999999",
			expectedValue:           time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			defaultValue:            time.Time{},
			expectedAfterResetValue: time.Time{},
			setDefault:              true,
		},
		{
			title:                   "reset to non-zero default value",
			value:                   "27/08/1980T14:22:20.999999999",
			expectedValue:           time.Date(1980, 8, 27, 14, 22, 20, 999999999, time.UTC),
			defaultValue:            time.Date(2020, 10, 17, 8, 9, 10, 11, time.UTC),
			expectedAfterResetValue: time.Date(2020, 10, 17, 8, 9, 10, 11, time.UTC),
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.Time("long", "usage")
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
