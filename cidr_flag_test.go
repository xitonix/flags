package flags_test

import (
	"errors"
	"fmt"
	"testing"

	"go.xitonix.io/flags"
	"go.xitonix.io/flags/core"
)

func TestCIDR(t *testing.T) {
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
			f := flags.CIDR(tc.long, tc.usage)
			checkFlagInitialState(t, f, "cidr", tc.expectedUsage, tc.expectedLong, "")
			checkCIDRFlagValues(t, core.CIDR{}, f.Get(), f.Var())
		})
	}
}

func TestCIDRP(t *testing.T) {
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
			f := flags.CIDRP(tc.long, tc.usage, tc.short)
			checkFlagInitialState(t, f, "cidr", tc.expectedUsage, tc.expectedLong, tc.expectedShort)
			checkCIDRFlagValues(t, core.CIDR{}, f.Get(), f.Var())
		})
	}
}

func TestCIDRFlag_WithKey(t *testing.T) {
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
			f := flags.CIDR("long", "usage").WithKey(tc.key)
			actual := f.Key().String()
			if actual != tc.expectedKey {
				t.Errorf("Expected Key: %s, Actual: %s", tc.expectedKey, actual)
			}
		})
	}
}

func TestCIDRFlag_WithDefault(t *testing.T) {
	c, _ := core.ParseCIDR("192.168.1.1/24")
	testCases := []struct {
		title                string
		defaultValue         core.CIDR
		expectedDefaultValue core.CIDR
	}{
		{
			title:                "empty default value",
			defaultValue:         core.CIDR{},
			expectedDefaultValue: core.CIDR{},
		},
		{
			title:                "none empty default value",
			defaultValue:         *c,
			expectedDefaultValue: *c,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.CIDR("long", "usage").WithDefault(tc.defaultValue)
			actual := f.Default()
			if !tc.expectedDefaultValue.Equal(actual.(core.CIDR)) {
				t.Errorf("Expected Default Value: %v, Actual: %s", tc.expectedDefaultValue, actual)
			}
		})
	}
}

