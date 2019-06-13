package flags_test

import (
	"testing"
	"time"

	"go.xitonix.io/flags"
	"go.xitonix.io/flags/test"
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
			if f.LongName() != tc.expectedLong {
				t.Errorf("Expected Long Name: %s, Actual: %s", tc.expectedLong, f.LongName())
			}
			if f.Usage() != tc.expectedUsage {
				t.Errorf("Expected Usage: %s, Actual: %s", tc.expectedUsage, f.Usage())
			}

			if f.IsDeprecated() {
				t.Error("The flag must not be marked as deprecated by default")
			}

			if f.IsHidden() {
				t.Error("The flag must not be marked as hidden by default")
			}

			if f.IsSet() {
				t.Error("The flag value must not be set initially")
			}

			if f.ShortName() != "" {
				t.Errorf("The short name was expected to be empty but it was %s", f.ShortName())
			}

			if f.Default() != nil {
				t.Errorf("The initial default value was expected to be nil, but it was %v", f.Default())
			}

			if f.Type() != "time" {
				t.Errorf("The flag type was expected to be 'time.Time', but it was %s", f.Type())
			}

			if !f.Get().IsZero() {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
		})
	}
}

func TestTimeP(t *testing.T) {
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
			short:         "Short",
			expectedShort: "Short",
		},
		{
			title:         "long and short names with white space",
			long:          " Long ",
			expectedLong:  "long",
			short:         " Short ",
			expectedShort: "Short",
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
			f := flags.TimeP(tc.long, tc.usage, tc.short)
			if f.LongName() != tc.expectedLong {
				t.Errorf("Expected Long Name: %s, Actual: %s", tc.expectedLong, f.LongName())
			}
			if f.Usage() != tc.expectedUsage {
				t.Errorf("Expected Usage: %s, Actual: %s", tc.expectedUsage, f.Usage())
			}

			if f.IsDeprecated() {
				t.Error("The flag must not be marked as deprecated by default")
			}

			if f.IsHidden() {
				t.Error("The flag must not be marked as hidden by default")
			}

			if f.IsSet() {
				t.Error("The flag value must not be set initially")
			}

			if f.ShortName() != tc.expectedShort {
				t.Errorf("The short name was expected to be %s but it was %s", tc.expectedShort, f.ShortName())
			}

			if f.Default() != nil {
				t.Errorf("The initial default value was expected to be nil, but it was %v", f.Default())
			}

			if f.Type() != "time" {
				t.Errorf("The flag type was expected to be 'time.Time', but it was %s", f.Type())
			}

			if !f.Get().IsZero() {
				t.Errorf("The flag value was expected to be empty, but it was %v", f.Get())
			}

			if f.Var() == nil {
				t.Error("The initial flag variable should not be nil")
			}
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
			if !test.ErrorContains(err, tc.expectedError) {
				t.Errorf("Expected to receive an error with '%s', but received %s", tc.expectedError, err)
			}
			actual := f.Get()
			if !actual.Equal(tc.expectedValue) {
				t.Errorf("Expected value: %v, Actual: %v", tc.expectedValue, actual)
			}

			if !tc.expectedValue.Equal(*fVar) {
				t.Errorf("Expected flag variable: %v, Actual: %v", tc.expectedValue, *fVar)
			}
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
			if !test.ErrorContains(err, tc.expectedError) {
				t.Errorf("Expected to receive an error with '%s', but received %s", tc.expectedError, err)
			}
			actual := f.Get()
			if !actual.Equal(tc.expectedValue) {
				t.Errorf("Expected value: %v, Actual: %v", tc.expectedValue, actual)
			}

			if !tc.expectedValue.Equal(*fVar) {
				t.Errorf("Expected flag variable: %v, Actual: %v", tc.expectedValue, *fVar)
			}

			f.ResetToDefault()

			if tc.setDefault && f.IsSet() {
				t.Error("IsSet() Expected: false, Actual: true")
			}

			actual = f.Get()
			if !actual.Equal(tc.expectedAfterResetValue) {
				t.Errorf("Expected value after reset: %v, Actual: %v", tc.expectedAfterResetValue, actual)
			}

			if !tc.expectedAfterResetValue.Equal(*fVar) {
				t.Errorf("Expected flag variable after reset: %v, Actual: %v", tc.expectedAfterResetValue, *fVar)
			}
		})
	}
}
