package core_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/xitonix/flags"
)

func TestIPAddressSlice(t *testing.T) {
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
			f := flags.IPAddressSlice(tc.long, tc.usage)
			checkFlagInitialState(t, f, "[]ip", tc.expectedUsage, tc.expectedLong, "")
			checkIPSliceFlagValues(t, []net.IP{}, f.Get(), f.Var())
		})
	}
}

func TestIPAddressSliceFlag_WithShort(t *testing.T) {
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
			f := flags.IPAddressSlice(tc.long, tc.usage).WithShort(tc.short)
			checkFlagInitialState(t, f, "[]ip", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkIPSliceFlagValues(t, []net.IP{}, f.Get(), f.Var())
		})
	}
}

func TestIPAddressSliceFlag_WithKey(t *testing.T) {
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
			f := flags.IPAddressSlice("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestIPAddressSliceFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         []net.IP
		expectedDefaultValue []net.IP
	}{
		{
			title:                "empty default value",
			defaultValue:         []net.IP{},
			expectedDefaultValue: []net.IP{},
		},
		{
			title:                "non empty IPv4 default value",
			defaultValue:         []net.IP{net.ParseIP("192.168.1.1")},
			expectedDefaultValue: []net.IP{net.ParseIP("192.168.1.1")},
		},
		{
			title:                "non empty IPv6 default value",
			defaultValue:         []net.IP{net.ParseIP("2001:db8::68")},
			expectedDefaultValue: []net.IP{net.ParseIP("2001:db8::68")},
		},
		{
			title:                "non empty mixed default values",
			defaultValue:         []net.IP{net.ParseIP("2001:db8::68"), net.ParseIP("192.168.1.1")},
			expectedDefaultValue: []net.IP{net.ParseIP("2001:db8::68"), net.ParseIP("192.168.1.1")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IPAddressSlice("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default().([]net.IP)
			for i, a := range actual {
				if !tc.expectedDefaultValue[i].Equal(a) {
					t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
				}
			}
		})
	}
}

func TestIPAddressSliceFlag_Hide(t *testing.T) {
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
			f := flags.IPAddressSlice("long", "usage")
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

func TestIPAddressSliceFlag_IsDeprecated(t *testing.T) {
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
			f := flags.IPAddressSlice("long", "usage")
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

func TestIPAddressSliceFlag_IsRequired(t *testing.T) {
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
			f := flags.IPAddressSlice("long", "usage")
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

func TestIPAddressSliceFlag_WithDelimiter(t *testing.T) {
	testCases := []struct {
		title         string
		value         string
		delimiter     string
		expectedValue []net.IP
	}{
		{
			title:         "IPv4 empty delimiter",
			value:         "192.168.1.1,192.168.1.2",
			expectedValue: []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("192.168.1.2")},
		},
		{
			title:         "IPv4 white space delimiter with white spaced input",
			value:         "192.168.1.1 192.168.1.2",
			delimiter:     " ",
			expectedValue: []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("192.168.1.2")},
		},
		{
			title:         "IPv4 none white space delimiter",
			value:         "192.168.1.1|192.168.1.2",
			delimiter:     "|",
			expectedValue: []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("192.168.1.2")},
		},
		{
			title:         "IPv4 no delimited input",
			value:         "192.168.1.1",
			delimiter:     "|",
			expectedValue: []net.IP{net.ParseIP("192.168.1.1")},
		},

		{
			title:         "IPv6 empty delimiter",
			value:         "2001:db8::68,2002:ab8::69",
			expectedValue: []net.IP{net.ParseIP("2001:db8::68"), net.ParseIP("2002:ab8::69")},
		},
		{
			title:         "IPv6 white space delimiter with white spaced input",
			value:         "2001:db8::68 2002:ab8::69",
			delimiter:     " ",
			expectedValue: []net.IP{net.ParseIP("2001:db8::68"), net.ParseIP("2002:ab8::69")},
		},
		{
			title:         "IPv6 none white space delimiter",
			value:         "2001:db8::68|2002:ab8::69",
			delimiter:     "|",
			expectedValue: []net.IP{net.ParseIP("2001:db8::68"), net.ParseIP("2002:ab8::69")},
		},
		{
			title:         "IPv6 no delimited input",
			value:         "2001:db8::68",
			delimiter:     "|",
			expectedValue: []net.IP{net.ParseIP("2001:db8::68")},
		},

		{
			title:         "mixed versions with empty delimiter",
			value:         "192.168.1.1,2002:ab8::69",
			expectedValue: []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("2002:ab8::69")},
		},
		{
			title:         "mixed versions white space delimiter with white spaced input",
			value:         "192.168.1.1 2002:ab8::69",
			delimiter:     " ",
			expectedValue: []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("2002:ab8::69")},
		},
		{
			title:         "mixed versions none white space delimiter",
			value:         "192.168.1.1|2002:ab8::69",
			delimiter:     "|",
			expectedValue: []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("2002:ab8::69")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IPAddressSlice("long", "usage").WithDelimiter(tc.delimiter)
			fVar := f.Var()
			err := f.Set(tc.value)
			actual := f.Get()
			checkIPSliceFlag(t, f, err, "", tc.expectedValue, actual, fVar)
			for _, act := range actual {
				if act == nil {
					t.Error("Did not expect a nil value")
				}
			}
		})
	}
}