func TestCIDRFlag_Hide(t *testing.T) {
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
			f := flags.CIDR("long", "usage")
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

func TestCIDRFlag_IsDeprecated(t *testing.T) {
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
			f := flags.CIDR("long", "usage")
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

func TestCIDRFlag_IsRequired(t *testing.T) {
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
			f := flags.CIDR("long", "usage")
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

func TestCIDRFlag_Set(t *testing.T) {
	const (
		ipV4Network = "192.168.1.1/24"
		ipV6Network = "2001:db8:85a3::8a2e:370:7334/16"
	)
	var (
		ipV4Cidr, _ = core.ParseCIDR(ipV4Network)
		ipV6Cidr, _ = core.ParseCIDR(ipV6Network)
		empty       = core.CIDR{}
	)
	testCases := []struct {
		title         string
		value         string
		expectedValue core.CIDR
		expectedError string
	}{
		{
			title:         "no value",
			expectedValue: empty,
		},
		{
			title:         "white space value",
			value:         "   ",
			expectedValue: empty,
		},
		{
			title:         "IPv4 value with white space",
			value:         "  " + ipV4Network + "  ",
			expectedValue: *ipV4Cidr,
		},
		{
			title:         "IPv6 value with white space",
			value:         "  " + ipV6Network + "  ",
			expectedValue: *ipV6Cidr,
		},
		{
			title:         "IPv4 value without white space",
			value:         ipV4Network,
			expectedValue: *ipV4Cidr,
		},
		{
			title:         "IPv6 value without white space",
			value:         ipV6Network,
			expectedValue: *ipV6Cidr,
		},
		{
			title:         "invalid value",
			value:         "abc",
			expectedError: "is not a valid cidr value",
			expectedValue: empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.CIDR("long", "usage")
			fVar := f.Var()
			err := f.Set(tc.value)
			checkCIDRFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestCIDRFlag_Validation(t *testing.T) {
	const (
		ipV4Network           = "192.168.1.1/24"
		ipV4Network2          = "192.168.1.2/24"
		ipV6Network           = "2001:db8:85a3::8a2e:370:7334/16"
		ipV6Network2          = "2001:db8:85a3::8a2e:370:8022/16"
		unacceptableV4Network = "10.10.8.8/8"
		unacceptableV6Network = "2001:db8::68/8"
	)
	var (
		ipV4Cidr, _  = core.ParseCIDR(ipV4Network)
		ipV4Cidr2, _ = core.ParseCIDR(ipV4Network2)
		ipV6Cidr, _  = core.ParseCIDR(ipV6Network)
		ipV6Cidr2, _ = core.ParseCIDR(ipV6Network2)
		empty        = core.CIDR{}
	)

	testCases := []struct {
		title             string
		value             string
		expectedValue     core.CIDR
		validationCB      func(in core.CIDR) error
		setValidationCB   bool
		validationList    []core.CIDR
		setValidationList bool
		expectedError     string
	}{
		{
			title:           "nil validation callback with IPv4 input",
			setValidationCB: true,
			value:           ipV4Network,
			expectedValue:   *ipV4Cidr,
			expectedError:   "",
		},
		{
			title:           "nil validation callback with IPv6 input",
			setValidationCB: true,
			value:           ipV6Network,
			expectedValue:   *ipV6Cidr,
			expectedError:   "",
		},
		{
			title:             "nil validation list with IPv4 input",
			setValidationList: true,
			value:             ipV4Network,
			expectedValue:     *ipV4Cidr,
			expectedError:     "",
		},
		{
			title:             "nil validation list with IPv6 input",
			setValidationList: true,
			value:             ipV6Network,
			expectedValue:     *ipV6Cidr,
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with IPv4 input",
			setValidationList: true,
			setValidationCB:   true,
			value:             ipV4Network,
			expectedValue:     *ipV4Cidr,
			expectedError:     "",
		},
		{
			title:             "nil validation list and callback with IPv6 input",
			setValidationList: true,
			setValidationCB:   true,
			value:             ipV6Network,
			expectedValue:     *ipV6Cidr,
			expectedError:     "",
		},
		{
			title:             "empty validation list with IPv4 input",
			validationList:    make([]core.CIDR, 0),
			setValidationList: true,
			value:             ipV4Network,
			expectedValue:     *ipV4Cidr,
			expectedError:     "",
		},
		{
			title:             "empty validation list with IPv6 input",
			validationList:    make([]core.CIDR, 0),
			setValidationList: true,
			value:             ipV6Network,
			expectedValue:     *ipV6Cidr,
			expectedError:     "",
		},
		{
			title:             "invalid item in the validation list",
			validationList:    []core.CIDR{{}},
			setValidationList: true,
			value:             ipV4Network,
			expectedValue:     *ipV4Cidr,
			expectedError:     "",
		},
		{
			title:             "IPv4 single item in the validation list",
			validationList:    []core.CIDR{*ipV4Cidr},
			setValidationList: true,
			value:             unacceptableV4Network,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected value is %s.", unacceptableV4Network, ipV4Network),
			expectedValue:     empty,
		},
		{
			title:             "IPv6 single item in the validation list",
			validationList:    []core.CIDR{*ipV6Cidr},
			setValidationList: true,
			value:             unacceptableV6Network,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected value is %s.", unacceptableV6Network, ipV6Network),
			expectedValue:     empty,
		},
		{
			title:             "two IPv4 items in the validation list",
			validationList:    []core.CIDR{*ipV4Cidr, *ipV4Cidr2},
			setValidationList: true,
			value:             unacceptableV4Network,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected values are %s,%s", unacceptableV4Network, ipV4Network, ipV4Cidr2),
			expectedValue:     empty,
		},
		{
			title:             "non unique IPv4 items in the validation list",
			validationList:    []core.CIDR{*ipV4Cidr, *ipV4Cidr},
			setValidationList: true,
			value:             unacceptableV4Network,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected value is %s.", unacceptableV4Network, ipV4Network),
			expectedValue:     empty,
		},
		{
			title:             "non unique IPv6 items in the validation list",
			validationList:    []core.CIDR{*ipV6Cidr, *ipV6Cidr},
			setValidationList: true,
			value:             unacceptableV6Network,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected value is %s.", unacceptableV6Network, ipV6Network),
			expectedValue:     empty,
		},
		{
			title:             "two IPv6 items in the validation list",
			validationList:    []core.CIDR{*ipV6Cidr, *ipV6Cidr2},
			setValidationList: true,
			value:             unacceptableV6Network,
			expectedError:     fmt.Sprintf("%s is not an acceptable value for --long. The expected values are %s,%s.", unacceptableV6Network, ipV6Network, ipV6Cidr2),
			expectedValue:     empty,
		},
		{
			title:             "acceptable IPv4 input with validation list",
			validationList:    []core.CIDR{*ipV4Cidr},
			setValidationList: true,
			value:             ipV4Network,
			expectedValue:     *ipV4Cidr,
			expectedError:     "",
		},
		{
			title:             "acceptable IPv6 input with validation list",
			validationList:    []core.CIDR{*ipV6Cidr},
			setValidationList: true,
			value:             ipV6Network,
			expectedValue:     *ipV6Cidr,
			expectedError:     "",
		},
		{
			title:             "acceptable IPv4 input with mix validation list",
			validationList:    []core.CIDR{*ipV4Cidr, *ipV6Cidr},
			setValidationList: true,
			value:             ipV4Network,
			expectedValue:     *ipV4Cidr,
			expectedError:     "",
		},
		{
			title:             "acceptable IPv6 input with mix validation list",
			validationList:    []core.CIDR{*ipV4Cidr, *ipV6Cidr},
			setValidationList: true,
			value:             ipV6Network,
			expectedValue:     *ipV6Cidr,
			expectedError:     "",
		},
		{
			title: "IPv4 input with validation callback with no validation error",
			validationCB: func(in core.CIDR) error {
				return nil
			},
			setValidationCB: true,
			value:           ipV4Network,
			expectedValue:   *ipV4Cidr,
		},
		{
			title: "IPv6 input with validation callback with no validation error",
			validationCB: func(in core.CIDR) error {
				return nil
			},
			setValidationCB: true,
			value:           ipV6Network,
			expectedValue:   *ipV6Cidr,
		},
		{
			title: "IPv4 input with failing validation callback",
			validationCB: func(in core.CIDR) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           ipV4Network,
			expectedError:   "validation callback failed",
			expectedValue:   empty,
		},
		{
			title: "IPv6 input with failing validation callback",
			validationCB: func(in core.CIDR) error {
				return errors.New("validation callback failed")
			},
			setValidationCB: true,
			value:           ipV6Network,
			expectedError:   "validation callback failed",
			expectedValue:   empty,
		},
		{
			title: "validation callback takes priority over validation list",
			validationCB: func(in core.CIDR) error {
				return errors.New("validation callback failed")
			},
			setValidationCB:   true,
			validationList:    []core.CIDR{*ipV4Cidr, *ipV6Cidr},
			setValidationList: true,
			value:             ipV4Network,
			expectedError:     "validation callback failed",
			expectedValue:     empty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.CIDR("long", "usage")
			fVar := f.Var()
			if tc.setValidationCB {
				f = f.WithValidationCallback(tc.validationCB)
			}
			if tc.setValidationList {
				f = f.WithValidRange(tc.validationList...)
			}
			err := f.Set(tc.value)
			checkCIDRFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)
		})
	}
}

func TestCIDRFlag_ResetToDefault(t *testing.T) {
	const (
		ipV4Network        = "192.168.1.1/24"
		ipV4DefaultNetwork = "192.168.1.2/16"
		ipV6Network        = "2001:db8:85a3::8a2e:370:7334/16"
		ipV6DefaultNetwork = "2001:db8:85a3::8a2e:260:6224/24"
	)
	var (
		ipV4, _        = core.ParseCIDR(ipV4Network)
		ipV4Default, _ = core.ParseCIDR(ipV4DefaultNetwork)
		ipV6, _        = core.ParseCIDR(ipV6Network)
		ipV6Default, _ = core.ParseCIDR(ipV6DefaultNetwork)
	)
	testCases := []struct {
		title                   string
		value                   string
		expectedValue           core.CIDR
		defaultValue            core.CIDR
		expectedAfterResetValue core.CIDR
		expectedError           string
		setDefault              bool
	}{
		{
			title:                   "empty value",
			expectedValue:           core.CIDR{},
			expectedAfterResetValue: core.CIDR{},
		},
		{
			title:                   "IPv4 reset without defining the default value",
			value:                   ipV4Network,
			expectedValue:           *ipV4,
			expectedAfterResetValue: *ipV4,
			setDefault:              false,
		},
		{
			title:                   "IPv6 reset without defining the default value",
			value:                   ipV6Network,
			expectedValue:           *ipV6,
			expectedAfterResetValue: *ipV6,
			setDefault:              false,
		},
		{
			title:                   "reset to non-zero IPv4 default value",
			value:                   ipV4Network,
			expectedValue:           *ipV4,
			defaultValue:            *ipV4Default,
			expectedAfterResetValue: *ipV4Default,
			setDefault:              true,
		},
		{
			title:                   "reset to non-zero IPv6 default value",
			value:                   ipV6Network,
			expectedValue:           *ipV6,
			defaultValue:            *ipV6Default,
			expectedAfterResetValue: *ipV6Default,
			setDefault:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := flags.CIDR("long", "usage")
			if tc.setDefault {
				f = f.WithDefault(tc.defaultValue)
			}
			fVar := f.Var()
			err := f.Set(tc.value)

			checkCIDRFlag(t, f, err, tc.expectedError, tc.expectedValue, f.Get(), fVar)

			f.ResetToDefault()

			if tc.setDefault && f.IsSet() {
				t.Error("IsSet() Expected: false, Actual: true")
			}

			checkCIDRFlagValues(t, tc.expectedAfterResetValue, f.Get(), fVar)
		})
	}
}
