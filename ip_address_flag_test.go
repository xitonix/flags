package flags_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"go.xitonix.io/flags"
)

func TestIPAddress(t *testing.T) {
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
			f := flags.IPAddress(tc.long, tc.usage)
			checkFlagInitialState(t, f, "IP", tc.expectedUsage, tc.expectedLong, "")
			checkIPFlagValues(t, nil, f.Get(), f.Var())
		})
	}
}

func TestIPAddressP(t *testing.T) {
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
			f := flags.IPAddressP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "IP", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkIPFlagValues(t, nil, f.Get(), f.Var())
		})
	}
}

func TestIPAddressFlag_WithKey(t *testing.T) {
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
			f := flags.IPAddress("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestIPAddressFlag_WithDefault(t *testing.T) {
	testCases := []struct {
		title                string
		defaultValue         net.IP
		expectedDefaultValue net.IP
	}{
		{
			title:                "nil default value",
			defaultValue:         nil,
			expectedDefaultValue: nil,
		},
		{
			title:                "zero IPv4 default value",
			defaultValue:         net.IPv4zero,
			expectedDefaultValue: net.IPv4zero,
		},
		{
			title:                "zero IPv6 default value",
			defaultValue:         net.IPv6zero,
			expectedDefaultValue: net.IPv6zero,
		},
		{
			title:                "non zero IPv4 default value",
			defaultValue:         net.ParseIP("192.168.1.1"),
			expectedDefaultValue: net.ParseIP("192.168.1.1"),
		},
		{
			title:                "non zero IPv6 default value",
			defaultValue:         net.ParseIP("2001:db8:85a3::8a2e:370:7334"),
			expectedDefaultValue: net.ParseIP("2001:db8:85a3::8a2e:370:7334"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IPAddress("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if !tc.expectedDefaultValue.Equal(actual.(net.IP)) {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestIPAddressFlag_Hide(t *testing.T) {
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
			f := flags.IPAddress("long", "usage")
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

func TestIPAddressFlag_IsDeprecated(t *testing.T) {
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
			f := flags.IPAddress("long", "usage")
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

func TestIPAddressFlag_Set(t *testing.T) {
	const (
		ipV4Address = "192.168.1.1"
		ipV6Address = "2001:db8:85a3::8a2e:370:7334"
	)
	var (
		ipV4 = net.ParseIP(ipV4Address)
		ipV6 = net.ParseIP(ipV6Address)
	)
	testCases := []struct {
		title         string
		value         string
		expectedValue net.IP
		expectedError string
	}{
		{
			title:         "no value",
			expectedValue: nil,
		},
		{
			title:         "white space value",
			value:         "   ",
			expectedValue: nil,
		},
		{
			title:         "IPv4 value with white space",
			value:         "  " + ipV4Address + "  ",
			expectedValue: ipV4,
		},
		{
			title:         "IPv6 value with white space",
			value:         "  " + ipV6Address + "  ",
			expectedValue: ipV6,
		},
		{
			title:         "IPv4 value without white space",
			value:         ipV4Address,
			expectedValue: ipV4,
		},
		{
			title:         "IPv6 value without white space",
			value:         ipV6Address,
			expectedValue: ipV6,
		},
		{
			title:         "invalid value",
			value:         "abc",
			expectedError: "is not a valid IP value",
			expectedValue: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IPAddress("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkIPFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestIPAddressFlag_Validation(t *testing.T) {
	const (
		ipV4Address           = "192.168.1.1"
		ipV6Address           = "2001:db8:85a3::8a2e:370:7334"
		unacceptableV4Address = "10.10.8.8"
		unacceptableV6Address = "2001:db8::68"
	)
	var (
		ipV4 = net.ParseIP(ipV4Address)
		ipV6 = net.ParseIP(ipV6Address)
	)

	testCases := []struct {
		title             string
		value             string
		expectedValue     net.IP
		validationCB      func(in net.IP) error
		setValidationCB   bool
		validationList    []net.IP
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with IPv4 input",
			setValidationCB: true,
			value:           ipV4Address,
			expectedValue:   ipV4,
			expectedError:   "",
		},
		{
			title:           "nil validation callback with IPv6 input",
			setValidationCB: true,
			value:           ipV6Address,
			expectedValue:   ipV6,
			expectedError:   "",
		},
		{
			title:             "nil validation list with IPv4 input",
			setValidationList: true,
			value:             ipV4Address,
			expectedValue:     ipV4,
			expectedError:     "",
		},
		{
			title:             "nil validation list with IPv6 input",
			setValidationList: true,
			value:             ipV6Address,
			expectedValue:     ipV6,
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with IPv4 input",
			setValidationList: true,
			setValidationCB:   true,
			value:             ipV4Address,
			expectedValue:     ipV4,
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with IPv6 input",
			setValidationList: true,
			setValidationCB:   true,
			value:             ipV6Address,
			expectedValue:     ipV6,
			expectedError:     "",
		},
		{
			title:             "empty validation list with IPv4 input",
			validationList:    make([]net.IP, 0),
			setValidationList: true,
			value:             ipV4Address,
			expectedValue:     ipV4,
			expectedError:     "",
		},
		{
			title:             "empty validation list with IPv6 input",
			validationList:    make([]net.IP, 0),
			setValidationList: true,
			value:             ipV6Address,
			expectedValue:     ipV6,
			expectedError:     "",
		},
		{
			title:             "IPv4 single item in the validation list",
			validationList:    []net.IP{ipV4},
			setValidationList: true,
			value:             unacceptableV4Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected value is %s.", unacceptableV4Address, ipV4Address),
			expectedValue:     nil,
		},
		{
			title:             "IPv6 single item in the validation list",
			validationList:    []net.IP{ipV6},
			setValidationList: true,
			value:             unacceptableV6Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected value is %s.", unacceptableV6Address, ipV6Address),
			expectedValue:     nil,
		},
		{
			title:             "duplicate IPv6 items in the validation list",
			validationList:    []net.IP{ipV6, ipV6},
			setValidationList: true,
			value:             unacceptableV6Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected value is %s.", unacceptableV6Address, ipV6Address),
			expectedValue:     nil,
		},
		{
			title:             "duplicate IPv4 items in the validation list",
			validationList:    []net.IP{ipV4, ipV4},
			setValidationList: true,
			value:             unacceptableV4Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected value is %s.", unacceptableV4Address, ipV4Address),
			expectedValue:     nil,
		},
		{
			title:             "two IPv4 items in the validation list",
			validationList:    []net.IP{ipV4, net.ParseIP("192.168.1.2")},
			setValidationList: true,
			value:             unacceptableV4Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected values are %s,192.168.1.2.", unacceptableV4Address, ipV4Address),
			expectedValue:     nil,
		},
		{
			title:             "two IPv6 items in the validation list",
			validationList:    []net.IP{ipV6, net.ParseIP("fe80:3::1ff:fe23:4567:890a")},
			setValidationList: true,
			value:             unacceptableV6Address,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected values are %s,fe80:3::1ff:fe23:4567:890a.", unacceptableV6Address, ipV6Address),
			expectedValue:     nil,
		},
		{
			title:             "invalid item in the validation list",
			validationList:    []net.IP{{}},
			setValidationList: true,
			value:             ipV4Address,
			expectedValue:     ipV4,
			expectedError:     "",
		},
		{
			title:             "acceptable IPv4 input with validation list",
			validationList:    []net.IP{ipV4},
			setValidationList: true,
			value:             ipV4Address,
			expectedValue:     ipV4,
			expectedError:     "",
		},
		{
			title:             "acceptable IPv6 input with validation list",
			validationList:    []net.IP{ipV6},
			setValidationList: true,
			value:             ipV6Address,
			expectedValue:     ipV6,
			expectedError:     "",
		},
		{
			title:             "acceptable IPv4 input with mix validation list",
			validationList:    []net.IP{ipV4, ipV6},
			setValidationList: true,
			value:             ipV4Address,
			expectedValue:     ipV4,
			expectedError:     "",
		},
		{
			title:             "acceptable IPv6 input with mix validation list",
			validationList:    []net.IP{ipV4, ipV6},
			setValidationList: true,
			value:             ipV6Address,
			expectedValue:     ipV6,
			expectedError:     "",
		},
		{
			title: "IPv4 input with validation callback with no validation error",
			validationCB: func(in net.IP) error {
				return nil
			},
			setValidationCB: true,
			value:           ipV4Address,
			expectedValue:   ipV4,
		},
		{
			title: "IPv6 input with validation callback with no validation error",
			validationCB: func(in net.IP) error {
				return nil
			},
			setValidationCB: true,
			value:           ipV6Address,
			expectedValue:   ipV6,
		},
		{
			title: "IPv4 input with failing validation callback",
			validationCB: func(in net.IP) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           ipV4Address,
			expectedError:   "validation callback failed",
			expectedValue:   nil,
		},
		{
			title: "IPv6 input with failing validation callback",
			validationCB: func(in net.IP) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           ipV6Address,
			expectedError:   "validation callback failed",
			expectedValue:   nil,
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in net.IP) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []net.IP{ipV4, ipV6},
			setValidationList: true,
			value:             ipV4Address,
			expectedError:     "validation callback failed",
			expectedValue:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IPAddress("long", "usage")
			fVar := f.Var()
			if tc.setValidationCB {
				f = f.WithValidationCallback(tc.validationCB)
			}
			if tc.setValidationList {
				f = f.WithValidRange(tc.validationList...)
			}
			err := f.Set(tc.value)
			checkIPFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestIPAddressFlag_ResetToDefault(t *testing.T) {
	const (
		ipV4Address        = "192.168.1.1"
		ipV4DefaultAddress = "192.168.1.2"
		ipV6Address        = "2001:db8:85a3::8a2e:370:7334"
		ipV6DefaultAddress = "2001:db8:85a3::8a2e:260:6224"
	)
	var (
		ipV4        = net.ParseIP(ipV4Address)
		ipV4Default = net.ParseIP(ipV4DefaultAddress)
		ipV6        = net.ParseIP(ipV6Address)
		ipV6Default = net.ParseIP(ipV6DefaultAddress)
	)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           net.IP
		defaultValue            net.IP
		expectedAfterResetValue net.IP
		expectedError           string
		setDefault              bool
	}{
		{
			title:                   "no value",
			expectedValue:           nil,
			expectedAfterResetValue: nil,
		},
		{
			title:                   "IPv4 reset without defining the default value",
			value:                   ipV4Address,
			expectedValue:           ipV4,
			expectedAfterResetValue: ipV4,
			setDefault:              false,
		},
		{
			title:                   "IPv6 reset without defining the default value",
			value:                   ipV6Address,
			expectedValue:           ipV6,
			expectedAfterResetValue: ipV6,
			setDefault:              false,
		},
		{
			title:                   "reset to IPv4 zero default value",
			value:                   ipV4Address,
			expectedValue:           ipV4,
			defaultValue:            net.IPv4zero,
			expectedAfterResetValue: net.IPv4zero,
			setDefault:              true,
		},
		{
			title:                   "reset to IPv6 zero default value",
			value:                   ipV6Address,
			expectedValue:           ipV6,
			defaultValue:            net.IPv6zero,
			expectedAfterResetValue: net.IPv6zero,
			setDefault:              true,
		},
		{
			title:                   "reset to non-zero IPv4 default value",
			value:                   ipV4Address,
			expectedValue:           ipV4,
			defaultValue:            ipV4Default,
			expectedAfterResetValue: ipV4Default,
			setDefault:              true,
		},
		{
			title:                   "reset to non-zero IPv6 default value",
			value:                   ipV6Address,
			expectedValue:           ipV6,
			defaultValue:            ipV6Default,
			expectedAfterResetValue: ipV6Default,
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.IPAddress("long", "usage")
			if tc.setDefault {
				f = f.WithDefault(tc.defaultValue)
			}
			fVar := f.Var()
			err := f.Set(tc.value)

			checkIPFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)

			f.ResetToDefault()

			if tc.setDefault && f.IsSet() {
				t.Error("IsSet() Expected: false, Actual: true")
			}

			checkIPFlagValues(t, tc.expectedAfterResetValue, f.Get(), fVar)
		})
	}
}