func TestIPAddressSliceFlag_Set(t *testing.T) {
	const (
		ipV4AddressOne = "192.168.1.1"
		ipV4AddressTwo = "192.168.1.2"
		ipV6AddressOne = "2001:db8::68"
		ipV6AddressTwo = "2002:ab8::69"
	)
	var (
		ipV4One = net.ParseIP(ipV4AddressOne)
		ipV4Two = net.ParseIP(ipV4AddressTwo)

		ipV6One = net.ParseIP(ipV6AddressOne)
		ipV6Two = net.ParseIP(ipV6AddressTwo)
	)
	empty := make([]net.IP, 0)
	testCases := []struct {
		title         string
		value         string
		expectedValue []net.IP
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
			title:         "single IPv4 value with white space",
			value:         " " + ipV4AddressOne + " ",
			expectedValue: []net.IP{ipV4One},
		},
		{
			title:         "single IPv4 value with no white space",
			value:         ipV4AddressOne,
			expectedValue: []net.IP{ipV4One},
		},
		{
			title:         "single IPv6 value with white space",
			value:         " " + ipV6AddressOne + " ",
			expectedValue: []net.IP{ipV6One},
		},
		{
			title:         "single IPv6 value with no white space",
			value:         ipV6AddressOne,
			expectedValue: []net.IP{ipV6One},
		},
		{
			title:         "comma separated IPv4 values with no white space",
			value:         ipV4AddressOne + "," + ipV4AddressTwo,
			expectedValue: []net.IP{ipV4One, ipV4Two},
		},
		{
			title:         "comma separated IPv6 values with no white space",
			value:         ipV6AddressOne + "," + ipV6AddressTwo,
			expectedValue: []net.IP{ipV6One, ipV6Two},
		},
		{
			title:         "comma separated mixed values with no white space",
			value:         ipV4AddressOne + "," + ipV6AddressOne,
			expectedValue: []net.IP{ipV4One, ipV6One},
		},

		{
			title:         "comma separated IPv4 values with white space",
			value:         "  " + ipV4AddressOne + "   ,   " + ipV4AddressTwo + "   ",
			expectedValue: []net.IP{ipV4One, ipV4Two},
		},
		{
			title:         "comma separated IPv6 values with white space",
			value:         "  " + ipV6AddressOne + "   ,   " + ipV6AddressTwo + "   ",
			expectedValue: []net.IP{ipV6One, ipV6Two},
		},
		{
			title:         "comma separated mixed values with white space",
			value:         "   " + ipV4AddressOne + "   ,   " + ipV6AddressOne + "   ",
			expectedValue: []net.IP{ipV4One, ipV6One},
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
			expectedError: "'invalid' is not a valid []ip value for --long",
			expectedValue: empty,
		},
		{
			title:         "partially invalid value",
			value:         ipV4AddressOne + ",invalid," + ipV6AddressOne,
			expectedError: "'invalid' is not a valid []ip value for --long",
			expectedValue: empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IPAddressSlice("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkIPSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestIPAddressSliceFlag_Validation(t *testing.T) {
	const (
		ipV4AddressOne          = "192.168.1.1"
		ipV4AddressTwo          = "192.168.1.2"
		ipV4AddressThree        = "192.168.1.3"
		unacceptableIPv4Address = "8.8.8.8"
		ipV6AddressOne          = "2001:db8::68"
		ipV6AddressTwo          = "2002:ab8::69"
		ipV6AddressThree        = "2002:ab8::70"
		unacceptableIPv6Address = "2002:aa8::80"
	)
	var (
		ipV4One   = net.ParseIP(ipV4AddressOne)
		ipV4Two   = net.ParseIP(ipV4AddressTwo)
		ipV4Three = net.ParseIP(ipV4AddressThree)

		ipV6One   = net.ParseIP(ipV6AddressOne)
		ipV6Two   = net.ParseIP(ipV6AddressTwo)
		ipV6Three = net.ParseIP(ipV6AddressThree)
	)
	empty := make([]net.IP, 0)
	testCases := []struct {
		title             string
		value             string
		expectedValue     []net.IP
		validationCB      func(in net.IP) error
		setValidationCB   bool
		validationList    []net.IP
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "IPv4 with nil validation callback",
			setValidationCB: true,
			value:           fmt.Sprintf("%s,%s", ipV4AddressOne, ipV4AddressTwo),
			expectedValue:   []net.IP{ipV4One, ipV4Two},
			expectedError:   "",
		},
		{
			title:           "IPv6 with nil validation callback",
			setValidationCB: true,
			value:           fmt.Sprintf("%s,%s", ipV6AddressOne, ipV6AddressTwo),
			expectedValue:   []net.IP{ipV6One, ipV6Two},
			expectedError:   "",
		},
		{
			title:           "mixed versions with nil validation callback",
			setValidationCB: true,
			value:           fmt.Sprintf("%s,%s", ipV6AddressOne, ipV6AddressTwo),
			expectedValue:   []net.IP{ipV6One, ipV6Two},
			expectedError:   "",
		},
		{
			title:             "IPv4 with nil validation list",
			setValidationList: true,
			value:             fmt.Sprintf("%s,%s", ipV4AddressOne, ipV4AddressTwo),
			expectedValue:     []net.IP{ipV4One, ipV4Two},
			expectedError:     "",
		},
		{
			title:             "IPv6 with nil validation list",
			setValidationList: true,
			value:             fmt.Sprintf("%s,%s", ipV6AddressOne, ipV6AddressTwo),
			expectedValue:     []net.IP{ipV6One, ipV6Two},
			expectedError:     "",
		},
		{
			title:             "mixed versions with nil validation list",
			setValidationList: true,
			value:             fmt.Sprintf("%s,%s", ipV4AddressOne, ipV6AddressTwo),
			expectedValue:     []net.IP{ipV4One, ipV6Two},
			expectedError:     "",
		},
		{
			title:             "IPv4 nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             fmt.Sprintf("%s,%s", ipV4AddressOne, ipV4AddressTwo),
			expectedValue:     []net.IP{ipV4One, ipV4Two},
			expectedError:     "",
		},
		{
			title:             "IPv6 nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             fmt.Sprintf("%s,%s", ipV6AddressOne, ipV6AddressTwo),
			expectedValue:     []net.IP{ipV6One, ipV6Two},
			expectedError:     "",
		},
		{
			title:             "mixed values with nil validation list and callback",
			setValidationList: true,
			setValidationCB:   true,
			value:             fmt.Sprintf("%s,%s", ipV4AddressOne, ipV6AddressTwo),
			expectedValue:     []net.IP{ipV4One, ipV6Two},
			expectedError:     "",
		},
		{
			title:             "empty validation list",
			validationList:    make([]net.IP, 0),
			setValidationList: true,
			value:             fmt.Sprintf("%s,%s", ipV4AddressOne, ipV6AddressTwo),
			expectedValue:     []net.IP{ipV4One, ipV6Two},
			expectedError:     "",
		},
		{
			title:             "unacceptable IPv4 input with single validation entry",
			validationList:    []net.IP{ipV4One},
			setValidationList: true,
			value:             ipV4AddressTwo,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --ip-addresses. The expected value is %s.", ipV4AddressTwo, ipV4AddressOne),
			expectedValue:     empty,
		},
		{
			title:             "invalid item in the validation list",
			validationList:    []net.IP{{}},
			setValidationList: true,
			value:             ipV4AddressOne,
			expectedValue:     []net.IP{ipV4One},
			expectedError:     "",
		},
		{
			title:             "duplicate IPv4 in the validation list",
			validationList:    []net.IP{ipV4One, ipV4One},
			setValidationList: true,
			value:             ipV4AddressTwo,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --ip-addresses. The expected value is %s.", ipV4AddressTwo, ipV4AddressOne),
			expectedValue:     empty,
		},
		{
			title:             "duplicate IPv6 in the validation list",
			validationList:    []net.IP{ipV6One, ipV6One},
			setValidationList: true,
			value:             ipV6AddressTwo,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --ip-addresses. The expected value is %s.", ipV6AddressTwo, ipV6AddressOne),
			expectedValue:     empty,
		},
		{
			title:             "unacceptable IPv6 input with single validation entry",
			validationList:    []net.IP{ipV6One},
			setValidationList: true,
			value:             ipV6AddressTwo,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --ip-addresses. The expected value is %s.", ipV6AddressTwo, ipV6AddressOne),
			expectedValue:     empty,
		},
		{
			title:             "unacceptable IPv4 input with two validation entries",
			validationList:    []net.IP{ipV4One, ipV4Two},
			setValidationList: true,
			value:             unacceptableIPv4Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --ip-addresses. The expected values are %s,%s.", unacceptableIPv4Address, ipV4AddressOne, ipV4AddressTwo),
			expectedValue:     empty,
		},
		{
			title:             "unacceptable IPv6 input with two validation entries",
			validationList:    []net.IP{ipV6One, ipV6Two},
			setValidationList: true,
			value:             unacceptableIPv6Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --ip-addresses. The expected values are %s,%s.", unacceptableIPv6Address, ipV6AddressOne, ipV6AddressTwo),
			expectedValue:     empty,
		},
		{
			title:             "unacceptable IPv4 input with three validation entries",
			validationList:    []net.IP{ipV4One, ipV4Two, ipV4Three},
			setValidationList: true,
			value:             unacceptableIPv4Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --ip-addresses. The expected values are %s,%s,%s.", unacceptableIPv4Address, ipV4AddressOne, ipV4AddressTwo, ipV4AddressThree),
			expectedValue:     empty,
		},
		{
			title:             "unacceptable IPv6 input with three validation entries",
			validationList:    []net.IP{ipV6One, ipV6Two, ipV6Three},
			setValidationList: true,
			value:             unacceptableIPv6Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --ip-addresses. The expected values are %s,%s,%s.", unacceptableIPv6Address, ipV6AddressOne, ipV6AddressTwo, ipV6AddressThree),
			expectedValue:     empty,
		},
		{
			title:             "empty value",
			validationList:    []net.IP{ipV4One, ipV4Two},
			setValidationList: true,
			value:             "",
			expectedError:     "",
			expectedValue:     empty,
		},
		{
			title:             "white space value",
			validationList:    []net.IP{ipV4One, ipV4Two},
			setValidationList: true,
			value:             "  ",
			expectedValue:     empty,
		},
		{
			title: "IPv4 validation callback with no validation error",
			validationCB: func(in net.IP) error {
				return nil
			},
			setValidationCB: true,
			value:           ipV4AddressOne,
			expectedValue:   []net.IP{ipV4One},
		},
		{
			title: "IPv6 validation callback with no validation error",
			validationCB: func(in net.IP) error {
				return nil
			},
			setValidationCB: true,
			value:           ipV6AddressOne,
			expectedValue:   []net.IP{ipV6One},
		},
		{
			title: "IPv4 validation callback with validation error",
			validationCB: func(in net.IP) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           ipV6AddressOne,
			expectedError:   "validation callback failed",
			expectedValue:   empty,
		},
		{
			title: "IPv6 validation callback with validation error",
			validationCB: func(in net.IP) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           ipV6AddressOne,
			expectedError:   "validation callback failed",
			expectedValue:   empty,
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in net.IP) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []net.IP{ipV4One, ipV6One},
			setValidationList: true,
			value:             unacceptableIPv4Address,
			expectedError:     "validation callback failed",
			expectedValue:     empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IPAddressSlice("ip-addresses", "usage")
			fVar := f.Var()
			if tc.setValidationCB {
				f = f.WithValidationCallback(tc.validationCB)
			}
			if tc.setValidationList {
				f = f.WithValidRange(tc.validationList...)
			}
			err := f.Set(tc.value)
			checkIPSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestIPAddressSliceFlag_ResetToDefault(t *testing.T) {
	const (
		ipV4AddressOne = "192.168.1.1"
		ipV4AddressTwo = "192.168.1.2"

		ipV6AddressOne = "2001:db8::68"
		ipV6AddressTwo = "2002:ab8::69"
	)
	var (
		ipV4One = net.ParseIP(ipV4AddressOne)
		ipV4Two = net.ParseIP(ipV4AddressTwo)

		ipV6One = net.ParseIP(ipV6AddressOne)
		ipV6Two = net.ParseIP(ipV6AddressTwo)
	)
	empty := make([]net.IP, 0)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           []net.IP
		defaultValue            []net.IP
		expectedAfterResetValue []net.IP
		expectedError           string
		setDefault              bool
		expectedIsSetAfterReset bool
	}{
		{
			title:                   "IPv4 reset without defining the default value",
			value:                   ipV4AddressOne,
			expectedValue:           []net.IP{ipV4One},
			expectedAfterResetValue: []net.IP{ipV4One},
			setDefault:              false,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "IPv6 reset without defining the default value",
			value:                   ipV6AddressOne,
			expectedValue:           []net.IP{ipV6One},
			expectedAfterResetValue: []net.IP{ipV6One},
			setDefault:              false,
			expectedIsSetAfterReset: true,
		},
		{
			title:                   "IPv4 reset to empty default value",
			value:                   ipV4AddressOne,
			expectedValue:           []net.IP{ipV4One},
			defaultValue:            empty,
			expectedAfterResetValue: empty,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "IPv6 reset to empty default value",
			value:                   ipV6AddressOne,
			expectedValue:           []net.IP{ipV6One},
			defaultValue:            empty,
			expectedAfterResetValue: empty,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "IPv4 reset to nil default value",
			value:                   ipV4AddressOne,
			expectedValue:           []net.IP{ipV4One},
			defaultValue:            nil,
			expectedAfterResetValue: nil,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "IPv6 reset to nil default value",
			value:                   ipV6AddressOne,
			expectedValue:           []net.IP{ipV6One},
			defaultValue:            nil,
			expectedAfterResetValue: nil,
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "IPv4 reset to non-empty default value",
			value:                   ipV4AddressOne,
			expectedValue:           []net.IP{ipV4One},
			defaultValue:            []net.IP{ipV4One, ipV4Two},
			expectedAfterResetValue: []net.IP{ipV4One, ipV4Two},
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
		{
			title:                   "IPv6 reset to non-empty default value",
			value:                   ipV6AddressOne,
			expectedValue:           []net.IP{ipV6One},
			defaultValue:            []net.IP{ipV6One, ipV6Two},
			expectedAfterResetValue: []net.IP{ipV6One, ipV6Two},
			setDefault:              true,
			expectedIsSetAfterReset: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IPAddressSlice("long", "usage")
			if tc.setDefault {
				f = f.WithDefault(tc.defaultValue)
			}
			fVar := f.Var()
			err := f.Set(tc.value)
			checkIPSliceFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)

			f.ResetToDefault()

			if f.IsSet() != tc.expectedIsSetAfterReset {
				t.Errorf("IsSet() Expected: %v, Actual: %v", tc.expectedIsSetAfterReset, f.IsSet())
			}

			checkIPSliceFlagValues(t, tc.expectedAfterResetValue, f.Get(), fVar)
		})
	}
}
